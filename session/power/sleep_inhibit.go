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

type sleepInhibitor struct {
	login1Manager *liblogin1.Manager
	fd            int
	what          string
	who           string
	why           string
	mode          string

	OnWakeup        func()
	OnBeforeSuspend func()
}

func newSleepInhibitor(login1Manager *liblogin1.Manager) *sleepInhibitor {
	inhibitor := &sleepInhibitor{
		login1Manager: login1Manager,
		fd:            -1,
		what:          "sleep",
		who:           "com.deepin.daemon.Power",
		why:           "run screenlock",
		mode:          "delay",
	}

	login1Manager.ConnectPrepareForSleep(func(before bool) {
		logger.Info("login1 PrepareForSleep", before)
		// signal `PrepareForSleep` true -> false
		if before {
			if inhibitor.OnBeforeSuspend != nil {
				inhibitor.OnBeforeSuspend()
			}
			inhibitor.unblock()
		} else {
			if inhibitor.OnWakeup != nil {
				inhibitor.OnWakeup()
			}
			inhibitor.block()
		}
	})
	return inhibitor
}

func (inhibitor *sleepInhibitor) block() error {
	logger.Debug("Block", inhibitor.what)
	if inhibitor.fd != -1 {
		return nil
	}
	fd, err := inhibitor.login1Manager.Inhibit(
		inhibitor.what, inhibitor.who, inhibitor.why, inhibitor.mode)
	if err != nil {
		logger.Error("inbhibit block failed:", err)
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
