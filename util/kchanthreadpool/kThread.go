package kchanthreadpool

import (
    "fmt"
    "sync"
)

// thread
type kThread struct {
    tid uint
    tp  *ThreadPool

    maxBuffsize int
    lock        sync.Mutex // lock for below object
    chTask      chan *Task
}

//
func newKThread(id uint, buffsize int, tp *ThreadPool) *kThread {
    var th = &kThread{}
    th.tid = id
    th.tp = tp
    th.maxBuffsize = buffsize
    th.chTask = make(chan *Task, buffsize)

    return th
}

//
func (th *kThread) GetTaskSize() int {
    if nil != th.chTask {
        return len(th.chTask)
    }

    return 0
}

/*
AddTask

@return busy bool : true -- buff is full, you may need to try again
*/
func (th *kThread) AddTask(elem *Task) (busy bool) {
    // lock
    th.lock.Lock()

    curBuffSize := len(th.chTask)
    if curBuffSize >= th.maxBuffsize-1 {
        // unlock
        th.lock.Unlock()

        return true
    }

    th.chTask <- elem

    // unlock
    th.lock.Unlock()

    return false
}

//
func (th *kThread) Start() {
    go threadRun(th)
}

//
func threadRun(th *kThread) {
    var (
        jobTodo *Task
        ok      bool
        tid     = th.tid
    )

    for {

        select {
        case <-th.tp.chExit:
            sEventMsg := fmt.Sprintf("thread-id:%d exit", th.tid)
            // fmt.Printf(sEventMsg)
            th.tp.handler.OnEvent(sEventMsg)
            return

        case jobTodo, ok = <-th.chTask:
            if !ok {
                // channel closed
                sEventMsg := fmt.Sprintf("thread-id:%d channel closed", th.tid)
                // fmt.Printf(sEventMsg)
                th.tp.handler.OnEvent(sEventMsg)
                return
            }

            if nil == jobTodo {
                sErrMsg := fmt.Sprintf("thread-id:%d have nil *Task", th.tid)
                // fmt.Printf(sErrMsg)
                th.tp.handler.OnError(sErrMsg)
                continue
            }

            // run job
            jobTodo.Data.Do(tid)
        }
    }
}
