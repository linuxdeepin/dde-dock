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

func (this *Manager) DisconnectDevice(path dbus.ObjectPath) error {
	if dev, err := nm.NewDevice(NMDest, path); err != nil {
		return err
	} else {
		dev.Disconnect()
		nm.DestroyDevice(dev)
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			dbus.NotifyChange(this, "WirelessConnections")
		case NM_DEVICE_TYPE_ETHERNET:
			LOGGER.Debug("DisconnectDevice...", path)
			dbus.NotifyChange(this, "WiredConnections")
		}
		return nil
	}
}

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
		// device maybe is repeated added
		return
	}
	LOGGER.Debug("addWirelessDevices:", wirelessDevice)

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = new_state
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(dev.Path, new_state)
		}
		// TODO remove
		// dbus.NotifyChange(this, "WirelessDevices")
	})

	// connect signal AccessPointAdded() and AccessPointRemoved()
	if aps, err := this.GetAccessPoints(dev.Path); err == nil {
		for _, ap := range aps {
			if this.AccessPointAdded != nil {
				this.AccessPointAdded(dev.Path, ap)
			}
		}
	}
	if devWireless, err := nm.NewDeviceWireless(NMDest, dev.Path); err == nil {
		devWireless.ConnectAccessPointAdded(func(apPath dbus.ObjectPath) {
			if this.AccessPointAdded != nil {
				if ap, err := newAccessPoint(apPath); err == nil {
					this.AccessPointAdded(dev.Path, ap)
				}
			}
		})
		devWireless.ConnectAccessPointRemoved(func(apPath dbus.ObjectPath) {
			if this.AccessPointRemoved != nil {
				this.AccessPointRemoved(dev.Path, apPath)
			}
		})
	}

	this.WirelessDevices = append(this.WirelessDevices, wirelessDevice)
	dbus.NotifyChange(this, "WirelessDevices")
}
func (this *Manager) addWiredDevice(dev *nm.Device) {
	wiredDevice := NewDevice(dev)
	if isDeviceExists(this.WiredDevices, wiredDevice) {
		// device maybe is repeated added
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		wiredDevice.State = new_state
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(dev.Path, new_state)
		}
		// TODO remove
		// dbus.NotifyChange(this, "WirelessDevices")
	})
	this.WiredDevices = append(this.WiredDevices, wiredDevice)
	dbus.NotifyChange(this, "WiredDevices")
}
func (this *Manager) addOtherDevice(dev *nm.Device) {
	this.OtherDevices = append(this.OtherDevices, NewDevice(dev))

	otherDevice := NewDevice(dev)
	if isDeviceExists(this.OtherDevices, otherDevice) {
		// may be repeated to add device
		return
	}

	// connect signal DeviceStateChanged()
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		otherDevice.State = new_state
		if this.DeviceStateChanged != nil {
			this.DeviceStateChanged(dev.Path, new_state)
		}
		// TODO remove
		// dbus.NotifyChange(this, "WirelessDevices")
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

func (this *Manager) GetAccessPoints(path dbus.ObjectPath) ([]AccessPoint, error) {
	aps := make([]AccessPoint, 0)
	if dev, err := nm.NewDeviceWireless(NMDest, path); err == nil {
		nmAps, err := dev.GetAccessPoints()
		if err != nil {
			LOGGER.Error("GetAccessPoints:", err) // TODO test
			return nil, err
		}
		for _, apPath := range nmAps {
			// TODO remove
			// if ap, err := nm.NewAccessPoint(NMDest, apPath); err == nil {
			// 	actived := dev.ActiveAccessPoint.Get() == apPath
			// 	aps = append(aps, AccessPoint{string(ap.Ssid.Get()),
			// 		parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()) != ApKeyNone,
			// 		calcStrength(ap.Strength.Get()),
			// 		ap.Path,
			// 		actived,
			// 	})
			// }
			if ap, err := newAccessPoint(apPath); err == nil {
				aps = append(aps, ap)
			}
		}
		return aps, nil
	} else {
		return nil, err
	}
}

func newAccessPoint(apPath dbus.ObjectPath) (ap AccessPoint, err error) {
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

func (this *Manager) ActivateConnection(uuid string, dev dbus.ObjectPath) {
	if cpath, err := _NMSettings.GetConnectionByUuid(uuid); err == nil {
		LOGGER.Debug("ActivateConnection:", uuid, dev)
		// TODO, ap path, "/"
		spath := dbus.ObjectPath("/")
		_NMManager.ActivateConnection(cpath, dev, spath)
	}
}
func (this *Manager) DeactivateConnection(cpath dbus.ObjectPath) error {
	return _NMManager.DeactivateConnection(cpath)
}
