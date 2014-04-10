package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type _ConnectionData map[string]map[string]dbus.Variant

func (this *Manager) initConnectionManage() {
	this.VPNConnections = make([]string, 0)
	this.WiredConnections = make([]string, 0)
	this.WirelessConnections = make([]string, 0)

	conns, err := _NMSettings.ListConnections()
	if err != nil {
		panic(err)
	}
	for _, c := range conns {
		this.handleConnectionChanged(OpAdded, c)
	}
	_NMSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		this.handleConnectionChanged(OpAdded, path)
	})
}

func (this *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	LOGGER.Debugf("handleConnectionChanged: operation %d, path %s", operation, path)
	switch operation {
	case OpAdded:
		nmConn, _ := nm.NewSettingsConnection(NMDest, path)
		nmConn.ConnectRemoved(func() {
			this.handleConnectionChanged(OpRemoved, path)
			nm.DestroySettingsConnection(nmConn)
		})
		cdata, err := nmConn.GetSettings()
		if err != nil {
			return
		}
		uuid := getSettingConnectionUuid(cdata)

		switch getSettingConnectionType(cdata) {
		case "802-11-wireless": // TODO
			this.WirelessConnections = append(this.WirelessConnections, uuid)
			dbus.NotifyChange(this, "WirelessConnections")
		case "802-3-ethernet":
			this.WiredConnections = append(this.WiredConnections, uuid)
			dbus.NotifyChange(this, "WiredConnections")
		case "pppoe":
		case "vpn":
			this.VPNConnections = append(this.VPNConnections, uuid)
			dbus.NotifyChange(this, "VPNConnections")
		case "cdma":
		}
	case OpRemoved:
		//TODO:
		//removed := false
		//if this.WirelessConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.WirelessConnections); removed {
		//dbus.NotifyChange(this, "WirelessConnections")
		//} else if this.WiredConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.WiredConnections); removed {
		//dbus.NotifyChange(this, "WiredConnections")
		//} else if this.VPNConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.VPNConnections); removed {
		//dbus.NotifyChange(this, "VPNConnections")
		//}
	}

}

func newWirelessConnection(id string, ssid []byte, keyFlag int) (uuid string) {
	uuid = newUUID()
	data := newWirelessConnectionData(id, uuid, ssid, keyFlag)

	LOGGER.Debug("new wireless connection data:", data) // TODO test

	// TODO
	_, err := _NMSettings.AddConnection(data)
	if err != nil {
		// panic(err) // TODO fixme
		LOGGER.Error(err)
	}

	return
}

// GetConnectionByAccessPoint return the connection's uuid of access point, return empty if none.
func (this *Manager) GetConnectionByAccessPoint(apPath dbus.ObjectPath) (uuid string, err error) {
	ap, err := nm.NewAccessPoint(NMDest, apPath)
	if err != nil {
		LOGGER.Error(err)
		return
	}

	conns, err := _NMSettings.ListConnections()
	if err != nil {
		LOGGER.Error(err)
		return
	}

	for _, cpath := range conns {
		if nmConn, err := nm.NewSettingsConnection(NMDest, cpath); err == nil {
			if cdata, err := nmConn.GetSettings(); err == nil {
				if isSettingWirelessSsidExists(cdata) && string(getSettingWirelessSsid(cdata)) == string(ap.Ssid.Get()) {
					uuid = getSettingConnectionUuid(cdata)
					LOGGER.Debug("connection is already exists", apPath, uuid)
					break
				}
			}
		}
	}
	// TODO remove
	// for _, conUuid := range this.WirelessConnections {
	// 	if cpath, err := _NMSettings.GetConnectionByUuid(conUuid); err == nil {
	// 		if nmConn, err := nm.NewSettingsConnection(NMDest, cpath); err == nil {
	// 			if cdata, err := nmConn.GetSettings(); err == nil {
	// 				if string(getSettingWirelessSsid(cdata)) == string(ap.Ssid.Get()) { // TODO
	// 					LOGGER.Debug("connection is already exists", apPath, conUuid)
	// 					uuid = conUuid
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// LOGGER.Debugf("GetConnectionByAccessPoint: apPath=%s, uuid=%s", apPath, uuid) // TODO test

	return
}

