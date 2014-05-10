package main

import (
	idbus "dbus/org/freedesktop/dbus/system"
	"dlib/dbus"
)

const (
	dbusBluezDest       = "org.bluez"
	dbusBluezPath       = "/org/bluez"
	dbusBluezIfsAdapter = "org.bluez.Adapter1"
	dbusBluezIfsDevice  = "org.bluez.Device1"

	dbusBluetoothDest = "com.deepin.daemon.Bluetooth"
	dbusBluetoothPath = "/com/deepin/daemon/Bluetooth"
	dbusBluetoothIfs  = "com.deepin.daemon.Bluetooth"
)

var bluezObjectManager *idbus.ObjectManager

type dbusObjectData map[string]dbus.Variant
type dbusInterfaceData map[string]map[string]dbus.Variant
type dbusInterfacesData map[dbus.ObjectPath]map[string]map[string]dbus.Variant

type Bluetooth struct {
	objects dbusInterfacesData

	// TODO adapter
	PrimaryAdapter string `access:"readwrite"`
	adapters       []*adapter
	Adapters       string // array of adapters that marshaled by json
	// Adapters       []string // array of adapter names

	// device
	devices []*device
	// Devices []dbus.ObjectPath
	Devices string // device objects that marshaled by json

	// alias properties for primary adapter
	Alias   string `access:"readwrite"`
	Powered bool   `access:"readwrite"`
	// Discovering         bool // TODO merge to Powered
	Discoverable        bool   `access:"readwrite"`
	DiscoverableTimeout uint32 `access:"readwrite"`

	// Pairable        bool   `access:"readwrite"`
	// PairableTimeout uint32 `access:"readwrite"`

	// TODO after switching primary adapter, power off other adapters

	// signals
	DeviceAdded   func(devJSON string)
	DeviceRemoved func(devJSON string)
}

func NewBluetooth() (bluettoth *Bluetooth) {
	bluettoth = &Bluetooth{}
	// TODO
	return
}

func (b *Bluetooth) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		dbusBluetoothDest,
		dbusBluetoothPath,
		dbusBluetoothIfs,
	}
}

func (b *Bluetooth) initBluetooth() {
	// initialize dbus object manager
	var err error
	bluezObjectManager, err = idbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		panic(err)
	}
	b.objects, err = bluezObjectManager.GetManagedObjects()
	if err != nil {
		panic(err)
	}

	// add exists adapters and devices
	for path, data := range b.objects {
		b.handleInterfacesAdded(path, data)
	}

	// connect signals
	bluezObjectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
	bluezObjectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)

	// TODO update properties
}
func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[dbusBluezIfsAdapter]; ok {
		b.addAdapter(path)
	}
	if _, ok := data[dbusBluezIfsDevice]; ok {
		b.addDevice(path, data[dbusBluezIfsDevice])
	}
}
func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(dbusBluezIfsAdapter, interfaces) {
		b.removeAdapter(path)
	}
	if isStringInArray(dbusBluezIfsDevice, interfaces) {
		b.removeDevice(path)
	}
}

// GetDevices return all devices object that marshaled by json.
func (b *Bluetooth) GetDevices() (devicesJSON string) {
	// TODO
	return
}

// TODO
func (b *Bluetooth) RemoveDevice(dpath dbus.ObjectPath) (err error) {
	return
}
