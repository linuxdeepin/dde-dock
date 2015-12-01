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
