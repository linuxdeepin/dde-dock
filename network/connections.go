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
	ConnectionType string
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

func tryRemoveConnection(path dbus.ObjectPath, conns []*Connection) ([]*Connection, bool) {
	var newConns []*Connection
	found := false
	for _, conn := range conns {
		if conn.Path != path {
			newConns = append(newConns, conn)
		} else {
			found = true
		}
	}
	return newConns, found
}

func (this *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	switch operation {
	case OpAdded:
		nmConn, _ := nm.NewSettingsConnection(path)
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

func NewConnection(core *nm.SettingsConnection) *Connection {
	c := &Connection{}
	settings, err := core.GetSettings()
	if err != nil {
		return c
	}
	c.Path = core.Path
	c.Name = settings["connection"]["id"].Value().(string)
	c.Uuid = settings["connection"]["uuid"].Value().(string)
	c.ConnectionType = settings["connection"]["type"].Value().(string)
	c.Data, err = core.GetSettings()
	return c
}

func newWirelessConnection(id string, ssid string, keyFlag int) *Connection {
	data := make(_ConnectionData)
	data[fieldConnection] = make(map[string]dbus.Variant)
	data[fieldIPv4] = make(map[string]dbus.Variant)
	data[fieldIPv6] = make(map[string]dbus.Variant)
	data[fieldWireless] = make(map[string]dbus.Variant)

	data[fieldConnection]["id"] = dbus.MakeVariant(id)
	uuid := newUUID()
	data[fieldConnection]["uuid"] = dbus.MakeVariant(uuid)
	data[fieldConnection]["type"] = dbus.MakeVariant(fieldWireless)

	data[fieldWireless]["ssid"] = dbus.MakeVariant([]uint8(ssid))

	if keyFlag != ApKeyNone {
		data[fieldWirelessSecurity] = make(map[string]dbus.Variant)
		data[fieldWireless]["security"] = dbus.MakeVariant(fieldWirelessSecurity)
		switch keyFlag {
		case ApKeyWep:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("none")
		case ApKeyPsk:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-psk")
			data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
		case ApKeyEap:
			data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-eap")
			data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
		}
	}

	data[fieldIPv4]["method"] = dbus.MakeVariant("auto")

	data[fieldIPv6]["method"] = dbus.MakeVariant("auto")

	newConn, err := _NMSettings.AddConnection(data)
	core, err := nm.NewSettingsConnection(newConn)
	if err != nil {
		panic(err)
	}
	return &Connection{data, core.Path, uuid, id, fieldWireless}
}

func (this *Manager) GetConnectionByAccessPoint(path dbus.ObjectPath) (*Connection, error) {
	if ap, err := nm.NewAccessPoint(path); err == nil {
		for _, c := range this.WirelessConnections {
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

func (this *Manager) GetActiveConnection(devPath dbus.ObjectPath) (ActiveConnection, error) {
	dev, err := nm.NewDevice(devPath)
	if err != nil {
		return ActiveConnection{}, err
	}
	ac, _ := nm.NewActiveConnection(dev.ActiveConnection.Get())
	name := ""
	if c, err := nm.NewSettingsConnection(ac.Connection.Get()); err == nil {
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
		_dev, _ := nm.NewDeviceWired(devPath)
		macaddress = _dev.HwAddress.Get()
		speed = fmt.Sprintf("%d", _dev.Speed.Get())
		nm.DestroyDeviceWired(_dev)
	case NM_DEVICE_TYPE_WIFI:
		_dev, _ := nm.NewDeviceWireless(devPath)
		macaddress = _dev.HwAddress.Get()
		speed = fmt.Sprintf("%d", _dev.Bitrate.Get()/1024)
		nm.DestroyDeviceWireless(_dev)
	}

	return ActiveConnection{
		Interface:    name,
		HWAddress:    macaddress,
		IPAddress:    ip,
		SubnetMask:   mask,
		RouteAddress: route,
		Speed:        speed,
	}, nil
}

func (this *Manager) UpdateConnection(data map[string]map[string]string) {
	/*func (this *Manager) UpdateConnection(data string) {*/
	fmt.Println("Update:", data)
}
