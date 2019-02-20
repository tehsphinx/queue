package queue

import (
	"sync"
	"testing"
)

func TestQueue_Push_Next_Pop(t *testing.T) {
	var q = NewQueue()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			_ = q.Push(struct{}{})
		}
		_ = q.Close()
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		var v interface{}
		for q.Next() {
			v, _ = q.Pop()
		}
		_ = v
		wg.Done()
	}(wg)

	wg.Wait()
}

func TestQueue_Has(t *testing.T) {
	var q = NewQueue()

	for i := 0; i < 10; i++ {
		err := q.Push(i)
		if err != nil {
			t.Fatal(err)
		}
	}

	if !q.Has(func(elem interface{}) bool {
		return elem.(int) == 8
	}) {
		t.Errorf("did not contain 8")
	}

	if q.Has(func(elem interface{}) bool {
		return elem.(int) == 12
	}) {
		t.Errorf("shouldn't contain 12")
	}

	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestQueue_PushUnique(t *testing.T) {
	var q = NewQueue()

	existFunc := func(i int) func(elem interface{}) bool {
		return func(elem interface{}) bool {
			return elem.(int) == i
		}
	}

	for i := 0; i < 10; i++ {
		ok, err := q.PushUnique(i, existFunc(i))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("could not add %d", i)
		}
	}

	for i := 0; i < 10; i++ {
		ok, err := q.PushUnique(i, existFunc(i))
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Errorf("was able to add %d", i)
		}
	}

	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
}

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

func BenchmarkQueue_Push_Next_Pop(b *testing.B) {
	var q = NewQueue()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < b.N; i++ {
			_ = q.Push(struct{}{})
		}
		_ = q.Close()
		wg.Done()
	}(wg)

	for i := 0; i < 1; i++ {
		go func(wg *sync.WaitGroup) {
			var v interface{}
			for q.Next() {
				v, _ = q.Pop()
			}
			_ = v
			wg.Done()
		}(wg)
	}

	wg.Wait()
}
