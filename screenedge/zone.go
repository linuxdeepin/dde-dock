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
	"strings"
)

type areaRange struct {
	X1 int32
	Y1 int32
	X2 int32
	Y2 int32
}

const (
	DISTANCE         = int32(10)
	EDGE_TOPLEFT     = "TopLeft"
	EDGE_BOTTOMLEFT  = "BottomLeft"
	EDGE_TOPRIGHT    = "TopRight"
	EDGE_BOTTOMRIGHT = "BottomRight"
	ACTION_WORKSPACE = "workspace"
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

	topLeftArea = areaRange{startX, startY, startX + DISTANCE, startY + DISTANCE}
	logger.Debug("TopLeft: ", topLeftArea)
	bottomLeftArea = areaRange{startX, endY - DISTANCE, startX + DISTANCE, endY}
	logger.Debug("BottomLeft: ", bottomLeftArea)
	topRightArea = areaRange{endX - DISTANCE, startY, endX, startY + DISTANCE}
	logger.Debug("TopRight: ", topRightArea)
	bottomRightArea = areaRange{endX - DISTANCE, endY - DISTANCE, endX, endY}
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
	keys := zoneSettings.ListKeys()

	for _, key := range keys {
		v := zoneSettings.GetString(key)
		if v == cmd {
			return key
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

func newManager() *Manager {
	m := &Manager{}

	registerZoneArea()
	m.listenSignal()

	return m
}
