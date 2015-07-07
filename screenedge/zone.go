/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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

	leftTopEdge     = "TopLeft"
	leftBottomEdge  = "BottomLeft"
	rightTopEdge    = "TopRight"
	rightBottomEdge = "BottomRight"

	edgeActionWorkspace = "workspace"
)

const (
	leftTopDelay     int32 = 0
	leftBottomDelay        = 0
	rightTopDelay          = 0
	rightBottomDelay       = 500
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
	logger.Debug("TopLeft: ", topLeftArea)
	bottomLeftArea = areaRange{startX, endY - edgeDistance, startX + edgeDistance, endY}
	logger.Debug("BottomLeft: ", bottomLeftArea)
	topRightArea = areaRange{endX - edgeDistance, startY, endX, startY + edgeDistance}
	logger.Debug("TopRight: ", topRightArea)
	bottomRightArea = areaRange{endX - edgeDistance, endY - edgeDistance, endX, endY}
	logger.Debug("BottomRight: ", bottomRightArea)

	logger.Debug("topLeft: ", topLeftArea)
	logger.Debug("bottomLeft: ", bottomLeftArea)
	logger.Debug("topRight: ", topRightArea)
	logger.Debug("bottomRight: ", bottomRightArea)

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
	if action, ok := edgeActionMap[edge]; ok {
		strs := strings.Split(action, " ")
		l := len(strs)
		if l < 0 {
			return
		}

		argv := []string{}
		for i := 1; i < l; i++ {
			argv = append(argv, strs[i])
		}
		go exec.Command(strs[0], argv...).Run()
	}
}

func isInArea(x, y int32, area areaRange) bool {
	if x >= (area.X1) && x < (area.X2) &&
		y >= (area.Y1) && y < (area.Y2) {
		return true
	}

	return false
}

func getEdgeForCommand(cmd string) string {
	keys := zoneSettings().ListKeys()

	for _, key := range keys {
		switch key {
		case "left-up", "left-down", "right-up", "right-down":
			v := zoneSettings().GetString(key)
			if v == cmd {
				return key
			}
		}
	}

	return ""
}

func enableOneEdge(edge string) {
	switch edge {
	case "left-up":
		edgeActionMap["BottomLeft"] = ""
		edgeActionMap["TopRight"] = ""
		edgeActionMap["BottomRight"] = ""
	case "left-down":
		edgeActionMap["TopLeft"] = ""
		edgeActionMap["TopRight"] = ""
		edgeActionMap["BottomRight"] = ""
	case "right-up":
		edgeActionMap["TopLeft"] = ""
		edgeActionMap["BottomLeft"] = ""
		edgeActionMap["BottomRight"] = ""
	case "right-down":
		edgeActionMap["TopLeft"] = ""
		edgeActionMap["BottomLeft"] = ""
		edgeActionMap["TopRight"] = ""
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
	if timeout == 0 {
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
