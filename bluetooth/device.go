/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	"dbus/org/bluez"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

const deviceRssiNotInRange = -1000 // -1000db means device not in range

const (
	deviceStateDisconnected = 0
	deviceStateConnecting   = 1
	deviceStateConnected    = 2
)

var (
	errInvaildDevicePath = fmt.Errorf("Invaild Device Path")
)

type device struct {
	bluezDevice *bluez.Device1

	Path        dbus.ObjectPath
	AdapterPath dbus.ObjectPath

	Name    string
	Alias   string
	Trusted bool
	Paired  bool
	State   uint32
	UUIDs   []string
	// optional
	Icon    string
	RSSI    int16
	Address string

	oldConnected bool
	connected    bool
	connecting   bool
	blockNotify  bool

	lk           sync.Mutex
	confirmation chan bool
}

func newDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) (d *device) {
	d = &device{Path: dpath}
	d.bluezDevice, _ = bluezNewDevice(dpath)
	d.AdapterPath = d.bluezDevice.Adapter.Get()
	d.Name = d.bluezDevice.Name.Get()
	d.Alias = d.bluezDevice.Alias.Get()
	d.Address = d.bluezDevice.Address.Get()
	d.Trusted = d.bluezDevice.Trusted.Get()
	d.Paired = d.bluezDevice.Paired.Get()
	d.connected = d.bluezDevice.Connected.Get()
	d.UUIDs = d.bluezDevice.UUIDs.Get()
	d.oldConnected = d.connected
	d.blockNotify = false
	d.notifyStateChanged()

	// optional properties, read from dbus object data in order to
	// check if property exists
	d.Icon = getDBusObjectValueString(data, "Icon")
	if isDBusObjectKeyExists(data, "RSSI") {
		d.RSSI = getDBusObjectValueInt16(data, "RSSI")
	}
	d.fixRssi()

	d.connectProperties()
	return
}
func destroyDevice(d *device) {
	bluezDestroyDevice(d.bluezDevice)
}

func (d *device) notifyDeviceAdded() {
	logger.Debug("DeviceAdded", marshalJSON(d))
	dbus.Emit(bluetooth, "DeviceAdded", marshalJSON(d))
	bluetooth.setPropState()
}
func (d *device) notifyDeviceRemoved() {
	logger.Debug("DeviceRemoved", marshalJSON(d))
	dbus.Emit(bluetooth, "DeviceRemoved", marshalJSON(d))
	bluetooth.setPropState()
}
func (d *device) notifyDevicePropertiesChanged() {
	logger.Debug("DevicePropertiesChanged", marshalJSON(d))
	dbus.Emit(bluetooth, "DevicePropertiesChanged", marshalJSON(d))
	bluetooth.setPropState()
}

