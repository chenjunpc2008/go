package kthreadpool

import (
    "fmt"
)

const (
    cGetTaskGap      = 1
    cMaxTaskInOneRun = 2000
)

// thread
type kThread struct {
    tid     uint
    tp      *ThreadPool
    taskQue tQueue
}

//
func newKThread(id uint, tp *ThreadPool) *kThread {
    var th = &kThread{}
    th.tid = id
    th.tp = tp
    th.taskQue.Init()

    return th
}

//
func (th *kThread) GetTaskSize() int {
    return th.taskQue.GetSize()
}

//
func (th *kThread) AddTask(elem *Task) {
    th.taskQue.PushBack(elem)
}

//
func (th *kThread) Start() {
    go threadRun(th)
}

//
func threadRun(th *kThread) {
    var (
        jobTodo *Task
        have    int
        i       int
        tid     = th.tid

        // this chan is only for wait, because we checked out any
        // time.Sleep() or time.After will() will cause a 15ms gap
        chWait chan int = make(chan int)
    )

    // close the channel, use the channel for wait a little while
    close(chWait)

    for {
        if th.tp.bExit {
            sEventMsg := fmt.Sprintf("thread-id:%d exit", th.tid)
            // fmt.Printf(sEventMsg)
            th.tp.handler.OnEvent(sEventMsg)
            return
        }

        if 0 == have {
            // wait for a while
            _ = <-chWait
        }

        for i = 0; i < cMaxTaskInOneRun; i++ {
            jobTodo, have = th.taskQue.PopFrontOne()
            if 0 == have {
                break
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
