//Package progress calculates an average item processing speed over a double sliding window
package progress

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/semaphore"
)

type Window struct {
	sem        *semaphore.Weighted
	total      time.Duration
	interval   time.Duration
	average    int64
	current    int64
	data       []int64
	historical []int64
	pos        int
	cancel     context.CancelFunc
}

//NewWindow returns a double sliding Window that calculates an average item processing speed
func NewWindow(total, interval time.Duration) (*Window, error) {
	if total == 0 {
		return nil, errors.New("progress: total cannot be 0")
	}
	if interval == 0 {
		return nil, errors.New("progress: interval cannot be 0")
	}
	if total <= interval || total%interval != 0 {
		return nil, errors.New("progress: total has to be a multiplier of interval")
	}

	cctx, cancel := context.WithCancel(context.Background())

	w := &Window{
		total:      total,
		interval:   interval,
		data:       make([]int64, int(total/interval)),
		historical: make([]int64, int(total/interval)),
		cancel:     cancel,
		sem:        semaphore.NewWeighted(1),
	}

	for i := range w.data {
		w.data[i] = -1
	}
	for i := range w.historical {
		w.historical[i] = -1
	}

	go w.tick(cctx)

	return w, nil
}

func (w *Window) tick(ctx context.Context) {
	ticker := time.NewTicker(w.interval)

	for {
		select {
		case <-ticker.C:
			w.move(ctx)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (w *Window) move(ctx context.Context) {
	w.sem.Acquire(ctx, 1)

	var total int64
	var num int64

	w.data[w.pos] = w.current

	for _, v := range w.data {
		if v == -1 {
			continue
		}
		total += v
		num++
	}

	if num == 0 {
		num = 1
	}

	firstAverage := total / num

	w.historical[w.pos] = firstAverage

	total = 0
	num = 0

	for _, v := range w.historical {
		if v == -1 {
			continue
		}
		total += v
		num++
	}

	if num == 0 {
		num = 1
	}

	secondAverage := total / num

	w.average = secondAverage * int64(len(w.historical))

	w.pos++

	if w.pos >= len(w.historical) {
		w.pos = 0
	}

	w.current = 0

	w.sem.Release(1)
}

//ItemCompleted adds a completed item to the Window
func (w *Window) ItemCompleted() {
	w.sem.Acquire(context.Background(), 1)

	w.current += 1

	w.sem.Release(1)
}

//Average gives you average item processing speed averaged over a double sliding Window
func (w *Window) Average() int64 {
	return w.average
}

//End stops the Window's ticking
func (w *Window) End() {
	w.cancel()
}
