package main

import (
	"errors"
	"fmt"

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
		WorkerCount: 5,
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
				wp.Abort()
			}
		}()
	}

	wp.Done()
	wp.WaitForCompletion()
}
