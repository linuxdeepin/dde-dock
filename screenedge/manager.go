/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"dbus/com/deepin/api/xmousearea"
	"dbus/com/deepin/daemon/display"
	"fmt"
	"github.com/BurntSushi/xgbutil"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/initializer/v2"
	"sync"
)

type Manager struct {
	settings   *Settings
	edges      map[string]*edge
	display    *display.Display
	xmousearea *xmousearea.XMouseArea
	areaId     string
	mutex      sync.Mutex
	timer      *edgeTimer
	xu         *xgbutil.XUtil
}

func NewManager() (*Manager, error) {
	var m = new(Manager)
	m.timer = new(edgeTimer)
	m.edges = make(map[string]*edge)
	m.settings = NewSettings()
	err := initializer.Do(
		m.initXUtil).Do(
		m.initDBusIFC).Do(
		func() error {
			m.handleDBusSignal()
			m.handleSettingsChanged()
			return m.initEdges()
		}).Do(
		m.registerEdgeAreas).Do(
		func() error {
			return dbus.InstallOnSession(m)
		}).GetError()

	if err != nil {
		m.destroy()
		return nil, err
	}
	return m, nil
}

func (m *Manager) initXUtil() error {
	var err error
	m.xu, err = xgbutil.NewConn()
	return err
}

func (m *Manager) initDBusIFC() error {
	var err error
	m.display, err = display.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		return err
	}

	m.xmousearea, err = xmousearea.NewXMouseArea("com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea")
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) finalizeDBusIFC() {
	if m.display != nil {
		display.DestroyDisplay(m.display)
		m.display = nil
	}

	if m.xmousearea != nil {
		xmousearea.DestroyXMouseArea(m.xmousearea)
		m.xmousearea = nil
	}
}

func (m *Manager) destroy() {
	if m.xu != nil {
		m.xu.Conn().Close()
		m.xu = nil
	}

	m.unregisterEdgeAreas()
	m.finalizeDBusIFC()
	dbus.UnInstallObject(m)
}

func (m *Manager) getPrimaryArea() (*areaRange, error) {
	rect := m.display.PrimaryRect.Get()
	if len(rect) != 4 {
		return nil, fmt.Errorf("Length of display PrimaryRect is not 4, rect: %v", rect)
	}
	var (
		startX int32
		startY int32
		endX   int32
		endY   int32
	)
	if v, ok := rect[0].(int16); ok {
		startX = int32(v)
	}
	if v, ok := rect[1].(int16); ok {
		startY = int32(v)
	}
	if v, ok := rect[2].(uint16); ok {
		endX = int32(v) + startX
	}
	if v, ok := rect[3].(uint16); ok {
		endY = int32(v) + startY
	}
	primaryArea := &areaRange{startX, startY, endX, endY}
	logger.Debugf("PrimaryArea: %v", primaryArea)
	return primaryArea, nil
}

func (m *Manager) getEdgeAreas() []*areaRange {
	var list []*areaRange
	for _, v := range m.edges {
		list = append(list, v.Area)
	}
	return list
}

func (m *Manager) registerEdgeAreas() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.areaId) > 0 {
		return nil
	}

	var err error
	m.areaId, err = m.xmousearea.RegisterAreas(m.getEdgeAreas(), 0)
	if err != nil {
		logger.Error("Register area failed: ", err)
		return err
	}
	logger.Debug("Register areas: ", m.areaId)
	return nil
}

func (m *Manager) unregisterEdgeAreas() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.areaId) == 0 {
		return nil
	}
	logger.Debug("Unregister areas: ", m.areaId)
	err := m.xmousearea.UnregisterArea(m.areaId)
	if err != nil {
		logger.Error("UnregisterArea failed:", err)
		return err
	}
	m.areaId = ""
	return nil
}

func (m *Manager) initEdges() error {
	m.edges[TopLeft] = &edge{}
	m.edges[BottomLeft] = &edge{}
	m.edges[TopRight] = &edge{}
	m.edges[BottomRight] = &edge{}

	err := m.setEdgeAreas()
	if err != nil {
		return err
	}

	// set edge commands
	for k, edge := range m.edges {
		edge.Command = m.settings.GetEdgeAction(k)
	}
	return nil
}

func (m *Manager) setEdgeAreas() error {
	logger.Debug("setEdgeAreas")
	primaryArea, err := m.getPrimaryArea()
	if err != nil {
		logger.Error("getPrimaryArea failed:", err)
		return err
	}
	const sideLength int32 = 5
	corners := primaryArea.GetCornerSquares(sideLength)
	cornerNames := []string{TopLeft, TopRight, BottomRight, BottomLeft}
	for index, name := range cornerNames {
		m.edges[name].Area = corners[index]
	}
	return nil
}
