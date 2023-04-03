/*
 * mastercoderk@gmail.com
 */

package gohelper

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	peroid = time.Second / 4
)

func TestTaskPoolBandwidthThrottled(t *testing.T) {
	pool := NewTaskPool[string](1, 3)

	var a, b time.Time
	task1 := func() {
		a = time.Now()
		time.Sleep(2 * peroid)
	}

	task2 := func() {
		b = time.Now()
		time.Sleep(2 * peroid)
	}

	pool.Run("a", task1)
	pool.Run("b", task2)

	pool.Join()
	diff := b.Sub(a).Abs()
	assert.Greater(t, diff, peroid)
}

func TestTaskPoolBandwidthEnough(t *testing.T) {
	pool := NewTaskPool[string](2, 3)

	var a, b time.Time
	task1 := func() {
		a = time.Now()
		time.Sleep(2 * peroid)
	}

	task2 := func() {
		b = time.Now()
		time.Sleep(2 * peroid)
	}

	pool.Run("a", task1)
	pool.Run("b", task2)

	pool.Join()
	diff := b.Sub(a).Abs()
	assert.Less(t, diff, peroid)
}

func TestTaskPoolRequestThrottled(t *testing.T) {
	pool := NewTaskPool[int](3, 1)
	var ts []time.Time
	lock := sync.Mutex{}

	task := func() {
		lock.Lock()
		ts = append(ts, time.Now())
		lock.Unlock()

		time.Sleep(2 * peroid)
	}

	pool.Run(1, task)
	pool.Run(1, task)

	pool.Join()

	assert.Equal(t, 2, len(ts))
	assert.Greater(t, ts[0].Sub(ts[1]).Abs(), peroid)
}

func TestTaskPoolRequestNotThrottled(t *testing.T) {
	pool := NewTaskPool[int](3, 2)
	var ts []time.Time
	lock := sync.Mutex{}

	task := func() {
		lock.Lock()
		ts = append(ts, time.Now())
		lock.Unlock()

		time.Sleep(2 * peroid)
	}

	pool.Run(1, task)
	pool.Run(1, task)

	pool.Join()

	assert.Equal(t, 2, len(ts))
	assert.Less(t, ts[0].Sub(ts[1]).Abs(), peroid)
}

func TestTaskPoolTryRun(t *testing.T) {
	pool := NewTaskPool[int](3, 1)

	task := func() {
		time.Sleep(peroid)
	}

	pool.Run(1, task)
	result := pool.TryRun(1, task)
	assert.False(t, result)

	pool.Join()
}

func TestTaskPoolRunBlock(t *testing.T) {
	pool := NewTaskPool[int](2, 2)

	t1 := time.Now()
	task := func() {
		time.Sleep(2 * peroid)
	}
	pool.BlockRun(1, task)
	t2 := time.Now()
	pool.Run(1, task)
	t3 := time.Now()

	pool.Join()
	assert.Greater(t, t2.Sub(t1).Abs(), peroid)
	assert.Less(t, t3.Sub(t2).Abs(), peroid)
}
