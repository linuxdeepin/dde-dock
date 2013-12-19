package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type _ConnectionData map[string]map[string]dbus.Variant

type Connection struct {
	core *nm.SettingsConnection
	data _ConnectionData

	Uuid           string
	Name           string
	ConnectionType string
}

func (this *Connection) GetDBusInfo_() dbus.DBusInfo {
	if this.core != nil {
		return dbus.DBusInfo{DBusDest, string(this.core.Path), DBusIFC + ".Connection"}
	} else {
		return dbus.DBusInfo{DBusDest, "/", DBusIFC + ".Connection"}
	}
}

func (this *Manager) initConnectionManage() {
	this.VPNConnections = make([]*Connection, 0)
	this.WiredConnections = make([]*Connection, 0)
	this.WirelessConnections = make([]*Connection, 0)

	for _, c := range _Settings.ListConnections() {
		this.handleConnectionChanged(OpAdded, string(c))
	}
	_Settings.ConnectNewConnection(func(path dbus.ObjectPath) {
		this.handleConnectionChanged(OpAdded, string(path))
	})
}

func (this *Manager) handleConnectionChanged(operation int32, path string) {
	switch operation {
	case OpAdded:
		c := NewConnection(nm.GetSettingsConnection(path))
		switch c.ConnectionType {
		case "802-11-wireless":
			this.WirelessConnections = append(this.WirelessConnections, c)
			dbus.NotifyChange(this, "WirelessConnections")
		case "802-3-ethernet":
			this.WiredConnections = append(this.WiredConnections, c)
			dbus.NotifyChange(this, "WiredConnections")
		case "pppoe":
		case "vpn":
			this.VPNConnections = append(this.VPNConnections, c)
			dbus.NotifyChange(this, "VPNConnections")
		case "cdma":
		}
	}
}

func NewConnection(core *nm.SettingsConnection) *Connection {
	c := &Connection{core: core}
	settings := core.GetSettings()
	c.Name = settings["connection"]["id"].Value().(string)
	c.Uuid = settings["connection"]["uuid"].Value().(string)
	c.ConnectionType = settings["connection"]["type"].Value().(string)
	c.data = core.GetSettings()
	return c
}

func NewWirelessConnection(id string, ssid string, keyFlag int) *Connection {
	data := make(_ConnectionData)
	data[fieldConnection] = make(map[string]dbus.Variant)
	data[fieldIPv4] = make(map[string]dbus.Variant)
	data[fieldIPv6] = make(map[string]dbus.Variant)
	data[fieldWireless] = make(map[string]dbus.Variant)

	data[fieldConnection]["id"] = dbus.MakeVariant(id)
	uuid := newUUID()
	data[fieldConnection]["uuid"] = dbus.MakeVariant(uuid)
	data[fieldConnection]["type"] = dbus.MakeVariant(fieldWireless)

	data[fieldWireless]["ssid"] = dbus.MakeVariant([]uint8(ssid))

	if keyFlag != ApKeyNone {
		data[fieldWirelessSecurity] = make(map[string]dbus.Variant)
		data[fieldWireless]["security"] = dbus.MakeVariant(fieldWirelessSecurity)
		switch keyFlag {
		case ApKeyWep:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("none")
		case ApKeyPsk:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-psk")
			data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
		case ApKeyEap:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-eap")
			data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
		}
	}

	data[fieldIPv4]["method"] = dbus.MakeVariant("auto")

	data[fieldIPv6]["method"] = dbus.MakeVariant("auto")

	core := nm.GetSettingsConnection(string(_Settings.AddConnection(data)))

	return &Connection{core, data, uuid, id, fieldWireless}
}

func (this *Manager) GetConnectionByAccessPoint(path dbus.ObjectPath) *Connection {
	ap := nm.GetAccessPoint(string(path))
	for _, c := range this.WirelessConnections {
		fmt.Println(c.data)
		if c.ConnectionType == fieldWireless && string(c.data[fieldWireless]["ssid"].Value().([]uint8)) == string(ap.Ssid.Get()) {
			return c
		}
	}
	fmt.Println("CCC:", path, string(ap.Ssid.Get()))
	return NewWirelessConnection(string(ap.Ssid.Get()), string(ap.Ssid.Get()), parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()))
}
