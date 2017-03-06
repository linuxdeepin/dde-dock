/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package gesture

import (
	"dbus/com/deepin/daemon/gesture"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
}

var (
	gs       *gesture.Gesture
	_manager *gestureManager
	logger   = log.NewLogger("gesture")
)

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("gesture", daemon, logger)
	return daemon
}

func init() {
	loader.Register(NewDaemon())
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	var err error
	_manager, err = newGestureManager()
	if err != nil {
		logger.Error("Failed to initialize gesture manager:", err)
		return err
	}

	gs, err = gesture.NewGesture("com.deepin.daemon.Gesture",
		"/com/deepin/daemon/Gesture")
	if err != nil {
		logger.Error("Failed to initialize gesture object:", err)
		_manager = nil
		return err
	}

	_manager.handleGSettingsChanged()
	gs.ConnectEvent(func(name, direction string, fingers int32) {
		err := _manager.Exec(name, direction, fingers)
		if err != nil {
			logger.Error("Exec failed:", err)
		}
	})

	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	gesture.DestroyGesture(gs)
	gs = nil
	_manager = nil
	return nil
}
