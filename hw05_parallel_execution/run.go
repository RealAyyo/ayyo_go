package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrPositiveWorkersCount = errors.New("number of workers must be positive")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrPositiveWorkersCount
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasksChan := make(chan Task)
	errorsChan := make(chan error, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskWorker(ctx, tasksChan, errorsChan)
		}()
	}

	var wgErrs sync.WaitGroup
	wgErrs.Add(1)
	go func() {
		defer wgErrs.Done()
		var errorCount int
		for err := range errorsChan {
			if err != nil {
				errorCount++
				if m > 0 && errorCount >= m {
					cancel()
					return
				}
			}
		}
	}()

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			break
		case tasksChan <- task:
		}
	}
	close(tasksChan)
	wg.Wait()
	close(errorsChan)
	wgErrs.Wait()

	if ctx.Err() != nil && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func taskWorker(ctx context.Context, tasksChan chan Task, errorsChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasksChan:
			if !ok {
				return
			}
			err := task()
			if err != nil {
				errorsChan <- err
			}
		}
	}
}
