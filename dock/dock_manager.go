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
	"dbus/com/deepin/daemon/display"
	"errors"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xrect"
	"pkg.deepin.io/lib/dbus"
)

type DockManager struct {
	hideStateManager *HideStateManager
	setting          *Setting
	dockProperty     *DockProperty
	entryManager     *EntryManager
	clientManager    *ClientManager

	// 共用部分
	hideMode    HideModeType
	displayMode DisplayModeType

	activeWindow       xproto.Window
	dpy                *display.Display
	displayPrimaryRect *xrect.XRect
	dockRect           *xrect.XRect

	dockHeight int
	dockWidth  int
}

func NewDockManager() (*DockManager, error) {
	m := new(DockManager)
	err := m.init()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *DockManager) initEntryManager() error {
	m.entryManager = NewEntryManager()
	m.entryManager.initRuntimeApps()
	m.entryManager.initDockedApps()
	err := dbus.InstallOnSession(m.entryManager.dockedAppManager)
	if err != nil {
		return err
	}

	err = dbus.InstallOnSession(m.entryManager)
	return err
}

func (m *DockManager) initHideStateManager() error {
	m.hideStateManager = NewHideStateManager()
	m.hideStateManager.dockRect = m.dockRect
	logger.Debug("initHideStateManager dockRect", m.hideStateManager.dockRect)
	m.hideStateManager.mode = m.hideMode
	m.hideStateManager.initHideState()
	err := dbus.InstallOnSession(m.hideStateManager)
	return err
}

func (m *DockManager) initDockProperty() error {
	m.dockProperty = NewDockProperty(m)
	err := dbus.InstallOnSession(m.dockProperty)
	return err
}

func (m *DockManager) initSetting() error {
	m.setting = NewSetting(m)
	if m.setting == nil {
		return errors.New("create setting failed")
	}
	err := dbus.InstallOnSession(m.setting)

	logger.Debug("init display and hide mode")
	m.displayMode = DisplayModeType(m.setting.core.GetEnum(DisplayModeKey))
	m.hideMode = HideModeType(m.setting.core.GetEnum(HideModeKey))
	m.dockHeight = getDockHeightByDisplayMode(m.displayMode)
	return err
}

func (m *DockManager) init() error {
	var err error

	err = m.initSetting()
	if err != nil {
		return err
	}
	logger.Info("initialize setting done")

	// ensure init display after init setting
	err = m.initDisplay()
	if err != nil {
		return err
	}
	logger.Info("initialize display done")

	err = m.initEntryManager()
	if err != nil {
		return err
	}
	logger.Info("initialize entry proxyer manager done")

	err = m.initHideStateManager()
	if err != nil {
		return err
	}
	logger.Info("initialize hide state manager done")

	err = m.initDockProperty()
	if err != nil {
		return err
	}
	logger.Info("initialize dock property done")

	m.setting.listenSettingsChanged()
	logger.Info("initialize settings done")

	// init client manager
	m.clientManager = NewClientManager()
	err = dbus.InstallOnSession(m.clientManager)
	if err != nil {
		return err
	}
	logger.Info("initialize client manager done")
	return nil
}

func (m *DockManager) destroy() {
	if m.dockProperty != nil {
		m.dockProperty.destroy()
		m.dockProperty = nil
	}

	if m.setting != nil {
		m.setting.destroy()
		m.setting = nil
	}

	if m.hideStateManager != nil {
		m.hideStateManager.destroy()
		m.hideStateManager = nil
	}

	if m.dpy != nil {
		display.DestroyDisplay(m.dpy)
		m.dpy = nil
	}

}

func getDockHeightByDisplayMode(mode DisplayModeType) int {
	switch mode {
	case DisplayModeModernMode:
		return 68
	case DisplayModeEfficientMode:
		return 48
	case DisplayModeClassicMode:
		return 32
	default:
		return 0
	}
}

func (m *DockManager) updateDockRect() {
	// calc dock rect
	primaryX, primaryY, primaryW, primaryH := m.displayPrimaryRect.Pieces()
	dockX := primaryX + (primaryW-m.dockWidth)/2
	dockY := primaryY + primaryH - m.dockHeight

	if m.dockRect == nil {
		m.dockRect = xrect.New(dockX, dockY, m.dockWidth, m.dockHeight)
	} else {
		m.dockRect.XSet(dockX)
		m.dockRect.YSet(dockY)
		m.dockRect.WidthSet(m.dockWidth)
		m.dockRect.HeightSet(m.dockHeight)
	}

	logger.Debug("primary rect:", m.displayPrimaryRect)
	logger.Debug("dock width:", m.dockWidth)
	logger.Debug("dock height:", m.dockHeight)
	logger.Debug("updateDockRect dock rect:", m.dockRect)

	if m.hideStateManager != nil {
		m.hideStateManager.updateStateWithoutDelay()
	} else {
		logger.Debug("m.hideStateManager is nil")
	}

	if m.dockProperty != nil {
		m.dockProperty.updateDockRect()
	} else {
		logger.Debug("m.dockProperty is nil")
	}
}
