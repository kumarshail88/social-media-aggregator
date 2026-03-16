package main

import "context"

type ApiConnector interface {
	Poll(ctx context.Context, ch chan string) error
}
