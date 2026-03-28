# delayqueue

[![](https://godoc.org/github.com/jaz303/delayqueue?status.svg)](http://godoc.org/github.com/jaz303/delayqueue)

This is a simple delay queue that sends each added value to a channel after a specified delay has elapsed.
Compared to the naive approach of spawning one goroutine per item, this implementation uses a constant two goroutines, a single auxiliary timer, and a priority queue, so should therefore comfortably handle a large number of items.

The delay queue makes a best effort attempt to provide in-order delivery of messages, but buffering at the output channel can cause out-of-order delivery. See `queue_test.go` code comment for further discussion.

## Install

```shell
go get github.com/jaz303/delayqueue
```

## Usage

```go
// Create a new delay queue of ints that will run until
// the provided context is cancelled. The second argument
// is the buffer size of the outgoing channel.
queue := delayqueue.New[int](context.Background(), 0)

go func() {
    // Read items from the queue
    for i := range queue.C {
        log.Printf("Read from queue: %+v", i)
    }
}()

now := time.Now()

// Add items to the queue
queue.Add(now.Add(100 * time.Millisecond), 2)
queue.Add(now.Add(500 * time.Millisecond), 3)
queue.Add(now, 1)
```

# Benchmarks

`benchmark_test.go` compares `delayqueue` to a goroutine per item approach:

```
goos: darwin
goarch: arm64
pkg: github.com/jaz303/delayqueue
cpu: Apple M3 Max
BenchmarkDelayQueueAsync-16              1040580              1142 ns/op             144 B/op          4 allocs/op
BenchmarkDelayQueueAsyncBuffered-16      1000000              1127 ns/op             144 B/op          4 allocs/op
BenchmarkDelayQueueSync-16                966073              1171 ns/op             144 B/op          4 allocs/op
BenchmarkDelayQueueSyncBuffered-16       1000000              1082 ns/op             144 B/op          4 allocs/op
BenchmarkNaive-16                        1000000              1021 ns/op             818 B/op          3 allocs/op
```

`delayqueue` is about 11% slower, performs 33% more allocations (1 more), but does use 5x less memory.
