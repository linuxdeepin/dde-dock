package main

import (
	"dlib/dbus"
)

func (m *Manager) updatePropActiveConnections() {
	m.activeConnections = make([]activeConnection, 0)
	for _, cpath := range nmGetActiveConnections() {
		if aconn, err := nmNewActiveConnection(cpath); err == nil {
			aconnObj := activeConnection{
				Devices: aconn.Devices.Get(),
				Uuid:    aconn.Uuid.Get(),
				State:   aconn.State.Get(),
			}
			m.activeConnections = append(m.activeConnections, aconnObj)
			// TODO
			// aconn.State.ConnectChanged(func() {
			// 	aconnObj.State = aconn.State.Get()
			// 	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
			// 	// dbus.NotifyChange(m, "ActiveConnections")
			// 	logger.Debug("ActiveConnection:", m.ActiveConnections)
			// })
		}
	}
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
	logger.Debug("ActiveConnection:", m.ActiveConnections) // TODO test
}

func (m *Manager) updatePropActivatingConnection() {
	// TODO
	dbus.NotifyChange(m, "ActivatingConnection")
}

func (m *Manager) updatePropState() {
	m.State = nmGetState()
	dbus.NotifyChange(m, "State")
}

func (m *Manager) updatePropWiredDevices() {
	dbus.NotifyChange(m, "WiredDevices")
}
func (m *Manager) updatePropWirelessDevices() {
	dbus.NotifyChange(m, "WirelessDevices")
}
func (m *Manager) updatePropDevices() {
	m.Devices, _ = marshalJSON(m.devices)
	dbus.NotifyChange(m, "Devices")
	logger.Debug("updatePropDevices", m.Devices) // TODO test
}

// TODO remove
// create connection for each wired device if not exists
func (m *Manager) updatePropWiredConnections() {
	m.WiredConnections = make([]string, 0)
	for _, wiredDev := range m.WiredDevices {
		uuid := m.GetWiredConnectionUuid(wiredDev.Path)
		m.WiredConnections = append(m.WiredConnections, uuid)
	}
	dbus.NotifyChange(m, "WiredConnections")
}
func (m *Manager) updatePropWirelessConnections() {
	dbus.NotifyChange(m, "WirelessConnections")
}
func (m *Manager) updatePropVpnConnections() {
	dbus.NotifyChange(m, "VPNConnections")
}
func (m *Manager) updatePropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	// logger.Debug(m.Connections) // TODO test
	dbus.NotifyChange(m, "Connections")
}
