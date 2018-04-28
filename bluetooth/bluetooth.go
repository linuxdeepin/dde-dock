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
	"errors"
	"sync"
	"time"

	"fmt"

	"github.com/linuxdeepin/go-dbus-factory/org.bluez"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	bluezDBusServiceName      = "org.bluez"
	bluezAdapterDBusInterface = "org.bluez.Adapter1"
	bluezDeviceDBusInterface  = "org.bluez.Device1"

	dbusServiceName = "com.deepin.daemon.Bluetooth"
	dbusPath        = "/com/deepin/daemon/Bluetooth"
	dbusInterface   = dbusServiceName
)

const (
	StateUnavailable = 0
	StateAvailable   = 1
	StateConnected   = 2
)

type dbusObjectData map[string]dbus.Variant

//go:generate dbusutil-gen -type Bluetooth bluetooth.go

type Bluetooth struct {
	service       *dbusutil.Service
	systemSigLoop *dbusutil.SignalLoop
	config        *config
	objectManager *bluez.ObjectManager
	sysDBusDaemon *ofdbus.DBus
	agent         *agent

	// adapter
	adaptersLock sync.Mutex
	adapters     map[dbus.ObjectPath]*adapter

	// device
	devicesLock sync.Mutex
	devices     map[dbus.ObjectPath][]*device

	PropsMu sync.RWMutex
	State   uint32 // StateUnavailable/StateAvailable/StateConnected

	methods *struct {
		DebugInfo                     func() `out:"info"`
		GetDevices                    func() `in:"adapter" out:"devicesJSON"`
		ConnectDevice                 func() `in:"device"`
		DisconnectDevice              func() `in:"device"`
		RemoveDevice                  func() `in:"adapter,device"`
		SetDeviceAlias                func() `in:"device,alias"`
		SetDeviceTrusted              func() `in:"device,trusted"`
		Confirm                       func() `in:"device,accept"`
		FeedPinCode                   func() `in:"device,accept,pinCode"`
		FeedPasskey                   func() `in:"device,accept,passkey"`
		GetAdapters                   func() `out:"adaptersJSON"`
		RequestDiscovery              func() `in:"adapter"`
		SetAdapterPowered             func() `in:"adapter,powered"`
		SetAdapterAlias               func() `in:"adapter,alias"`
		SetAdapterDiscoverable        func() `in:"adapter,discoverable"`
		SetAdapterDiscovering         func() `in:"adapter,discovering"`
		SetAdapterDiscoverableTimeout func() `in:"adapter,timeout"`
	}

	signals *struct {
		// adapter/device properties changed signals
		AdapterAdded, AdapterRemoved, AdapterPropertiesChanged struct {
			adapterJSON string
		}

		DeviceAdded, DeviceRemoved, DevicePropertiesChanged struct {
			devJSON string
		}

		// pair request signals
		DisplayPinCode struct {
			device  dbus.ObjectPath
			pinCode string
		}
		DisplayPasskey struct {
			device  dbus.ObjectPath
			passkey uint32
			entered uint32
		}

		// RequestConfirmation you should call Confirm with accept
		RequestConfirmation struct {
			device  dbus.ObjectPath
			passkey string
		}

		// RequestAuthorization you should call Confirm with accept
		RequestAuthorization struct {
			device dbus.ObjectPath
		}

		// RequestPinCode you should call FeedPinCode with accept and key
		RequestPinCode struct {
			device dbus.ObjectPath
		}

		// RequestPasskey you should call FeedPasskey with accept and key
		RequestPasskey struct {
			device dbus.ObjectPath
		}

		Cancelled struct {
			device dbus.ObjectPath
		}
	}
}

func newBluetooth(service *dbusutil.Service) (b *Bluetooth) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return nil
	}

	b = &Bluetooth{
		service:       service,
		systemSigLoop: dbusutil.NewSignalLoop(systemConn, 10),
	}
	b.adapters = make(map[dbus.ObjectPath]*adapter)
	return
}

func (b *Bluetooth) destroy() {
	b.agent.destroy()

	b.objectManager.RemoveHandler(proxy.RemoveAllHandlers)
	b.sysDBusDaemon.RemoveHandler(proxy.RemoveAllHandlers)

	b.devicesLock.Lock()
	for _, devices := range b.devices {
		for _, device := range devices {
			device.destroy()
		}
	}
	b.devicesLock.Unlock()

	b.adaptersLock.Lock()
	for _, adapter := range b.adapters {
		adapter.destroy()
	}
	b.adaptersLock.Unlock()

	err := b.service.StopExport(b)
	if err != nil {
		logger.Warning(err)
	}
	b.systemSigLoop.Stop()
}

func (*Bluetooth) GetInterfaceName() string {
	return dbusInterface
}

func (b *Bluetooth) init() {
	b.systemSigLoop.Start()
	b.config = newConfig()
	b.config.save()
	b.devices = make(map[dbus.ObjectPath][]*device)

	systemConn := b.systemSigLoop.Conn()

	b.sysDBusDaemon = ofdbus.NewDBus(systemConn)
	b.sysDBusDaemon.InitSignalExt(b.systemSigLoop, true)
	b.sysDBusDaemon.ConnectNameOwnerChanged(b.handleDBusNameOwnerChanged)

	// initialize dbus object manager
	b.objectManager = bluez.NewObjectManager(systemConn)

	// connect signals
	b.objectManager.InitSignalExt(b.systemSigLoop, true)
	b.objectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
	b.objectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)

	b.agent.init()
	b.loadObjects()
}

