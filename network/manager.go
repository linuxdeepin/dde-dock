package main

import "dlib/dbus"
import "dlib/dbus/property"
import nm "dbus/org/freedesktop/networkmanager"

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
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"`
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`

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

	this.updateDeviceManage()
	this.updateConnectionManage()
}

func NewManager() (m *Manager) {
	this := &Manager{}
	this.updateManager()
	this.WirelessEnabled = property.NewWrapProperty(this, "WirelessEnabled", _Manager.WirelessEnabled)
	this.NetworkingEnabled = property.NewWrapProperty(this, "NetworkingEnabled", _Manager.NetworkingEnabled)

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
