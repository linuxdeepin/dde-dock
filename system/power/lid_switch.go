/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"pkg.deepin.io/lib/arch"
	"pkg.deepin.io/lib/dbus"
)

func (m *Manager) initLidSwitch() {
	if arch.Get() == arch.Sunway {
		m.initLidSwitchSW()
	} else {
		m.initLidSwitchCommon()
	}
	logger.Debug("hasLidSwitch:", m.HasLidSwitch)
}

func (m *Manager) handleLidSwitchEvent(closed bool) {
	if closed {
		logger.Info("Lid Closed")
		dbus.Emit(m, "LidClosed")
	} else {
		logger.Info("Lid Opened")
		dbus.Emit(m, "LidOpened")
	}
}
