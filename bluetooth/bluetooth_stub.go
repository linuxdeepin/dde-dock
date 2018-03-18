/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package bluetooth

import (
	"fmt"

	"pkg.deepin.io/lib/dbus1"
)

func (b *Bluetooth) DebugInfo() (string, *dbus.Error) {
	info := fmt.Sprintf("adapters: %s\ndevices: %s", marshalJSON(b.adapters), marshalJSON(b.devices))
	return info, nil
}

func (b *Bluetooth) updateState() {
	newState := StateUnavailable
	if len(b.adapters) > 0 {
		newState = StateAvailable
	}

	for _, devices := range b.devices {
		for _, d := range devices {
			if d.connected {
				newState = StateConnected
				break
			}
		}
	}

	b.PropsMu.Lock()
	b.setPropState(uint32(newState))
	b.PropsMu.Unlock()
}

//ClearUnpairedDevice will remove all device in unpaired list
func (b *Bluetooth) ClearUnpairedDevice() *dbus.Error {
	var removeDevices []*device
	for _, devices := range b.devices {
		for _, d := range devices {
			if !d.Paired {
				logger.Info("remove unpaired device", d)
				removeDevices = append(removeDevices, d)
			}
		}
	}

	for _, d := range removeDevices {
		bluezRemoveDevice(d.AdapterPath, d.Path)
	}
	return nil
}
