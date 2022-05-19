package kchanthreadpool

// TaskJobIF job interface
type TaskJobIF interface {
    Do(threadID uint)
}

// Task task for thread pool
type Task struct {
    Data TaskJobIF
    Key  uint // key value for serial job
}

// NewTask new task object
func NewTask() *Task {
    return &Task{}
}
