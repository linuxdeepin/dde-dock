package network

import (
	"dlib/dbus"
)

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
	// logger.Debug("ActiveConnections:", m.ActiveConnections) // TODO test
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
	// logger.Debug("updatePropDevices", m.Devices) // TODO test
}

// TODO
func (m *Manager) updatePropAccessPoints() {
	// m.AccessPoints, _ = marshalJSON(m.accessPoints)
	// dbus.NotifyChange(m, "AccessPoints")
	// testJSON, _ := marshalJSON(m.accessPoints)
	// logger.Debug("updatePropAccessPoints", testJSON) // TODO test
}

func (m *Manager) updatePropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	logger.Debug(m.Connections) // TODO test
	dbus.NotifyChange(m, "Connections")
}
