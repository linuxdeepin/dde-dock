package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type connection struct {
	Path   dbus.ObjectPath
	Uuid   string
	Id     string
	HwAddr string // if not empty, only works for special device
	Ssid   string // only used for wireless connection
}

type activeConnectionInfo struct {
	Interface    string
	HWAddress    string
	IPAddress    string
	SubnetMask   string
	RouteAddress string
	Speed        string
}

func (m *Manager) initConnectionManage() {
	m.connections = make(map[string][]connection)
	m.VPNConnections = make([]string, 0)
	m.WiredConnections = make([]string, 0) // TODO remove
	m.WirelessConnections = make([]string, 0)

	// create special wired connection if need
	m.updatePropWiredConnections() // TODO remove

	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(opAdded, c)
	}
	nmSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(opAdded, path)
	})
}

func (m *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	logger.Debugf("handleConnectionChanged: operation %d, path %s", operation, path) // TODO test
	conn := connection{Path: path}
	switch operation {
	case opAdded:
		nmConn, _ := nmNewSettingsConnection(path)
		nmConn.ConnectRemoved(func() { // TODO is still need?
			m.handleConnectionChanged(opRemoved, path)
			nm.DestroySettingsConnection(nmConn)
		})

		cdata, err := nmConn.GetSettings()
		if err != nil {
			return
		}
		uuid := getSettingConnectionUuid(cdata)
		conn.Uuid = uuid
		conn.Id = getSettingConnectionId(cdata)

		switch getSettingConnectionType(cdata) {
		case NM_SETTING_WIRED_SETTING_NAME:
			// wired connection will be treatment specially
			// TODO remove
			// m.WiredConnections = append(m.WiredConnections, uuid)
			// dbus.NotifyChange(m, "WiredConnections")
		case NM_SETTING_WIRELESS_SETTING_NAME: // TODO
			m.WirelessConnections = append(m.WirelessConnections, uuid)
			m.updatePropWirelessConnections()

			conn.Ssid = string(getSettingWirelessSsid(cdata))
			if isSettingWirelessMacAddressExists(cdata) {
				conn.HwAddr = convertMacAddressToString(getSettingWirelessMacAddress(cdata))
			}
			switch generalGetConnectionType(cdata) {
			case typeWireless:
				m.connections[typeWireless] = m.addConnection(m.connections[typeWireless], conn)
			case typeWirelessAdhoc:
				m.connections[typeWirelessAdhoc] = m.addConnection(m.connections[typeWirelessAdhoc], conn)
			case typeWirelessHotspot:
				m.connections[typeWirelessHotspot] = m.addConnection(m.connections[typeWirelessHotspot], conn)
			}
		case NM_SETTING_PPPOE_SETTING_NAME:
			m.connections[typePppoe] = m.addConnection(m.connections[typePppoe], conn)
		case NM_SETTING_GSM_SETTING_NAME:
			m.connections[typeMobile] = m.addConnection(m.connections[typeMobile], conn)
		case NM_SETTING_VPN_SETTING_NAME:
			m.VPNConnections = append(m.VPNConnections, uuid)
			m.updatePropVpnConnections()
			m.connections[typeVpn] = m.addConnection(m.connections[typeVpn], conn)
		}
		m.updatePropConnections()
	case opRemoved:
		for k, conns := range m.connections {
			if m.isConnectionExists(conns, conn) {
				m.connections[k] = m.removeConnection(conns, conn)
			}
		}
		m.updatePropConnections()
		//TODO: remove
		//removed := false
		//if m.WirelessConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), m.WirelessConnections); removed {
		//dbus.NotifyChange(m, "WirelessConnections")
		//} else if m.WiredConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), m.WiredConnections); removed {
		//dbus.NotifyChange(m, "WiredConnections")
		//} else if m.VPNConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), m.VPNConnections); removed {
		//dbus.NotifyChange(m, "VPNConnections")
		//}
	}
}
func (m *Manager) addConnection(conns []connection, conn connection) []connection {
	if m.isConnectionExists(conns, conn) {
		return conns
	}
	conns = append(conns, conn)
	return conns
}
func (m *Manager) removeConnection(conns []connection, conn connection) []connection {
	i := m.getConnectionIndex(conns, conn)
	if i < 0 {
		return conns
	}
	copy(conns[i:], conns[i+1:])
	conns = conns[:len(conns)-1]
	return conns
}
func (m *Manager) isConnectionExists(conns []connection, conn connection) bool {
	if m.getConnectionIndex(conns, conn) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getConnectionIndex(conns []connection, conn connection) int {
	for i, c := range conns {
		if c.Path == conn.Path {
			return i
		}
	}
	return -1
}

// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() (typesJSON string) {
	typesJSON, _ = marshalJSON(supportedConnectionTypesInfo)
	return
}

