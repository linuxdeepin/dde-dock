package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type AccessPoint struct {
	Ssid     string
	NeedKey  bool
	Strength uint8
	Path     dbus.ObjectPath
}

type Device struct {
	Path  dbus.ObjectPath
	State uint32
}

func NewDevice(core *nm.Device) *Device {
	return &Device{core.Path, core.State.Get()}
}

func NewAccessPoint(apPath dbus.ObjectPath) (ap AccessPoint, err error) {
	calcStrength := func(s uint8) uint8 {
		switch {
		case s <= 10:
			return 0
		case s <= 25:
			return 25
		case s <= 50:
			return 50
		case s <= 75:
			return 75
		case s <= 100:
			return 100
		}
		return 0
	}

	nmAp, err := nm.NewAccessPoint(NMDest, apPath)
	if err != nil {
		return
	}

	ap = AccessPoint{string(nmAp.Ssid.Get()),
		parseFlags(nmAp.Flags.Get(), nmAp.WpaFlags.Get(), nmAp.RsnFlags.Get()) != ApKeyNone,
		calcStrength(nmAp.Strength.Get()),
		nmAp.Path,
	}
	return
}

// DisconnectDevice will disconnect all connection in target device.
func (this *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	dev, err := nm.NewDevice(NMDest, devPath)
	if err != nil {
		LOGGER.Error(err)
		return
	}
	err = dev.Disconnect()
	if err != nil {
		LOGGER.Error(err)
		return
	}
	return
}

// TODO remove
// func (this *Manager) DisconnectDevice(path dbus.ObjectPath) error {
// 	if dev, err := nm.NewDevice(NMDest, path); err != nil {
// 		return err
// 	} else {
// 		dev.Disconnect()
// 		nm.DestroyDevice(dev)
// 		switch dev.DeviceType.Get() {
// 		case NM_DEVICE_TYPE_WIFI:
// 			dbus.NotifyChange(this, "WirelessConnections")
// 		case NM_DEVICE_TYPE_ETHERNET:
// 			LOGGER.Debug("DisconnectDevice...", path)
// 			dbus.NotifyChange(this, "WiredConnections")
// 		}
// 		return nil
// 	}
// }

func (this *Manager) initDeviceManage() {
	_NMManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		this.handleDeviceChanged(OpAdded, path)
	})
	_NMManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		this.handleDeviceChanged(OpRemoved, path)
	})
	devs, err := _NMManager.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, p := range devs {
		this.handleDeviceChanged(OpAdded, p)
	}
}

func (this *Manager) addWirelessDevice(dev *nm.Device) {
	wirelessDevice := NewDevice(dev)
	if isDeviceExists(this.WirelessDevices, wirelessDevice) {
		// device maybe repeat added
		return
	}
	LOGGER.Debug("addWirelessDevices:", wirelessDevice)

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = newState
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(this, "WirelessDevices")
	})

	// connect signal AccessPointAdded() and AccessPointRemoved()
	if devWireless, err := nm.NewDeviceWireless(NMDest, dev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			if this.AccessPointAdded != nil {
				if ap, err := NewAccessPoint(apPath); err == nil {
					// LOGGER.Debug("AccessPointAdded:", ap.Ssid, apPath) // TODO test
					this.AccessPointAdded(string(dev.Path), string(ap.Path))
				}
			}
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			if this.AccessPointRemoved != nil {
				// LOGGER.Debug("AccessPointRemoved:", apPath) // TODO test
				this.AccessPointRemoved(string(dev.Path), string(apPath))
			}
		})
	}

	this.WirelessDevices = append(this.WirelessDevices, wirelessDevice)
	dbus.NotifyChange(this, "WirelessDevices")
}
func (this *Manager) addWiredDevice(dev *nm.Device) {
	wiredDevice := NewDevice(dev)
	if isDeviceExists(this.WiredDevices, wiredDevice) {
		// device maybe repeat added
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		wiredDevice.State = newState
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(this, "WirelessDevices")
	})
	this.WiredDevices = append(this.WiredDevices, wiredDevice)
	dbus.NotifyChange(this, "WiredDevices")
}
func (this *Manager) addOtherDevice(dev *nm.Device) {
	this.OtherDevices = append(this.OtherDevices, NewDevice(dev))

	otherDevice := NewDevice(dev)
	if isDeviceExists(this.OtherDevices, otherDevice) {
		// device maybe repeat added
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(newState uint32, old_state uint32, reason uint32) {
		otherDevice.State = newState
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(string(dev.Path), newState)
		}
		// TODO remove
		dbus.NotifyChange(this, "WirelessDevices")
	})
	this.OtherDevices = append(this.OtherDevices, otherDevice)
	dbus.NotifyChange(this, "OtherDevices")
}
func isDeviceExists(devs []*Device, dev *Device) bool {
	for _, d := range devs {
		if d.Path == dev.Path {
			return true
		}
	}
	return false
}

