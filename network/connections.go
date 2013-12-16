package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

type _ConnectionData map[string]map[string]dbus.Variant

type Connection struct {
	core           *nm.SettingsConnection
	data           _ConnectionData
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
	return c
}
func NewWirelessConnection(id string, ssid string) *Connection {
	data := make(_ConnectionData)
	data[fieldConnection] = make(map[string]dbus.Variant)
	data[fieldIPv4] = make(map[string]dbus.Variant)
	data[fieldIPv6] = make(map[string]dbus.Variant)
	data[fieldWireless] = make(map[string]dbus.Variant)
	data[fieldWirelessSecurity] = make(map[string]dbus.Variant)

	data[fieldConnection]["id"] = dbus.MakeVariant(id)
	data[fieldConnection]["uuid"] = dbus.MakeVariant(newUUID())
	data[fieldConnection]["type"] = dbus.MakeVariant(fieldWireless)

	data[fieldWireless]["ssid"] = dbus.MakeVariant([]uint8(ssid))

	data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("none")

	data[fieldIPv4]["method"] = dbus.MakeVariant("auto")

	data[fieldIPv6]["method"] = dbus.MakeVariant("auto")

	return &Connection{data: data}
}

func (this *Manager) GetConnectionByAccessPoint(ap AccessPoint) *Connection {
	for _, c := range this.WirelessConnections {
		if string(c.data[fieldWireless]["ssid"].Value().([]uint8)) == ap.Ssid {
			return c
		}
	}
	return nil
}
