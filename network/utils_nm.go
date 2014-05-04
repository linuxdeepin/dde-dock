package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"
import "strings"

func nmNewDevice(devPath dbus.ObjectPath) (dev *nm.Device, err error) {
	dev, err = nm.NewDevice(nmDest, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewDeviceWired(devPath dbus.ObjectPath) (dev *nm.DeviceWired, err error) {
	dev, err = nm.NewDeviceWired(nmDest, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewDeviceWireless(devPath dbus.ObjectPath) (dev *nm.DeviceWireless, err error) {
	dev, err = nm.NewDeviceWireless(nmDest, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewAccessPoint(apPath dbus.ObjectPath) (ap *nm.AccessPoint, err error) {
	ap, err = nm.NewAccessPoint(nmDest, apPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewActiveConnection(apath dbus.ObjectPath) (ac *nm.ActiveConnection, err error) {
	ac, err = nm.NewActiveConnection(nmDest, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewAgentManager() (manager *nm.AgentManager, err error) {
	manager, err = nm.NewAgentManager(nmDest, "/org/freedesktop/NetworkManager/AgentManager")
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewDHCP4Config(path dbus.ObjectPath) (dhcp4 *nm.DHCP4Config, err error) {
	dhcp4, err = nm.NewDHCP4Config(nmDest, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetDevices() (devPaths []dbus.ObjectPath, err error) {
	devPaths, err = nmManager.GetDevices()
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetWiredDeviceHwAddr(devPath dbus.ObjectPath) (hwAddr string, err error) {
	wiredDev, err := nmNewDeviceWired(devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	hwAddr = strings.ToUpper(wiredDev.HwAddress.Get())
	return
}

func nmGetWirelessDeviceHwAddr(devPath dbus.ObjectPath) (hwAddr string, err error) {
	wirelessDev, err := nmNewDeviceWireless(devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	hwAddr = strings.ToUpper(wirelessDev.HwAddress.Get())
	return
}

func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nm.SettingsConnection, err error) {
	conn, err = nm.NewSettingsConnection(nmDest, cpath)
	if err != nil {
		logger.Error(err)
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

func nmAddAndActivateConnection(data connectionData, devPath dbus.ObjectPath) (cpath, apath dbus.ObjectPath, err error) {
	if len(devPath) == 0 {
		devPath = "/"
	}
	spath := dbus.ObjectPath("/")
	cpath, apath, err = nmManager.AddAndActivateConnection(data, devPath, spath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmActivateConnection(cpath, devPath dbus.ObjectPath) (apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	apath, err = nmManager.ActivateConnection(cpath, devPath, spath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmDeactivateConnection(apath dbus.ObjectPath) (err error) {
	err = nmManager.DeactivateConnection(apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetActiveConnections() (apaths []dbus.ObjectPath) {
	apaths = nmManager.ActiveConnections.Get()
	return
}

func nmGetState() (state uint32) {
	state = nmManager.State.Get()
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

func nmGetConnectionData(cpath dbus.ObjectPath) (data connectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(nmDest, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	data, err = nmConn.GetSettings()
	if err != nil {
		logger.Error(err)
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
		logger.Error("get uuid of connection failed, uuid is empty")
	}
	return
}

func nmGetConnectionType(cpath dbus.ObjectPath) (ctype string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	ctype = getCustomConnectinoType(data)
	if len(ctype) == 0 {
		logger.Error("get type of connection failed, type is empty")
	}
	return
}

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := nmSettings.ListConnections()
	if err != nil {
		logger.Error(err)
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
	cpath, err = nmSettings.GetConnectionByUuid(uuid)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// get wireless connection by ssid, the connection with special hardware address is priority
func nmGetWirelessConnection(ssid []byte, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, ok bool) {
	var hwAddr string
	if len(devPath) != 0 {
		hwAddr, _ = nmGetWirelessDeviceHwAddr(devPath)
	}
	ok = false
	for _, p := range nmGetWirelessConnectionListBySsid(ssid) {
		data, err := nmGetConnectionData(p)
		if err != nil {
			continue
		}
		if isSettingWirelessMacAddressExists(data) {
			if hwAddr == convertMacAddressToString(getSettingWirelessMacAddress(data)) {
				cpath = p
				ok = true
				return
			}
		} else if !ok {
			cpath = p
			ok = true
		}
	}
	return
}

func nmGetWirelessConnectionListBySsid(ssid []byte) (cpaths []dbus.ObjectPath) {
	for _, p := range nmGetConnectionList() {
		data, err := nmGetConnectionData(p)
		if err != nil {
			continue
		}
		if getCustomConnectinoType(data) != typeWireless {
			continue
		}
		if isSettingWirelessSsidExists(data) && string(getSettingWirelessSsid(data)) == string(ssid) {
			cpaths = append(cpaths, p)
		}
	}
	return
}

func nmAddConnection(data connectionData) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.AddConnection(data)
	if err != nil {
		logger.Error(err)
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
