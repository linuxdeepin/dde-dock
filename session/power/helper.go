/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	libdisplay "dbus/com/deepin/daemon/display"
	libsessionwatcher "dbus/com/deepin/daemon/sessionwatcher"
	liblockfront "dbus/com/deepin/dde/lockfront"
	libsessionmanager "dbus/com/deepin/sessionmanager"
	libpower "dbus/com/deepin/system/power"
	liblogin1 "dbus/org/freedesktop/login1"
	libscreensaver "dbus/org/freedesktop/screensaver"
	"pkg.deepin.io/lib/notify"

	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgbutil"
)

type Helper struct {
	Power          *libpower.Power
	Notification   *notify.Notification
	SessionManager *libsessionmanager.SessionManager
	SessionWatcher *libsessionwatcher.SessionWatcher
	ScreenSaver    *libscreensaver.ScreenSaver
	Display        *libdisplay.Display
	LockFront      *liblockfront.LockFront
	Login1Manager  *liblogin1.Manager

	xu *xgbutil.XUtil
}

func NewHelper() (*Helper, error) {
	h := &Helper{}
	err := h.init()
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Helper) init() error {
	var err error
	h.Power, err = libpower.NewPower("com.deepin.system.Power", "/com/deepin/system/Power")
	if err != nil {
		logger.Warning("init Power failed:", err)
		return err
	}

	notify.Init(dbusDest)
	h.Notification = notify.NewNotification("", "", "")

	h.SessionManager, err = libsessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager")
	if err != nil {
		logger.Warning("init SessionManager failed:", err)
		return err
	}

	h.ScreenSaver, err = libscreensaver.NewScreenSaver("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
	if err != nil {
		logger.Warning("init ScreenSaver failed:", err)
		return err
	}

	h.Display, err = libdisplay.NewDisplay(dbusDisplayDest, dbusDisplayPath)
	if err != nil {
		logger.Warning("init Display failed:", err)
		return err
	}

	h.LockFront, err = liblockfront.NewLockFront("com.deepin.dde.lockFront", "/com/deepin/dde/lockFront")
	if err != nil {
		logger.Warning("init LockFront failed:", err)
		return err
	}

	h.SessionWatcher, err = libsessionwatcher.NewSessionWatcher("com.deepin.daemon.SessionWatcher", "/com/deepin/daemon/SessionWatcher")
	if err != nil {
		logger.Warning("init SessionWatcher failed:", err)
		return err
	}

	h.Login1Manager, err = liblogin1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		logger.Warning("init login1 manager failed:", err)
		return err
	}

	// init X conn
	h.xu, err = xgbutil.NewConn()
	if err != nil {
		return err
	}
	dpms.Init(h.xu.Conn())
	return nil
}

func (h *Helper) Destroy() {
	if h.Power != nil {
		libpower.DestroyPower(h.Power)
		h.Power = nil
	}

	if h.Notification != nil {
		h.Notification.Destroy()
		h.Notification = nil
		notify.Destroy()
	}

	if h.SessionManager != nil {
		libsessionmanager.DestroySessionManager(h.SessionManager)
		h.SessionManager = nil
	}

	if h.ScreenSaver != nil {
		libscreensaver.DestroyScreenSaver(h.ScreenSaver)
		h.ScreenSaver = nil
	}

	if h.Display != nil {
		libdisplay.DestroyDisplay(h.Display)
		h.Display = nil
	}

	if h.LockFront != nil {
		h.LockFront = nil
	}

	if h.SessionWatcher != nil {
		libsessionwatcher.DestroySessionWatcher(h.SessionWatcher)
		h.SessionWatcher = nil
	}

	if h.Login1Manager != nil {
		liblogin1.DestroyManager(h.Login1Manager)
		h.Login1Manager = nil
	}

	// NOTE: Don't close x conn, because the bug of lib xgbutil.
	// [xgbutil] eventloop.go:27: BUG: Could not read an event or an error.
	if h.xu != nil {
		//h.xu.Conn().Close()
		h.xu = nil
	}
}
