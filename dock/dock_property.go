package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"sync"
)

type DockProperty struct {
	heightLock sync.RWMutex
	Height     int32

	panelLock  sync.RWMutex
	PanelWidth int32
}

func NewDockProperty() *DockProperty {
	return &DockProperty{}
}

func (e *DockProperty) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.Property",
		"/dde/dock/Property",
		"dde.dock.Property",
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
