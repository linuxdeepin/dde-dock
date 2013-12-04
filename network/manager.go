package main

import "dlib/dbus"
import nm "networkmanager"

const (
	DBUS_DEST = "com.deepin.daemon.Network"
	DBUS_PATH = "/com/deepin/daemon/Network"
	DBUS_IFC  = "com.deepin.daemon.Network"
)

const (
	OP_ADDED = iota
	OP_REMOVED
)

var (
	_Manager  = nm.GetManager("/org/freedesktop/NetworkManager")
	_Settings = nm.GetSettings("/org/freedesktop/NetworkManager/Settings")
)

type AccessPoint struct {
	Uuid string
}

type Manager struct {
	//update by manager.go
	WiredEnabled      bool `access:"readwrite"`
	WirelessEnabled   bool `access:"readwrite"`
	VPNEnabled        bool `access:"readwrite"`
	NetworkingEnabled bool `access:"readwrite"`

	//update by devices.go
	APs         []AccessPoint
	HasWireless bool
	HasWired    bool
	devices     map[string]*nm.Device

	//update by connections.go
	WiredConnections    []*Connection
	WirelessConnections []*Connection
	VPNConnections      []*Connection
}

func (this *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{DBUS_DEST, DBUS_PATH, DBUS_IFC}
}

func (this *Manager) updateManager() {
	this.WiredEnabled = true
	this.WirelessEnabled = _Manager.GetWirelessEnabled()
	this.NetworkingEnabled = _Manager.GetNetworkingEnabled()

	this.updateDeviceManage()
	this.updateConnectionManage()
}

func NewManager() (m *Manager) {
	this := &Manager{}
	this.updateManager()
	_Manager.ConnectPropertiesChanged(func(props map[string]dbus.Variant) {
		this.updateManager()
		if _, ok := props["WirelessEnabled"]; ok {
			dbus.NotifyChange(this, "WirelessEnabled")
		}
		if _, ok := props["NetworkingEnabled"]; ok {
			dbus.NotifyChange(this, "NetworkingEnabled")
		}
	})

	_Settings.ConnectNewConnection(func(path dbus.ObjectPath) {
		this.handleConnectionChanged(OP_ADDED, string(path))
	})
	return this
}

func main() {
	manager := NewManager()
	dbus.InstallOnSession(manager)

	select {}
}
