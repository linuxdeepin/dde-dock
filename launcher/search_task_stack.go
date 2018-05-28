/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package launcher

import (
	"sync"
)

type searchTaskStack struct {
	tasks   []*searchTask
	items   map[string]*Item
	manager *Manager
	mu      sync.Mutex
}

func newSearchTaskStack(manager *Manager) *searchTaskStack {
	return &searchTaskStack{
		items:   manager.items,
		manager: manager,
	}
}

func (sts *searchTaskStack) Clear() {
	sts.mu.Lock()

	for _, task := range sts.tasks {
		task.Cancel()
	}
	sts.tasks = nil

	sts.mu.Unlock()
}

func (sts *searchTaskStack) Pop() {
	sts.mu.Lock()

	// cancel top task
	top := sts.topTask()
	if top != nil {
		top.Cancel()
		logger.Debug("Pop", top)
		sts.tasks = sts.tasks[:len(sts.tasks)-1]
	}

	sts.mu.Unlock()
}

func (sts *searchTaskStack) Push(c rune) {
	sts.mu.Lock()

	logger.Debugf("Push %c", c)
	prev := sts.topTask()
	task := newSearchTask(c, sts, prev)
	sts.tasks = append(sts.tasks, task)

	sts.mu.Unlock()

	task.search(prev)
}

func (sts *searchTaskStack) topTask() *searchTask {
	length := len(sts.tasks)
	if length == 0 {
		return nil
	}
	return sts.tasks[length-1]
}

func (sts *searchTaskStack) indexOf(task *searchTask) int {
	for idx, t := range sts.tasks {
		if t == task {
			return idx
		}
	}
	return -1
}

func (sts *searchTaskStack) GetNext(task *searchTask) *searchTask {
	sts.mu.Lock()
	defer sts.mu.Unlock()

	idx := sts.indexOf(task)
	if idx == -1 {
		return nil
	}
	idx++
	if idx < len(sts.tasks) {
		return sts.tasks[idx]
	}
	return nil
}
