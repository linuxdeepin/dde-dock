package power

import (
	"time"
)

type TimeAfterTaskState uint

const (
	TimeAfterTaskStateReady TimeAfterTaskState = iota
	TimeAfterTaskStateRunning
	TimeAfterTaskStateDone
)

type TimeAfterTask struct {
	State      TimeAfterTaskState
	fn         func()
	cancelable bool
	timer      *time.Timer
}

func NewTimeAfterTask(delay time.Duration, fn func()) *TimeAfterTask {
	t := &TimeAfterTask{
		State:      TimeAfterTaskStateReady,
		cancelable: true,
	}
	t.timer = time.AfterFunc(delay, func() {
		t.cancelable = false
		t.State = TimeAfterTaskStateRunning
		fn()
		t.State = TimeAfterTaskStateDone
	})
	return t
}

func (t *TimeAfterTask) Cancel() {
	if t.cancelable {
		t.timer.Stop()
		t.State = TimeAfterTaskStateDone
	}
}

type TimeAfterTasks []*TimeAfterTask

func (tasks TimeAfterTasks) Wait(delay time.Duration, countMax int) {
	count := 0
	for {
		allDone := true
		for _, task := range tasks {
			if task.State != TimeAfterTaskStateDone {
				allDone = false
				break
			}
		}
		if allDone || count >= countMax {
			break
		}
		time.Sleep(delay)
		count++
	}
}

func (tasks TimeAfterTasks) CancelAll() {
	for _, task := range tasks {
		task.Cancel()
	}
}
