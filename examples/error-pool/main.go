package main

import (
	"errors"
	"fmt"
	"time"

	goworker "github.com/anshal21/go-worker"
)

func doWork(i int) error {
	fmt.Println(i)
	if i == 29 {
		return errors.New("some error")
	}
	return nil
}

func main() {

	// instantiate the workerpool
	wp := goworker.NewWorkerPool(&goworker.WorkerPoolInput{
		WorkerCount: 10,
	})

	// starts the execution of queued tasks
	wp.Start()

	// scatter the task to multiple workers
	for i := 0; i < 100; i++ {
		val := i
		future := wp.Add(&goworker.Task{
			F: func() (interface{}, error) {
				return nil, doWork(val)
			},
		})

		// abort worker pool in the case of error
		go func() {
			err := future.Error()
			if err != nil {
				if err != goworker.ErrorInactiveWorkerPool && err != goworker.ErrorWorkerPoolAborted {
					wp.Abort()
				}
			}
		}()
	}

	wp.Done()
	wp.WaitForCompletion()
	time.Sleep(10 * time.Second)
}
