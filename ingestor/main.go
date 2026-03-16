package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// token const token = "jWD0oGjXAyMmNbfBchMF23CZ-UDozjjsSvYYV_5xJBA"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mastodonToken := os.Getenv("MASTODON_TOKEN")

	mode := flag.String("mode", "api", "Connector mode: api or stream")
	flag.Parse()

	if *mode == "stream" {
		streamConnector(ctx, mastodonToken)
	} else if *mode == "api" {
		apiConnector(ctx, mastodonToken)
	} else {
		log.Fatal("Please select one of the following connector options: stream, api")
	}

}

func streamConnector(ctx context.Context, token string) {

	//Init connector
	ch := make(chan string)

	//TODO Read from config
	const mastodonURL = "https://mastodon.social/api/v1/streaming/public"

	connector := MastodonStreamConnector{mastodonURL, token}

	go connector.Stream(ctx, ch)

	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
		case <-ctx.Done():
			fmt.Println("shutting down")
			return
		}
	}
}

func apiConnector(ctx context.Context, token string) {

	ch := make(chan MastodonMessage)

	const mastodonApiURL = "https://mastodon.social/api/v1/timelines/home"

	connector := MastodonApiConnector{mastodonApiURL, token}

	go connector.Poll(ctx, ch)

	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
		case <-ctx.Done():
			fmt.Println("shutting down")
			return
		}
	}

}
