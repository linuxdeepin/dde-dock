/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bluetooth

import (
	"fmt"
	"sync"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const deviceRssiNotInRange = -1000 // -1000db means device not in range

const (
	deviceStateDisconnected = 0
	deviceStateConnecting   = 1
	deviceStateConnected    = 2
)

var (
	errInvalidDevicePath = fmt.Errorf("invalid device path")
)

type device struct {
	bluezDevice *bluez.Device

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

func newDevice(systemSigLoop *dbusutil.SignalLoop, dpath dbus.ObjectPath,
	data map[string]dbus.Variant) (d *device) {
	d = &device{Path: dpath}
	d.bluezDevice, _ = bluezNewDevice(dpath)
	d.AdapterPath, _ = d.bluezDevice.Adapter().Get(0)
	d.Name, _ = d.bluezDevice.Name().Get(0)
	d.Alias, _ = d.bluezDevice.Alias().Get(0)
	d.Address, _ = d.bluezDevice.Address().Get(0)
	d.Trusted, _ = d.bluezDevice.Trusted().Get(0)
	d.Paired, _ = d.bluezDevice.Paired().Get(0)
	d.connected, _ = d.bluezDevice.Connected().Get(0)
	d.UUIDs, _ = d.bluezDevice.UUIDs().Get(0)
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

	d.bluezDevice.InitSignalExt(systemSigLoop, true)
	d.connectProperties()
	return
}

func (d *device) destroy() {
	d.bluezDevice.RemoveHandler(proxy.RemoveAllHandlers)
}

func (d *device) notifyDeviceAdded() {
	logger.Debug("DeviceAdded", marshalJSON(d))
	globalBluetooth.service.Emit(globalBluetooth, "DeviceAdded", marshalJSON(d))
	globalBluetooth.updateState()
}
func (d *device) notifyDeviceRemoved() {
	logger.Debug("DeviceRemoved", marshalJSON(d))
	globalBluetooth.service.Emit(globalBluetooth, "DeviceRemoved", marshalJSON(d))
	globalBluetooth.updateState()
}
func (d *device) notifyDevicePropertiesChanged() {
	logger.Debug("DevicePropertiesChanged", marshalJSON(d))
	globalBluetooth.service.Emit(globalBluetooth, "DevicePropertiesChanged", marshalJSON(d))
	globalBluetooth.updateState()
}

func (d *device) connectProperties() {
	d.bluezDevice.Connected().ConnectChanged(func(hasValue bool, value bool) {
		d.handleConnected()
	})

	d.bluezDevice.Name().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Name = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Alias().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Alias = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Address().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Address = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Trusted().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		d.Trusted = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Paired().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		d.Paired = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Icon().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Icon = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.UUIDs().ConnectChanged(func(hasValue bool, value []string) {
		if !hasValue {
			return
		}
		d.UUIDs = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.RSSI().ConnectChanged(func(hasValue bool, value int16) {
		rssi, err := d.bluezDevice.RSSI().Get(0)
		if err != nil && !d.Paired {
			logger.Debug("Get dbus property RSSI failed", d.Path, err)
			bluezRemoveDevice(d.AdapterPath, d.Path)
			return
		}
		d.RSSI = rssi
		d.fixRssi()
		d.notifyDevicePropertiesChanged()
	})
}
func (d *device) handleConnected() {
	d.connected, _ = d.bluezDevice.Connected().Get(0)
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
	d := newDevice(b.systemSigLoop, dpath, data)
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
	d.destroy()
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
		return nil, errInvalidDevicePath
	}
	return b.devices[apath][index], nil
}

// GetDevices return all device objects that marshaled by json.
func (b *Bluetooth) GetDevices(apath dbus.ObjectPath) (devicesJSON string, err *dbus.Error) {
	devices := b.devices[apath]
	devicesJSON = marshalJSON(devices)
	return
}

func (b *Bluetooth) ConnectDevice(dpath dbus.ObjectPath) *dbus.Error {
	// mark device is connecting
	d, err := b.getDevice(dpath)
	if err != nil {
		return dbusutil.ToError(err)
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
	return nil
}

func (b *Bluetooth) DisconnectDevice(dpath dbus.ObjectPath) *dbus.Error {
	// mark disconnected in time
	d, err := b.getDevice(dpath)
	if err != nil {
		return dbusutil.ToError(err)
	}
	d.connected = false
	d.notifyStateChanged()
	b.config.setDeviceConfigConnected(d.connectAddress(), false)
	go bluezDisconnectDevice(dpath)
	return nil
}

func (b *Bluetooth) RemoveDevice(apath, dpath dbus.ObjectPath) *dbus.Error {
	// TODO
	go bluezRemoveDevice(apath, dpath)
	return nil
}

func (b *Bluetooth) SetDeviceAlias(dpath dbus.ObjectPath, alias string) *dbus.Error {
	go bluezSetDeviceAlias(dpath, alias)
	return nil
}

func (b *Bluetooth) SetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) *dbus.Error {
	go bluezSetDeviceTrusted(dpath, trusted)
	return nil
}
