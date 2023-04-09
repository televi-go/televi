package util

import "testing"

func TestNewLinkedList(t *testing.T) {
	l := NewLinkedList[int]()
	if l.Len() != 0 {
		t.Errorf("Expected length to be %d, but got %d", 0, l.Len())
	}
}
func TestPushFront(t *testing.T) {
	l := NewLinkedList[int]()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)
	if l.Len() != 3 {
		t.Errorf("Expected length to be %d, but got %d", 3, l.Len())
	}
}
func TestWaitPushFront(t *testing.T) {
	l := NewLinkedList[int]()
	go func() {
		for i := 0; i < 5; i++ {
			l.PushFront(i)
		}
	}()
	for i := 4; i <= 0; i-- {
		value := l.WaitPopFront()
		if value != i {
			t.Errorf("Expected value to be %d, but got %d", i, value)
		}
	}
	if l.Len() != 0 {
		t.Errorf("Expected length to be %d, but got %d", 0, l.Len())
	}
}
func TestDelete(t *testing.T) {
	l := NewLinkedList[int]()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)
	l.Delete(1)
	if l.Len() != 2 {
		t.Errorf("Expected length to be %d, but got %d", 2, l.Len())
	}
	i := 0
	for item := range l.Iter() {
		if i == 0 && item != 3 {
			t.Errorf("Expected value to be %d, but got %d", 3, item)
		}
		if i == 1 && item != 1 {
			t.Errorf("Expected value to be %d, but got %d", 1, item)
		}
		i++
	}
}
func TestIter(t *testing.T) {
	l := NewLinkedList[int]()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)
	i := 0
	for item := range l.Iter() {
		if i == 0 && item != 3 {
			t.Errorf("Expected value to be %d, but got %d", 3, item)
		}
		if i == 1 && item != 2 {
			t.Errorf("Expected value to be %d, but got %d", 2, item)
		}
		if i == 2 && item != 1 {
			t.Errorf("Expected value to be %d, but got %d", 1, item)
		}
		i++
	}
}
