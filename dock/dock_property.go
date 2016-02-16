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
	heightLock sync.RWMutex
	// Height是前端dock的高度。
	Height int32

	panelLock sync.RWMutex
	// PanelWidth是前端dock底板的宽度。
	PanelWidth int32
}

func NewDockProperty() *DockProperty {
	return &DockProperty{}
}

func (e *DockProperty) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "dde.dock.Property",
		ObjectPath: "/dde/dock/Property",
		Interface:  "dde.dock.Property",
	}
}

func (p *DockProperty) updateDockHeight(mode DisplayModeType) int32 {
	p.heightLock.Lock()
	defer p.heightLock.Unlock()
	switch mode {
	case DisplayModeModernMode:
		p.Height = 68
		return p.Height
	case DisplayModeEfficientMode:
		p.Height = 48
		return p.Height
	case DisplayModeClassicMode:
		p.Height = 32
		return p.Height
	}

	return 0
}

// SetPanelWidth由前端界面调用，为后端设置底板的宽度。
func (p *DockProperty) SetPanelWidth(width int32) int32 {
	p.panelLock.Lock()
	defer p.panelLock.Unlock()
	if p.PanelWidth != width {
		p.PanelWidth = width
	}
	return p.PanelWidth
}

func (p *DockProperty) destroy() {
	dbus.UnInstallObject(p)
}
