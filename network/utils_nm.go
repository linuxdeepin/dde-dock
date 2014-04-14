package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

func nmNewDevice(devPath dbus.ObjectPath) (dev *nm.Device, ok bool) {
	dev, err := nm.NewDevice(NMDest, devPath)
	if err != nil {
		LOGGER.Error(err)
		ok = false
		return
	}
	ok = true
	return
}

func nmGetDeviceInterface(devPath dbus.ObjectPath) (devInterface string) {
	dev, ok := nmNewDevice(devPath)
	if !ok {
		return
	}
	devInterface = dev.Interface.Get()
	return
}

func nmGetConnectionData(cpath dbus.ObjectPath) (data _ConnectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(NMDest, cpath)
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

func nmGetConnectionUuid(cpath dbus.ObjectPath) (uuid string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	uuid = getSettingConnectionUuid(data)
	if len(uuid) == 0 {
		LOGGER.Error("get uuid of connection failed, uuid is empty")
	}
	return
}

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := _NMSettings.ListConnections()
	if err != nil {
		LOGGER.Error(err)
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

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if getSettingConnectionUuid(data) == uuid {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmAddConnection(data _ConnectionData) {
	_, err := _NMSettings.AddConnection(data)
	if err != nil {
		LOGGER.Error(err)
	}
	return
}
