/*
Package kchanthreadpool 协程池
说明
使用channel做任务队列的协程池
每个协程仅从自身任务队列中取任务，这样可以减少锁的争抢
每个任务可分配至最少任务的队列中，也可根据任务key值分配到相同的任务池，以此来保证任务的严格串行
*/
package kchanthreadpool

import (
    "errors"
    "fmt"
)

/*
 */

// PoolStatus pool status
type PoolStatus int

const (
    // PoolStatusClosed closed
    PoolStatusClosed PoolStatus = 0
    // PoolStatusInitialized initialized
    PoolStatusInitialized PoolStatus = 1
    // PoolStatusStarting starting
    PoolStatusStarting PoolStatus = 2
    // PoolStatusRunning running
    PoolStatusRunning PoolStatus = 3
    // PoolStatusStopping stopping
    PoolStatusStopping PoolStatus = 4
)

/*
PoolHandlerIF handler for upper apps
*/
type PoolHandlerIF interface {
    OnError(msg string)
    OnEvent(msg string)
}

// ThreadPool thread pool
type ThreadPool struct {
    chExit chan int // channel for exit

    numThreads uint

    threads []*kThread
    status  PoolStatus
    handler PoolHandlerIF

    iBuffSizePerThread int
}

// NewThreadPool new thread pool
func NewThreadPool(num uint, buffsize int, cb PoolHandlerIF) (*ThreadPool, error) {
    var (
        tp  *ThreadPool = &ThreadPool{}
        err error
    )

    tp.handler = cb
    tp.chExit = make(chan int)
    tp.status = PoolStatusClosed
    tp.numThreads = num
    tp.threads = make([]*kThread, 0)
    tp.iBuffSizePerThread = buffsize

    err = tp.init()

    return tp, err
}

// init
func (tp *ThreadPool) init() error {
    if PoolStatusClosed != tp.status {
        return errors.New("init failed status not closed")
    }

    var childTh *kThread

    for i := uint(0); i < tp.numThreads; i++ {
        childTh = newKThread(i, tp.iBuffSizePerThread, tp)
        tp.threads = append(tp.threads, childTh)
    }

    tp.status = PoolStatusInitialized

    return nil
}

// Start start
func (tp *ThreadPool) Start() error {
    if PoolStatusInitialized != tp.status {
        return errors.New("init failed status not initialized")
    }

    tp.status = PoolStatusStarting

    for k, v := range tp.threads {
        if nil == v {
            return fmt.Errorf("nil threadHld in pos:%v", k)
        }

        v.Start()
    }

    tp.status = PoolStatusRunning

    return nil
}

// Stop stop
func (tp *ThreadPool) Stop() error {
    if PoolStatusClosed == tp.status || PoolStatusStopping == tp.status {
        return errors.New("stop failed status in closed or stopping")
    }

    close(tp.chExit)

    tp.status = PoolStatusClosed

    return nil
}

/*
AddTaskByMini add task

@return busy bool : true -- buff is full, you may need to try again
@return retErr error : error
*/
func (tp *ThreadPool) AddTaskByMini(elem *Task) (busy bool, retErr error) {
    if PoolStatusRunning != tp.status {
        return false, fmt.Errorf("add task initialed when status:%d", tp.status)
    }

    var (
        miniTaskSize    int
        currentTaskSize int
        miniThread      *kThread
    )

    for k, v := range tp.threads {
        if nil == v {
            sErrMsg := fmt.Sprintf("nil threadHld in pos:%v", k)
            tp.handler.OnError(sErrMsg)
            continue
        }

        currentTaskSize = v.GetTaskSize()
        if 0 == currentTaskSize {
            miniThread = v
            break
        }

        if 0 == miniTaskSize || miniTaskSize > currentTaskSize {
            miniTaskSize = currentTaskSize
            miniThread = v
        }
    }

    if nil == miniThread {
        sErrMsg := "nil threadHld to add task"
        tp.handler.OnError(sErrMsg)
        return false, errors.New(sErrMsg)
    }

    busy = miniThread.AddTask(elem)
    if busy {
        return true, nil
    }

    return false, nil
}

/*
AddTaskByKey add task

@return busy bool : true -- buff is full, you may need to try again
@return retErr error : error
*/
func (tp *ThreadPool) AddTaskByKey(elem *Task) (busy bool, retErr error) {
    if PoolStatusRunning != tp.status {
        return false, fmt.Errorf("add task initialed when status:%d", tp.status)
    }

    var (
        dest       uint
        destThread *kThread
    )

    dest = elem.Key % tp.numThreads

    destThread = tp.threads[dest]

    if nil == destThread {
        sErrMsg := "nil threadHld to add task"
        tp.handler.OnError(sErrMsg)
        return false, errors.New(sErrMsg)
    }

    busy = destThread.AddTask(elem)
    if busy {
        return true, nil
    }

    return false, nil
}

/*
GetTaskSize task size
*/
func (tp *ThreadPool) GetTaskSize() int {
    var (
        currentTaskSize int
    )

    for k, v := range tp.threads {

        if nil == v {
            sErrMsg := fmt.Sprintf("nil threadHld in pos:%v", k)
            tp.handler.OnError(sErrMsg)
            continue
        }

        currentTaskSize += v.GetTaskSize()
    }

    return currentTaskSize
}

/*
GetTaskInfo task info
*/
func (tp *ThreadPool) GetTaskInfo() string {

    var (
        sOut            string
        sTmp            string
        currentTaskSize int
    )

    sOut = "["
    for k, v := range tp.threads {
        sOut += "{"
        if nil == v {
            sErrMsg := fmt.Sprintf("nil threadHld in pos:%v", k)
            tp.handler.OnError(sErrMsg)

            sTmp = fmt.Sprintf("tid:%d nil", k)
            sOut += sTmp
            sOut += "} "
            continue
        }

        currentTaskSize = v.GetTaskSize()

        sTmp = fmt.Sprintf("tid:%d tsize:%d", k, currentTaskSize)
        sOut += sTmp
        sOut += "} "
    }

    sOut += "]"

    return sOut
}
