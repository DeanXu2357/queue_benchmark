package implbyrb

import (
	"context"
	"sync"

	"qbench/contract"
)

/**
 * This is the implementation of the Queue interface using a ring buffer.
 */
type rbImpl struct {
	lock      sync.Locker
	readable  *sync.Cond
	writeable *sync.Cond
	buf       []string

	cap  int64
	head int64
	tail int64
}

func (r *rbImpl) Put(ctx context.Context, s string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.readable.Signal()

	for r.isFull() {
		r.writeable.Wait()
	}

	r.buf[r.head] = s
	r.head = (r.head + 1) % r.cap
	return nil
}

func (r *rbImpl) isFull() bool {
	return (r.tail+1)%r.cap == r.head
}

func (r *rbImpl) Consume(ctx context.Context, f func(jobs []string) error) error {
	data := r.get()

	if err := f([]string{data}); err != nil {
		return err
	}

	return nil
}

func (r *rbImpl) get() string {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.writeable.Signal()

	for r.isEmpty() {
		r.readable.Wait()
	}

	d := r.buf[r.tail]
	r.tail = (r.tail + 1) % r.cap
	return d
}

func (r *rbImpl) isEmpty() bool {
	return r.head == r.tail
}

func (r *rbImpl) Close(ctx context.Context) error {
	return nil
}

func NewRingBufferImpl(s int64) contract.Queue {
	b := make([]string, s)
	m := sync.Mutex{}

	return &rbImpl{
		cap:       s,
		buf:       b,
		lock:      &m,
		readable:  sync.NewCond(&m),
		writeable: sync.NewCond(&m),
	}
}
