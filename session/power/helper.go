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
	"github.com/godbus/dbus"
	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	"pkg.deepin.io/lib/dbusutil/proxy"

	// system bus
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	libpower "github.com/linuxdeepin/go-dbus-factory/com.deepin.system.power"
	"github.com/linuxdeepin/go-dbus-factory/net.hadess.sensorproxy"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"

	// session bus
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.sessionwatcher"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.screensaver"
	"github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbusutil"
)

type Helper struct {
	Notifications *notifications.Notifications

	Power         *libpower.Power // sig
	LoginManager  *login1.Manager // sig
	SensorProxy   *sensorproxy.SensorProxy
	SysDBusDaemon *ofdbus.DBus
	Daemon        *daemon.Daemon

	SessionManager *sessionmanager.SessionManager
	SessionWatcher *sessionwatcher.SessionWatcher
	ScreenSaver    *screensaver.ScreenSaver // sig
	Display        *display.Display

	xConn *x.Conn
}

func newHelper(systemBus, sessionBus *dbus.Conn) (*Helper, error) {
	h := &Helper{}
	err := h.init(systemBus, sessionBus)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Helper) init(sysBus, sessionBus *dbus.Conn) error {
	var err error

	h.Notifications = notifications.NewNotifications(sessionBus)

	h.Power = libpower.NewPower(sysBus)
	h.LoginManager = login1.NewManager(sysBus)
	h.SensorProxy = sensorproxy.NewSensorProxy(sysBus)
	h.SysDBusDaemon = ofdbus.NewDBus(sysBus)
	h.Daemon = daemon.NewDaemon(sysBus)
	h.SessionManager = sessionmanager.NewSessionManager(sessionBus)
	h.ScreenSaver = screensaver.NewScreenSaver(sessionBus)
	h.Display = display.NewDisplay(sessionBus)
	h.SessionWatcher = sessionwatcher.NewSessionWatcher(sessionBus)

	// init X conn
	h.xConn, err = x.NewConn()
	if err != nil {
		return err
	}
	return nil
}

func (h *Helper) initSignalExt(systemSigLoop, sessionSigLoop *dbusutil.SignalLoop) {
	// sys
	h.SysDBusDaemon.InitSignalExt(systemSigLoop, true)
	h.LoginManager.InitSignalExt(systemSigLoop, true)
	h.Power.InitSignalExt(systemSigLoop, true)
	h.SensorProxy.InitSignalExt(systemSigLoop, true)
	h.Daemon.InitSignalExt(systemSigLoop, true)
	// session
	h.ScreenSaver.InitSignalExt(sessionSigLoop, true)
	h.SessionWatcher.InitSignalExt(sessionSigLoop, true)
	h.Display.InitSignalExt(sessionSigLoop, true)
}

func (h *Helper) Destroy() {
	h.SysDBusDaemon.RemoveHandler(proxy.RemoveAllHandlers)
	h.LoginManager.RemoveHandler(proxy.RemoveAllHandlers)
	h.Power.RemoveHandler(proxy.RemoveAllHandlers)
	h.SensorProxy.RemoveHandler(proxy.RemoveAllHandlers)
	h.Daemon.RemoveHandler(proxy.RemoveAllHandlers)

	h.ScreenSaver.RemoveHandler(proxy.RemoveAllHandlers)
	h.SessionWatcher.RemoveHandler(proxy.RemoveAllHandlers)

	if h.xConn != nil {
		h.xConn.Close()
		h.xConn = nil
	}
}
