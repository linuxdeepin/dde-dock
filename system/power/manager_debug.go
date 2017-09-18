// +build debug

/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"gir/gudev-1.0"
)

func (m *Manager) Debug(cmd string) {
	logger.Debug("Debug", cmd)
	switch cmd {
	case "init-batteries":
		devices := m.gudevClient.QueryBySubsystem("power_supply")
		for _, dev := range devices {
			m.addBattery(dev)
		}
		logger.Debug("initBatteries done")
		for _, dev := range devices {
			dev.Unref()
		}

	case "remove-all-batteries":
		var devices []*gudev.Device
		for _, bat := range m.batteries {
			devices = append(devices, bat.newDevice())
		}

		for _, dev := range devices {
			m.removeBattery(dev)
			dev.Unref()
		}

	case "destroy":
		m.destroy()

	default:
		logger.Warning("Command not support")
	}
}
