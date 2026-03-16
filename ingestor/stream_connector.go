package main

import "context"

type StreamConnector interface {
	Stream(ctx context.Context, ch chan string) error
}
