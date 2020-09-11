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
	dbus "github.com/godbus/dbus"
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
type TouchType int32
type TouchDirection int32

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

	TouchTypeRightButton = TouchType(C.TOUCH_TYPE_RIGHT_BUTTON)

	ButtonTypeDown = TouchType(C.BUTTON_TYPE_DOWN)
	ButtonTypeUp   = TouchType(C.BUTTON_TYPE_UP)

	//handleTouchScreenEvent
	TouchDirectionNone  = TouchDirection(C.DIR_NONE)
	TouchDirectionTop   = TouchDirection(C.DIR_TOP)
	TouchDirectionRight = TouchDirection(C.DIR_RIGHT)
	TouchDirectionBot   = TouchDirection(C.DIR_BOT)
	TouchDirectionLeft  = TouchDirection(C.DIR_LEFT)
	TouchTypeNone       = TouchType(C.GT_NONE)
	TouchTypeTap        = TouchType(C.GT_TAP)
	TouchTypeMovement   = TouchType(C.GT_MOVEMENT)
	TouchTypeEdge       = TouchType(C.GT_EDGE)
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

func (t TouchType) String() string {
	switch t {
	case TouchTypeRightButton:
		return "touch right button"
	case ButtonTypeDown:
		return "down"
	case ButtonTypeUp:
		return "up"
	case TouchTypeNone:
		return "touch none"
	case TouchTypeTap:
		return "touch tap"
	case TouchTypeMovement:
		return "touch movement"
	case TouchTypeEdge:
		return "touch edge"

	}
	return "Unknown"
}

func (t TouchDirection) String() string {
	switch t {
	case TouchDirectionNone:
		return "none"
	case TouchDirectionTop:
		return "top"
	case TouchDirectionRight:
		return "right"
	case TouchDirectionBot:
		return "bot"
	case TouchDirectionLeft:
		return "left"
	}
	return "Unknown"
}

