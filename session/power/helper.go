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
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/notify"

	// system bus
	libpower "github.com/linuxdeepin/go-dbus-factory/com.deepin.system.power"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"

	// session bus
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.sessionwatcher"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.screensaver"
	"github.com/linuxdeepin/go-x11-client"
)

type Helper struct {
	Notification *notify.Notification

	Power        *libpower.Power // sig
	LoginManager *login1.Manager // sig

	SessionManager *sessionmanager.SessionManager
	SessionWatcher *sessionwatcher.SessionWatcher
	ScreenSaver    *screensaver.ScreenSaver // sig
	Display        *display.Display

	xConn *x.Conn
}

func newHelper(systemConn, sessionConn *dbus.Conn) (*Helper, error) {
	h := &Helper{}
	err := h.init(systemConn, sessionConn)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Helper) init(systemConn, sessionConn *dbus.Conn) error {
	var err error

	notify.Init(dbusServiceName)
	h.Notification = notify.NewNotification("", "", "")

	h.Power = libpower.NewPower(systemConn)
	h.LoginManager = login1.NewManager(systemConn)
	h.SessionManager = sessionmanager.NewSessionManager(sessionConn)
	h.ScreenSaver = screensaver.NewScreenSaver(sessionConn)
	h.Display = display.NewDisplay(sessionConn)
	h.SessionWatcher = sessionwatcher.NewSessionWatcher(sessionConn)

	// init X conn
	h.xConn, err = x.NewConn()
	if err != nil {
		return err
	}
	return nil
}

func (h *Helper) Destroy() {
	h.Power.RemoveHandler(proxy.RemoveAllHandlers)
	h.LoginManager.RemoveHandler(proxy.RemoveAllHandlers)
	h.ScreenSaver.RemoveHandler(proxy.RemoveAllHandlers)

	if h.Notification != nil {
		h.Notification.Destroy()
		h.Notification = nil
		notify.Destroy()
	}

	if h.xConn != nil {
		h.xConn.Close()
		h.xConn = nil
	}
}
