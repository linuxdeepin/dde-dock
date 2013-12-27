package main

import "testing"
import "dlib/dbus"
import nm "dbus/com/deepin/daemon/network"
import "fmt"

func init() {
	dbus.InstallOnSession(_Manager)
}

func TestInvalid(t *testing.T) {
	_Manager.ActiveAccessPoint("/", "/")

	_Manager.DisconnectDevice("/")

	_Manager.ActiveWiredDevice(false, "/")
	_Manager.ActiveWiredDevice(true, "/")

	_Manager.GetAccessPoints("/")

	_Manager.GetActiveConnection("/")

	_Manager.GetConnectionByAccessPoint("/")

	_Manager.SetKey("xxoo", "sd")
	_Manager.GetDBusInfo()

	dumy := make(map[string]map[string]string)
	_Manager.UpdateConnection(dumy)
}

func TestDBusFailed(t *testing.T) {
	if _, err := nm.NewNetworkManager("/"); err == nil {
		t.Fail()
	}

	m, err := nm.NewNetworkManager("/com/deepin/daemon/Network")
	if err != nil {
		t.Fatal("NewNetworkManager1")
	}
	if _, err := nm.NewNetworkManager("/com/deepin/daemon/Networkxx"); err == nil {
		t.Fatal("NewNetworkManager2")
	}
	if err = m.ActiveAccessPoint("/", "/"); err == nil {
		t.Fatal("ActiveAccessPoint")
	}

	if err = m.DisconnectDevice("/"); err == nil {
		t.Fatal("DisconnectDevice")
	}
	if err = m.ActiveWiredDevice(false, "/"); err == nil {
		t.Fatal("ActiveWiredDevice false")
	}
	if err = m.ActiveWiredDevice(true, "/"); err == nil {
		t.Fatal("ActiveWiredDevice true")
	}
	if _, err = m.GetAccessPoints("/"); err == nil {
		t.Fatal("GetAccessPoints")
	}
	if _, err = m.GetActiveConnection("/"); err == nil {
		t.Fatal("GetActiveConnection")
	}
	if _, err = m.GetConnectionByAccessPoint("/"); err == nil {
		t.Fatal("GetConnectionByAccessPoint")
	}
	if err = m.SetKey("xxoo", "sd"); err == nil {
		//This wouldn't failed
	}
}

func TestDBusSuccess(t *testing.T) {
	m, err := nm.NewNetworkManager("/com/deepin/daemon/Network")
	if err != nil {
		t.Fatal("Create NetworkManager failed")
	}
	for _, d := range m.WirelessDevices.Get() {
		fmt.Println(d)
	}
}
