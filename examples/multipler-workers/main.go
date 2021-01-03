package main

import (
	"fmt"

	goworker "github.com/anshal21/go-worker"
)

func main() {
	// instantiate the workerpool
	wp := goworker.NewWorkerPool(&goworker.WorkerPoolInput{
		WorkerCount: 5,
	})

	// starts the execution of queued tasks
	wp.Start()

	// add tasks to the pool
	for i := 0; i < 100; i++ {
		val := i
		wp.Add(&goworker.Task{
			F: func() (interface{}, error) {
				fmt.Println(val)
				return nil, nil
			},
		})
	}

	// tells goworker that all the tasks are added
	wp.Done()

	wp.WaitForCompletion()
}
