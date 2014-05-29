package network

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

// TODO different connection structures for different types
type connection struct {
	Path dbus.ObjectPath
	Uuid string
	Id   string

	// if not empty, the connection will only apply to special device,
	// works for wired, wireless, infiniband, wimax devices
	HwAddress string

	// works for wireless, olpc-mesh connections
	Ssid string
}

type activeConnection struct {
	nmaconn *nm.ActiveConnection
	path    dbus.ObjectPath

	Devices []dbus.ObjectPath
	// SpecificObject dbus.ObjectPath // TODO
	Uuid  string
	State uint32
	Vpn   bool
	// VpnState uint32 // TODO
}

type activeConnectionInfo struct {
	DeviceType   string
	Interface    string
	HwAddress    string
	IpAddress    string
	SubnetMask   string
	RouteAddress string
	Dns1         string
	Dns2         string
	Speed        string
}

func (m *Manager) initConnectionManage() {
	m.connections = make(map[string][]*connection)

	// create special wired connection if need
	/*m.updatePropWiredConnections() // TODO remove*/

	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(opAdded, c)
	}
	nmSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(opAdded, path)
	})
}

func (m *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	// logger.Debugf("handleConnectionChanged: operation %d, path %s", operation, path) // TODO test
	conn := &connection{Path: path}
	switch operation {
	case opAdded:
		nmConn, _ := nmNewSettingsConnection(path)
		nmConn.ConnectRemoved(func() {
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
		case NM_SETTING_WIRELESS_SETTING_NAME:
			conn.Ssid = string(getSettingWirelessSsid(cdata))
			if isSettingWirelessMacAddressExists(cdata) {
				conn.HwAddress = convertMacAddressToString(getSettingWirelessMacAddress(cdata))
			}
			switch getCustomConnectionType(cdata) {
			case connectionWireless:
				m.connections[connectionWireless] = m.addConnection(m.connections[connectionWireless], conn)
			case connectionWirelessAdhoc:
				m.connections[connectionWirelessAdhoc] = m.addConnection(m.connections[connectionWirelessAdhoc], conn)
			case connectionWirelessHotspot:
				m.connections[connectionWirelessHotspot] = m.addConnection(m.connections[connectionWirelessHotspot], conn)
			}
		case NM_SETTING_PPPOE_SETTING_NAME:
			m.connections[connectionPppoe] = m.addConnection(m.connections[connectionPppoe], conn)
		case NM_SETTING_GSM_SETTING_NAME, NM_SETTING_CDMA_SETTING_NAME:
			m.connections[connectionMobile] = m.addConnection(m.connections[connectionMobile], conn)
		case NM_SETTING_VPN_SETTING_NAME:
			m.connections[connectionVpn] = m.addConnection(m.connections[connectionVpn], conn)
		}
		m.updatePropConnections()
	case opRemoved:
		for k, conns := range m.connections {
			if m.isConnectionExists(conns, conn) {
				m.connections[k] = m.removeConnection(conns, conn)
			}
		}
		m.updatePropConnections()
	}
}
func (m *Manager) addConnection(conns []*connection, conn *connection) []*connection {
	if m.isConnectionExists(conns, conn) {
		return conns
	}
	conns = append(conns, conn)
	return conns
}
func (m *Manager) removeConnection(conns []*connection, conn *connection) []*connection {
	i := m.getConnectionIndex(conns, conn)
	if i < 0 {
		return conns
	}
	copy(conns[i:], conns[i+1:])
	conns = conns[:len(conns)-1]
	return conns
}
func (m *Manager) isConnectionExists(conns []*connection, conn *connection) bool {
	if m.getConnectionIndex(conns, conn) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getConnectionIndex(conns []*connection, conn *connection) int {
	for i, c := range conns {
		if c.Path == conn.Path {
			return i
		}
	}
	return -1
}

// TODO remove
// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() (typesJSON string) {
	typesJSON, _ = marshalJSON(supportedConnectionTypesInfo)
	return
}

// TODO remove
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

func (m *Manager) GetActiveConnectionInfo(devPath dbus.ObjectPath) (acinfoJSON string, err error) {
	// get connection data
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devName := getDeviceName(nmDev.DeviceType.Get())

	aconn := nmDev.ActiveConnection.Get()
	if aconn == "/" {
		acinfoJSON = ""
		return
	}
	nmAConn, err := nmNewActiveConnection(aconn)
	if err != nil {
		return
	}
	nmConn, err := nmNewSettingsConnection(nmAConn.Connection.Get())
	if err != nil {
		return
	}

	// query connection data
	cdata, err := nmConn.GetSettings()
	name := ""
	dns2 := ""
	if err == nil {
		name = getSettingConnectionId(cdata)
		dns2 = getSettingVkIp4ConfigDns(cdata)
	}

	// query dhcp4
	ip, mask, route, dns1 := nmGetDHCP4Info(nmDev.Dhcp4Config.Get())

	// get hardware address
	hwAddress, err := nmGeneralGetDeviceHwAddr(devPath)
	if err != nil {
		hwAddress = "00:00:00:00:00:00"
	}

	// get network speed (Mb/s)
	var speed = "-"
	switch nmDev.DeviceType.Get() {
	case NM_DEVICE_TYPE_ETHERNET:
		devWired, _ := nmNewDeviceWired(devPath)
		speed = fmt.Sprintf("%d", devWired.Speed.Get())
	case NM_DEVICE_TYPE_WIFI:
		devWireless, _ := nmNewDeviceWireless(devPath)
		speed = fmt.Sprintf("%d", devWireless.Bitrate.Get()/1024)
	}

	acinfo := &activeConnectionInfo{
		DeviceType:   devName,
		Interface:    name,
		HwAddress:    hwAddress,
		IpAddress:    ip,
		SubnetMask:   mask,
		RouteAddress: route,
		Dns1:         dns1,
		Dns2:         dns2,
		Speed:        speed,
	}
	acinfoJSON, _ = marshalJSON(acinfo)
	return
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

// EditConnection open a connection through uuid, return ConnectionSession's dbus object path if success.
func (m *Manager) EditConnection(uuid string, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
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

// TODO remove dbus interface
// GetConnectionPathByUuid return connection setting dbus path by uuid
func (m *Manager) getConnectionPathByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmGetConnectionByUuid(uuid)
	return
}

// TODO remove
// GetActiveConnectionState get current state of the active connection.
func (m *Manager) getActiveConnectionState(apath dbus.ObjectPath) (state uint32) {
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
	_, err = nmActivateConnection(cpath, devPath)
	return
}

// TODO
func (m *Manager) DeactivateConnection(uuid string) (err error) {
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
