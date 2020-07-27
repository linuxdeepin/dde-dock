/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"syscall"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	login1 "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	"pkg.deepin.io/dde/daemon/appearance"
	"pkg.deepin.io/dde/daemon/bluetooth"
	"pkg.deepin.io/dde/daemon/network"
)

type sleepInhibitor struct {
	loginManager *login1.Manager
	fd           int

	OnWakeup        func()
	OnBeforeSuspend func()
}

func newSleepInhibitor(login1Manager *login1.Manager, daemon *daemon.Daemon) *sleepInhibitor {
	inhibitor := &sleepInhibitor{
		loginManager: login1Manager,
		fd:           -1,
	}

	_, err := daemon.ConnectHandleForSleep(func(before bool) {
		logger.Info("login1 HandleForSleep", before)
		// signal `HandleForSleep` true -> false
		if !_manager.sessionActive {
			//如果此用户此时不是活跃状态,则不处理待机唤醒信号.
			return
		}

		if before {
			// TODO(jouyouyun): implement 'HandleForSleep' register
			appearance.HandlePrepareForSleep(true)
			network.HandlePrepareForSleep(true)
			bluetooth.HandlePrepareForSleep(true)
			if inhibitor.OnBeforeSuspend != nil {
				inhibitor.OnBeforeSuspend()
			}
			suspendPulseSinks(1)
			suspendPulseSources(1)
			err := inhibitor.unblock()
			if err != nil {
				logger.Warning(err)
			}
		} else {
			suspendPulseSources(0)
			suspendPulseSinks(0)
			if inhibitor.OnWakeup != nil {
				inhibitor.OnWakeup()
			}
			network.HandlePrepareForSleep(false)
			bluetooth.HandlePrepareForSleep(false)
			appearance.HandlePrepareForSleep(false)
			err := inhibitor.block()
			if err != nil {
				logger.Warning(err)
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	return inhibitor
}

func (inhibitor *sleepInhibitor) block() error {
	logger.Debug("block sleep")
	if inhibitor.fd != -1 {
		return nil
	}
	fd, err := inhibitor.loginManager.Inhibit(0,
		"sleep", dbusServiceName, "run screen lock", "delay")
	if err != nil {
		return err
	}
	inhibitor.fd = int(fd)
	return nil
}

func (inhibitor *sleepInhibitor) unblock() error {
	if inhibitor.fd == -1 {
		return nil
	}
	logger.Debug("unblock sleep")
	err := syscall.Close(inhibitor.fd)
	inhibitor.fd = -1
	if err != nil {
		logger.Warning("failed to close fd:", err)
		return err
	}
	return nil
}
