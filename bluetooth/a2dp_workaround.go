/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

func (b *Bluetooth) disconnectA2DPDeviceExcept(d *device) {
	for _, devices := range b.devices {
		for _, device := range devices {
			if device.Path == d.Path {
				continue
			}
			for _, uuid := range device.UUIDs {
				if uuid == A2DP_SINK_UUID {
					bluezDisconnectDevice(device.Path)
				}
			}
		}
	}
}
