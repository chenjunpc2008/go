package kchanthreadpool

import (
    "fmt"
    "runtime"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

type myHld struct {
}

func (h *myHld) OnError(msg string) {
    const ftag = "myHld::OnError()"
    fmt.Println(ftag, msg)
}

func (h *myHld) OnEvent(msg string) {
    const ftag = "myHld::OnEvent()"
    fmt.Println(ftag, msg)
}

type myJob1 struct {
    data string
}

func (j *myJob1) Do(tid uint) {
    const ftag = "myJob1::Do()"
    fmt.Println(ftag, "ct", tid, j.data)
}

func Test_Pool_1(t *testing.T) {

    runtime.GOMAXPROCS(20)

    hld := new(myHld)

    pool, err := NewThreadPool(3, hld, 200)
    assert.Equal(t, nil, err)

    err = pool.Start()
    assert.Equal(t, nil, err)
    select {
    case <-time.After(2 * time.Second):
    }

    err = pool.Stop()
    assert.Equal(t, nil, err)

    select {
    case <-time.After(2 * time.Second):
    }
}

func Test_Pool_2(t *testing.T) {

    runtime.GOMAXPROCS(20)

    hld := new(myHld)

    pool, err := NewThreadPool(1, hld, 200)
    assert.Equal(t, nil, err)

    err = pool.Start()
    assert.Equal(t, nil, err)

    var (
        localJob *myJob1
        taskHold *Task
    )

    {
        localJob = &myJob1{data: "task 1"}
        taskHold = NewTask()
        taskHold.Data = localJob
        err = pool.AddTaskByMini(taskHold)
        assert.Equal(t, nil, err)
    }

    {
        localJob = &myJob1{data: "task 2"}
        taskHold = NewTask()
        taskHold.Data = localJob
        err = pool.AddTaskByMini(taskHold)
        assert.Equal(t, nil, err)
    }

    {
        localJob = &myJob1{data: "task 3"}
        taskHold = NewTask()
        taskHold.Data = localJob
        err = pool.AddTaskByMini(taskHold)
        assert.Equal(t, nil, err)
    }

    select {
    case <-time.After(10 * time.Second):
    }

    err = pool.Stop()
    assert.Equal(t, nil, err)

    select {
    case <-time.After(10 * time.Second):
    }
}

func Test_Pool_3(t *testing.T) {

    runtime.GOMAXPROCS(20)

    hld := new(myHld)
    pool, err := NewThreadPool(3, hld, 200)
    assert.Equal(t, nil, err)

    err = pool.Start()
    assert.Equal(t, nil, err)

    var (
        localJob *myJob1
        taskHold *Task
    )

    for i := 0; i < 50; i++ {
        msg := fmt.Sprintf("task %d", i)
        localJob = &myJob1{data: msg}
        taskHold = NewTask()
        taskHold.Data = localJob
        err = pool.AddTaskByMini(taskHold)
        assert.Equal(t, nil, err)

    }

    select {
    case <-time.After(5 * time.Second):
    }

    for i := 0; i < 50; i++ {
        msg := fmt.Sprintf("task %d", i+50)
        localJob = &myJob1{data: msg}
        taskHold = NewTask()
        taskHold.Data = localJob
        err = pool.AddTaskByMini(taskHold)
        assert.Equal(t, nil, err)

    }

    select {
    case <-time.After(5 * time.Second):
    }

    err = pool.Stop()
    assert.Equal(t, nil, err)

    select {
    case <-time.After(5 * time.Second):
    }
}
