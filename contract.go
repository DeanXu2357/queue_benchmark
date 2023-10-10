package main

import "context"

type Queue interface {
	Put(ctx context.Context, s string) error
	Consume(ctx context.Context, f func(jobs []string) error) error
	Close(ctx context.Context) error
}
