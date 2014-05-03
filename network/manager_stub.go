package main

import (
	"dlib/dbus"
)

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections = make([]string, 0)
	for _, cpath := range nmGetActiveConnections() {
		if aconn, err := nmNewActiveConnection(cpath); err == nil {
			m.ActiveConnections = append(m.ActiveConnections, aconn.Uuid.Get())
			logger.Debugf("ActiveConnections, uuid=%s, state=%d", aconn.Uuid.Get(), aconn.State.Get()) // TODO test
		}
	}
	dbus.NotifyChange(m, "ActiveConnections")
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
func (m *Manager) updatePropOtherDevices() {
	dbus.NotifyChange(m, "OtherDevices")
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
