package kbuffthreadpool

import "sync"

type tQueue struct {
    // lock
    lock sync.Mutex

    buff []*Task
}

//
func (q *tQueue) Init() {
    q.buff = make([]*Task, 0)
}

//
func (q *tQueue) GetSize() int {
    return len(q.buff)
}

//
func (q *tQueue) PushBack(elem *Task) {
    // lock
    q.lock.Lock()

    q.buff = append(q.buff, elem)

    // unlock
    q.lock.Unlock()
}

/*
pop front one task

@return *Task  : task object
@return int : 0 -- don't have any more task in queue
*/
func (q *tQueue) PopFrontOne() (*Task, int) {
    // lock
    q.lock.Lock()

    if 0 == len(q.buff) {
        // unlock
        q.lock.Unlock()
        return nil, 0
    }

    var ta *Task

    ta = q.buff[0]
    q.buff = q.buff[1:]

    // unlock
    q.lock.Unlock()
    return ta, 1
}
