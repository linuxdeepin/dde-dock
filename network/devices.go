package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type AccessPoint struct {
	Ssid     string
	NeedKey  bool
	Strength uint8
}

type Device struct {
	Path  dbus.ObjectPath
	State uint32
}

func NewDevice(core *nm.Device) *Device {
	return &Device{core.Path, core.State.Get()}
}

func (this *Manager) ActiveWiredDevice(active bool, path dbus.ObjectPath) {
	dev := nm.GetDevice(string(path))
	if active && dev.State.Get() == NM_DEVICE_STATE_DISCONNECTED {
		for _, c := range dev.AvailableConnections.Get() {
			_Manager.ActivateConnection(c, path, dbus.ObjectPath("/"))
		}
	} else if !active && dev.State.Get() == NM_DEVICE_STATE_ACTIVATED {
		nm.GetDevice(string(path)).Disconnect()
	}
}

func (this *Manager) initDeviceManage() {
	_Manager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		this.handleDeviceChanged(OpAdded, path)
	})
	_Manager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		this.handleDeviceChanged(OpRemoved, path)
	})
	for _, p := range _Manager.GetDevices() {
		this.handleDeviceChanged(OpAdded, p)
	}
}

func tryRemoveDevice(path dbus.ObjectPath, devices []*Device) ([]*Device, bool) {
	var newDevices []*Device
	found := false
	for _, dev := range devices {
		if dev.Path != path {
			newDevices = append(newDevices, dev)
		} else {
			found = true
		}
	}
	return newDevices, found
}

func (this *Manager) addWirelessDevice(path dbus.ObjectPath) {
	dev := nm.GetDevice(string(path))
	wirelessDevice := NewDevice(dev)
	dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
		wirelessDevice.State = new_state
		dbus.NotifyChange(this, "WirelessDevices")
	})
	this.WirelessDevices = append(this.WirelessDevices, wirelessDevice)
	dbus.NotifyChange(this, "WirelessDevices")

	nmWirelessDev := nm.GetDeviceWireless(string(path))
	nmWirelessDev.ConnectAccessPointAdded(func(p dbus.ObjectPath) {
		fmt.Println("Ap add...", p, this.GetAccessPoints(path))
		dbus.NotifyChange(this, "WirelessDevices")
	})
	nmWirelessDev.ConnectAccessPointRemoved(func(p dbus.ObjectPath) {
		fmt.Println("Ap removed...", p)
		dbus.NotifyChange(this, "WirelessDevices")
		dbus.NotifyChange(this, "WirelessDevices")
	})
}

func (this *Manager) handleDeviceChanged(operation int32, path dbus.ObjectPath) {
	switch operation {
	case OpAdded:
		dev := nm.GetDevice(string(path))
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			this.addWirelessDevice(path)
		case NM_DEVICE_TYPE_ETHERNET:
			wiredDevice := NewDevice(dev)
			dev.ConnectStateChanged(func(new_state uint32, old_state uint32, reason uint32) {
				wiredDevice.State = dev.State.Get()
				dbus.NotifyChange(this, "WiredDevices")
			})
			this.WiredDevices = append(this.WiredDevices, wiredDevice)
			dbus.NotifyChange(this, "WiredDevices")
		default:
			this.OtherDevices = append(this.OtherDevices, NewDevice(dev))
		}
	case OpRemoved:
		var removed bool
		if this.WirelessDevices, removed = tryRemoveDevice(path, this.WirelessDevices); removed {
			dbus.NotifyChange(this, "WirelessDevices")
			fmt.Println("WirelessRemoved..")
		}
		if this.WiredDevices, removed = tryRemoveDevice(path, this.WiredDevices); removed {
			dbus.NotifyChange(this, "WiredDevices")
		}
	default:
		panic("Didn't support operation")
	}
}

func (this *Manager) GetAccessPoints(path dbus.ObjectPath) []AccessPoint {
	aps := make([]AccessPoint, 0)
	dev := nm.GetDeviceWireless(string(path))
	for _, apPath := range dev.GetAccessPoints() {
		ap := nm.GetAccessPoint(string(apPath))
		aps = append(aps, AccessPoint{string(ap.Ssid.Get()), false, ap.Strength.Get()})
	}
	fmt.Println("APS:", aps)
	return aps
}
