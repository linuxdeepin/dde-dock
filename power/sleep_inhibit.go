/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	liblogin1 "dbus/org/freedesktop/login1"
	"syscall"
)

func init() {
	submoduleList = append(submoduleList, newSleepInhibitor)
}

type sleepInhibitor struct {
	fd      int
	login1  *liblogin1.Manager
	manager *Manager
	what    string
	who     string
	why     string
	mode    string
}

func newSleepInhibitor(m *Manager) (string, submodule, error) {
	name := "SleepInhibitor"
	inhibitor := &sleepInhibitor{
		fd:      -1,
		manager: m,
		what:    "sleep",
		who:     "lock screen",
		why:     "run screenlock",
		mode:    "delay",
	}
	var err error
	inhibitor.login1, err = liblogin1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		return name, nil, err
	}
	err = inhibitor.block()
	if err != nil {
		return name, nil, err
	}
	return name, inhibitor, nil
}

// 处理待机之前，唤醒时事件
func (inhibitor *sleepInhibitor) Start() error {
	// signal `PrepareForSleep` true -> false
	m := inhibitor.manager
	inhibitor.login1.ConnectPrepareForSleep(func(before bool) {
		if before {
			// sleep is blocked
			m.handleBeforeSuspend()
			inhibitor.unblock()
		} else {
			m.handleWeakup()
			inhibitor.block()
		}
	})
	return nil
}

func (inhibitor *sleepInhibitor) Destroy() {
	// close fd
	inhibitor.unblock()
	if inhibitor.login1 != nil {
		liblogin1.DestroyManager(inhibitor.login1)
		inhibitor.login1 = nil
	}
}

func (inhibitor *sleepInhibitor) block() error {
	logger.Debug("Block", inhibitor.what)
	if inhibitor.fd != -1 {
		return nil
	}
	fd, err := inhibitor.login1.Inhibit(
		inhibitor.what, inhibitor.who, inhibitor.why, inhibitor.mode)
	if err != nil {
		logger.Error("inbhibit failed", err)
		return err
	}
	inhibitor.fd = int(fd)
	return nil
}

func (inhibitor *sleepInhibitor) unblock() error {
	if inhibitor.fd == -1 {
		return nil
	}
	logger.Debug("Unblock", inhibitor.what)
	err := syscall.Close(inhibitor.fd)
	inhibitor.fd = -1
	if err != nil {
		logger.Error("close fd error:", err)
		return err
	}
	return nil
}
