package utilitariosgorms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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
