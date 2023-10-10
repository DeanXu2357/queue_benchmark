package main_test

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"qbench"
)

func Test_chImpl(t *testing.T) {
	q := main.NewChImpl(500)
	ctx := context.Background()

	assertFIFOInSPSC(t, ctx, q)
}

func Test_chImpl_MPSC(t *testing.T) {
	q := main.NewChImpl(500)
	ctx := context.Background()

	assertMPSC(t, ctx, q)
}

func Test_chImpl_SPMC(t *testing.T) {
	q := main.NewChImpl(500)
	ctx := context.Background()

	assertSPMC(t, ctx, q)
}

// assertFIFOInSPSC asserts that the queue is FIFO in a single producer, single consumer scenario.
func assertFIFOInSPSC(t *testing.T, ctx context.Context, q main.Queue) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	size := 1000
	go func() {
		defer wg.Done()

		for i := 0; i < size; i++ {
			if err := q.Put(ctx, strconv.Itoa(i)); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	output := make([]string, 0, size)
	m := sync.Mutex{}
	go func() {
		defer wg.Done()

		for i := 0; i < size; i++ {
			if err := q.Consume(ctx, func(jobs []string) error {
				if len(jobs) >= 1 {
					for j := 0; j < len(jobs); j++ {
						m.Lock()

						output = append(output, jobs[j])
						m.Unlock()
					}
				}
				return nil
			}); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	wg.Wait()

	if len(output) != size {
		t.Errorf("Expected %d, got %d", size, len(output))
		return
	}

	prev := 0
	for i := range output {
		if prev != 0 {
			n, _ := strconv.Atoi(output[i])
			if n < prev {
				t.Errorf("Expected %d to be greater than %d", n, prev)
				return
			}

			prev = n
		} else {
			prev, _ = strconv.Atoi(output[i])
		}
	}

	if err := q.Close(ctx); err != nil {
		t.Error(err)
		return
	}
}

func assertMPSC(t *testing.T, ctx context.Context, q main.Queue) {
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 0; i < 100000; i++ {
			if err := q.Put(ctx, strconv.Itoa(i)); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 100000; i < 200000; i++ {
			if err := q.Put(ctx, strconv.Itoa(i)); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	output := make([]string, 0, 200000)
	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 0; i < 200000; i++ {
			if err := q.Consume(ctx, func(s []string) error {
				if len(s) >= 1 {
					for i := 0; i < len(s); i++ {
						output = append(output, s[i])
					}
				}
				return nil
			}); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	wg.Wait()

	check := make(map[string]struct{}, 200000)
	for i := 0; i < 200000; i++ {
		check[strconv.Itoa(i)] = struct{}{}
	}

	for i := range output {
		if _, ok := check[output[i]]; !ok {
			t.Errorf("Expected %s to be in the map", output[i])
			return
		}
		delete(check, output[i])
	}

	if len(check) != 0 {
		t.Errorf("Expected check map to be empty, got %d", len(check))
		return
	}

	if err := q.Close(ctx); err != nil {
		t.Error(err)
		return
	}
}

func assertSPMC(t *testing.T, ctx context.Context, q main.Queue) {
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 0; i < 200000; i++ {
			if err := q.Put(ctx, strconv.Itoa(i)); err != nil {
				t.Error(err)
				return
			}
		}
		if err := q.Close(ctx); err != nil {
			t.Error(err)
			return
		}
	}()

	output := make([]string, 0, 200000)
	m := sync.Mutex{}
	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 0; i < 100000; i++ {
			if err := q.Consume(ctx, func(s []string) error {
				m.Lock()
				defer m.Unlock()

				if len(s) >= 1 {
					for i := 0; i < len(s); i++ {
						output = append(output, s[i])
					}
				}
				return nil
			}); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		for i := 0; i < 100000; i++ {
			if err := q.Consume(ctx, func(s []string) error {
				m.Lock()
				defer m.Unlock()

				if len(s) >= 1 {
					for i := 0; i < len(s); i++ {
						output = append(output, s[i])
					}
				}
				return nil
			}); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	wg.Wait()

	if len(output) != 200000 {
		t.Errorf("Expected %d, got %d", 200000, len(output))
		return
	}

	check := make(map[string]struct{}, 200000)
	for i := range output {
		if _, ok := check[output[i]]; ok {
			t.Errorf("Expected %s to not be duplicate", output[i])
		}
		check[output[i]] = struct{}{}
	}
}
