package implbyrb_test

import (
	"context"
	"testing"

	"qbench/contract"
	"qbench/implbyrb"
)

func Test_rbImpl(t *testing.T) {
	q := implbyrb.NewRingBufferImpl(500)
	ctx := context.Background()

	contract.AssertFIFOInSPSC(t, ctx, q)
}

func Test_rbImpl_MPSC(t *testing.T) {
	q := implbyrb.NewRingBufferImpl(500)
	ctx := context.Background()

	contract.AssertMPSC(t, ctx, q)
}

func Test_rbImpl_SPMC(t *testing.T) {
	q := implbyrb.NewRingBufferImpl(500)
	ctx := context.Background()

	contract.AssertSPMC(t, ctx, q)
}
