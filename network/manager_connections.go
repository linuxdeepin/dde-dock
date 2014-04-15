package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type _ConnectionData map[string]map[string]dbus.Variant

type ActiveConnection struct {
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

	m.initWiredConnections()

	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(OpAdded, c)
	}
	NMSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(OpAdded, path)
	})
}

// create connection for each wired device if not exists
func (m *Manager) initWiredConnections() {
	for _, wiredDev := range m.WiredDevices {
		uuid := m.GetWiredConnectionUuid(wiredDev.Path)
		m.WiredConnections = append(m.WiredConnections, uuid)
	}
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
	Logger.Debugf("handleConnectionChanged: operation %d, path %s", operation, path)
	switch operation {
	case OpAdded:
		nmConn, _ := nmNewSettingsConnection(path)
		nmConn.ConnectRemoved(func() {
			m.handleConnectionChanged(OpRemoved, path)
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
			dbus.NotifyChange(m, "WirelessConnections")
		case "802-3-ethernet": // wired connection will be treatment specially
			// TODO remove
			// m.WiredConnections = append(m.WiredConnections, uuid)
			// dbus.NotifyChange(m, "WiredConnections")
		case "pppoe":
		case "vpn":
			m.VPNConnections = append(m.VPNConnections, uuid)
			dbus.NotifyChange(m, "VPNConnections")
		case "cdma":
		}
	case OpRemoved:
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

func newWiredConnection(id string) (uuid string) {
	Logger.Debugf("new wired connection, id=%s", id)
	uuid = newUUID()
	data := newWiredConnectionData(id, uuid)
	nmAddConnection(data)
	return
}

func newWirelessConnection(id string, ssid []byte, keyFlag int) (uuid string) {
	Logger.Debugf("new wireless connection, id=%s, ssid=%s, keyFlag=%d", id, ssid, keyFlag)
	uuid = newUUID()
	data := newWirelessConnectionData(id, uuid, ssid, keyFlag)
	nmAddConnection(data)
	return
}

func newPppoeConnection(id string) (uuid string) {
	Logger.Debugf("new pppoe connection, id=%s", id)
	uuid = newUUID()
	data := newPppoeConnectionData(id, uuid)
	nmAddConnection(data)
	return
}

// TODO [remove or rename] GetConnectionByAccessPoint return the connection's uuid of access point, return empty if none.
func (m *Manager) GetConnectionByAccessPoint(apPath dbus.ObjectPath) (uuid string, err error) {
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	cpath, ok := nmGetWirelessConnectionBySsid(ap.Ssid.Get())
	if !ok {
		return
	}

	uuid = nmGetConnectionUuid(cpath)
	return

	// TODO remove
	// conns, err := NMSettings.ListConnections()
	// if err != nil {
	// 	Logger.Error(err)
	// 	return
	// }

	// for _, cpath := range conns {
	// 	if nmConn, err := nmNewSettingsConnection(cpath); err == nil {
	// 		if cdata, err := nmConn.GetSettings(); err == nil {
	// 			if isSettingWirelessSsidExists(cdata) && string(getSettingWirelessSsid(cdata)) == string(ap.Ssid.Get()) {
	// 				uuid = getSettingConnectionUuid(cdata)
	// 				Logger.Debug("connection is already exists", apPath, uuid)
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// TODO remove
	// for _, conUuid := range m.WirelessConnections {
	// 	if cpath, err := NMSettings.GetConnectionByUuid(conUuid); err == nil {
	// 		if nmConn, err := nmNewSettingsConnection(cpath); err == nil {
	// 			if cdata, err := nmConn.GetSettings(); err == nil {
	// 				if string(getSettingWirelessSsid(cdata)) == string(ap.Ssid.Get()) { // TODO
	// 					Logger.Debug("connection is already exists", apPath, conUuid)
	// 					uuid = conUuid
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// Logger.Debugf("GetConnectionByAccessPoint: apPath=%s, uuid=%s", apPath, uuid) // TODO test

	return
}

// CreateConnectionByAccessPoint create connection for access point and return the uuid.
func (m *Manager) CreateConnectionForAccessPoint(apPath dbus.ObjectPath) (uuid string, err error) {
	Logger.Debug("CreateConnectionForAccessPoint: apPath", apPath)
	uuid, err = m.GetConnectionByAccessPoint(apPath)
	if len(uuid) != 0 {
		// connection already exists
		return
	}

	// create connection
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	// TODO FIXME
	keyFlag := parseFlags(ap)
	if keyFlag == ApKeyEap {
		Logger.Debug("ignore wireless connection:", string(ap.Ssid.Get()))
		return "", dbus.NewNoObjectError(apPath)
	}

	uuid = newWirelessConnection(string(ap.Ssid.Get()), []byte(ap.Ssid.Get()), parseFlags(ap))
	return
}

// TODO rename to GetActiveConnectionInfo
func (m *Manager) GetActiveConnection(devPath dbus.ObjectPath) (ret *ActiveConnection, err error) {
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

	ip, mask, route := parseDHCP4(dev.Dhcp4Config.Get())
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

	return &ActiveConnection{
		Interface:    name,
		HWAddress:    macaddress,
		IPAddress:    ip,
		SubnetMask:   mask,
		RouteAddress: route,
		Speed:        speed,
	}, nil
}

// CreateConnection create a new connection, return ConnectionSession's dbus object path if success.
func (m *Manager) CreateConnection(connType string) (session *ConnectionSession, err error) {
	session, err = NewConnectionSessionByCreate(connType)
	if err != nil {
		Logger.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		Logger.Error(err)
		return
	}

	return
}

// OpenConnection open a connection through uuid, return ConnectionSession's dbus object path if success.
func (m *Manager) EditConnection(uuid string) (session *ConnectionSession, err error) {
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
		Logger.Debug("read only connection, create new")
		return m.CreateConnection(getSettingConnectionType(connData))
	}

	session, err = NewConnectionSessionByOpen(uuid)
	if err != nil {
		Logger.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		Logger.Error(err)
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

// GetActiveConnectionState get current state of the active connection.
func (m *Manager) GetActiveConnectionState(apath dbus.ObjectPath) (state uint32) {
	conn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}
	state = conn.State.Get()
	return
}

// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() []string {
	return supportedConnectionTypes
}

func (m *Manager) ActivateConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (err error) {
	Logger.Debugf("ActivateConnectionForAccessPoint: apPath=%s, devPath=%s", apPath, devPath)
	// if there is no connection for current access point, create one
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	cpath, ok := nmGetWirelessConnectionBySsid(ap.Ssid.Get())
	if ok {
		_, err = nmActivateConnection(cpath, devPath)
	} else {
		uuid := newUUID()
		data := newWirelessConnectionData(string(ap.Ssid.Get()), uuid, []byte(ap.Ssid.Get()), parseFlags(ap))
		_, _, err = nmAddAndActivateConnection(data, devPath)
	}
	return
}

func (m *Manager) ActivateConnection(uuid string, devPath dbus.ObjectPath) (err error) {
	Logger.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, devPath)
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	// TODO fixme
	// if only one access point connection, do nothing for it will be
	// activate by network manager automatic
	if nmGetConnectionType(cpath) == typeWireless {
		count := 0
		for _, tmpcpath := range nmGetConnectionList() {
			ctype := nmGetConnectionType(tmpcpath)
			if ctype == typeWireless {
				count++
			}
		}
		if count <= 1 {
			Logger.Debug("only one access point connection, will be activate by network manager automatic")
			return
		}
	}

	_, err = nmActivateConnection(cpath, devPath)
	return
}

// TODO remove
func (m *Manager) DeactivateConnection(uuid string) (err error) {
	apath, ok := nmGetActiveConnectionByUuid(uuid)
	if !ok {
		Logger.Error("not found active connection with uuid", uuid)
		return
	}
	Logger.Debug("DeactivateConnection:", uuid, apath)
	err = nmDeactivateConnection(apath)
	return
}
