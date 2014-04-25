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
	ActiveConnections []string      // uuid collection of connections that activated
	State             uint32        // networking state

	//update by devices.go
	WiredDevices    []*device
	WirelessDevices []*device // TODO is "device" struct still needed?
	OtherDevices    []*device

	//update by connections.go
	WiredConnections    []string
	WirelessConnections []string
	VPNConnections      []string
	uuid2connectionType map[string]string // TODO remove

	//signals
	NeedSecrets                  func(string, string, string)
	DeviceStateChanged           func(devPath string, newState uint32)
	AccessPointAdded             func(devPath string, apPath string)
	AccessPointRemoved           func(devPath string, apPath string)
	AccessPointPropertiesChanged func(devPath string, apPath string)

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

	// update property "ActiveConnections" after initilizing device
	m.updatePropActiveConnections()
	nmManager.ActiveConnections.ConnectChanged(func() {
		m.updatePropActiveConnections()
	})

	// update property "State"
	m.updatePropState()
	nmManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent("org.snyh.agent")
}
