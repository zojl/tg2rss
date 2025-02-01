package fetcher

import (
	"net/http"
	"io"
)

func FetchHTML(url string) (string, error) {
	response, err := FetchStream(url)

	return string(response), err
}

func FetchStream(url string) ([]byte, error) {
	resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return respBody, nil
}