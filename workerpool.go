package goworker

import (
	"sync"
)

// Request contains the parameters to initialize a worker pool
// Default values are used for missing parameters
// Following are the default values
//	- WorkerCount: 1
//	- Buffer: 100
type WorkerPoolInput struct {
	// WorkerCount tells how many workers are to be initiated
	// in the pool, default value is 1
	WorkerCount int
	// Buffer specifies how many tasks can be add to the queue, before
	// the wp.Add becomes a blocking call
	Buffer int
}

const (
	_defaultBufferSize  = 100
	_defaultWorkerCount = 1
)

// WorkerPool is an interface to a workerpool implementation
type WorkerPool interface {
	// Add adds a task to the worker pool
	// the function returns a future object that holds the response of the execution
	// currently it takes Task as input so only error response is supported
	Add(task *Task) *Future
	// Start starts the execution of the tasks in the workerpool
	Start()
	// Done notifies the workerpool that there are no more tasks to be added
	// if Done is never called, the workerpool will keep waiting for tasks forever
	Done()
	// Abort notifies workerpool to abort the execution of all remaning tasks in the pool
	Abort()
	// WaitForCompletion is a blocking function that waits till workerpool finishes all the tasks
	WaitForCompletion()
}

type taskWrapper struct {
	task     *Task
	response *Future
}

// NewWorkerPool returns a new instance of WorkerPool with provided specs
func NewWorkerPool(request *WorkerPoolInput) WorkerPool {
	request = populateDefaults(request)

	return &workerpool{
		workerCount: request.WorkerCount,
		abortChan:   make(chan struct{}, request.WorkerCount),
		taskChan:    make(chan taskWrapper, request.Buffer),
	}
}

func populateDefaults(request *WorkerPoolInput) *WorkerPoolInput {
	if request.WorkerCount == 0 {
		request.WorkerCount = _defaultWorkerCount
	}
	if request.Buffer == 0 {
		request.Buffer = _defaultBufferSize
	}
	return request
}

type workerpool struct {
	wg          sync.WaitGroup
	workerCount int
	abortChan   chan struct{}
	taskChan    chan taskWrapper
}

func (w *workerpool) Start() {
	w.wg.Add(w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		go func() {
			defer w.wg.Done()
			for {
				select {
				case task, ok := <-w.taskChan:
					if !ok {
						return
					}
					res, err := task.task.Run()
					task.response.NotifyResult(res)
					task.response.NotifyError(err)
				case <-w.abortChan:
					w.taskChan = nil
					return

				}
				if w.taskChan == nil {
					return
				}
			}
		}()
	}
}

func (w *workerpool) Add(task *Task) *Future {
	res := NewFuture()
	workUnit := taskWrapper{
		task:     task,
		response: res,
	}

	// return if workerpool is in an inactive state
	if w.taskChan == nil {
		res.NotifyResult(nil)
		res.NotifyError(ErrorInactiveWorkerPool)
		return res
	}

	w.taskChan <- workUnit
	return res
}

func (w *workerpool) Done() {
	if w.taskChan != nil {
		close(w.taskChan)
	}
}

func (w *workerpool) Abort() {
	if w.taskChan != nil {
		for i := 0; i < w.workerCount; i++ {
			w.abortChan <- struct{}{}
		}
	}
}

func (w *workerpool) WaitForCompletion() {
	w.wg.Wait()
}
