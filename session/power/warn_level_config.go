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

package power

import (
	"time"

	gio "pkg.deepin.io/gir/gio-2.0"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/gsettings"
)

type warnLevelConfig struct {
	UsePercentageForPolicy bool
	LowTime                uint64
	DangerTime             uint64
	CriticalTime           uint64
	ActionTime             uint64
	LowPercentage          float64
	DangerPercentage       float64
	CriticalPercentage     float64
	ActionPercentage       float64
}

func (c *warnLevelConfig) isValid() bool {
	if c.LowTime > c.DangerTime &&
		c.DangerTime > c.CriticalTime &&
		c.CriticalTime > c.ActionTime &&

		c.LowPercentage > c.DangerPercentage &&
		c.DangerPercentage > c.CriticalPercentage &&
		c.CriticalPercentage > c.ActionPercentage {
		return true
	}
	return false
}

type WarnLevelConfigManager struct {
	UsePercentageForPolicy gsprop.Bool `prop:"access:rw"`

	LowTime      gsprop.Int `prop:"access:rw"`
	DangerTime   gsprop.Int `prop:"access:rw"`
	CriticalTime gsprop.Int `prop:"access:rw"`
	ActionTime   gsprop.Int `prop:"access:rw"`

	LowPercentage      gsprop.Int `prop:"access:rw"`
	DangerPercentage   gsprop.Int `prop:"access:rw"`
	CriticalPercentage gsprop.Int `prop:"access:rw"`
	ActionPercentage   gsprop.Int `prop:"access:rw"`

	settings    *gio.Settings
	changeTimer *time.Timer
	changeCb    func()
}

func NewWarnLevelConfigManager(gs *gio.Settings) *WarnLevelConfigManager {

	m := &WarnLevelConfigManager{
		settings: gs,
	}

	m.UsePercentageForPolicy.Bind(gs, settingKeyUsePercentageForPolicy)
	m.LowTime.Bind(gs, settingKeyLowTime)
	m.DangerTime.Bind(gs, settingKeyDangerTime)
	m.CriticalTime.Bind(gs, settingKeyCriticalTime)
	m.ActionTime.Bind(gs, settingKeyActionTime)

	m.LowPercentage.Bind(gs, settingKeyLowPercentage)
	m.DangerPercentage.Bind(gs, settingKeyDangerlPercentage)
	m.CriticalPercentage.Bind(gs, settingKeyCriticalPercentage)
	m.ActionPercentage.Bind(gs, settingKeyActionPercentage)

	m.connectSettingsChanged()
	return m
}

func (m *WarnLevelConfigManager) getWarnLevelConfig() *warnLevelConfig {
	return &warnLevelConfig{
		UsePercentageForPolicy: m.UsePercentageForPolicy.Get(),
		LowTime:                uint64(m.LowTime.Get()),
		DangerTime:             uint64(m.DangerTime.Get()),
		CriticalTime:           uint64(m.CriticalTime.Get()),
		ActionTime:             uint64(m.ActionTime.Get()),

		LowPercentage:      float64(m.LowPercentage.Get()),
		DangerPercentage:   float64(m.DangerPercentage.Get()),
		CriticalPercentage: float64(m.CriticalPercentage.Get()),
		ActionPercentage:   float64(m.ActionPercentage.Get()),
	}
}

func (m *WarnLevelConfigManager) setChangeCallback(fn func()) {
	m.changeCb = fn
}

func (m *WarnLevelConfigManager) delayCheckValid() {
	logger.Debug("delayCheckValid")
	if m.changeTimer != nil {
		m.changeTimer.Stop()
	}
	m.changeTimer = time.AfterFunc(20*time.Second, func() {
		logger.Debug("checkValid")
		wlc := m.getWarnLevelConfig()
		if !wlc.isValid() {
			logger.Info("Warn level config is invalid, reset")
			err := m.Reset()
			if err != nil {
				logger.Warning(err)
			}
		}
	})
}

func (m *WarnLevelConfigManager) notifyChange() {
	if m.changeCb != nil {
		logger.Debug("WarnLevelConfig change")
		m.changeCb()
	}
	m.delayCheckValid()
}

func (m *WarnLevelConfigManager) connectSettingsChanged() {
	gsettings.ConnectChanged(gsSchemaPower, "*", func(key string) {
		switch key {
		case settingKeyUsePercentageForPolicy,
			settingKeyLowPercentage,
			settingKeyDangerlPercentage,
			settingKeyCriticalPercentage,
			settingKeyActionPercentage,

			settingKeyLowTime,
			settingKeyDangerTime,
			settingKeyCriticalTime,
			settingKeyActionTime:

			logger.Debug("key changed", key)
			m.notifyChange()
		}
	})

}

func (m *WarnLevelConfigManager) Reset() *dbus.Error {
	s := m.settings
	s.Reset(settingKeyUsePercentageForPolicy)
	s.Reset(settingKeyLowPercentage)
	s.Reset(settingKeyDangerlPercentage)
	s.Reset(settingKeyCriticalPercentage)
	s.Reset(settingKeyActionPercentage)
	s.Reset(settingKeyLowTime)
	s.Reset(settingKeyDangerTime)
	s.Reset(settingKeyCriticalTime)
	s.Reset(settingKeyActionTime)
	return nil
}

func (*WarnLevelConfigManager) GetInterfaceName() string {
	return dbusInterface + ".WarnLevelConfig"
}
