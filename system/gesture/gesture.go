/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package gesture

// #cgo pkg-config: libinput glib-2.0
// #cgo LDFLAGS: -ludev -lm
// #include <stdlib.h>
// #include "core.h"
import "C"

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.daemon.Gesture"
	dbusPath        = "/com/deepin/daemon/Gesture"
	dbusInterface   = "com.deepin.daemon.Gesture"
)

type GestureType int32

var (
	GestureTypeSwipe = GestureType(C.GESTURE_TYPE_SWIPE)
	GestureTypePinch = GestureType(C.GESTURE_TYPE_PINCH)
	GestureTypeTap   = GestureType(C.GESTURE_TYPE_TAP)

	GestureDirectionNone  = GestureType(C.GESTURE_DIRECTION_NONE)
	GestureDirectionUp    = GestureType(C.GESTURE_DIRECTION_UP)
	GestureDirectionDown  = GestureType(C.GESTURE_DIRECTION_DOWN)
	GestureDirectionLeft  = GestureType(C.GESTURE_DIRECTION_LEFT)
	GestureDirectionRight = GestureType(C.GESTURE_DIRECTION_RIGHT)
	GestureDirectionIn    = GestureType(C.GESTURE_DIRECTION_IN)
	GestureDirectionOut   = GestureType(C.GESTURE_DIRECTION_OUT)
)

func (t GestureType) String() string {
	switch t {
	case GestureTypeSwipe:
		return "swipe"
	case GestureTypePinch:
		return "pinch"
	case GestureTypeTap:
		return "tap"
	case GestureDirectionNone:
		return "none"
	case GestureDirectionUp:
		return "up"
	case GestureDirectionDown:
		return "down"
	case GestureDirectionLeft:
		return "left"
	case GestureDirectionRight:
		return "right"
	case GestureDirectionIn:
		return "in"
	case GestureDirectionOut:
		return "out"
	}
	return "Unknown"
}

type Manager struct {
	service *dbusutil.Service
	signals *struct {
		Event struct {
			name      string
			direction string
			fingers   int32
		}
	}
}

var (
	_m     *Manager
	logger = log.NewLogger(dbusServiceName)
)

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewDaemon())
}

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("gesture", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

//export handleGestureEvent
func handleGestureEvent(ty, direction, fingers C.int) {
	logger.Debug("Emit gesture event:", GestureType(ty).String(),
		GestureType(direction).String(),
		int32(fingers))
	_m.service.Emit(_m, "Event", GestureType(ty).String(),
		GestureType(direction).String(),
		int32(fingers))
}

func (*Daemon) Start() error {
	logger.BeginTracing()
	logger.Info("Start gesture daemon")
	service := loader.GetService()
	_m = &Manager{
		service: service,
	}
	err := service.Export(dbusPath, _m)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	// TODO: debug level
	go C.start_loop()
	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}
	C.quit_loop()
	service := loader.GetService()
	err := service.StopExport(_m)
	if err != nil {
		return err
	}

	_m = nil
	return nil
}
