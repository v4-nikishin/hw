package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var rng []Task
	var cntTask int
	var cntErr int32

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 || n >= len(tasks) {
		rng = tasks
		cntTask = len(tasks)
	}

	for {
		if cntTask != len(tasks) {
			if n < len(tasks)-cntTask {
				rng = tasks[cntTask : cntTask+n]
				cntTask += n
			} else {
				rng = tasks[cntTask:]
				cntTask = len(tasks)
			}
		}

		wg := sync.WaitGroup{}
		for _, t := range rng {
			t := t
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := t(); err != nil {
					atomic.AddInt32(&cntErr, 1)
				}
			}()
		}
		wg.Wait()

		if atomic.LoadInt32(&cntErr) >= int32(m) {
			return ErrErrorsLimitExceeded
		}
		if cntTask == len(tasks) {
			break
		}
	}
	return nil
}
