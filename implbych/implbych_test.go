package implbych_test

import (
	"context"
	"testing"

	"qbench/contract"
	"qbench/implbych"
)

func Test_chImpl(t *testing.T) {
	q := implbych.NewChImpl(500)
	ctx := context.Background()

	contract.AssertFIFOInSPSC(t, ctx, q)
}

func Test_chImpl_MPSC(t *testing.T) {
	q := implbych.NewChImpl(500)
	ctx := context.Background()

	contract.AssertMPSC(t, ctx, q)
}

func Test_chImpl_SPMC(t *testing.T) {
	q := implbych.NewChImpl(500)
	ctx := context.Background()

	contract.AssertSPMC(t, ctx, q)
}
