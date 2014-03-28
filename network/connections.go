package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "fmt"

type _ConnectionData map[string]map[string]dbus.Variant

type Connection struct {
	Data _ConnectionData

	Path           dbus.ObjectPath
	Uuid           string
	Name           string
	ConnectionType string // TODO
}

func (this *Manager) initConnectionManage() {
	this.VPNConnections = make([]*Connection, 0)
	this.WiredConnections = make([]*Connection, 0)
	this.WirelessConnections = make([]*Connection, 0)

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
		c := NewConnection(nmConn)
		switch c.ConnectionType {
		case "802-11-wireless":
			this.WirelessConnections = append(this.WirelessConnections, c)
			dbus.NotifyChange(this, "WirelessConnections")
		case "802-3-ethernet":
			this.WiredConnections = append(this.WiredConnections, c)
			dbus.NotifyChange(this, "WiredConnections")
		case "pppoe":
		case "vpn":
			this.VPNConnections = append(this.VPNConnections, c)
			dbus.NotifyChange(this, "VPNConnections")
		case "cdma":
		}
	case OpRemoved:
		removed := false
		if this.WirelessConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.WirelessConnections); removed {
			dbus.NotifyChange(this, "WirelessConnections")
		} else if this.WiredConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.WiredConnections); removed {
			dbus.NotifyChange(this, "WiredConnections")
		} else if this.VPNConnections, removed = tryRemoveConnection(dbus.ObjectPath(path), this.VPNConnections); removed {
			dbus.NotifyChange(this, "VPNConnections")
		}
	}

}

// TODO
func NewConnection(core *nm.SettingsConnection) *Connection {
	c := &Connection{}
	settings, err := core.GetSettings()
	if err != nil {
		return c // TODO still return even error occured?
	}
	c.Path = core.Path
	// TODO remove
	// c.Name = settings["connection"]["id"].Value().(string)
	// c.Uuid = settings["connection"]["uuid"].Value().(string)
	// c.ConnectionType = settings["connection"]["type"].Value().(string)
	c.Name = getSettingConnectionId(settings)
	c.Uuid = getSettingConnectionUuid(settings)
	c.ConnectionType = getSettingConnectionType(settings)
	c.Data, err = core.GetSettings() // TODO need GetSettings() again?
	return c
}

func newWirelessConnection(id string, ssid string, keyFlag int) *Connection {
	data := make(_ConnectionData)

	uuid := newUUID()
	initWirelessConnection(data, id, uuid, ssid, keyFlag)
	LOGGER.Debugf("%v", data)
	// map[connection:map[uuid:"9c60bc6f-d2ac-4571-a06f-58b0d5d22eac" type:"802-11-wireless"] 802-11-wireless:map[security:"802-11-wireless-security"] 802-11-wireless-security:map[auth-alg:"open"] ipv4:map[] ipv6:map[]]
	// map[connection:map[id:"CMCC-AUTO" uuid:"8b135d1c-6d3b-4797-a992-304a82b95d07" type:"802-11-wireless"] ipv4:map[method:"auto"] ipv6:map[method:"auto"] 802-11-wireless:map[ssid:@ay [0x43, 0x4d, 0x43, 0x43, 0x2d, 0x41, 0x55, 0x54, 0x4f] security:"802-11-wireless-security"] 802-11-wireless-security:map[key-mgmt:"wpa-eap" auth-alg:"open"]]

	// TODO
	newConn, err := _NMSettings.AddConnection(data)
	if err != nil {
		panic(err)
	}
	core, err := nm.NewSettingsConnection(NMDest, newConn)
	if err != nil {
		panic(err)
	}
	return &Connection{data, core.Path, uuid, id, fieldWireless}
}

func (this *Manager) GetConnectionByAccessPoint(path dbus.ObjectPath) (*Connection, error) {
	if ap, err := nm.NewAccessPoint(NMDest, path); err == nil {
		for _, c := range this.WirelessConnections {
			// TODO
			if c.ConnectionType == fieldWireless && string(c.Data[fieldWireless]["ssid"].Value().([]uint8)) == string(ap.Ssid.Get()) {
				return c, nil
			}
		}
		fmt.Println("CCC:", path, string(ap.Ssid.Get()))
		return newWirelessConnection(string(ap.Ssid.Get()), string(ap.Ssid.Get()), parseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())), nil
	} else {
		return nil, dbus.NewNoObjectError(path)
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
		name = NewConnection(c).Name
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

func (this *Manager) getConnectionPathByUUID(uuid string) (path dbus.ObjectPath, ok bool) {
	ok = false
	for _, c := range this.WiredConnections {
		if uuid == c.Uuid {
			path = c.Path
			ok = true
		}
	}
	for _, c := range this.WirelessConnections {
		if uuid == c.Uuid {
			path = c.Path
			ok = true
		}
	}
	for _, c := range this.VPNConnections {
		if uuid == c.Uuid {
			path = c.Path
			ok = true
		}
	}
	return
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
	// TODO
	return
}