func (this *Manager) handleDeviceChanged(operation int32, path dbus.ObjectPath) {
	LOGGER.Debugf("handleDeviceChanged: operation %d, path %s", operation, path)
	switch operation {
	case OpAdded:
		dev, err := nm.NewDevice(NMDest, path)
		if err != nil {
			panic(err)
		}
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			this.addWirelessDevice(dev)
		case NM_DEVICE_TYPE_ETHERNET:
			this.addWiredDevice(dev)
		default:
			this.addOtherDevice(dev)
		}
	case OpRemoved:
		var removed bool
		if this.WirelessDevices, removed = tryRemoveDevice(path, this.WirelessDevices); removed {
			dbus.NotifyChange(this, "WirelessDevices")
			LOGGER.Debug("WirelessRemoved..")
		} else if this.WiredDevices, removed = tryRemoveDevice(path, this.WiredDevices); removed {
			dbus.NotifyChange(this, "WiredDevices")
		}
	default:
		panic("Didn't support operation")
	}
}

const (
	ApKeyNone = iota
	ApKeyWep
	ApKeyPsk
	ApKeyEap
)

func parseFlags(flags, wpaFlags, rsnFlags uint32) int {
	r := ApKeyNone

	if (flags&NM_802_11_AP_FLAGS_PRIVACY != 0) && (wpaFlags == NM_802_11_AP_SEC_NONE) && (rsnFlags == NM_802_11_AP_SEC_NONE) {
		r = ApKeyWep
	}
	if wpaFlags != NM_802_11_AP_SEC_NONE {
		r = ApKeyPsk
	}
	if rsnFlags != NM_802_11_AP_SEC_NONE {
		r = ApKeyPsk
	}
	if (wpaFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) || (rsnFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) {
		r = ApKeyEap
	}
	return r
}

// GetAccessPoints return all access point's dbus path of target device.
func (this *Manager) GetAccessPoints(path dbus.ObjectPath) (aps []dbus.ObjectPath, err error) {
	dev, err := nm.NewDeviceWireless(NMDest, path)
	if err != nil {
		return
	}
	aps, err = dev.GetAccessPoints()
	return
}

// GetAccessPointProperty return access point's detail information.
func (this *Manager) GetAccessPointProperty(apPath dbus.ObjectPath) (ap AccessPoint, err error) {
	ap, err = NewAccessPoint(apPath)
	return
}

func (this *Manager) getDeviceAddress(devPath dbus.ObjectPath, devType uint32) string {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nm.NewDeviceWired(NMDest, devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWired(dev) }()
		return dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nm.NewDeviceWireless(NMDest, devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWireless(dev) }()
		return dev.HwAddress.Get()
	}
	return ""
}

func (this *Manager) ActivateConnection(uuid string, dev dbus.ObjectPath) (err error) {
	LOGGER.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, dev)
	cpath, err := _NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		LOGGER.Error(err)
		return
	}
	// TODO, ap path, "/"
	spath := dbus.ObjectPath("/")
	_, err = _NMManager.ActivateConnection(cpath, dev, spath)
	if err != nil {
		LOGGER.Error(err)
	}
	return
}

// TODO remove
func (this *Manager) DeactivateConnection(uuid string) (err error) {
	cpath := this.getActiveConnectionByUuid(uuid)
	if len(cpath) == 0 {
		return
	}
	LOGGER.Debug("DeactivateConnection:", uuid, cpath)
	err = _NMManager.DeactivateConnection(cpath)
	return
}
func (this *Manager) getActiveConnectionByUuid(uuid string) (cpath dbus.ObjectPath) {
	for _, path := range _NMManager.ActiveConnections.Get() {
		if conn, err := nm.NewActiveConnection(NMDest, path); err == nil {
			if conn.Uuid.Get() == uuid {
				cpath = path
			}
		}
	}
	return
}
