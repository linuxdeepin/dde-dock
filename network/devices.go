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
		dbus.NotifyChange(this, "WirelessDevices")
	})
	nmWirelessDev.ConnectAccessPointRemoved(func(p dbus.ObjectPath) {
		dbus.NotifyChange(this, "WirelessDevices")
		dbus.NotifyChange(this, "WirelessDevices")
	})
	this.GetAccessPoints((path))
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

func parseFlags_(flag uint32) string {
	r := ""
	if flag&NM_802_11_AP_SEC_PAIR_WEP40 != 0 {
		r += " PAIR_WEP40"
	}
	if flag&NM_802_11_AP_SEC_PAIR_WEP104 != 0 {
		r += " PAIR_WEP104"
	}
	if flag&NM_802_11_AP_SEC_PAIR_TKIP != 0 {
		r += " PAIR_TKIP"
	}
	if flag&NM_802_11_AP_SEC_PAIR_CCMP != 0 {
		r += " PAIR_CCMP"
	}
	if flag&NM_802_11_AP_SEC_GROUP_WEP40 != 0 {
		r += " GROUP_WEP40"
	}
	if flag&NM_802_11_AP_SEC_GROUP_WEP104 != 0 {
		r += " GROUP_WEP40"
	}
	if flag&NM_802_11_AP_SEC_GROUP_TKIP != 0 {
		r += " GROUP_WEP40"
	}
	if flag&NM_802_11_AP_SEC_GROUP_CCMP != 0 {
		r += " GROUP_WEP40"
	}
	if flag&NM_802_11_AP_SEC_KEY_MGMT_PSK != 0 {
		r += " MGMT_PSK"
	}
	if flag&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0 {
		r += " MGMT_802.1X"
	}
	return r
}

func (this *Manager) GetAccessPoints(path dbus.ObjectPath) []AccessPoint {
	aps := make([]AccessPoint, 0)
	dev := nm.GetDeviceWireless(string(path))
	/*ac := nm.GetActiveConnection(string(nm.GetDevice(string(path)).ActiveConnection.Get())).Connection.Get()*/
	for i, apPath := range dev.GetAccessPoints() {
		ap := nm.GetAccessPoint(string(apPath))
		/*actived := string(nm.GetSettingsConnection(string(ac)).GetSettings()[fieldWireless]["ssid"].Value().([]uint8)) == string(ap.Ssid.Get())*/
		actived := false
		aps = append(aps, AccessPoint{string(ap.Ssid.Get()),
			parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()) != ApKeyNone,
			ap.Strength.Get(),
			ap.Path,
			actived,
		})
		fmt.Printf("%s %s %s %s\n", ap.Path, aps[i].Ssid, parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()), parseFlags_(ap.WpaFlags.Get()))
	}
	return aps
}
