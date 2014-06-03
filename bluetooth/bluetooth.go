package bluetooth

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
	// adapter
	PrimaryAdapter string `access:"readwrite"` // TODO dbus.ObjectPath
	adapters       []*adapter
	Adapters       string // array of adapters that marshaled by json

	// device
	devices map[dbus.ObjectPath][]*device
	Devices string // device objects that marshaled by json

	// alias properties for primary adapter
	Alias               dbus.Property `access:"readwrite"`
	Powered             dbus.Property `access:"readwrite"`
	Discoverable        dbus.Property `access:"readwrite"`
	DiscoverableTimeout dbus.Property `access:"readwrite"`
	// Alias               string `access:"readwrite"`
	// Powered             bool   `access:"readwrite"`
	// Discoverable        bool   `access:"readwrite"`
	// DiscoverableTimeout uint32 `access:"readwrite"`

	// signals
	DeviceAdded      func(devJSON string)
	DeviceRemoved    func(devJSON string)
	RequestPinCode   func(devJSON string)
	AuthorizeService func(devJSON string, uuid string)
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
	b.devices = make(map[dbus.ObjectPath][]*device)

	// initialize dbus object manager
	var err error
	bluezObjectManager, err = idbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		panic(err)
	}
	objects, err := bluezObjectManager.GetManagedObjects()
	if err != nil {
		panic(err)
	}

	// add exists adapters and devices
	for path, data := range objects {
		b.handleInterfacesAdded(path, data)
	}

	// connect signals
	bluezObjectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
	bluezObjectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)
}
func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[dbusBluezIfsAdapter]; ok {
		b.addAdapter(path)
		if len(b.PrimaryAdapter) == 0 {
			b.updatePropPrimaryAdapter(path)
		}
	}
	if _, ok := data[dbusBluezIfsDevice]; ok {
		b.addDevice(path, data[dbusBluezIfsDevice])
	}
}
func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(dbusBluezIfsAdapter, interfaces) {
		b.removeAdapter(path)
		if dbus.ObjectPath(b.PrimaryAdapter) == path {
			if len(b.adapters) > 0 {
				b.updatePropPrimaryAdapter(b.adapters[0].Path)
			} else {
				b.updatePropPrimaryAdapter("")
			}
		}
	}
	if isStringInArray(dbusBluezIfsDevice, interfaces) {
		b.removeDevice(path)
	}
}
