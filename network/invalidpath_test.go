package main

import "testing"
import "dlib/dbus"

func init() {
	dbus.InstallOnSession(_Manager)
}

func TestInvalid(t *testing.T) {
	_Manager.ActiveAccessPoint("/", "/")
	_Manager.ActiveWiredDevice(false, "/")
	_Manager.ActiveWiredDevice(true, "/")
	_Manager.DisconnectDevice("/")
	_Manager.GetAccessPoints("/")
	_Manager.SetKey("xxoo", "sd")
	_Manager.GetActiveConnection("/")
	_Manager.GetDBusInfo()
}
