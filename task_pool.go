/*
 * mastercoderk@gmail.com
 */

package gohelper

import "sync"

type Task func()

type RequesterID interface {
	comparable
}

type TaskPool[ID RequesterID] interface {
	Run(ID, Task) // block until your turn
	TryRun(ID, Task) bool
	BlockRun(ID, Task)
	Join()
}

func NewTaskPool[ID RequesterID](
	totalConcurrentTask int,
	concurrentTaskForEach int,
) TaskPool[ID] {

	pool := &taskPool[ID]{
		totaActiveTaskAllowed:   totalConcurrentTask,
		runnerActiveTaskAllowed: concurrentTaskForEach,
		worker:                  make(chan struct{}, totalConcurrentTask),
		clientWorkers:           make(map[ID]chan struct{}),
		mapLock:                 sync.Mutex{},
	}

	for i := 0; i < totalConcurrentTask; i++ {
		pool.worker <- struct{}{}
	}
	return pool
}

type taskPool[ID RequesterID] struct {
	totaActiveTaskAllowed   int
	runnerActiveTaskAllowed int
	worker                  chan struct{}
	clientWorkers           map[ID]chan struct{}
	mapLock                 sync.Mutex
}

func (p *taskPool[ID]) Run(id ID, t Task) {
	ch := p.getRequestWorkers(id)
	permit := <-ch

	labour := <-p.worker

	go func() {
		defer func() {
			p.worker <- labour
			ch <- permit
		}()
		defer func() {
			_ = recover()
		}()
		t()
	}()
}

func (p *taskPool[ID]) TryRun(id ID, t Task) bool {
	ch := p.getRequestWorkers(id)
	var permit struct{}
	select {
	case permit = <-ch:
		// oops
	default:
		return false
	}

	var labour struct{}
	select {
	case labour = <-p.worker:
	default:
		ch <- permit
		return false
	}

	go func() {
		defer func() {
			p.worker <- labour
			ch <- permit
		}()
		defer func() {
			_ = recover()
		}()
		t()
	}()

	return true
}

func (p *taskPool[ID]) BlockRun(id ID, t Task) {
	ch := p.getRequestWorkers(id)
	permit := <-ch
	defer func() {
		ch <- permit
	}()

	labour := <-p.worker
	defer func() {
		p.worker <- labour
	}()

	defer func() {
		_ = recover()
	}()
	t()
}

func (p *taskPool[ID]) getRequestWorkers(id ID) chan struct{} {
	p.mapLock.Lock()
	defer p.mapLock.Unlock()

	ch, exists := p.clientWorkers[id]
	if !exists {
		ch = make(chan struct{}, p.runnerActiveTaskAllowed)
		for i := 0; i < p.runnerActiveTaskAllowed; i++ {
			ch <- struct{}{}
		}
		p.clientWorkers[id] = ch
	}

	return ch
}

func (p *taskPool[ID]) Join() {
	for i := 0; i < p.totaActiveTaskAllowed; i++ {
		<-p.worker
	}
}
