package delayqueue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const epsilon = 5 * time.Millisecond

func TestElementsYieldedInTimeOrder(t *testing.T) {
	q := New[int](context.Background(), 3)
	assert.Nil(t, q.Add(time.Now().Add(100*time.Millisecond), 2))
	assert.Nil(t, q.Add(time.Now().Add(200*time.Millisecond), 3))
	assert.Nil(t, q.Add(time.Now().Add(50*time.Millisecond), 1))

	collect := []int{}
	for i := 0; i < 3; i++ {
		collect = append(collect, <-q.C)
	}

	assert.Equal(t, 1, collect[0])
	assert.Equal(t, 2, collect[1])
	assert.Equal(t, 3, collect[2])
}

func TestElementsYieldedAtCorrectTime(t *testing.T) {
	now := time.Now()

	vals := []int{1, 2, 3, 4, 5}
	dues := []time.Time{
		now.Add(100 * time.Millisecond),
		now.Add(300 * time.Millisecond),
		now.Add(600 * time.Millisecond),
		now.Add(650 * time.Millisecond),
		now.Add(1 * time.Second),
	}

	q := New[int](context.Background(), 5)
	for i := range vals {
		assert.Nil(t, q.Add(dues[i], vals[i]))
	}

	for i := range vals {
		v := <-q.C
		assert.Equal(t, vals[i], v)
		dt := time.Since(dues[i])
		if dt < 0 {
			dt = -dt
		}
		assert.True(t, dt < epsilon,
			"Time between scheduled deadline and channel received time was %s, more than expected %s", dt, epsilon)
	}
}

func TestResort(t *testing.T) {
	q := New[int](context.Background(), 0)
	now := time.Now()
	assert.Nil(t, q.Add(now.Add(100*time.Millisecond), 2))
	assert.Nil(t, q.Add(now.Add(200*time.Millisecond), 4))

	time.Sleep(500 * time.Millisecond)
	assert.Nil(t, q.Add(now.Add(150*time.Millisecond), 3))
	assert.Nil(t, q.Add(now.Add(80*time.Millisecond), 1))

	time.Sleep(500 * time.Millisecond)
	collect := []int{}
	for i := 0; i < 4; i++ {
		collect = append(collect, <-q.C)
	}

	// Note: demonstrates the limitations of in order delivery.
	// Nothing can be done about the 100ms entry arriving first; it's already been offered to the
	// channel by the time new messages with earlier deadlines are added. However, for the
	// queue of "ready" items just waiting to be received on channel, we make a best effort attempt
	// to offer those items in order.
	//
	// This also demonstrates how creating a delay queue with a non-zero channel buffer allows
	// for a greater chance of out-of-order delivery. The greater the channel buffer size,
	// the relative higher chance of out-of-order all else being equal (whether there's any out
	// of order delivery depends on a number of factors in the producer and consumer pattern).
	assert.Equal(t, 2, collect[0])
	assert.Equal(t, 1, collect[1])
	assert.Equal(t, 3, collect[2])
	assert.Equal(t, 4, collect[3])
	assert.Equal(t, 2, q.Resorts())
}
