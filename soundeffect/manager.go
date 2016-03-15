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
	"gir/gio-2.0"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"sync"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	settingKeyEnabled = "enabled"

	DBusDest = "com.deepin.daemon.SoundEffect"
	dbusPath = "/com/deepin/daemon/SoundEffect"
	dbusIFC  = DBusDest
)

type Manager struct {
	Enabled *property.GSettingsBoolProperty `access:"readwrite"`
	setting *gio.Settings
	count   int
	mutex   sync.Mutex
}

func NewManager() *Manager {
	var m = new(Manager)

	m.setting = gio.NewSettings(soundEffectSchema)
	m.Enabled = property.NewGSettingsBoolProperty(m, "Enabled", m.setting, settingKeyEnabled)
	return m
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) PlaySystemSound(event string) {
	if len(event) == 0 {
		return
	}

	go func() {
		m.mutex.Lock()
		m.count++
		logger.Debug("start", m.count)
		m.mutex.Unlock()

		err := soundutils.PlaySystemSound(event, "", true)
		if err != nil {
			logger.Error(err)
		}

		m.mutex.Lock()
		logger.Debug("end", m.count)
		m.count--
		m.mutex.Unlock()
	}()
}
