package utilitariosgorms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// Verifica com base no tamanho em MB passado se o arquivo é do tamanho correspondente.
func FileSizeFromURLVerifyUsingRange(fileURL string, maxSizeMB int) error {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("URL inválida: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	// Pedimos somente o primeiro byte do arquivo
	req.Header.Set("Range", "bytes=0-0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer request GET Range: %w", err)
	}
	defer resp.Body.Close()

	// Se o servidor não suportar Range, pode devolver 200 OK inteiro,
	// o que não ajuda muito. Então conferimos se obtemos status 206.
	if resp.StatusCode != http.StatusPartialContent {
		return errors.New("o servidor não suportou a requisição parcial (range request)")
	}

	// Exemplo de cabeçalho: "Content-Range: bytes 0-0/1500000"
	contentRange := resp.Header.Get("Content-Range")
	if contentRange == "" {
		return errors.New("não foi possível obter Content-Range no cabeçalho")
	}

	// contentRange deve conter algo como "bytes 0-0/1500000"
	parts := strings.Split(contentRange, "/")
	if len(parts) != 2 {
		return errors.New("formato de Content-Range inesperado")
	}

	totalSizeStr := parts[1] // "1500000"
	totalSize, err := strconv.ParseInt(totalSizeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("erro ao converter tamanho: %w", err)
	}

	// converte para MB
	sizeMB := totalSize / (1024 * 1024)
	if sizeMB > int64(maxSizeMB) {
		return fmt.Errorf("erro: arquivo excede limite de %dMB (tamanho: %dMB)", maxSizeMB, sizeMB)
	}

	// Tudo certo
	return nil
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
