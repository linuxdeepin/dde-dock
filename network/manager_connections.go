package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type activeConnectionInfo struct {
	Interface    string
	HWAddress    string
	IPAddress    string
	SubnetMask   string
	RouteAddress string
	Speed        string
}

func (m *Manager) initConnectionManage() {
	m.VPNConnections = make([]string, 0)
	m.WiredConnections = make([]string, 0)
	m.WirelessConnections = make([]string, 0)

	m.updatePropWiredConnections()

	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(opAdded, c)
	}
	nmSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(opAdded, path)
	})
}

// GetWiredConnectionUuid return connection uuid for target wired device.
func (m *Manager) GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string) {
	// check if target wired connection exists, if not, create one
	id := "wired-connection-" + nmGetDeviceInterface(wiredDevPath)
	cpath, ok := nmGetConnectionById(id)
	if ok {
		uuid = nmGetConnectionUuid(cpath)
	} else {
		uuid = newWiredConnection(id)
	}
	return
}

func (m *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	logger.Debugf("handleConnectionChanged: operation %d, path %s", operation, path)
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

		switch getSettingConnectionType(cdata) {
		case "802-11-wireless": // TODO
			m.WirelessConnections = append(m.WirelessConnections, uuid)
			m.updatePropWirelessConnections()
		case "802-3-ethernet": // wired connection will be treatment specially
			// TODO remove
			// m.WiredConnections = append(m.WiredConnections, uuid)
			// dbus.NotifyChange(m, "WiredConnections")
		case "pppoe":
		case "vpn":
			m.VPNConnections = append(m.VPNConnections, uuid)
			m.updatePropVpnConnections()
		case "cdma":
		}
	case opRemoved:
		//TODO:
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

// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() (typesJSON string) {
	typesJSON, _ = marshalJSON(supportedConnectionTypesInfo)
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

// GetConnectionByUuid return connection setting dbus path by uuid
func (m *Manager) GetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
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

// use disconnect device instead
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