type Manager struct {
	service *dbusutil.Service

	// nolint
	methods *struct {
		SetShortPressDuration   func() `in:"duration"`
		SetEdgeMoveStopDuration func() `in:"duration"`
	}

	// nolint
	signals *struct {
		Event struct {
			name      string
			direction string
			fingers   int32
		}

		//gesture double click down
		DbclickDown struct {
			fingers int32
		}

		//gesture swipe moving info
		SwipeMoving struct {
			fingers        int32
			accelX, accelY float64
		}

		//gesture swipe stop or interrupted
		SwipeStop struct {
			fingers int32
		}

		TouchEdgeEvent struct {
			direction      string
			scaleX, scaleY float64
		}

		TouchSinglePressTimeout struct {
			time           int32
			scaleX, scaleY float64
		}

		TouchPressTimeout struct {
			fingers, time  int32
			scaleX, scaleY float64
		}

		TouchUpOrCancel struct {
			scaleX, scaleY float64
		}

		TouchEdgeMoveStop struct {
			direction      string
			scaleX, scaleY float64
			duration       int
		}

		TouchEdgeMoveStopLeave struct {
			direction      string
			scaleX, scaleY float64
			duration       int
		}

		TouchMoving struct {
			scalex, scaley float64
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

//duration unit ms
func (*Manager) SetShortPressDuration(duration int) *dbus.Error {
	C.set_timer_short_duration(C.int(duration))
	return nil
}

//duration unit ms
func (*Manager) SetEdgeMoveStopDuration(duration int) *dbus.Error {
	C.set_edge_move_stop_time(C.int(duration))
	return nil
}

//touchpad gesture
//export handleGestureEvent
func handleGestureEvent(ty, direction, fingers C.int) {
	err := _m.service.Emit(_m, "Event", GestureType(ty).String(),
		GestureType(direction).String(), int32(fingers))
	if err != nil {
		logger.Error("handleGestureEvent failed:", err)
	}
}

//export handleDbclickDown
func handleDbclickDown(fingers C.int) {
	err := _m.service.Emit(_m, "DbclickDown", int32(fingers))
	if err != nil {
		logger.Error("handleDbclickDown failed:", err)
	}
}

//export handleSwipeMoving
func handleSwipeMoving(fingers C.int, accelX, accelY C.double) {
	logger.Debug("emit SwipeMoving:", float64(accelX), float64(accelY))
	err := _m.service.Emit(_m, "SwipeMoving", int32(fingers), float64(accelX), float64(accelY))
	if err != nil {
		logger.Error("handleSwipeMoving failed:", err)
	}
}

//export handleSwipeStop
func handleSwipeStop(fingers C.int) {
	logger.Debug("emit SwipeStop:", int32(fingers))
	err := _m.service.Emit(_m, "SwipeStop", int32(fingers))
	if err != nil {
		logger.Error("handleSwipeStop failed:", err)
	}
}

//touchscreen gesture
//export handleTouchEvent
func handleTouchEvent(ty, btn C.int) {
	err := _m.service.Emit(_m, "Event", TouchType(ty).String(),
		TouchType(btn).String(), 0)
	if err != nil {
		logger.Error("handleTouchEvent failed:", err)
	}
}

//export handleTouchScreenEvent
func handleTouchScreenEvent(ty, direction, fingers C.int, scaleX, scaleY C.double) {
	if int(ty) == int(C.get_edge_type()) {
		err := _m.service.Emit(_m, "TouchEdgeEvent", TouchDirection(direction).String(), float64(scaleX), float64(scaleY))
		if err != nil {
			logger.Error("handleTouchScreenEvent failed:", err)
		}
	}
}

//export handleTouchEdgeMoveStop
func handleTouchEdgeMoveStop(direction C.int, x, y C.double, duration C.int) {
	err := _m.service.Emit(_m, "TouchEdgeMoveStop", TouchDirection(direction).String(), float64(x), float64(y), int(duration))
	if err != nil {
		logger.Error("handleTouchEdgeMoveStop failed:", err)
	}
}

//export handleTouchEdgeMoveStopLeave
func handleTouchEdgeMoveStopLeave(direction C.int, x, y C.double, duration C.int) {
	err := _m.service.Emit(_m, "TouchEdgeMoveStopLeave", TouchDirection(direction).String(), float64(x), float64(y), int(duration))
	if err != nil {
		logger.Error("handleTouchEdgeMoveStopLeave failed:", err)
	}
}

//export handleTouchMoving
func handleTouchMoving(scalex, scaley C.double) {
	err := _m.service.Emit(_m, "TouchMoving", float64(scalex), float64(scaley))
	if err != nil {
		logger.Error("handleTouchMoving failed:", err)
	}
}

//export handleTouchShortPress
func handleTouchShortPress(time C.int, scalex, scaley C.double) {
	err := _m.service.Emit(_m, "TouchSinglePressTimeout", int32(time), float64(scalex), float64(scaley))
	if err != nil {
		logger.Error("handleTouchShortPress failed:", err)
	}
}

//export handleTouchPressTimeout
func handleTouchPressTimeout(fingers, time C.int, scalex, scaley C.double) {
	err := _m.service.Emit(_m, "TouchPressTimeout", int32(fingers), int(time), float64(scalex), float64(scaley))
	if err != nil {
		logger.Error("handleTouchPressTimeout", err)
	}
}

//export handleTouchUpOrCancel
func handleTouchUpOrCancel(scalex, scaley C.double) {
	err := _m.service.Emit(_m, "TouchUpOrCancel", float64(scalex), float64(scaley))
	if err != nil {
		logger.Error("handleTouchUpOrCancel failed:", err)
	}
}

func (*Daemon) Start() error {
	logger.BeginTracing()
	logger.Info("start gesture daemon")
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

	conf, err := loadConfig(getConfigPath())
	if err != nil {
		logger.Warning("Failed to load gesture config:", err)
		conf = &Config{}
	}
	go C.start_loop(C.int(conf.Verbose), C.double(conf.LongPressDistance))
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

func (*Daemon) SetLongPressDuration(duration int) {
	C.set_timer_duration(C.int(duration))
}
