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

	tasksChan := make(chan Task, len(tasks))
	errorsChan := make(chan error, len(tasks))
	doneChan := make(chan bool)

	go func() {
		for _, task := range tasks {
			tasksChan <- task
		}
		close(tasksChan)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				select {
				case <-doneChan:
					return
				default:
					errorsChan <- task()
				}
			}
		}()
	}

	errorCount := 0
	for i := 0; i < len(tasks); i++ {
		err := <-errorsChan
		if err != nil {
			errorCount++
			if errorCount >= m {
				close(doneChan)
				wg.Wait()
				close(errorsChan)
				return ErrErrorsLimitExceeded
			}
		}
	}

	wg.Wait()
	close(errorsChan)
	return nil
}
