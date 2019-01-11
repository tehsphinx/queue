package queue

import (
	"sync"
	"testing"
)

func BenchmarkQueue_Push(b *testing.B) {
	var q = NewQueue()

	for i := 0; i < b.N; i++ {
		_ = q.Push(struct{}{})
	}
}

func BenchmarkQueue_Pop(b *testing.B) {
	var q = NewQueue()

	for i := 0; i < b.N; i++ {
		_ = q.Push(struct{}{})
	}

	b.ResetTimer()

	var v interface{}
	for i := 0; i < b.N; i++ {
		v, _ = q.Pop()
	}
	_ = v
}

func BenchmarkQueue_PopBlocking(b *testing.B) {
	var q = NewQueue()

	for i := 0; i < b.N; i++ {
		_ = q.Push(struct{}{})
	}

	b.ResetTimer()

	var v interface{}
	for i := 0; i < b.N; i++ {
		v, _ = q.PopBlocking()
	}
	_ = v
}

func BenchmarkQueue_Push_PopBlocking(b *testing.B) {
	var q = NewQueue()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < b.N; i++ {
			_ = q.Push(struct{}{})
		}
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		var v interface{}
		for i := 0; i < b.N; i++ {
			v, _ = q.PopBlocking()
		}
		_ = v
		wg.Done()
	}(wg)

	wg.Wait()
}
