/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"time"
)

type areaRange struct {
	X1 int32
	Y1 int32
	X2 int32
	Y2 int32
}

const (
	edgeDistance int32 = 5

	leftTopEdge     = "left-up"
	leftBottomEdge  = "left-down"
	rightTopEdge    = "right-up"
	rightBottomEdge = "right-down"

	edgeActionWorkspace = "workspace"
)

var (
	areaId string

	topLeftArea     areaRange
	bottomLeftArea  areaRange
	topRightArea    areaRange
	bottomRightArea areaRange
)

func registerZoneArea() {
	var (
		startX int32
		startY int32
		endX   int32
		endY   int32
	)

	rect := dspObj.PrimaryRect.Get()
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
	logger.Debugf("PrimaryRect: %d, %d, %d, %d\n", startX, endX, startY, endY)

	if endX <= startX || endY <= startY {
		return
	}

	topLeftArea = areaRange{startX, startY, startX + edgeDistance, startY + edgeDistance}
	logger.Debug("left-up: ", topLeftArea)
	bottomLeftArea = areaRange{startX, endY - edgeDistance, startX + edgeDistance, endY}
	logger.Debug("left-down: ", bottomLeftArea)
	topRightArea = areaRange{endX - edgeDistance, startY, endX, startY + edgeDistance}
	logger.Debug("right-up: ", topRightArea)
	bottomRightArea = areaRange{endX - edgeDistance, endY - edgeDistance, endX, endY}
	logger.Debug("right-down: ", bottomRightArea)

	var err error
	areaId, err = areaObj.RegisterAreas([]areaRange{
		topLeftArea,
		bottomLeftArea,
		topRightArea,
		bottomRightArea,
	}, 0)
	if err != nil {
		logger.Warning("Register area failed: ", err)
		return
	}

	logger.Debug("MouseArea Id: ", areaId)
}

func unregisterZoneArea() {
	areaObj.UnregisterArea(areaId)
}

func execEdgeAction(edge string) {
	action, ok := edgeActionMap[edge]
	if !ok {
		return
	}

	go exec.Command("/bin/sh", "-c", action).Run()
}

func isInArea(x, y int32, area areaRange) bool {
	if x >= (area.X1) && x < (area.X2) &&
		y >= (area.Y1) && y < (area.Y2) {
		return true
	}

	return false
}

func getEdgeForCommand(cmd string) []string {
	var ret []string

	keys := zoneSettings().ListKeys()
	for _, key := range keys {
		switch key {
		case "left-up", "left-down", "right-up", "right-down":
			v := zoneSettings().GetString(key)
			if strings.Contains(v, cmd) {
				ret = append(ret, key)
				continue
			}
		}
	}

	return ret
}

func enableSpecialEdges(edges []string) {
	mutex.Lock()
	defer mutex.Unlock()

	for k, _ := range edgeActionMap {
		if isItemInList(k, edges) {
			continue
		}
		edgeActionMap[k] = ""
	}
}

func (m *Manager) destroy() {
	unregisterZoneArea()
	finalizeDBusIFC()
	dbus.UnInstallObject(m)
}

func newManager() *Manager {
	m := &Manager{}

	m.lTopTimer = &edgeTimer{}
	m.lBottomTimer = &edgeTimer{}
	m.rTopTimer = &edgeTimer{}
	m.rBottomTimer = &edgeTimer{}

	registerZoneArea()
	m.listenSignal()

	return m
}

type edgeTimer struct {
	timer *time.Timer
}

func (eTimer *edgeTimer) DoAction(edge string, timeout int32) {
	action, ok := edgeActionMap[edge]
	if !ok {
		return
	}

	logger.Debug("Exec action:", edge, action)
	if !strings.Contains(action, "com.deepin.dde.ControlCenter") {
		execEdgeAction(edge)
		return
	}

	eTimer.timer = time.NewTimer(time.Millisecond *
		time.Duration(timeout))
	go func() {
		if eTimer.timer == nil {
			return
		}

		<-eTimer.timer.C
		execEdgeAction(edge)
		eTimer.timer = nil
	}()
}

func (eTimer *edgeTimer) StopTimer() {
	if eTimer.timer == nil {
		return
	}

	eTimer.timer.Stop()
	eTimer.timer = nil
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
