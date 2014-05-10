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

type bluezObjectData map[string]map[string]dbus.Variant
type bluezObjectsData map[dbus.ObjectPath]map[string]map[string]dbus.Variant

type Bluetooth struct {
	objects bluezObjectsData

	// TODO adapter
	PrimaryAdapter dbus.ObjectPath `access:"readwrite"`
	adapters       []*adapter
	Adapters       []dbus.ObjectPath // hci object paths

	// device
	devices []*device
	// Devices []dbus.ObjectPath
	Devices string // device objects that marshaled by json

	// property for primary adapter
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
	// TODO
	var err error
	bluezObjectManager, err = idbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		panic(err)
	}
	b.objects, err = bluezObjectManager.GetManagedObjects()
	if err != nil {
		panic(err)
	}
	for k, v := range b.objects {
		logger.Debug(k, ":")
		for childKey, childValue := range v {
			logger.Debug("-->", childKey, ":", childValue)
		}
		logger.Debug()
		// adapter
	}

	// connect signals
	bluezObjectManager.ConnectInterfacesAdded(func(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
		// TODO
		logger.Debug(path, ":", data)
	})
	bluezObjectManager.ConnectInterfacesRemoved(func(path dbus.ObjectPath, interfaces []string) {
		// TODO
		logger.Debug(path, ":", interfaces)
	})
}

// TODO
func (b *Bluetooth) RemoveDevice(dpath dbus.ObjectPath) (err error) {
	return
}
