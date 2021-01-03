package goworker

import "errors"

var (
	// ErrorInactiveWorkerPool represents error due to inactive workerpool
	ErrorInactiveWorkerPool = errors.New("operation failed, workerpool is in-active")
)
