package queue

import (
	"container/list"
	"errors"
	"sync"
)

// NewQueue returns a new FIFO queue.
func NewQueue() *Queue {
	q := &Queue{
		list: list.New(),
		lock: sync.Mutex{},
	}
	q.cond = sync.NewCond(&q.lock)
	return q
}

// Queue is a basic FIFO queue.
type Queue struct {
	list   *list.List
	lock   sync.Mutex
	cond   *sync.Cond
	closed bool
}

// Push adds a node to the queue.
// Will fail if the queue has been closed.
func (s *Queue) Push(n interface{}) error {
	s.lock.Lock()
	if s.closed {
		s.lock.Unlock()
		return errors.New("queue.Queue: pushing item into a closed queue")
	}
	s.list.PushBack(n)
	s.lock.Unlock()

	s.cond.Broadcast()
	return nil
}

// Pop removes and returns a node from the queue.
// The bool value is false if the queue was empty.
func (s *Queue) Pop() (interface{}, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	e := s.list.Front()
	if e == nil {
		return nil, false
	}
	s.list.Remove(e)
	return e.Value, true
}

// Next blocks until an element is available or the queue is closed.
// Reports false if the queue has been emptied and is closed.
// Beware that if multiple goroutines read from the queue, Pop can return false after Next was true.
func (s *Queue) Next() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for s.list.Len() == 0 {
		if s.closed {
			return false
		}
		s.cond.Wait()
	}
	return true
}

// PopBlocking removes and returns a node from the queue.
// It blocks until an element is available or the queue is closed.
// The bool value is false if the queue has been emptied and was closed.
func (s *Queue) PopBlocking() (interface{}, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for s.list.Len() == 0 {
		s.cond.Wait()
		if s.closed && s.list.Len() == 0 {
			return nil, false
		}
	}

	e := s.list.Front()
	s.list.Remove(e)
	return e.Value, true
}

// Len reports the current length of the queue.
func (s *Queue) Len() int {
	return s.list.Len()
}

// Clear empties the buffer
func (s *Queue) Clear() {
	s.list.Init()
}

// Close closes the queue. Use this to release all PopBlocking calls.
// The queue will be emptied before PopBlocking returns false.
func (s *Queue) Close() error {
	s.lock.Lock()
	s.closed = true
	s.lock.Unlock()
	s.cond.Broadcast()
	return nil
}
