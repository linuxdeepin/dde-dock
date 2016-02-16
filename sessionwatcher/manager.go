/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import "time"

type Manager struct {
	taskList *taskInfos
	quit     chan struct{}
}

func newManager() *Manager {
	var m = new(Manager)
	m.quit = make(chan struct{})
	m.taskList = new(taskInfos)
	return m
}

func (m *Manager) AddTask(task *taskInfo) {
	if m.IsTaskExist(task.Name) {
		logger.Debugf("Task '%s' has exist", task.Name)
		return
	}

	*m.taskList = append(*m.taskList, task)
}

func (m *Manager) IsTaskExist(name string) bool {
	for _, task := range *m.taskList {
		if name == task.Name {
			return true
		}
	}
	return false
}

func (m *Manager) HasRunning() bool {
	for _, task := range *m.taskList {
		if !task.Over() {
			return true
		}
	}
	return false
}

func (m *Manager) LaunchAll() {
	for _, task := range *m.taskList {
		err := task.Launch()
		if err != nil {
			logger.Warningf("Launch '%s' failed: %v",
				task.Name, err)
		}
	}
}

func (m *Manager) StartLoop() {
	for {
		select {
		case <-m.quit:
			return
		case <-time.After(loopDuration):
			if !m.HasRunning() {
				logger.Debug("All program has launched failure")
				m.QuitLoop()
				return
			}

			m.LaunchAll()
		}
	}
}

func (m *Manager) QuitLoop() {
	if m.quit == nil {
		return
	}
	close(m.quit)
	m.quit = nil
}
