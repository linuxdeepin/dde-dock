/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package sessionwatcher

import (
	"dbus/org/freedesktop/dbus"
	"os/exec"
	"pkg.deepin.io/lib/log"
	"time"
)

const (
	dockDest   = "com.deepin.dde.dock"
	dockLaunch = "dde-dock"

	dockAppletsDest   = "dde.dock.entry.AppletManager"
	dockAppletsLaunch = "dde-dock-applets"

	maxDuration   = time.Second * 5
	deltaDuration = time.Second * 1
)

var logger = log.NewLogger("daemon/sessionwatcher")

type Manager struct {
	quit chan struct{}

	launchDockFailed    bool
	launchAppletsFailed bool
}

func NewManager() *Manager {
	var m = Manager{}
	m.quit = make(chan struct{})
	return &m
}

func (m *Manager) canLaunchDock() bool {
	if m.launchDockFailed {
		return false
	}

	exist, err := isDBusDestExist(dockDest)
	if err != nil {
		logger.Debugf("Check '%s' exist failed: %v", dockDest, err)
	}
	if exist {
		return false
	}
	return true
}

func (m *Manager) canLaunchDockApplets() bool {
	if m.launchAppletsFailed {
		return false
	}

	exist, err := isDBusDestExist(dockAppletsDest)
	if err != nil {
		logger.Debugf("Check '%s' exist failed: %v", dockAppletsDest, err)
	}
	if exist {
		return false
	}
	return true
}

func (m *Manager) restartDock() {
	err := doAction("killall", []string{dockLaunch})
	if err != nil {
		logger.Debugf("killall '%s' failed: %v", dockLaunch, err)
	}

	err = doLaunchCommand(dockLaunch, nil)
	if err != nil {
		m.launchDockFailed = true
		logger.Warningf("Launch '%s' failed: %v", dockLaunch, err)
		return
	}
	logger.Debug("Restart dde-dock over")

	return
}

func (m *Manager) restartDockApplets() {
	err := doAction("killall", []string{dockAppletsLaunch})
	if err != nil {
		logger.Debugf("killall '%s' failed: %v", dockAppletsLaunch, err)
	}

	err = doLaunchCommand(dockAppletsLaunch, nil)
	if err != nil {
		m.launchAppletsFailed = true
		logger.Warningf("Launch '%s' failed: %v", dockAppletsLaunch, err)
		return
	}
	logger.Debug("Restart dde-dock-applets over")
	return
}

func doLaunchCommand(cmd string, args []string) error {
	var (
		err          error
		waitDuration = time.Second * 0
	)
	for waitDuration < maxDuration {
		err = doAction(cmd, args)
		if err == nil {
			return nil
		}

		waitDuration += deltaDuration
		<-time.After(waitDuration)
	}
	return err
}

func doAction(cmd string, args []string) error {
	// Run() block, why?
	//return exec.Command("/bin/sh", "-c", cmd).Run()
	return exec.Command(cmd, args...).Start()
}

func isDBusDestExist(dest string) (bool, error) {
	daemon, err := dbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return false, err
	}
	defer dbus.DestroyDBusDaemon(daemon)

	names, err := daemon.ListNames()
	if err != nil {
		return false, err
	}
	return isItemInList(dest, names), nil
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