// GetWiredConnectionUuid return connection uuid for target wired device.
func (m *Manager) GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string) {
	// check if target wired connection exists, if not, create one
	id := "wired-connection-" + nmGetDeviceInterface(wiredDevPath)
	// TODO check connection type, read only
	cpath, ok := nmGetConnectionById(id)
	if ok {
		uuid = nmGetConnectionUuid(cpath)
	} else {
		uuid = newWiredConnection(id)
	}
	return
}

// GetActiveConnectionInfo
func (m *Manager) GetActiveConnectionInfo(devPath dbus.ObjectPath) (ret *activeConnectionInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
		}
	}()
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return nil, err
	}
	ac, err := nmNewActiveConnection(dev.ActiveConnection.Get())
	if err != nil {
		return nil, err
	}
	name := ""
	if c, err := nmNewSettingsConnection(ac.Connection.Get()); err != nil {
		return nil, err
	} else {
		if cdata, err := c.GetSettings(); err == nil {
			name = getSettingConnectionId(cdata)
		}
	}

	ip, mask, route := nmGetDHCP4Info(dev.Dhcp4Config.Get())
	defer func() {
		nm.DestroyDevice(dev)
		nm.DestroyActiveConnection(ac)
	}()

	var macaddress = "0:0:0:0:0:0"
	var speed = "-"
	switch dev.DeviceType.Get() {
	case NM_DEVICE_TYPE_ETHERNET:
		_dev, _ := nmNewDeviceWired(devPath)
		macaddress = _dev.HwAddress.Get()
		speed = fmt.Sprintf("%d", _dev.Speed.Get())
		nm.DestroyDeviceWired(_dev)
	case NM_DEVICE_TYPE_WIFI:
		_dev, _ := nmNewDeviceWireless(devPath)
		macaddress = _dev.HwAddress.Get()
		speed = fmt.Sprintf("%d", _dev.Bitrate.Get()/1024)
		nm.DestroyDeviceWireless(_dev)
	}

	return &activeConnectionInfo{
		Interface:    name,
		HWAddress:    macaddress,
		IPAddress:    ip,
		SubnetMask:   mask,
		RouteAddress: route,
		Speed:        speed,
	}, nil
}

// CreateConnection create a new connection, return ConnectionSession's dbus object path if success.
func (m *Manager) CreateConnection(connType string, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	logger.Debug("CreateConnection", connType, devPath)
	session, err = NewConnectionSessionByCreate(connType, devPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}

	return
}

// OpenConnection open a connection through uuid, return ConnectionSession's dbus object path if success.
func (m *Manager) EditConnection(uuid string, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	// if is read only connection(default system connection created by
	// network manager), create a new connection
	// TODO
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	connData, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	if getSettingConnectionReadOnly(connData) {
		logger.Debug("read only connection, create new")
		return m.CreateConnection(generalGetConnectionType(connData), devPath)
	}

	session, err = NewConnectionSessionByOpen(uuid, devPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}

	return
}

// DeleteConnection delete a connection through uuid.
func (m *Manager) DeleteConnection(uuid string) (err error) {
	//TODO: remove(uninstall dbus) editing connection_session object
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	conn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return err
	}
	return conn.Delete()
}

// TODO hide dbus interface
// GetConnectionPathByUuid return connection setting dbus path by uuid
func (m *Manager) GetConnectionPathByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmGetConnectionByUuid(uuid)
	return
}

// TODO
// GetActiveConnectionState get current state of the active connection.
func (m *Manager) GetActiveConnectionState(apath dbus.ObjectPath) (state uint32) {
	conn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}
	state = conn.State.Get()
	return
}

func (m *Manager) ActivateConnection(uuid string, devPath dbus.ObjectPath) (err error) {
	logger.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, devPath)
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	if nmGetConnectionType(cpath) == typeWireless {
		count := 0
		for _, tmpcpath := range nmGetConnectionList() {
			ctype := nmGetConnectionType(tmpcpath)
			if ctype == typeWireless {
				count++
			}
		}
		if count <= 1 {
			logger.Debug("only one access point connection, will be activate by network manager automatic")
			return
		}
	}

	_, err = nmActivateConnection(cpath, devPath)
	return
}

// TODO remove
// use disconnect device instead
func (m *Manager) deactivateConnection(uuid string) (err error) {
	apath, ok := nmGetActiveConnectionByUuid(uuid)
	if !ok {
		logger.Error("not found active connection with uuid", uuid)
		return
	}
	logger.Debug("DeactivateConnection:", uuid, apath)
	err = nmDeactivateConnection(apath)
	return
}

// DisconnectDevice will disconnect all connection in target device.
func (m *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	err = dev.Disconnect()
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
