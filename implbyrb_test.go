package main_test

import (
	"context"
	"testing"

	"qbench"
)

func Test_rbImpl(t *testing.T) {
	q := main.NewRingBufferImpl(500)
	ctx := context.Background()

	assertFIFOInSPSC(t, ctx, q)
}

func Test_rbImpl_MPSC(t *testing.T) {
	q := main.NewRingBufferImpl(500)
	ctx := context.Background()

	assertMPSC(t, ctx, q)
}

func Test_rbImpl_SPMC(t *testing.T) {
	q := main.NewRingBufferImpl(500)
	ctx := context.Background()

	assertSPMC(t, ctx, q)
}
