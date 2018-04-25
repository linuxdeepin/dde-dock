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
	"sync"
	"time"

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
