package util

import "sync"

type Queue[T any] struct {
	items []T
	lock  sync.Mutex
	cond  *sync.Cond
}

func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{
		items: make([]T, 0),
	}
	q.cond = sync.NewCond(&q.lock)
	return q
}
func (q *Queue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}
func (q *Queue[T]) Dequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == 0 {
		var t T
		return t, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
func (q *Queue[T]) WaitDequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
func (q *Queue[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.items)
}
func (q *Queue[T]) Iter() <-chan T {
	c := make(chan T)
	go func() {
		q.lock.Lock()
		defer q.lock.Unlock()
		defer close(c)
		for _, item := range q.items {
			c <- item
		}
	}()
	return c
}
func (q *Queue[T]) Peek() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == 0 {
		var t T
		return t, false
	}
	return q.items[0], true
}
