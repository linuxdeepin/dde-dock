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
		case "802-11-wireless":
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

func newWirelessConnection(id string, ssid []byte, keyFlag int) string {
	// TODO FIXME
	if keyFlag == ApKeyEap {
		LOGGER.Debug("ignore wireless connection:", id)
		return ""
	}

	uuid := newUUID()
	data := newWirelessConnectionData(id, uuid, ssid, keyFlag)

	// LOGGER.Debug("new wireless connection data:", data) // TODO test

	_, err := _NMSettings.AddConnection(data)
	if err != nil {
		// panic(err) // TODO fixme
		LOGGER.Error(err)
	}

	return uuid
}

func (this *Manager) GetConnectionByAccessPoint(path dbus.ObjectPath) (string, error) {
	if ap, err := nm.NewAccessPoint(NMDest, path); err == nil {
		for _, c := range this.WirelessConnections {
			if nmConn, err := nm.NewSettingsConnection(NMDest, path); err == nil {
				if cdata, err := nmConn.GetSettings(); err == nil {
					if string(getSettingWirelessSsid(cdata)) == string(ap.Ssid.Get()) { // TODO
						return c, nil
					}
				}
			}

		}
		LOGGER.Debug("CCC:", path, string(ap.Ssid.Get()))
		return newWirelessConnection(string(ap.Ssid.Get()), []byte(ap.Ssid.Get()), parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())), nil
	} else {
		return "", dbus.NewNoObjectError(path)
	}
}

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
