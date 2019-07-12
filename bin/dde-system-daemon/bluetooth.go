package main

import (
	"path/filepath"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/keyfile"
)

const (
	bluetoothPrefixDir = "/var/lib/bluetooth"

	kfSectionGeneral  = "General"
	kfKeyTechnologies = "SupportedTechnologies"
)

func (*Daemon) BluetoothGetDeviceTechnologies(adapter, device string) ([]string, *dbus.Error) {
	var filename = filepath.Join(bluetoothPrefixDir, adapter, device, "info")
	technologies, err := doBluetoothGetDeviceTechnologies(filename)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	return technologies, nil
}

func doBluetoothGetDeviceTechnologies(filename string) ([]string, error) {
	var kf = keyfile.NewKeyFile()
	err := kf.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}
	return kf.GetStringList(kfSectionGeneral, kfKeyTechnologies)
}
