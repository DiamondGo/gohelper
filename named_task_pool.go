/*
 * mastercoderk@gmail.com
 */

package gohelper

import "container/list"

type Task func()

type RunnerID interface {
	comparable
}

type NamedTaskPool[ID RunnerID] interface {
	Run(ID, Task) // block until your turn
	TryRun(ID, Task) bool
	RunUnblock(ID, Task)
}

func NewNamedTaskPool[ID RunnerID](
	totalConcurrentTask int,
	concurrentTaskForEach int,
) NamedTaskPool[ID] {
	return &namedTaskPool[ID]{
		totaActiveTaskAllowed:   totalConcurrentTask,
		runnerActiveTaskAllowed: concurrentTaskForEach,
		runningTasks:            make(map[ID]int),
		pendingTask:             list.New(),
	}
}

type namedTaskPool[ID RunnerID] struct {
	totaActiveTaskAllowed   int
	runnerActiveTaskAllowed int
	runningTasks            map[ID]int
	pendingTask             *list.List
}

func (p *namedTaskPool[ID]) Run(id ID, t Task) {
}

func (p *namedTaskPool[ID]) TryRun(id ID, t Task) bool {
	return false
}

func (p *namedTaskPool[ID]) RunUnblock(id ID, t Task) {
}
