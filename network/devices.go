package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type AccessPoint struct {
	Ssid     string
	NeedKey  bool
	Strength uint8
	Path     dbus.ObjectPath
	Actived  bool
}

type Device struct {
	Path  dbus.ObjectPath
	State uint32
}

func NewDevice(core *nm.Device) *Device {
	return &Device{core.Path, core.State.Get()}
}

func (this *Manager) DisconnectDevice(path dbus.ObjectPath) error {
	if dev, err := nm.NewDevice(path); err != nil {
		return err
	} else {
		dev.Disconnect()
		nm.DestroyDevice(dev)
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			dbus.NotifyChange(this, "WirelessConnections")
		case NM_DEVICE_TYPE_ETHERNET:
			fmt.Println("DisconnectDevice...", path)
			dbus.NotifyChange(this, "WiredConnections")
		}
		return nil
	}
}

func (this *Manager) ActiveWiredDevice(path dbus.ObjectPath) error {
	if dev, err := nm.NewDevice(path); err != nil {
		return err
	} else {
		if dev.State.Get() == NM_DEVICE_STATE_DISCONNECTED {
			for _, c := range dev.AvailableConnections.Get() {
				_NMManager.ActivateConnection(c, path, dbus.ObjectPath("/"))
			}
			dbus.NotifyChange(this, "WiredConnections")
		} else {
			return fmt.Errorf("WiredDeveice %v has been actived already.", path)
		}
		nm.DestroyDevice(dev)
	}
	return nil
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
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = new_state
		dbus.NotifyChange(this, "WirelessDevices")
	})
	this.WirelessDevices = append(this.WirelessDevices, wirelessDevice)
	dbus.NotifyChange(this, "WirelessDevices")
}
func (this *Manager) addWiredDevice(dev *nm.Device) {
	wiredDevice := NewDevice(dev)
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		wiredDevice.State = new_state
		dbus.NotifyChange(this, "WiredDevices")
	})
	this.WiredDevices = append(this.WiredDevices, wiredDevice)
	dbus.NotifyChange(this, "WiredDevices")
}
func (this *Manager) addOtherDevice(dev *nm.Device) {
	this.OtherDevices = append(this.OtherDevices, NewDevice(dev))

	otherDevice := NewDevice(dev)
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		otherDevice.State = new_state
		dbus.NotifyChange(this, "OtherDevices")
	})
	this.OtherDevices = append(this.OtherDevices, otherDevice)
	dbus.NotifyChange(this, "OtherDevices")
}

func (this *Manager) handleDeviceChanged(operation int32, path dbus.ObjectPath) {
	switch operation {
	case OpAdded:
		dev, err := nm.NewDevice(path)
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
			fmt.Println("WirelessRemoved..")
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
	aps := make([]AccessPoint, 0)
	if dev, err := nm.NewDeviceWireless(path); err == nil {
		nmAps, err := dev.GetAccessPoints()
		if err != nil {
			return nil, err
		}
		for _, apPath := range nmAps {
			if ap, err := nm.NewAccessPoint(apPath); err == nil {
				actived := dev.ActiveAccessPoint.Get() == apPath
				aps = append(aps, AccessPoint{string(ap.Ssid.Get()),
					parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()) != ApKeyNone,
					calcStrength(ap.Strength.Get()),
					ap.Path,
					actived,
				})
			}
		}
		return aps, nil
	} else {
		return nil, err
	}
}

func (this *Manager) getDeviceAddress(devPath dbus.ObjectPath, devType uint32) string {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		dev, err := nm.NewDeviceWired(devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWired(dev) }()
		return dev.HwAddress.Get()
	case NM_DEVICE_TYPE_WIFI:
		dev, err := nm.NewDeviceWireless(devPath)
		if err != nil {
			panic(err)
		}
		defer func() { nm.DestroyDeviceWireless(dev) }()
		return dev.HwAddress.Get()
	}
	return ""
}

func (this *Manager) ActiveAccessPoint(dev dbus.ObjectPath, ap dbus.ObjectPath) error {
	con, err := this.GetConnectionByAccessPoint(ap)
	if err != nil {
		return err
	}
	_NMManager.ActivateConnection(con.Path, dev, ap)
	return nil
}
