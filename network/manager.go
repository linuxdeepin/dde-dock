package main

import "dlib/dbus"
import "dlib/dbus/property"
import nm "dbus/org/freedesktop/networkmanager"

const NMDest = "org.freedesktop.NetworkManager"

const (
	OpAdded = iota
	OpRemoved
)

var (
	NMManager, _  = nm.NewManager(NMDest, "/org/freedesktop/NetworkManager")
	NMSettings, _ = nm.NewSettings(NMDest, "/org/freedesktop/NetworkManager/Settings")
)

type Manager struct {
	//update by manager.go
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"`
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`
	ActiveConnections []string      // uuid collection of connections that activated
	State             uint32        // networking state

	//update by devices.go
	WiredDevices    []*Device
	WirelessDevices []*Device // TODO is "Device" struct still needed?
	OtherDevices    []*Device

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
	m.WirelessEnabled = property.NewWrapProperty(m, "WirelessEnabled", NMManager.WirelessEnabled)
	m.NetworkingEnabled = property.NewWrapProperty(m, "NetworkingEnabled", NMManager.NetworkingEnabled)

	m.initDeviceManage()
	m.initConnectionManage()

	// update property "ActiveConnections" after initilizing device
	m.updatePropActiveConnections()
	NMManager.ActiveConnections.ConnectChanged(func() {
		m.updatePropActiveConnections()
	})

	// update property "State"
	m.updatePropState()
	NMManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent("org.snyh.agent")
}
