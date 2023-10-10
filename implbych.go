package main

import (
	"context"
	"time"
)

/**
 * This is the implementation of the Queue interface using a channel.
 */
type chImpl struct {
	ch chan string
}

func (c *chImpl) Put(ctx context.Context, s string) error {
	c.ch <- s
	return nil
}

func (c *chImpl) Consume(ctx context.Context, f func(jobs []string) error) error {
	if err := f([]string{<-c.ch}); err != nil {
		return err
	}
	return nil
}

func (c *chImpl) Close(ctx context.Context) error {
	time.Sleep(3 * time.Second)
	close(c.ch)
	return nil
}

func NewChImpl(s int64) Queue {
	return &chImpl{ch: make(chan string, s)}
}
