//Package progress calculates item processing speed over a double sliding window
package progress

import (
	"errors"
	"sync"
	"time"
)

type Window struct {
	mu         sync.RWMutex
	total      time.Duration
	interval   time.Duration
	current    int64
	data       []int64
	historical []int64
	pos        int
	end        chan struct{}
}

//NewWindow returns a double sliding Window that can help you calculate item processing speed
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

	w := &Window{
		total:      total,
		interval:   interval,
		data:       make([]int64, int(total/interval)),
		historical: make([]int64, int(total/interval)),
		end:        make(chan struct{}, 1),
	}

	for i := range w.data {
		w.data[i] = -1
	}
	for i := range w.historical {
		w.historical[i] = -1
	}

	go w.tick()

	return w, nil
}

func (w *Window) tick() {
	ticker := time.NewTicker(w.interval)

	for {
		select {
		case <-ticker.C:
			w.move()
		case <-w.end:
			ticker.Stop()
			return
		}
	}
}

func (w *Window) move() {
	w.mu.Lock()

	w.data[w.pos] = w.current

	w.historical[w.pos] = w.currentAverage()

	w.pos++

	if w.pos >= len(w.data) {
		w.pos = 0
	}

	w.current = 0

	w.mu.Unlock()
}

func (w *Window) currentAverage() int64 {
	var total int64
	var num int64

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

	return total / num
}

//ItemCompleted adds an item to the Window
func (w *Window) ItemCompleted() {
	w.mu.Lock()

	w.current += 1

	w.mu.Unlock()
}

//Average gives you the average amount of items processed in the window, averaged over another sliding window
func (w *Window) Average() int64 {
	var total int64
	var num int64

	w.mu.Lock()

	for _, v := range w.historical {
		if v == -1 {
			continue
		}
		total += v
		num++
	}

	w.mu.Unlock()

	if num == 0 {
		num = 1
	}

	average := total / num

	return average * int64(len(w.historical))
}

//End stops the Window's ticking
func (w *Window) End() {
	select {
	case w.end <- struct{}{}:
	default:
		return
	}
}
