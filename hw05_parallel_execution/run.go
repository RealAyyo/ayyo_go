package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup
	tasksChan := make(chan Task)
	errorsChan := make(chan error, n)
	doneChan := make(chan struct{})

	go func() {
		for _, task := range tasks {
			tasksChan <- task
		}
		close(tasksChan)
	}()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				select {
				case <-doneChan:
					return
				default:
					if err := task(); err != nil {
						errorsChan <- err
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneChan)
		close(errorsChan)
	}()

	var errorCount int
	for err := range errorsChan {
		if err != nil {
			errorCount++

			if errorCount >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}
