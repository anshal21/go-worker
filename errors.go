package goworker

import "errors"

var (
	// ErrorInactiveWorkerPool represents error due to inactive workerpool
	ErrorInactiveWorkerPool = errors.New("operation failed, workerpool is in-active")

	// ErrorWorkerPoolAborted is returned for the pending tasks in the queue after the pool is aborted
	ErrorWorkerPoolAborted = errors.New("workerpool was aborted before the task could be scheduled")
)
