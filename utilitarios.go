package utilitariosgorms

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// -- //

// Descriçã oda Function
func DownloadFileFromURLToBinary(fileURL string) ([]byte, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return data, nil
}

// -- //

// Retorna a Quantidade em Bytes
func GetFileSizeFromURL(fileURL string) (int64, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return 0, fmt.Errorf("invalid URL: %w", err)
	}

	headResp, err := http.Head(parsedURL.String())
	if err != nil {
		return 0, fmt.Errorf("failed to make HEAD request: %w", err)
	}
	defer headResp.Body.Close()

	if headResp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HEAD request failed: status code %d", headResp.StatusCode)
	}

	size := headResp.ContentLength
	if size <= 0 {
		return 0, fmt.Errorf("unable to determine file size")
	}

	return size, nil
}

//---//

// DownloadResponse armazena o resultado do download.
type DownloadResponse struct {
	Data       []byte // Conteúdo baixado
	RemoteIP   string // IP do servidor remoto
	SizeInMB   int64  // Tamanho do conteúdo em MB
	StatusCode int    // Código de status HTTP, se quiser usar
}

// DownloadWithTimeout baixa todo o conteúdo de "fileURL" com timeout
// de "timeout" (por exemplo, 30s). Retorna:
// - DownloadResponse (com Data, RemoteIP, SizeInMB, StatusCode)
// - error, caso haja falha (timeout, statuscode ruim, tamanho acima do permitido, etc.)
func DownloadWithTimeout(fileURL string, maxSizeMB int, timeout time.Duration) (*DownloadResponse, error) {

	// 1) Faz parse do URL para verificar se é válida.
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, fmt.Errorf("URL inválida: %w", err)
	}

	// 2) Criamos um dialer customizado para capturar o IP remoto.
	var remoteIP string
	dialer := &net.Dialer{}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			// address normalmente vem como "host:porta"
			conn, err := dialer.DialContext(ctx, network, address)
			if err == nil {
				remoteIP = conn.RemoteAddr().String()
				// remoteIP geralmente vem no formato "IP:port", ex.: "192.168.0.10:443"
				// Se quiser só o IP sem porta, pode dar um split:
				if idx := strings.LastIndex(remoteIP, ":"); idx != -1 {
					remoteIP = remoteIP[:idx]
				}
			}
			return conn, err
		},
	}

	// 3) Cria um *client* com o transport e com timeout total.
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout, // Cancela toda a operação (DNS + connect + download)
	}

	// 4) Cria uma request com contexto de 30s (ou "timeout" escolhido)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição GET: %w", err)
	}

	// 5) Executa a requisição
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar arquivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("status de resposta inválido: %d", resp.StatusCode)
	}

	// 6) Lê todo o conteúdo do Body em memória.
	//    CUIDADO: se for muito grande, pode estourar memória.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta HTTP: %w", err)
	}

	// 7) Verifica o tamanho do que foi lido
	sizeInBytes := int64(len(data))
	sizeInMB := sizeInBytes / (1024 * 1024)

	if sizeInMB > int64(maxSizeMB) {
		return nil, errors.New(
			fmt.Sprintf("arquivo excede limite de %dMB (baixados: %dMB)", maxSizeMB, sizeInMB))
	}

	// 8) Monta e retorna a estrutura de resposta
	return &DownloadResponse{
		Data:       data,
		RemoteIP:   remoteIP,
		SizeInMB:   sizeInMB,
		StatusCode: resp.StatusCode,
	}, nil
}

// -- //

// FormatDate converte uma string de data para um dos 10 formatos especificados.
// dateStr: string da data no formato RFC3339 (ex: "2023-04-05T14:30:00Z")
// formatOption: inteiro de 1 a 10 que especifica o formato de saída
func FormatDate(dateStr string, formatOption int) (string, error) {
	// Define os layouts de data
	const inputLayout = time.RFC3339

	// Parse a data de entrada
	parsedTime, err := time.Parse(inputLayout, dateStr)
	if err != nil {
		return "", fmt.Errorf("erro ao parsear a data: %w", err)
	}

	var outputLayout string

	// Seleciona o layout de saída baseado no formatOption
	switch formatOption {
	case 1:
		// Exemplo: "02-01-2006"
		outputLayout = "02-01-2006"
	case 2:
		// Exemplo: "January 2, 2006"
		outputLayout = "January 2, 2006"
	case 3:
		// Exemplo: "02 Jan 06 15:04 MST"
		outputLayout = "02 Jan 06 15:04 MST"
	case 4:
		// Exemplo: "2006/01/02"
		outputLayout = "2006/01/02"
	case 5:
		// Exemplo: "02-01-2006 15:04"
		outputLayout = "02-01-2006 15:04"
	case 6:
		// Exemplo: "Mon, 02 Jan 2006 15:04:05 MST"
		outputLayout = "Mon, 02 Jan 2006 15:04:05 MST"
	case 7:
		// Exemplo: "02-Jan-2006"
		outputLayout = "02-Jan-2006"
	case 8:
		// Exemplo: "02/01/2006"
		outputLayout = "02/01/2006"
	case 9:
		// Exemplo: "2006.01.02"
		outputLayout = "2006.01.02"
	case 10:
		// Exemplo: "02 Jan 2006 03:04 PM"
		outputLayout = "02 Jan 2006 03:04 PM"
	default:
		return "", errors.New("opção de formato inválida. Escolha um número de 1 a 10")
	}

	// Formata a data
	formattedDate := parsedTime.Format(outputLayout)
	return formattedDate, nil
}

// -- //
