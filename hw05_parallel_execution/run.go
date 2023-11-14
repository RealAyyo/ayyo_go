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
	startWorkers(ctx, &wg, n, tasksChan, errorsChan)
	processTasks(ctx, tasks, tasksChan)
	wg.Wait()
	close(errorsChan)

	return processErrors(errorsChan, m)
}

func startWorkers(ctx context.Context, wg *sync.WaitGroup, wCount int, tasksChan <-chan Task, errorsChan chan<- error) {
	for i := 0; i < wCount; i++ {
		wg.Add(1)
		go worker(ctx, wg, tasksChan, errorsChan)
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, tasksChan <-chan Task, errorsChan chan<- error) {
	defer wg.Done()
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

func processTasks(ctx context.Context, tasks []Task, tasksChan chan<- Task) {
	for _, task := range tasks {
		select {
		case <-ctx.Done():
			break
		case tasksChan <- task:
		}
	}
	close(tasksChan)
}

func processErrors(errorsChan <-chan error, m int) error {
	var errorCount int
	for err := range errorsChan {
		if err != nil {
			errorCount++
			if m > 0 && errorCount >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}
	return nil
}
