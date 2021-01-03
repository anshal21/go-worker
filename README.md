
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/anshal21/go-worker/blob/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/anshal21/go-worker)](https://goreportcard.com/report/github.com/anshal21/go-worker) [![Build Status](https://travis-ci.com/anshal21/go-worker.svg?branch=main)](https://travis-ci.com/anshal21/go-worker)

# go-worker
go-worker provides an implementation of thread pool pattern. It provides the *WorkerPool* interface which exposes following methods

### Add 
**Add** method adds the provided task to the task queue of workerpool. It takes a **Task** object as input and returns a **Future** object containing the response from the task

### Start
**Start** method starts the task execution in the workerpool

### Done
**Done** sends a signal to the workerpool, notifying that no more tasks are to be added to the pool

### Abort
**Abort** sends a signal to the workerpool, notifying it that further task execution should be aborted

### WaitForCompletion
**WaitForCompletion** is a blocking method, it waits for all the tasks in the pool to complete the execution before returning


