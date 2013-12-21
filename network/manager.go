package main

import "dlib/dbus"
import "dlib/dbus/property"
import nm "dbus/org/freedesktop/networkmanager"

const (
	DBusDest = "com.deepin.daemon.Network"
	DBusPath = "/com/deepin/daemon/Network"
	DBusIFC  = "com.deepin.daemon.Network"
)

const (
	OpAdded = iota
	OpRemoved
)

var (
	_NMManager  = nm.GetManager("/org/freedesktop/NetworkManager")
	_NMSettings = nm.GetSettings("/org/freedesktop/NetworkManager/Settings")
	_Manager    = _NewManager()
)

type Manager struct {
	//update by manager.go
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"`
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`

	//update by devices.go
	APs             []*AccessPoint
	WirelessDevices []*Device
	WiredDevices    []*Device
	OtherDevices    []*Device

	//update by connections.go
	WiredConnections    []*Connection
	WirelessConnections []*Connection
	VPNConnections      []*Connection

	NeedMoreConfigure func(string, string)

	agent *Agent
}

func (this *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{DBusDest, DBusPath, DBusIFC}
}

func (this *Manager) initManager() {
	this.WiredEnabled = true
	this.WirelessEnabled = property.NewWrapProperty(this, "WirelessEnabled", _NMManager.WirelessEnabled)
	this.NetworkingEnabled = property.NewWrapProperty(this, "NetworkingEnabled", _NMManager.NetworkingEnabled)
	this.initDeviceManage()
	this.initConnectionManage()
}

func _NewManager() (m *Manager) {
	this := &Manager{}
	this.initManager()
	this.agent = NewAgent("org.snyh.agent")
	return this
}

func main() {
	dbus.InstallOnSession(_Manager)
	select {}
}
