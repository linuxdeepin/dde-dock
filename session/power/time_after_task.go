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
