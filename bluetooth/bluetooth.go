package main

import (
	"dlib/dbus"
)

type Bluetooth struct {
	// TODO
	EnableBluetooth bool `access:"readwrite"`
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
}
