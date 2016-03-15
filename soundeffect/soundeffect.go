/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package soundeffect

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/soundeffect")
var manager *Manager

func Start() error {
	if manager != nil {
		return nil
	}
	logger.BeginTracing()
	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return err
	}
	return nil
}

func IsPlaying() bool {
	return manager.count > 0
}

func Stop() error {
	if manager == nil {
		return nil
	}

	manager.setting.Unref()
	manager = nil
	logger.EndTracing()
	return nil
}
