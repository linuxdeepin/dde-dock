/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/lib/dbus"
)

type dbusWatcher struct {
	isSystemBus bool
	dbusObj     *dbus.Conn
	sigChan     <-chan *dbus.Signal
	callbacks   []func(*dbus.Signal)
}

func newDbusWatcher(isSystemBus bool) (dw *dbusWatcher) {
	dw = &dbusWatcher{isSystemBus: isSystemBus}
	var err error
	if dw.isSystemBus {
		dw.dbusObj, err = dbus.SystemBus()
	} else {
		dw.dbusObj, err = dbus.SessionBus()
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
	dw.dbusObj.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, expression)
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
	dw.sigChan = dw.dbusObj.Signal()
	go func() {
		for s := range dw.sigChan {
			for _, cb := range dw.callbacks {
				cb(s)
			}
		}
	}()
}

func (dw *dbusWatcher) stop() {
	dw.dbusObj.DetachSignal(dw.sigChan)
}
