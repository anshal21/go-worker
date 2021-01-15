
<img src="https://github.com/anshal21/images/blob/master/go-worker/gopher.png" height=300 ></img>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/anshal21/go-worker/blob/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/anshal21/go-worker)](https://goreportcard.com/report/github.com/anshal21/go-worker) [![Build Status](https://travis-ci.com/anshal21/go-worker.svg?branch=main)](https://travis-ci.com/anshal21/go-worker)

# go-worker
go-worker is an implementation of thread pool pattern. It exposes the *WorkerPool* interface which provides following methods  

### Methods
`Add`: Adds the provided task to the task queue of workerpool. It takes a *Task* object as input and returns a *Future* object containing the response from the task 

`Start`: Starts the task execution in the workerpool  

`Done`: Sends a signal to the workerpool, notifying that no more tasks are to be added to the pool 

`Abort`: Sends a signal to the workerpool, notifying it that further task execution should be aborted  

`WaitForCompletion`: It is a blocking method, it waits for all the tasks in the pool to complete the execution before returning 

 ### Models
 `Future`: A future is returned for every task that is added to the workerpool. The user can use *future.Result()* and *future.Error()* methods to get the *result* and *error* from the task
 
 ### Errors
 `ErrorInactiveWorkerPool`: The error is returned when there is an attempt to add a task to already completed or aborted workerpool
 
 `ErrorWorkerPoolAborted`: The error is returned for all the tasks that couldn't be scheduled as the workerpool was aborted
 
 Link: https://anshal21.github.io/go-worker/
 
 ## Installation
 ```bash
 go get github.com/anshal21/go-worker
 ```

 ## Usage Patterns

 ### 1. [ Task Distribution / Parallel Processing ](https://github.com/anshal21/go-worker/blob/main/examples/multipler-workers/main.go)  
 ```go
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
 ``` 

  ### 2. [ Scatter-Gather ](https://github.com/anshal21/go-worker/blob/main/examples/scatter-gather/main.go)       
 ```go
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
 ```  

 ### 3. [ Exit on Error ]( https://github.com/anshal21/go-worker/blob/main/examples/error-pool/main.go )
  ```go
    // instantiate the workerpool
	wp := goworker.NewWorkerPool(&goworker.WorkerPoolInput{
		WorkerCount: 1,
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
				if err == goworker.ErrorInactiveWorkerPool || err == goworker.ErrorWorkerPoolAborted {
					return
				}
				wp.Abort()
			}
		}()
	}

	wp.Done()
	wp.WaitForCompletion()
 ``` 

### 4. [ Run Forever ]( https://github.com/anshal21/go-worker/blob/main/examples/forever-running/main.go )
  ```go

func addworkPeriodically() {
	for {
		for i := 0; i < 10; i++ {
			wp.Add(&goworker.Task{
				F: func() (interface{}, error) {
					fmt.Println("did something")
					return nil, nil
				},
			})
		}
	}
}

func main() {
	wp.Start()
	go addworkPeriodically()
	wp.WaitForCompletion()
}
 ``` 

