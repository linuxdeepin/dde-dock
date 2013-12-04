package main

import nm "networkmanager"
import "dlib/dbus"

type Connection struct {
	core           *nm.SettingsConnection
	Uuid           string
	Name           string
	ConnectionType string
}

func (this *Connection) GetDBusInfo_() dbus.DBusInfo {
	if this.core != nil {
		return dbus.DBusInfo{DBUS_DEST, string(this.core.Path), DBUS_IFC + ".Connection"}
	} else {
		return dbus.DBusInfo{DBUS_DEST, "/", DBUS_IFC + ".Connection"}
	}
}

func (this *Connection) ActiveAccessPoint(path string) {
	settings := this.core.GetSettings()
	settings["802-11-wireless"]["ssid"] = dbus.MakeVariant(nm.GetAccessPoint(path).GetSsid())
}

func (this *Manager) updateConnectionManage() {
	this.VPNConnections = make([]*Connection, 0)
	this.WiredConnections = make([]*Connection, 0)
	this.WirelessConnections = make([]*Connection, 0)

	for _, c := range _Settings.ListConnections() {
		this.handleConnectionChanged(OP_ADDED, string(c))
	}
}

func (this *Manager) handleConnectionChanged(operation int32, path string) {
	switch operation {
	case OP_ADDED:
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
