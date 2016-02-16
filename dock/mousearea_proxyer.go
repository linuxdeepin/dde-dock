/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"pkg.deepin.io/lib/dbus"
	"sync"
)

type coordinateRange struct {
	X0 int32
	Y0 int32
	X1 int32
	Y1 int32
}

type XMouseAreaInterface interface {
	ConnectCursorInto(func(int32, int32, string)) func()
	ConnectCursorOut(func(int32, int32, string)) func()
	UnregisterArea(string) error
	RegisterAreas(interface{}, int32) (string, error)
	RegisterFullScreen() (string, error)
}

// XMouseAreaProxyer为dde-api中XMouseAreaProxy接口的简单封装，用于触发隐藏dock的显示。
// 由于之前在C后端调用DBus不方便，因此特意实现了次接口，没有存在的意义，将会被废弃。
type XMouseAreaProxyer struct {
	lock    sync.RWMutex
	area    XMouseAreaInterface
	areaId  string
	idValid bool

	// InvalidId信号会在需要相应事件，但目前所持有的鼠标响应区域非法时被触发。
	InvalidId func()
}

func (a *XMouseAreaProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/XMouseAreaProxyer",
		Interface:  "dde.dock.XMouseAreaProxyer",
	}
}

func NewXMouseAreaProxyer(area XMouseAreaInterface, err error) (*XMouseAreaProxyer, error) {
	if err != nil {
		return nil, err
	}
	return &XMouseAreaProxyer{area: area, idValid: false}, nil
}

func (a *XMouseAreaProxyer) connectHandler(callback func(int32, int32, string)) func(x, y int32, id string) {
	return func(x, y int32, id string) {
		a.lock.Lock()
		if !a.idValid || id != a.areaId {
			if !a.idValid {
				dbus.Emit(a, "InvalidId")
			}
			logger.Debugf("valid: %v, event id: %v, areaId: %v", a.idValid, id, a.areaId)
			a.lock.Unlock()
			return
		}
		a.lock.Unlock()
		callback(x, y, id)
	}
}

func (a *XMouseAreaProxyer) connectMotionInto(callback func(int32, int32, string)) func() {
	return a.area.ConnectCursorInto(a.connectHandler(callback))
}

func (a *XMouseAreaProxyer) connectMotionOut(callback func(int32, int32, string)) func() {
	return a.area.ConnectCursorOut(a.connectHandler(callback))
}

func (a *XMouseAreaProxyer) unregister() {
	if a.idValid {
		a.area.UnregisterArea(a.areaId)
		a.idValid = false
	}
}

func (a *XMouseAreaProxyer) registerArea(registerHandler func() (string, error)) {
	a.lock.Lock()
	defer a.lock.Unlock()

	newAreaId, err := registerHandler()
	if err != nil {
		logger.Warning("register mousearea failed:", err)
		return
	}

	if a.areaId != newAreaId {
		a.unregister()
	}
	a.idValid = true
	a.areaId = newAreaId
}

// RegisterAreas注册多个区域为鼠标可响应区域。
// coordinateRange类型为{x0, y0, x1, y1}所表示的矩形区域。
// (x0, y0)表示矩形的左上角，(x1, y1)表示矩形的右下角。
func (a *XMouseAreaProxyer) RegisterAreas(areas []coordinateRange, eventMask int32) {
	a.registerArea(func() (string, error) {
		return a.area.RegisterAreas(areas, eventMask)
	})
}

// RegisterFullScreen将全屏注册为有效的鼠标响应区域。
func (a *XMouseAreaProxyer) RegisterFullScreen() {
	a.registerArea(a.area.RegisterFullScreen)
}

func (a *XMouseAreaProxyer) destroy() {
	dbus.UnInstallObject(a)
}