func (d *device) connectProperties() {
	d.bluezDevice.Connected.ConnectChanged(d.handleConnected)

	d.bluezDevice.Name.ConnectChanged(func() {
		d.Name = d.bluezDevice.Name.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.Alias.ConnectChanged(func() {
		d.Alias = d.bluezDevice.Alias.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.Address.ConnectChanged(func() {
		d.Address = d.bluezDevice.Address.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.Trusted.ConnectChanged(func() {
		d.Trusted = d.bluezDevice.Trusted.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.Paired.ConnectChanged(func() {
		d.Paired = d.bluezDevice.Paired.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.Icon.ConnectChanged(func() {
		d.Icon = d.bluezDevice.Icon.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.UUIDs.ConnectChanged(func() {
		d.UUIDs = d.bluezDevice.UUIDs.Get()
		d.notifyDevicePropertiesChanged()
	})
	d.bluezDevice.RSSI.ConnectChanged(func() {
		_, err := d.bluezDevice.RSSI.GetValue()
		if nil != err && !d.Paired {
			logger.Debug("Get dbus property RSSI failed", d.Path, err)
			bluezRemoveDevice(d.AdapterPath, d.Path)
			return
		}
		d.RSSI = d.bluezDevice.RSSI.Get()
		d.fixRssi()
		d.notifyDevicePropertiesChanged()
	})
}
func (d *device) handleConnected() {
	d.connected = d.bluezDevice.Connected.Get()
	if d.oldConnected != d.connected {
		d.oldConnected = d.connected
		if !d.blockNotify {
			if d.connected {
				notifyBluetoothConnected(d.Alias)
			} else {
				notifyBluetoothDisconnected(d.Alias)
			}
		}
	}
	d.notifyStateChanged()
}
func (d *device) notifyConnectFailed() {
	notifyBluetoothConnectFailed(d.Alias)
}
func (d *device) notifyIgnored() {
	notifyBluetoothDeviceIgnored(d.Alias)
}
func (d *device) notifyStateChanged() {
	if d.connected {
		d.connecting = false
		d.State = deviceStateConnected
	} else if d.connecting {
		d.State = deviceStateConnecting
	} else {
		d.State = deviceStateDisconnected
	}
	logger.Debugf("notifyStateChanged: %#v", d) // TODO test
	d.notifyDevicePropertiesChanged()
}

func (d *device) connectAddress() string {
	adapterAddress := bluezGetAdapterAddress(d.AdapterPath)
	return adapterAddress + "/" + d.Address
}

func (d *device) fixRssi() {
	if d.RSSI == 0 {
		if d.Trusted {
			d.RSSI = deviceRssiNotInRange / 2
		} else {
			d.RSSI = deviceRssiNotInRange
		}
	}
}

func (b *Bluetooth) addDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) {
	if b.isDeviceExists(dpath) {
		logger.Error("repeat add device", dpath)
		return
	}

	b.devicesLock.Lock()
	d := newDevice(dpath, data)
	b.devices[d.AdapterPath] = append(b.devices[d.AdapterPath], d)
	b.config.addDeviceConfig(d.connectAddress())

	d.notifyDeviceAdded()
	b.devicesLock.Unlock()

	connected := b.config.getDeviceConfigConnected(d.connectAddress())
	if connected {
		d.blockNotify = true
		go func() {
			time.Sleep(25 * time.Second)
			if d, _ := b.getDevice(dpath); nil != d && !d.connected {
				logger.Infof("auto connect device %v", d)
				bluezConnectDevice(dpath)
			}
			d.blockNotify = false
		}()
	}
}

func (b *Bluetooth) removeDevice(dpath dbus.ObjectPath) {
	apath, i := b.getDeviceIndex(dpath)
	if i < 0 {
		logger.Error("repeat remove device", dpath)
		return
	}

	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	b.devices[apath] = b.doRemoveDevice(b.devices[apath], i)
	return
}

func (b *Bluetooth) doRemoveDevice(devices []*device, i int) []*device {
	d := devices[i]
	b.config.removeDeviceConfig(d.connectAddress())
	d.notifyDeviceRemoved()
	destroyDevice(d)
	copy(devices[i:], devices[i+1:])
	devices[len(devices)-1] = nil
	devices = devices[:len(devices)-1]
	return devices
}

func (b *Bluetooth) isDeviceExists(dpath dbus.ObjectPath) bool {
	_, i := b.getDeviceIndex(dpath)
	if i >= 0 {
		return true
	}
	return false
}

func (b *Bluetooth) findDeviceIndex(dpath dbus.ObjectPath) (apath dbus.ObjectPath, index int) {
	for p, devices := range b.devices {
		for i, d := range devices {
			if d.Path == dpath {
				return p, i
			}
		}
	}
	return "", -1
}

func (b *Bluetooth) getDeviceIndex(dpath dbus.ObjectPath) (apath dbus.ObjectPath, index int) {
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	return b.findDeviceIndex(dpath)
}

func (b *Bluetooth) getDevice(dpath dbus.ObjectPath) (*device, error) {
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	apath, index := b.findDeviceIndex(dpath)
	if index < 0 {
		return nil, errInvaildDevicePath
	}
	return b.devices[apath][index], nil
}

// GetDevices return all device objects that marshaled by json.
func (b *Bluetooth) GetDevices(apath dbus.ObjectPath) (devicesJSON string, err error) {
	devices := b.devices[apath]
	devicesJSON = marshalJSON(devices)
	return
}

func (b *Bluetooth) ConnectDevice(dpath dbus.ObjectPath) (err error) {
	// mark device is connecting
	d, err := b.getDevice(dpath)
	if err != nil {
		return
	}
	d.connecting = true
	d.notifyStateChanged()

	go func() {
		d.blockNotify = true
		defer func() { d.blockNotify = false }()
		if !bluezGetDevicePaired(dpath) {
			bluezPairDevice(dpath)
			time.Sleep(200 * time.Millisecond)
		}

		// TODO: remove work code if bluez a2dp is ok
		// bluez do not support muti a2dp devices
		// disconnect a2dp device before connect
		for _, uuid := range d.UUIDs {
			if uuid == A2DP_SINK_UUID {
				b.disconnectA2DPDeviceExcept(d)
			}
		}

		err = bluezConnectDevice(dpath)
		if err == nil {
			b.config.setDeviceConfigConnected(d.connectAddress(), true)
			// trust device when connecting success
			if !bluezGetDeviceTrusted(dpath) {
				bluezSetDeviceTrusted(dpath, true)
			}
		} else {
			b.config.setDeviceConfigConnected(d.connectAddress(), false)
			logger.Warning("ConnectDevice failed:", dpath, err)
			if err.Error() == bluezErrorInvalidKey.Error() || !d.Paired {
				// we do not want to pop notify for device because we will remove it.
				bluezRemoveDevice(d.AdapterPath, d.Path)
				d.notifyIgnored()
				return
			}
			d.notifyConnectFailed()
		}
		d.handleConnected()

		if d.connecting {
			d.connecting = false
			d.notifyStateChanged()
		}
	}()
	return
}

func (b *Bluetooth) DisconnectDevice(dpath dbus.ObjectPath) (err error) {
	// mark disconnected in time
	d, err := b.getDevice(dpath)
	if err != nil {
		return
	}
	d.connected = false
	d.notifyStateChanged()
	b.config.setDeviceConfigConnected(d.connectAddress(), false)
	go bluezDisconnectDevice(dpath)
	return
}

func (b *Bluetooth) RemoveDevice(apath, dpath dbus.ObjectPath) (err error) {
	// TODO
	go bluezRemoveDevice(apath, dpath)
	return
}

func (b *Bluetooth) SetDeviceAlias(dpath dbus.ObjectPath, alias string) (err error) {
	go bluezSetDeviceAlias(dpath, alias)
	return
}

func (b *Bluetooth) SetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	go bluezSetDeviceTrusted(dpath, trusted)
	return
}
