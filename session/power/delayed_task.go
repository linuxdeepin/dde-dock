/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package power

import (
	"sync"
	"time"
)

type delayedTaskState uint

const (
	delayedTaskStateReady delayedTaskState = iota
	delayedTaskStateRunning
	delayedTaskStateDone
)

type delayedTask struct {
	state      delayedTaskState
	mu         sync.Mutex
	cancelable bool
	timer      *time.Timer
	name       string
}

func newDelayedTask(name string, delay time.Duration, fn func()) *delayedTask {
	t := &delayedTask{
		state:      delayedTaskStateReady,
		cancelable: true,
		name:       name,
	}
	t.timer = time.AfterFunc(delay, func() {
		t.mu.Lock()
		t.cancelable = false
		t.state = delayedTaskStateRunning
		t.mu.Unlock()

		fn()

		t.mu.Lock()
		t.state = delayedTaskStateDone
		t.mu.Unlock()
	})
	return t
}

func (t *delayedTask) Cancel() {
	t.mu.Lock()
	if t.cancelable {
		t.timer.Stop()
		t.state = delayedTaskStateDone
		logger.Debugf("delayedTask %s cancelled", t.name)
	}
	t.mu.Unlock()
}

type delayedTasks []*delayedTask

func (tasks delayedTasks) Wait(delay time.Duration, countMax int) {
	count := 0
	for {
		allDone := true
		for _, task := range tasks {
			task.mu.Lock()
			state := task.state
			task.mu.Unlock()

			if state != delayedTaskStateDone {
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

func (tasks delayedTasks) CancelAll() {
	for _, task := range tasks {
		task.Cancel()
	}
}
