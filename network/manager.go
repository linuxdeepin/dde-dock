package main

import "dlib/dbus"
import "dlib/dbus/property"
import nm "dbus/org/freedesktop/networkmanager"

const (
	dbusNmDest      = "org.freedesktop.NetworkManager"
	dbusNetworkDest = "com.deepin.daemon.Network"
	dbusNetworkPath = "/com/deepin/daemon/Network"
	dbusNetworkIfs  = "com.deepin.daemon.Network"
)

// TODO refactor code
const (
	opAdded = iota
	opRemoved
)

var (
	nmManager, _  = nm.NewManager(dbusNmDest, "/org/freedesktop/NetworkManager")
	nmSettings, _ = nm.NewSettings(dbusNmDest, "/org/freedesktop/NetworkManager/Settings")
)

type connectionData map[string]map[string]dbus.Variant

type Manager struct {
	//update by manager.go
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"` // TODO
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`
	State             uint32        // networking state

	activeConnections []*activeConnection
	ActiveConnections string // array of connections that activated and marshaled by json

	//update by devices.go
	WiredDevices    []*deviceOld
	WirelessDevices []*deviceOld
	devices         map[string][]*device
	Devices         string // array of device objects and marshaled by json
	accessPoints    map[dbus.ObjectPath][]*accessPoint
	// AccessPoints    string // array of access point objects and marshaled by json

	//update by connections.go
	WiredConnections    []string
	WirelessConnections []string
	VPNConnections      []string // TODO remove
	connections         map[string][]*connection
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
		dbusNetworkDest,
		dbusNetworkPath,
		dbusNetworkIfs,
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
	m.updateActiveConnections()
	nmManager.ActiveConnections.ConnectChanged(func() {
		m.updateActiveConnections()
	})

	// update property "State"
	m.updatePropState()
	nmManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent("org.snyh.agent")
}

func (m *Manager) updateActiveConnections() {
	// reset all exists active connection objects
	for i, _ := range m.activeConnections {
		// TODO how to disconnect BaseObserver.ConnectChanged()
		nm.DestroyActiveConnection(m.activeConnections[i].nmaconn)
		m.activeConnections[i] = nil
	}
	m.activeConnections = make([]*activeConnection, 0)
	for _, acpath := range nmGetActiveConnections() {
		if nmaconn, err := nmNewActiveConnection(acpath); err == nil {
			aconn := &activeConnection{
				nmaconn: nmaconn,
				path:    acpath,
				Devices: nmaconn.Devices.Get(),
				Uuid:    nmaconn.Uuid.Get(),
				State:   nmaconn.State.Get(),
				Vpn:     nmaconn.Vpn.Get(),
			}
			// go func() {
			// 	time.Sleep(500 * time.Millisecond)
			// 	logger.Debug("state changed:", aconn.Uuid, aconn.State, nmaconn.State.Get()) // TODO test
			// 	aconn.State = nmaconn.State.Get()
			// }()
			// TODO connect properties
			// nmaconn.Devices.ConnectChanged(func() {
			// 	// if m.isActiveConnectionExists(aconn) {
			// 	logger.Debug("state changed:", aconn.Devices, nmaconn.Devices.Get()) // TODO test
			// 	aconn.Devices = nmaconn.Devices.Get()
			// 	// m.updatePropActiveConnections()
			// 	// }
			// })
			nmaconn.State.ConnectChanged(func() {
				// TODO fix dbus property issue
				// if m.isActiveConnectionExists(aconn) {
				logger.Debug("state changed:", aconn.State, nmaconn.State.Get()) // TODO test
				aconn.State = nmaconn.State.Get()
				m.updatePropActiveConnections()
				// }
			})
			m.activeConnections = append(m.activeConnections, aconn)
		}
	}
	m.updatePropActiveConnections()
	logger.Debug("active connections changed:", m.ActiveConnections) // TODO test
}
func (m *Manager) isActiveConnectionExists(aconn *activeConnection) bool {
	if aconn == nil {
		return false
	}
	for _, a := range m.activeConnections {
		if &a == &aconn {
			return true
		}
	}
	return false
}
