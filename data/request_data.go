package data

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func Give_data(localhost_endpoint string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(localhost_endpoint + "/liveclientdata/allgamedata")
	if err != nil {
		return "", fmt.Errorf("Error accessing endpoint: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading body: %w", err)
	}
	return string(body), nil
}
