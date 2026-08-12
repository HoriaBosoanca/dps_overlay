package data

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func GiveData(localhostEndpoint string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(localhostEndpoint + "/liveclientdata/allgamedata")
	if err != nil {
		return "", fmt.Errorf("error accessing endpoint: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}
	return string(body), nil
}
