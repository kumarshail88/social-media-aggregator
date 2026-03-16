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

func main() {

	mode := flag.String("mode", "api", "Connector mode: api or stream")
	flag.Parse()

	if *mode == "stream" {
		streamConnector()
	} else if *mode == "api" {
		apiConnector()
	} else {
		log.Fatal("Please select one of the following connector options: stream, api")
	}
}

func streamConnector() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	//Init connector
	ch := make(chan string)

	//TODO Read from config
	const mastodonURL = "https://mastodon.social/api/v1/streaming/public"

	//Must be secured in KV or so
	const token = "jWD0oGjXAyMmNbfBchMF23CZ-UDozjjsSvYYV_5xJBA"

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

func apiConnector() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ch := make(chan MastodonMessage)

	const mastodonApiURL = "https://mastodon.social/api/v1/timelines/home"

	//Must be secured in KV or so
	const token = "jWD0oGjXAyMmNbfBchMF23CZ-UDozjjsSvYYV_5xJBA"

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
