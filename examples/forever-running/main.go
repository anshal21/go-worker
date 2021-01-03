package main

import (
	"fmt"
	"time"

	goworker "github.com/anshal21/go-worker"
)

var wp = goworker.NewWorkerPool(&goworker.WorkerPoolInput{
	WorkerCount: 1,
})

func addworkPeriodically() {

	cnt := 0
	batchCount := 0
	for {
		batchCount++
		for i := 0; i < 10; i++ {
			cnt++
			val := cnt
			wp.Add(&goworker.Task{
				F: func() (interface{}, error) {
					fmt.Printf("Batch: %v, Val: %v\n", batchCount, val)
					return nil, nil
				},
			})
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	wp.Start()
	go addworkPeriodically()
	wp.WaitForCompletion()
}
