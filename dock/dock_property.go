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

// DockProperty存储dock前端界面相关的一些属性，包括dock的高度以及底板的宽度。
type DockProperty struct {
	dockManager *DockManager
	Height      int32
	panelLock   sync.RWMutex
	// PanelWidth是前端dock底板的宽度。
	PanelWidth int32
}

func NewDockProperty(dockManager *DockManager) *DockProperty {
	return &DockProperty{
		dockManager: dockManager,
		PanelWidth:  int32(dockManager.dockWidth),
		Height:      int32(dockManager.dockHeight),
	}
}

func (e *DockProperty) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "dde.dock.Property",
		ObjectPath: "/dde/dock/Property",
		Interface:  "dde.dock.Property",
	}
}

func (p *DockProperty) updateDockRect() {
	rect := p.dockManager.dockRect
	p.Height = int32(rect.Height())
	p.PanelWidth = int32(rect.Width())
}

// SetPanelWidth由前端界面调用，为后端设置底板的宽度。
func (p *DockProperty) SetPanelWidth(width int32) int32 {
	p.panelLock.Lock()
	defer p.panelLock.Unlock()

	if p.dockManager.dockWidth != int(width) {
		p.dockManager.dockWidth = int(width)
		p.dockManager.updateDockRect()
	}
	return width
}

func (p *DockProperty) destroy() {
	dbus.UnInstallObject(p)
}
