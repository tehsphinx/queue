package queue

import (
	"sync"
	"testing"
	"time"
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

func TestQueue_PopBlocking(t *testing.T) {
	var q = NewQueue()

	for i := 0; i < 10; i++ {
		err := q.Push(i)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		interf, ok := q.PopBlocking()
		if !ok {
			t.Error("did not get a value")
		}
		if interf.(int) != i {
			t.Errorf("expected %d got %v", i, interf)
		}
	}
}

func TestQueue_Wait(t *testing.T) {
	ch := make(chan time.Time)

	// check if a new time unblocks Wait
	q := NewQueue()
	go func() {
		q.Wait()
		t := time.Now()
		ch <- t
	}()

	time.Sleep(100 * time.Millisecond)
	if err := q.Push(1); err != nil {
		t.Fatal(err)
	}

	expectedTime := time.Now()
	gotTime := <-ch

	if expectedTime.After(gotTime) {
		if 1*time.Millisecond < expectedTime.Sub(gotTime) {
			t.Error("time should be closer")
		}
	}
	if gotTime.After(expectedTime) {
		if 1*time.Millisecond < gotTime.Sub(expectedTime) {
			t.Error("time should be closer")
		}
	}

	// check if closing unblocks Wait
	q = NewQueue()
	go func() {
		q.Wait()
		t := time.Now()
		ch <- t
	}()

	time.Sleep(100 * time.Millisecond)
	if err := q.Close(); err != nil {
		t.Fatal(err)
	}

	expectedTime = time.Now()
	gotTime = <-ch

	if expectedTime.After(gotTime) {
		if 1*time.Millisecond < expectedTime.Sub(gotTime) {
			t.Error("time should be closer", expectedTime.Sub(gotTime))
		}
	}
	if gotTime.After(expectedTime) {
		if 1*time.Millisecond < gotTime.Sub(expectedTime) {
			t.Error("time should be closer", gotTime.Sub(expectedTime))
		}
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
