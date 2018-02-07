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
	"pkg.deepin.io/lib/dbus"
)

func (b *Bluetooth) DebugInfo() (info string) {
	info = fmt.Sprintf("adapters: %s\ndevices: %s", marshalJSON(b.adapters), marshalJSON(b.devices))
	return
}

func (b *Bluetooth) setPropState() {
	b.State = StateUnavailable
	if len(b.adapters) > 0 {
		b.State = StateAvailable
	}
	for _, devs := range b.devices {
		for _, d := range devs {
			if d.connected {
				b.State = StateConnected
				break
			}
		}
	}
	dbus.NotifyChange(b, "State")
}

//ClearUnpairedDevice will remove all device in unpaired list
func (b *Bluetooth) ClearUnpairedDevice() {
	removeDevices := [](*device){}
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
}
