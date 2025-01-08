package utilitariosgorms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
func FileSizeFromURLVerify(fileURL string, maxSizeMB int) error {
	// Parseia a URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		fmt.Printf("Erro: URL inválida (%v)\n", err)
		return errors.New("Não foi possível parserar a URL")
	}

	// Faz a requisição HEAD para obter o tamanho do arquivo
	headResp, err := http.Head(parsedURL.String())
	if err != nil {
		fmt.Printf("Erro: Falha ao fazer requisição HEAD (%v)\n", err)
		return errors.New("Não foi possível realizar o request")
	}
	defer headResp.Body.Close()

	// Verifica o status da resposta
	if headResp.StatusCode > 299 || headResp.StatusCode <= 200 {
		return errors.New("Erro: Status da resposta inválido")
	}

	// Obtém o tamanho do arquivo
	size := headResp.ContentLength
	if size <= 0 {
		return errors.New("Erro: Não foi possível determinar o tamanho do arquivo")
	}

	// Converte o tamanho para megabytes
	sizeMB := size / (1024 * 1024)

	// Verifica se o tamanho está dentro do limite
	if sizeMB > int64(maxSizeMB) {
		return errors.New(fmt.Sprintf("Erro: Tamanho do arquivo excede o limite (%d MB > %d MB)\n", sizeMB, maxSizeMB))
	}

	// Tamanho válido
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
