package main

import (
	"fmt"

	goworker "github.com/anshal21/go-worker"
)

func doWork(i int) int {
	return i * 2
}

func main() {

	data := make([]int, 0)
	for i := 0; i < 100; i++ {
		data = append(data, i)
	}

	// instantiate the workerpool
	wp := goworker.NewWorkerPool(&goworker.WorkerPoolInput{
		WorkerCount: 5,
	})

	// starts the execution of queued tasks
	wp.Start()

	results := make([]*goworker.Future, len(data))

	// scatter the task to multiple workers
	for index := range data {
		val := data[index]
		results[index] = wp.Add(&goworker.Task{
			F: func() (interface{}, error) {
				return doWork(val), nil
			},
		})
	}

	wp.Done()
	wp.WaitForCompletion()

	// gather results from workers
	for _, res := range results {
		fmt.Println(res.Result())
	}
}
