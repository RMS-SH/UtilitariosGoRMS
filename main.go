package utilitariosgorms

import (
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
