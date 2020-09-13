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
	"os"
	"time"

	dbus "github.com/godbus/dbus"
	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

type adapter struct {
	core    *bluez.HCI
	address string

	Path                dbus.ObjectPath
	Name                string
	Alias               string
	Powered             bool
	Discovering         bool
	Discoverable        bool
	DiscoverableTimeout uint32
	// discovering timer, when time is up, stop discovering until start button is clicked next time
	discoveringTimeout *time.Timer
	//Scan timeout flag
	discoveringTimeoutFlag              bool
	scanReadyToConnectDeviceTimeout     *time.Timer
	scanReadyToConnectDeviceTimeoutFlag bool
}

var defaultDiscoveringTimeout = 1 * time.Minute
var defaultFindDeviceTimeout = 1 * time.Second

func newAdapter(systemSigLoop *dbusutil.SignalLoop, apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	systemConn := systemSigLoop.Conn()
	a.core, _ = bluez.NewHCI(systemConn, apath)
	a.core.InitSignalExt(systemSigLoop, true)
	a.connectProperties()
	a.address, _ = a.core.Address().Get(0)
	// 用于定时停止扫描
	a.discoveringTimeout = time.AfterFunc(defaultDiscoveringTimeout, func() {
		logger.Debug("discovery time out, stop discovering")
		//扫描结束后添加备份
		for adapterpath, devices := range globalBluetooth.devices {
			for _, device := range devices {
				bd := newBackupDevice(device)
				globalBluetooth.backupDeviceLock.Lock()
				globalBluetooth.backupDevices[adapterpath] = append(globalBluetooth.backupDevices[adapterpath], bd)
				globalBluetooth.backupDeviceLock.Unlock()
			}
		}
		//Scan timeout
		a.discoveringTimeoutFlag = true
		if err := a.core.StopDiscovery(0); err != nil {
			logger.Warningf("stop discovery failed, err:%v", err)
		}
		globalBluetooth.prepareToConnectedDevice = ""
	})
	//扫描1S钟，未扫描到该设备弹出通知
	a.scanReadyToConnectDeviceTimeout = time.AfterFunc(defaultFindDeviceTimeout, func() {
		a.scanReadyToConnectDeviceTimeoutFlag = false
		_, err := globalBluetooth.getDevice(globalBluetooth.prepareToConnectedDevice)
		if err != nil {
			backupdevice, err1 := globalBluetooth.getBackupDevice(globalBluetooth.prepareToConnectedDevice)
			if err1 != nil {
				logger.Debug("getBackupDevice Failed:", err1)
			}
			notifyConnectFailedHostDown(backupdevice.Alias)
		}
		//清空备份
		globalBluetooth.backupDeviceLock.Lock()
		globalBluetooth.backupDevices = make(map[dbus.ObjectPath][]*backupDevice)
		globalBluetooth.backupDeviceLock.Unlock()
	})
	// stop timer at first
	a.discoveringTimeout.Stop()
	a.scanReadyToConnectDeviceTimeout.Stop()
	// fix alias
	alias, _ := a.core.Alias().Get(0)
	if alias == "first-boot-hostname" {
		hostname, err := os.Hostname()
		if err == nil {
			if hostname != "first-boot-hostname" {
				// reset alias
				err = a.core.Alias().Set(0, "")
				if err != nil {
					logger.Warning(err)
				}
			}
		} else {
			logger.Warning("failed to get hostname:", err)
		}
	}

	a.Alias, _ = a.core.Alias().Get(0)
	a.Name, _ = a.core.Name().Get(0)
	a.Powered, _ = a.core.Powered().Get(0)
	a.Discovering, _ = a.core.Discovering().Get(0)
	a.Discoverable, _ = a.core.Discoverable().Get(0)
	a.DiscoverableTimeout, _ = a.core.DiscoverableTimeout().Get(0)
	return
}

func (a *adapter) destroy() {
	a.core.RemoveHandler(proxy.RemoveAllHandlers)
}

func (a *adapter) String() string {
	return fmt.Sprintf("adapter %s [%s]", a.Alias, a.address)
}

func (a *adapter) notifyAdapterAdded() {
	logger.Info("AdapterAdded", a)
	err := globalBluetooth.service.Emit(globalBluetooth, "AdapterAdded", marshalJSON(a))
	if err != nil {
		logger.Warning(err)
	}
	globalBluetooth.updateState()
}

func (a *adapter) notifyAdapterRemoved() {
	logger.Info("AdapterRemoved", a)
	err := globalBluetooth.service.Emit(globalBluetooth, "AdapterRemoved", marshalJSON(a))
	if err != nil {
		logger.Warning(err)
	}
	globalBluetooth.updateState()
}

func (a *adapter) notifyPropertiesChanged() {
	err := globalBluetooth.service.Emit(globalBluetooth, "AdapterPropertiesChanged", marshalJSON(a))
	if err != nil {
		logger.Warning(err)
	}
	globalBluetooth.updateState()
}

func (a *adapter) connectProperties() {
	err := a.core.Name().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Name = value
		logger.Debugf("%s Name: %v", a, value)
		a.notifyPropertiesChanged()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = a.core.Alias().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Alias = value
		logger.Debugf("%s Alias: %v", a, value)
		a.notifyPropertiesChanged()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = a.core.Powered().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Powered = value
		logger.Debugf("%s Powered: %v", a, value)

		if a.Powered {
			err := a.core.Discoverable().Set(0, globalBluetooth.config.Discoverable)
			if err != nil {
				logger.Warningf("failed to set discoverable for %s: %v", a, err)
			}
			go func() {
				time.Sleep(1 * time.Second)
				err = a.core.StopDiscovery(0)
				// in case auto connect to device failed, only when signal power on is received, try to auto connect device
				globalBluetooth.tryConnectPairedDevices()
				// start discovery
				err = a.core.StartDiscovery(0)
				if err != nil {
					logger.Warningf("failed to start discovery for %s: %v", a, err)
				}
				// dont need to start discovering, according to blueZ, scan will be called, when power is set on
				a.discoveringTimeout.Reset(defaultDiscoveringTimeout)
			}()
		} else {
			// if power off, stop discovering time out
			a.discoveringTimeout.Stop()
		}
		// Sleep for 1s and wait for bluez to set the attributes before sending the attribute change signal
		time.Sleep(1 * time.Second)
		a.notifyPropertiesChanged()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = a.core.Discovering().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discovering = value
		logger.Debugf("%s Discovering: %v", a, value)
		//Scan timeout and send attribute change signal directly
		if a.discoveringTimeoutFlag {
			a.notifyPropertiesChanged()
		} else {
			if value != a.Powered {
				return
			}
			a.notifyPropertiesChanged()
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	err = a.core.Discoverable().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discoverable = value
		logger.Debugf("%s Discoverable: %v", a, value)
		a.notifyPropertiesChanged()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = a.core.DiscoverableTimeout().ConnectChanged(func(hasValue bool, value uint32) {
		if !hasValue {
			return
		}
		a.DiscoverableTimeout = value
		logger.Debugf("%s DiscoverableTimeout: %v", a, value)
		a.notifyPropertiesChanged()
	})
	if err != nil {
		logger.Warning(err)
	}
}
func (a *adapter) startDiscovery() {
	err := a.core.StartDiscovery(0)
	if err != nil {
		logger.Warningf("failed to start discovery for %s: %v", a, err)
	} else {
		logger.Debug("reset timer for stop scan")
		// start discovering success, reset discovering timer
		a.discoveringTimeout.Reset(defaultDiscoveringTimeout)
	}
}
