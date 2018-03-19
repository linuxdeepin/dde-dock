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

package network

import (
	"pkg.deepin.io/lib/dbus1"
)

type dbusWatcher struct {
	isSystemBus bool
	dbusConn    *dbus.Conn
	sigChan     chan *dbus.Signal
	callbacks   []func(*dbus.Signal)
}

func newDbusWatcher(isSystemBus bool) (dw *dbusWatcher) {
	dw = &dbusWatcher{isSystemBus: isSystemBus}
	var err error
	if dw.isSystemBus {
		dw.dbusConn, err = dbus.SystemBus()
	} else {
		dw.dbusConn, err = dbus.SessionBus()
	}
	if err != nil {
		logger.Error(err)
		return
	}
	dw.start()
	return
}

func destroyDbusWatcher(dw *dbusWatcher) {
	dw.reset()
	dw.stop()
}

func (dw *dbusWatcher) watch(expression string) {
	dw.dbusConn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, expression)
}

func (dw *dbusWatcher) connect(cb func(*dbus.Signal)) {
	dw.callbacks = append(dw.callbacks, cb)
}

func (dw *dbusWatcher) reset() {
	dw.callbacks = nil
	dw.stop()
	dw.start()
}

func (dw *dbusWatcher) start() {
	dw.sigChan = make(chan *dbus.Signal, 10)
	dw.dbusConn.Signal(dw.sigChan)
	go func() {
		for s := range dw.sigChan {
			for _, cb := range dw.callbacks {
				cb(s)
			}
		}
	}()
}

func (dw *dbusWatcher) stop() {
	dw.dbusConn.RemoveSignal(dw.sigChan)
}