func (b *Bluetooth) loadObjects() {
	// add exists adapters and devices
	objects, err := b.objectManager.GetManagedObjects(0)
	if err != nil {
		logger.Error(err)
		return
	}
	for path, data := range objects {
		b.handleInterfacesAdded(path, data)
	}
}

func (b *Bluetooth) removeAllObjects() {
	b.devicesLock.Lock()
	for _, devices := range b.devices {
		for _, device := range devices {
			device.notifyDeviceRemoved()
			device.destroy()
		}
	}
	b.devices = make(map[dbus.ObjectPath][]*device)
	b.devicesLock.Unlock()

	b.adaptersLock.Lock()
	for _, adapter := range b.adapters {
		adapter.notifyAdapterRemoved()
		adapter.destroy()
	}
	b.adapters = make(map[dbus.ObjectPath]*adapter)
	b.adaptersLock.Unlock()
}

func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[bluezAdapterDBusInterface]; ok {
		requestUnblockBluetoothDevice()
		b.addAdapter(dbus.ObjectPath(path))
	}
	if _, ok := data[bluezDeviceDBusInterface]; ok {
		b.addDevice(dbus.ObjectPath(path), data[bluezDeviceDBusInterface])
	}
}

func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(bluezAdapterDBusInterface, interfaces) {
		b.removeAdapter(dbus.ObjectPath(path))
	}
	if isStringInArray(bluezDeviceDBusInterface, interfaces) {
		b.removeDevice(dbus.ObjectPath(path))
	}
}

func (b *Bluetooth) handleDBusNameOwnerChanged(name, oldOwner, newOwner string) {
	// if a new dbus session was installed, the name and newOwner
	// will be not empty, if a dbus session was uninstalled, the
	// name and oldOwner will be not empty
	if name != bluezDBusServiceName {
		return
	}
	if newOwner != "" {
		logger.Info("bluetooth is starting")
		time.AfterFunc(1*time.Second, func() {
			b.loadObjects()
			b.agent.registerDefaultAgent()
		})
	} else {
		logger.Info("bluetooth stopped")
		b.removeAllObjects()
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
		time.AfterFunc(25*time.Second, func() {
			d, _ := b.getDevice(dpath)
			if d != nil && !d.connected {
				logger.Infof("auto connect %s", d)
				d.Connect()
			}
		})
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

func (b *Bluetooth) findDevice(dpath dbus.ObjectPath) (apath dbus.ObjectPath, index int) {
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
	return b.findDevice(dpath)
}

func (b *Bluetooth) getDevice(dpath dbus.ObjectPath) (*device, error) {
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	apath, index := b.findDevice(dpath)
	if index < 0 {
		return nil, errInvalidDevicePath
	}
	return b.devices[apath][index], nil
}

func (b *Bluetooth) addAdapter(apath dbus.ObjectPath) {
	if b.isAdapterExists(apath) {
		logger.Warning("repeat add adapter", apath)
		return
	}

	a := newAdapter(b.systemSigLoop, apath)
	// initialize adapter power state
	b.config.addAdapterConfig(bluezGetAdapterAddress(apath))
	oldPowered := b.config.getAdapterConfigPowered(bluezGetAdapterAddress(apath))

	err := a.bluezAdapter.Powered().Set(0, oldPowered)
	if err != nil {
		logger.Warning(err)
	}

	err = a.bluezAdapter.Discoverable().Set(0, false)
	if err != nil {
		logger.Warning(err)
	}

	b.adaptersLock.Lock()
	b.adapters[apath] = a
	b.adaptersLock.Unlock()

	a.notifyAdapterAdded()
}

func (b *Bluetooth) removeAdapter(apath dbus.ObjectPath) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	if b.adapters[apath] == nil {
		logger.Warning("repeat remove adapter", apath)
		return
	}

	b.doRemoveAdapter(apath)
}

func (b *Bluetooth) doRemoveAdapter(apath dbus.ObjectPath) {
	removeAdapter := b.adapters[apath]
	delete(b.adapters, apath)

	removeAdapter.notifyAdapterRemoved()
	removeAdapter.destroy()
}

func (b *Bluetooth) getAdapter(apath dbus.ObjectPath) (a *adapter, err error) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	a = b.adapters[apath]
	if a == nil {
		err = fmt.Errorf("adapter not exists %s", apath)
		logger.Error(err)
		return
	}
	return
}
func (b *Bluetooth) isAdapterExists(apath dbus.ObjectPath) bool {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()
	if b.adapters[apath] != nil {
		return true
	}
	return false
}

func (b *Bluetooth) feed(devPath dbus.ObjectPath, accept bool, key string) (err error) {
	_, err = b.getDevice(devPath)
	if nil != err {
		logger.Warningf("FeedRequest can not find device: %v, %v", devPath, err)
		return err
	}

	b.agent.mu.Lock()
	if b.agent.requestDevice != devPath {
		b.agent.mu.Unlock()
		logger.Warningf("FeedRequest can not find match device: %q, %q", b.agent.requestDevice, devPath)
		return errBluezCanceled
	}
	b.agent.mu.Unlock()

	select {
	case b.agent.rspChan <- authorize{path: devPath, accept: accept, key: key}:
		return nil
	default:
		return errors.New("rspChan no reader")
	}
}

func (b *Bluetooth) updateState() {
	newState := StateUnavailable
	if len(b.adapters) > 0 {
		newState = StateAvailable
	}

	for _, devices := range b.devices {
		for _, d := range devices {
			if d.connected {
				newState = StateConnected
				break
			}
		}
	}

	b.PropsMu.Lock()
	b.setPropState(uint32(newState))
	b.PropsMu.Unlock()
}
