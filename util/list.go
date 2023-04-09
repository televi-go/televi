package util

import "sync"

type node[T any] struct {
	value T
	next  *node[T]
}
type LinkedList[T any] struct {
	head *node[T]
	len  int
	lock sync.Mutex
	cond *sync.Cond
}

func NewLinkedList[T any]() *LinkedList[T] {
	l := &LinkedList[T]{}
	l.cond = sync.NewCond(&l.lock)
	return l
}
func (l *LinkedList[T]) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.len
}
func (l *LinkedList[T]) PushFront(value T) {
	l.lock.Lock()
	defer l.lock.Unlock()
	n := &node[T]{value: value}
	n.next = l.head
	l.head = n
	l.len++
	l.cond.Signal()
}
func (l *LinkedList[T]) WaitPopFront() T {
	l.lock.Lock()
	defer l.lock.Unlock()
	for l.head == nil {
		l.cond.Wait()
	}
	n := l.head
	l.head = n.next
	l.len--
	return n.value
}
func (l *LinkedList[T]) Delete(index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if index < 0 || index >= l.len {
		return
	}
	if index == 0 {
		l.head = l.head.next
		l.len--
		return
	}
	prev := l.head
	curr := l.head.next
	for i := 1; i < index; i++ {
		prev = curr
		curr = curr.next
	}
	prev.next = curr.next
	l.len--
}
func (l *LinkedList[T]) Iter() <-chan T {
	c := make(chan T)
	go func() {
		l.lock.Lock()
		defer l.lock.Unlock()
		defer close(c)
		for n := l.head; n != nil; n = n.next {
			c <- n.value
		}
	}()
	return c
}

func (l *LinkedList[T]) ToArray() []T {
	var result []T
	l.lock.Lock()
	defer l.lock.Unlock()
	for n := l.head; n != nil; n = n.next {
		result = append(result, n.value)
	}
	return result
}
