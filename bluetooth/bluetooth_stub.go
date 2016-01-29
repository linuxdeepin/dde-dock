/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
)

func (b *Bluetooth) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	logger.Debug("OnPropertiesChanged()", name)
}

func (b *Bluetooth) DebugInfo() (info string) {
	info = fmt.Sprintf("adapters: %s\ndevices: %s", marshalJSON(b.adapters), marshalJSON(b.devices))
	return
}

func (b *Bluetooth) setPropAdapters() {
	b.Adapters = marshalJSON(b.adapters)
	dbus.NotifyChange(b, "Adapters")
}

func (b *Bluetooth) setPropDevices() {
	b.Devices = marshalJSON(b.devices)
	dbus.NotifyChange(b, "Devices")
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

func (b *Bluetooth) clearUnpairedDevice() {
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