// CreateConnectionByAccessPoint create connection for access point and return the uuid.
func (this *Manager) CreateConnectionForAccessPoint(apPath dbus.ObjectPath) (uuid string, err error) {
	LOGGER.Debug("CreateConnectionForAccessPoint: apPath", apPath)
	uuid, err = this.GetConnectionByAccessPoint(apPath)
	if len(uuid) != 0 {
		// connection already exists
		return
	}

	// create connection
	if ap, err := nm.NewAccessPoint(NMDest, apPath); err == nil {
		// TODO FIXME
		keyFlag := parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())
		if keyFlag == ApKeyEap {
			LOGGER.Debug("ignore wireless connection:", string(ap.Ssid.Get()))
			return "", dbus.NewNoObjectError(apPath)
		}

		uuid = newWirelessConnection(string(ap.Ssid.Get()), []byte(ap.Ssid.Get()), parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get()))
	} else {
		err = dbus.NewNoObjectError(apPath)
	}
	return
}

// TODO
func (this *Manager) GetActiveConnection(devPath dbus.ObjectPath) (ret *ActiveConnection, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
		}
	}()
	dev, err := nm.NewDevice(NMDest, devPath)
	if err != nil {
		return nil, err
	}
	ac, err := nm.NewActiveConnection(NMDest, dev.ActiveConnection.Get())
	if err != nil {
		return nil, err
	}
	name := ""
	if c, err := nm.NewSettingsConnection(NMDest, ac.Connection.Get()); err != nil {
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
		_dev, _ := nm.NewDeviceWired(NMDest, devPath)
		macaddress = _dev.HwAddress.Get()
		speed = fmt.Sprintf("%d", _dev.Speed.Get())
		nm.DestroyDeviceWired(_dev)
	case NM_DEVICE_TYPE_WIFI:
		_dev, _ := nm.NewDeviceWireless(NMDest, devPath)
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
func (this *Manager) CreateConnection(connType string) (session *ConnectionSession, err error) {
	session, err = NewConnectionSessionByCreate(connType)
	if err != nil {
		LOGGER.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		LOGGER.Error(err)
		return
	}

	return
}

// OpenConnection open a connection through uuid, return ConnectionSession's dbus object path if success.
func (this *Manager) EditConnection(uuid string) (session *ConnectionSession, err error) {
	// if is read only connection(default system connection created by
	// network manager), create a new connection
	connData := this.getConnectionDataByUuid(uuid)
	if getSettingConnectionReadOnly(connData) {
		LOGGER.Debug("read only connection, create new")
		return this.CreateConnection(getSettingConnectionType(connData))
	}

	session, err = NewConnectionSessionByOpen(uuid)
	if err != nil {
		LOGGER.Error(err)
		return
	}

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		LOGGER.Error(err)
		return
	}

	return
}

func (this *Manager) getConnectionDataByUuid(uuid string) (data _ConnectionData) {
	coreObjPath, err := _NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		LOGGER.Error(err)
		return
	}
	nmConn, err := nm.NewSettingsConnection(NMDest, coreObjPath)
	if err != nil {
		LOGGER.Error(err)
		return
	}
	data, err = nmConn.GetSettings()
	if err != nil {
		LOGGER.Error(err)
		return
	}
	return
}

// DeleteConnection delete a connection through uuid.
func (this *Manager) DeleteConnection(uuid string) (err error) {
	//TODO: remove(uninstall dbus) editing connection_session object
	cpath, err := _NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		return err
	}
	conn, err := nm.NewSettingsConnection(NMDest, cpath)
	if err != nil {
		return err
	}
	return conn.Delete()
}

// GetConnectionByUuid return connection setting dbus path by uuid
func (this *Manager) GetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = _NMSettings.GetConnectionByUuid(uuid)
	return
}

// GetActiveConnectionState get current state of the active connection.
func (this *Manager) GetActiveConnectionState(uuid string) (state uint32) {
	cpath, err := _NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		LOGGER.Error("GetActiveConnectionState,", err)
		return
	}
	conn, err := nm.NewActiveConnection(NMDest, cpath)
	if err != nil {
		LOGGER.Error("GetActiveConnectionState,", err)
		return
	}
	state = conn.State.Get()
	return
}

// GetSupportedConnectionTypes return all supported connection types
func (this *Manager) GetSupportedConnectionTypes() []string {
	return supportedConnectionTypes
}
