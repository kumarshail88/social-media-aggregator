package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MastodonMessage struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
	URL       string `json:"url"`
}

type MastodonApiConnector struct {
	url   string
	token string
}

func (m *MastodonApiConnector) Poll(ctx context.Context, ch chan MastodonMessage) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	var sinceId string
	for {
		select {
		case <-ticker.C:
			err, lastId := fetchApi(ctx, m.url, m.token, ch, sinceId)
			if err != nil {
				return err
			}
			sinceId = lastId
			fmt.Println("lastId " + lastId)
			fmt.Println("sinceId " + sinceId)
		case <-ctx.Done():
			fmt.Println("shutting down")
			return nil
		}
	}
}

func fetchApi(ctx context.Context, url string, token string, ch chan MastodonMessage, sinceId string) (error error, lastId string) {

	// Prevent duplicate posts
	if sinceId != "" {
		url = url + "?since_id=" + sinceId
		fmt.Println("fetching " + url)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err, ""
	}

	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err, ""
	}

	var messages []MastodonMessage
	json.NewDecoder(response.Body).Decode(&messages)

	lastId = messages[len(messages)-1].ID

	for _, msg := range messages {
		ch <- msg
	}

	return nil, lastId
}
