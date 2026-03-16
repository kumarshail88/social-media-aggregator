package main

import (
	"bufio"
	"context"
	"net/http"
)

type MastodonStreamConnector struct {
	url   string
	token string
}

func (m *MastodonStreamConnector) Stream(ctx context.Context, ch chan string) error {
	err := fetchStream(ctx, m.url, m.token, ch)
	if err != nil {
		return err
	}
	return nil
}

func fetchStream(ctx context.Context, url string, token string, ch chan string) error {

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		ch <- scanner.Text()
	}

	return nil
}
