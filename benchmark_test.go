package delayqueue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func BenchmarkDelayQueueAsync(b *testing.B) {
	b.ReportAllocs()
	queue := New[int](context.Background(), 0)

	now := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			<-queue.C
		}
		wg.Done()
	}()

	for i := 0; i < b.N; i++ {
		queue.Add(now, i)
	}

	wg.Wait()
}

func BenchmarkDelayQueueAsyncBuffered(b *testing.B) {
	b.ReportAllocs()
	queue := New[int](context.Background(), 16)

	now := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			<-queue.C
		}
		wg.Done()
	}()

	for i := 0; i < b.N; i++ {
		queue.Add(now, i)
	}

	wg.Wait()
}

func BenchmarkDelayQueueSync(b *testing.B) {
	b.ReportAllocs()
	queue := New[int](context.Background(), 0)

	now := time.Now()

	for i := 0; i < b.N; i++ {
		queue.Add(now, i)
	}

	for i := 0; i < b.N; i++ {
		<-queue.C
	}
}

func BenchmarkDelayQueueSyncBuffered(b *testing.B) {
	b.ReportAllocs()
	queue := New[int](context.Background(), 16)

	now := time.Now()

	for i := 0; i < b.N; i++ {
		queue.Add(now, i)
	}

	for i := 0; i < b.N; i++ {
		<-queue.C
	}
}

func BenchmarkNaive(b *testing.B) {
	b.ReportAllocs()
	results := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		time.AfterFunc(time.Duration(0), func() {
			results <- i
			wg.Done()
		})
	}

	for i := 0; i < b.N; i++ {
		<-results
	}

	wg.Wait()
}
