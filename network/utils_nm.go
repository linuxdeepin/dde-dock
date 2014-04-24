package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

func nmNewDevice(devPath dbus.ObjectPath) (dev *nm.Device, err error) {
	dev, err = nm.NewDevice(NMDest, devPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewDeviceWired(devPath dbus.ObjectPath) (dev *nm.DeviceWired, err error) {
	dev, err = nm.NewDeviceWired(NMDest, devPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewDeviceWireless(devPath dbus.ObjectPath) (dev *nm.DeviceWireless, err error) {
	dev, err = nm.NewDeviceWireless(NMDest, devPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewAccessPoint(apPath dbus.ObjectPath) (ap *nm.AccessPoint, err error) {
	ap, err = nm.NewAccessPoint(NMDest, apPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewActiveConnection(apath dbus.ObjectPath) (ac *nm.ActiveConnection, err error) {
	ac, err = nm.NewActiveConnection(NMDest, apath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewAgentManager() (manager *nm.AgentManager, err error) {
	manager, err = nm.NewAgentManager(NMDest, "/org/freedesktop/NetworkManager/AgentManager")
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewDHCP4Config(path dbus.ObjectPath) (dhcp4 *nm.DHCP4Config, err error) {
	dhcp4, err = nm.NewDHCP4Config(NMDest, path)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetDevices() (devPaths []dbus.ObjectPath, err error) {
	devPaths, err = NMManager.GetDevices()
	if err != nil {
		Logger.Error(err)
	}
	return
}

func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nm.SettingsConnection, err error) {
	conn, err = nm.NewSettingsConnection(NMDest, cpath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetDeviceInterface(devPath dbus.ObjectPath) (devInterface string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devInterface = dev.Interface.Get()
	return
}

func nmAddAndActivateConnection(data _ConnectionData, devPath dbus.ObjectPath) (cpath, apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	cpath, apath, err = NMManager.AddAndActivateConnection(data, devPath, spath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmActivateConnection(cpath, devPath dbus.ObjectPath) (apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	apath, err = NMManager.ActivateConnection(cpath, devPath, spath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmDeactivateConnection(apath dbus.ObjectPath) (err error) {
	err = NMManager.DeactivateConnection(apath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetActiveConnections() (apaths []dbus.ObjectPath) {
	apaths = NMManager.ActiveConnections.Get()
	return
}

func nmGetState() (state uint32) {
	state = NMManager.State.Get()
	return
}

func nmGetActiveConnectionByUuid(uuid string) (apath dbus.ObjectPath, ok bool) {
	for _, apath = range nmGetActiveConnections() {
		if ac, err := nmNewActiveConnection(apath); err == nil {
			if ac.Uuid.Get() == uuid {
				ok = true
				return
			}
		}
	}
	ok = false
	return
}

func nmGetConnectionData(cpath dbus.ObjectPath) (data _ConnectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(NMDest, cpath)
	if err != nil {
		Logger.Error(err)
		return
	}
	data, err = nmConn.GetSettings()
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetConnectionUuid(cpath dbus.ObjectPath) (uuid string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	uuid = getSettingConnectionUuid(data)
	if len(uuid) == 0 {
		Logger.Error("get uuid of connection failed, uuid is empty")
	}
	return
}

func nmGetConnectionType(cpath dbus.ObjectPath) (ctype string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	ctype = generalGetConnectionType(data)
	if len(ctype) == 0 {
		Logger.Error("get type of connection failed, type is empty")
	}
	return
}

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := NMSettings.ListConnections()
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetConnectionById(id string) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if getSettingConnectionId(data) == id {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetWirelessConnectionBySsid(ssid []byte) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if isSettingWirelessSsidExists(data) && string(getSettingWirelessSsid(data)) == string(ssid) {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmAddConnection(data _ConnectionData) (cpath dbus.ObjectPath, err error) {
	cpath, err = NMSettings.AddConnection(data)
	if err != nil {
		Logger.Error(err)
	}
	return
}

func nmGetDHCP4Info(path dbus.ObjectPath) (ip string, mask string, route string) {
	dhcp4, err := nmNewDHCP4Config(path)
	if err != nil {
		return
	}
	options := dhcp4.Options.Get()
	if ipData, ok := options["ip_address"]; ok {
		ip, _ = ipData.Value().(string)
	}
	if maskData, ok := options["subnet_mask"]; ok {
		mask, _ = maskData.Value().(string)
	}
	if routeData, ok := options["routers"]; ok {
		route, _ = routeData.Value().(string)
	}
	return
}
