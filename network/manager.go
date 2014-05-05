package main

import "dlib/dbus"
import "dlib/dbus/property"
import nm "dbus/org/freedesktop/networkmanager"

const nmDest = "org.freedesktop.NetworkManager"

// TODO
const (
	opAdded = iota
	opRemoved
)

var (
	nmManager, _  = nm.NewManager(nmDest, "/org/freedesktop/NetworkManager")
	nmSettings, _ = nm.NewSettings(nmDest, "/org/freedesktop/NetworkManager/Settings")
)

type connectionData map[string]map[string]dbus.Variant

type Manager struct {
	//update by manager.go
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"` // TODO
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`
	State             uint32        // networking state

	activeConnections []activeConnection
	ActiveConnections string // array of connections that activated and marshaled by json

	ActivatingConnection string // TODO

	//update by devices.go
	WiredDevices    []*deviceOld
	WirelessDevices []*deviceOld
	devices         map[string][]*device
	Devices         string // array of device objects and marshaled by json

	//update by connections.go
	WiredConnections    []string
	WirelessConnections []string
	VPNConnections      []string // TODO remove
	connections         map[string][]connection
	Connections         string // array of connection information and marshaled by json

	//signals
	NeedSecrets                  func(string, string, string)
	DeviceStateChanged           func(devPath string, newState uint32)
	AccessPointAdded             func(devPath, apJSON string)
	AccessPointRemoved           func(devPath, apJSON string)
	AccessPointPropertiesChanged func(devPath, apJSON string)

	agent *Agent
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Network",
		"/com/deepin/daemon/Network",
		"com.deepin.daemon.Network",
	}
}

func NewManager() (m *Manager) {
	m = &Manager{}
	return
}

func (m *Manager) initManager() {
	m.WiredEnabled = true
	m.WirelessEnabled = property.NewWrapProperty(m, "WirelessEnabled", nmManager.WirelessEnabled)
	m.NetworkingEnabled = property.NewWrapProperty(m, "NetworkingEnabled", nmManager.NetworkingEnabled)

	m.initDeviceManage()
	m.initConnectionManage()

	// update property "ActiveConnections" after initilizing devices
	m.updatePropActiveConnections()
	nmManager.ActiveConnections.ConnectChanged(func() {
		m.updatePropActiveConnections()
	})

	// TODO need update dbus-factory about network-manager
	// update property "ActivatingConnection"
	// m.updatePropActivatingConnection()
	// nmManager.

	// update property "State"
	m.updatePropState()
	nmManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent("org.snyh.agent")
}
