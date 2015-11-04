package bluetooth

func (b *Bluetooth) disconnectA2DPDevice() {
	for _, devices := range b.devices {
		for _, device := range devices {
			for _, uuid := range device.UUIDs {
				if uuid == A2DP_SINK_UUID {
					bluezDisconnectDevice(device.Path)
				}
			}
		}
	}
}
