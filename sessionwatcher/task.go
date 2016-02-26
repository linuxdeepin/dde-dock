/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	"time"
)

const (
	loopDuration       = time.Second * 10
	admissibleDuration = time.Second * 2

	maxLaunchTimes = 5
)

type taskInfo struct {
	Name  string
	Times int // continuous launch times

	failed        bool
	prevTimestamp int64 // previous launch timestamp

	isRunning func() bool
	launcher  func() error
}
type taskInfos []*taskInfo

func newTaskInfo(name string,
	isRunning func() bool, launcher func() error) *taskInfo {
	if isRunning == nil || launcher == nil {
		return nil
	}

	var task = &taskInfo{
		Name:          name,
		Times:         0,
		failed:        false,
		prevTimestamp: time.Now().Unix(),
		isRunning:     isRunning,
		launcher:      launcher,
	}

	return task
}

func (task *taskInfo) Launch() error {
	if !task.CanLaunch() {
		task.Times = 0
		return nil
	}

	duration := time.Now().Unix() - task.prevTimestamp
	if duration < int64(loopDuration+admissibleDuration) {
		task.Times += 1
	} else {
		task.Times = 0
	}

	if task.Times == maxLaunchTimes {
		task.failed = true
		logger.Debugf("Launch '%s' failed: over max launch times",
			task.Name)
	}

	task.prevTimestamp = time.Now().Unix()
	return task.launcher()
}

func (task *taskInfo) CanLaunch() bool {
	if task.failed {
		return false
	}

	return (task.isRunning() == false)
}

func (task *taskInfo) Over() bool {
	return task.failed
}
