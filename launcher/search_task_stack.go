/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

import ()

type searchTaskStack struct {
	tasks   []*searchTask
	items   map[string]*Item
	manager *Manager
}

func newSearchTaskStack(manager *Manager) *searchTaskStack {
	return &searchTaskStack{
		items:   manager.items,
		manager: manager,
	}
}

func (sts *searchTaskStack) Clear() {
	for _, task := range sts.tasks {
		task.isCanceled = true
	}
	sts.tasks = nil
}

func (sts *searchTaskStack) Pop() {
	// cancel top task
	top := sts.topTask()
	if top != nil {
		top.Cancel()
		logger.Debug("Pop", top)
		sts.tasks = sts.tasks[:len(sts.tasks)-1]
	}
}

func (sts *searchTaskStack) Push(c rune) {
	logger.Debugf("Push %c", c)
	prev := sts.topTask()
	task := newSearchTask(c, sts, prev)
	sts.tasks = append(sts.tasks, task)
	task.doSearch(prev)
}

func (sts *searchTaskStack) topTask() *searchTask {
	length := len(sts.tasks)
	if length == 0 {
		return nil
	}
	return sts.tasks[length-1]
}
