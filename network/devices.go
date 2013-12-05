package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

const (
	NM_DEVICE_TYPE_UNKNOWN    = 0
	NM_DEVICE_TYPE_ETHERNET   = 1
	NM_DEVICE_TYPE_WIFI       = 2
	NM_DEVICE_TYPE_UNUSED1    = 3
	NM_DEVICE_TYPE_UNUSED2    = 4
	NM_DEVICE_TYPE_BT         = 5
	NM_DEVICE_TYPE_OLPC_MESH  = 6
	NM_DEVICE_TYPE_WIMAX      = 7
	NM_DEVICE_TYPE_MODEM      = 8
	NM_DEVICE_TYPE_INFINIBAND = 9
	NM_DEVICE_TYPE_BOND       = 10
	NM_DEVICE_TYPE_VLAN       = 11
	NM_DEVICE_TYPE_ADSL       = 12
	NM_DEVICE_TYPE_BRIDGE     = 13
)

const (
	NM_DEVICE_STATE_UNKNOWN      = 0
	NM_DEVICE_STATE_UNMANAGED    = 10
	NM_DEVICE_STATE_UNAVAILABLE  = 20
	NM_DEVICE_STATE_DISCONNECTED = 30
	NM_DEVICE_STATE_PREPARE      = 40
	NM_DEVICE_STATE_CONFIG       = 50
	NM_DEVICE_STATE_NEED_AUTH    = 60
	NM_DEVICE_STATE_IP_CONFIG    = 70
	NM_DEVICE_STATE_IP_CHECK     = 80
	NM_DEVICE_STATE_SECONDARIES  = 90
	NM_DEVICE_STATE_ACTIVATED    = 100
	NM_DEVICE_STATE_DEACTIVATING = 110
	NM_DEVICE_STATE_FAILED       = 120
)

func (this *Manager) updateDeviceManage() {
	this.devices = make(map[string]*nm.Device)
	_Manager.ConnectDeviceAdded(func(path string) {
		this.handleDeviceChanged(OP_ADDED, string(path))
	})
	_Manager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		this.handleDeviceChanged(OP_REMOVED, string(path))
	})
	for _, p := range _Manager.GetDevices() {
		this.handleDeviceChanged(OP_ADDED, string(p))
	}
	this.updateDeviceInfo()
}

func (this *Manager) handleDeviceChanged(operation int32, path string) {
	switch operation {
	case OP_ADDED:
		dev := nm.GetDevice(path)
		if dev.DeviceType.Get() == NM_DEVICE_TYPE_WIFI {
		}
		this.devices[path] = dev
		this.updateDeviceInfo()
	case OP_REMOVED:
		delete(this.devices, path)
		this.updateDeviceInfo()
	default:
		panic("Didn't support operation")
	}
}

func (this *Manager) updateDeviceInfo() {
	hasWired := false
	hasWireless := false
	this.APs = this.APs[0:0]
	for _, dev := range this.devices {
		switch dev.DeviceType.Get() {
		case NM_DEVICE_TYPE_WIFI:
			hasWireless = true
			this.updateAccessPoint(nm.GetDeviceWireless(string(dev.Path)))
		case NM_DEVICE_TYPE_ETHERNET:
			hasWired = true
		}
	}
	dbus.NotifyChange(this, "APs")
	if hasWired != this.HasWired {
		this.HasWired = hasWired
		dbus.NotifyChange(this, "HasWired")
	}
	if hasWireless != this.HasWireless {
		this.HasWireless = hasWireless
		dbus.NotifyChange(this, "HasWireless")
	}
}

func (this *Manager) updateAccessPoint(dev *nm.DeviceWireless) {
	for _, d := range dev.GetAccessPoints() {
		ssid := nm.GetAccessPoint(string(d)).Ssid.Get().(string)
		this.APs = append(this.APs, AccessPoint{ssid})
	}
	return
}
