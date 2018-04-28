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

	"strings"

	"github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	deviceStateDisconnected  = 0
	deviceStateConnecting    = 1
	deviceStateConnected     = 2
	deviceStateDisconnecting = 1
)

type deviceState uint32

func (s deviceState) String() string {
	switch s {
	case deviceStateDisconnected:
		return "Disconnected"
	case deviceStateConnecting:
		return "doing"
	case deviceStateConnected:
		return "Connected"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

var (
	errInvalidDevicePath = fmt.Errorf("invalid device path")
)

type device struct {
	bluezDevice *bluez.Device

	Path        dbus.ObjectPath
	AdapterPath dbus.ObjectPath

	Alias            string
	Trusted          bool
	Paired           bool
	State            deviceState
	ServicesResolved bool
	// optional
	UUIDs   []string
	Name    string
	Icon    string
	RSSI    int16
	Address string

	connected    bool
	connecting   bool
	agentWorking bool

	connectPhase      connectPhase
	disconnectPhase   disconnectPhase
	disconnectChan    chan struct{}
	mu                sync.Mutex
	confirmation      chan bool
	pairingFailedTime time.Time
}

type connectPhase uint32

const (
	connectPhaseNone = iota
	connectPhaseStart
	connectPhasePairStart
	connectPhasePairEnd
	connectPhaseConnectProfilesStart
	connectPhaseConnectProfilesEnd
)

type disconnectPhase uint32

const (
	disconnectPhaseNone = iota
	disconnectPhaseStart
	disconnectPhaseDisconnectStart
	disconnectPhaseDisconnectEnd
)

func (d *device) setDisconnectPhase(value disconnectPhase) {
	d.mu.Lock()
	d.disconnectPhase = value
	d.mu.Unlock()

	switch value {
	case disconnectPhaseDisconnectStart:
		logger.Debugf("%s disconnect start", d)
	case disconnectPhaseDisconnectEnd:
		logger.Debugf("%s disconnect end", d)
	}
	d.updateState()
	d.notifyDevicePropertiesChanged()
}

func (d *device) getDisconnectPhase() disconnectPhase {
	d.mu.Lock()
	value := d.disconnectPhase
	d.mu.Unlock()
	return value
}

func (d *device) setConnectPhase(value connectPhase) {
	d.mu.Lock()
	d.connectPhase = value
	d.mu.Unlock()

	switch value {
	case connectPhasePairStart:
		logger.Debugf("%s pair start", d)
	case connectPhasePairEnd:
		logger.Debugf("%s pair end", d)

	case connectPhaseConnectProfilesStart:
		logger.Debugf("%s connect profiles start", d)
	case connectPhaseConnectProfilesEnd:
		logger.Debugf("%s connect profiles end", d)
	}

	d.updateState()
	d.notifyDevicePropertiesChanged()
}

func (d *device) getConnectPhase() connectPhase {
	d.mu.Lock()
	value := d.connectPhase
	d.mu.Unlock()
	return value
}

func (d *device) agentWorkStart() {
	logger.Debugf("%s agent work start", d)
	d.agentWorking = true
	d.updateState()
	d.notifyDevicePropertiesChanged()
}

func (d *device) agentWorkEnd() {
	logger.Debugf("%s agent work end", d)
	d.agentWorking = false
	d.updateState()
	d.notifyDevicePropertiesChanged()
}

func (d *device) String() string {
	return fmt.Sprintf("device [%s] %s", d.Address, d.Alias)
}

func newDevice(systemSigLoop *dbusutil.SignalLoop, dpath dbus.ObjectPath,
	data map[string]dbus.Variant) (d *device) {
	d = &device{Path: dpath}
	systemConn := systemSigLoop.Conn()
	d.bluezDevice, _ = bluez.NewDevice(systemConn, dpath)
	d.AdapterPath, _ = d.bluezDevice.Adapter().Get(0)
	d.Name, _ = d.bluezDevice.Name().Get(0)
	d.Alias, _ = d.bluezDevice.Alias().Get(0)
	d.Address, _ = d.bluezDevice.Address().Get(0)
	d.Trusted, _ = d.bluezDevice.Trusted().Get(0)
	d.Paired, _ = d.bluezDevice.Paired().Get(0)
	d.connected, _ = d.bluezDevice.Connected().Get(0)
	d.UUIDs, _ = d.bluezDevice.UUIDs().Get(0)
	d.ServicesResolved, _ = d.bluezDevice.ServicesResolved().Get(0)

	d.disconnectChan = make(chan struct{})
	d.updateState()

	// optional properties, read from dbus object data in order to
	// check if property exists
	d.Icon = getDBusObjectValueString(data, "Icon")
	if isDBusObjectKeyExists(data, "RSSI") {
		d.RSSI = getDBusObjectValueInt16(data, "RSSI")
	}

	d.bluezDevice.InitSignalExt(systemSigLoop, true)
	d.connectProperties()
	return
}

func (d *device) destroy() {
	d.bluezDevice.RemoveHandler(proxy.RemoveAllHandlers)
}

func (d *device) notifyDeviceAdded() {
	logger.Debug("DeviceAdded", d)
	if d.Name == "" {
		logger.Debugf("%s is hidden", d)
		return
	}

	globalBluetooth.service.Emit(globalBluetooth, "DeviceAdded", marshalJSON(d))
	globalBluetooth.updateState()
}

func (d *device) notifyDeviceRemoved() {
	logger.Debug("DeviceRemoved", d)
	globalBluetooth.service.Emit(globalBluetooth, "DeviceRemoved", marshalJSON(d))
	globalBluetooth.updateState()
}

func (d *device) notifyDevicePropertiesChanged() {
	globalBluetooth.service.Emit(globalBluetooth, "DevicePropertiesChanged", marshalJSON(d))
	globalBluetooth.updateState()
}

func (d *device) connectProperties() {
	d.bluezDevice.Connected().ConnectChanged(func(hasValue bool, connected bool) {
		if !hasValue {
			return
		}
		logger.Debugf("%s Connected: %v", d, connected)
		d.connected = connected

		needNotify := true
		if !connected {
			select {
			case d.disconnectChan <- struct{}{}:
				logger.Debugf("%s disconnectChan send done", d)
				needNotify = false
			default:
			}
		}

		d.updateState()
		d.notifyDevicePropertiesChanged()

		if needNotify {
			d.notifyConnectedChanged()
		}
		return
	})

	d.bluezDevice.Name().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		logger.Debugf("%s Name: %v", d, value)
		d.Name = value
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Alias().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Alias = value
		logger.Debugf("%s Alias: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Address().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Address = value
		logger.Debugf("%s Address: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Trusted().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		d.Trusted = value
		logger.Debugf("%s Trusted: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Paired().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		d.Paired = value
		logger.Debugf("%s Paired: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.ServicesResolved().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		d.ServicesResolved = value
		logger.Debugf("%s ServicesResolved: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.Icon().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		d.Icon = value
		logger.Debugf("%s Icon: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.UUIDs().ConnectChanged(func(hasValue bool, value []string) {
		if !hasValue {
			return
		}
		d.UUIDs = value
		logger.Debugf("%s UUIDs: %v", d, value)
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.RSSI().ConnectChanged(func(hasValue bool, value int16) {
		if !hasValue {
			d.RSSI = 0
			logger.Debugf("%s RSSI invalidated", d)
		} else {
			d.RSSI = value
			logger.Debugf("%s RSSI: %v", d, value)
		}
		d.notifyDevicePropertiesChanged()
	})

	d.bluezDevice.LegacyPairing().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		logger.Debugf("%s LegacyPairing: %v", d, value)
	})

	d.bluezDevice.Blocked().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		logger.Debugf("%s Blocked: %v", d, value)
	})
}

func (d *device) notifyConnectedChanged() {
	connectPhase := d.getConnectPhase()
	if connectPhase != connectPhaseNone {
		// connect is in progress
		logger.Debugf("%s handleNotifySend: connect is in progress", d)
		return
	}

	disconnectPhase := d.getDisconnectPhase()
	if disconnectPhase != disconnectPhaseNone {
		// disconnect is in progress
		logger.Debugf("%s handleNotifySend: disconnect is in progress", d)
		return
	}

	if d.connected {
		notifyConnected(d.Alias)
	} else {
		if time.Since(d.pairingFailedTime) < 2*time.Second {
			return
		}
		notifyDisconnected(d.Alias)
	}
}

func (d *device) updateState() {
	newState := d.getState()
	if d.State != newState {
		d.State = newState
		logger.Debugf("%s State: %s", d, d.State)
	}
}

func (d *device) getState() deviceState {
	if d.agentWorking {
		return deviceStateConnecting
	}

	if d.connectPhase != connectPhaseNone {
		return deviceStateConnecting

	} else if d.disconnectPhase != connectPhaseNone {
		return deviceStateDisconnecting

	} else {
		if d.connected {
			return deviceStateConnected
		} else {
			return deviceStateDisconnected
		}
	}
}

func (d *device) connectAddress() string {
	adapterAddress := bluezGetAdapterAddress(d.AdapterPath)
	return adapterAddress + "/" + d.Address
}

func (d *device) Connect() {
	logger.Debug(d, "call Connect()")
	connectPhase := d.getConnectPhase()
	disconnectPhase := d.getDisconnectPhase()
	if connectPhase != connectPhaseNone {
		logger.Warningf("%s connect is in progress", d)
		return
	} else if disconnectPhase != disconnectPhaseNone {
		logger.Debugf("%s disconnect is in progress", d)
		return
	}

	d.setConnectPhase(connectPhaseStart)
	defer d.setConnectPhase(connectPhaseNone)

	blocked, err := d.bluezDevice.Blocked().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if blocked {
		err := d.bluezDevice.Blocked().Set(0, false)
		if err != nil {
			logger.Warning(err)
			return
		}
	}

	paired, err := d.bluezDevice.Paired().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if !paired {
		d.setConnectPhase(connectPhasePairStart)
		err := d.bluezDevice.Pair(0)
		d.setConnectPhase(connectPhasePairEnd)

		if err != nil {
			logger.Warningf("%s pair failed: %v", d, err)
			notifyConnectFailedPairing(d.Alias)
			d.pairingFailedTime = time.Now()
			d.setConnectPhase(connectPhaseNone)
			return
		} else {
			logger.Warningf("%s pair succeeded", d)
		}
	} else {
		logger.Debugf("%s already paired", d)
	}

	// TODO: remove work code if bluez a2dp is ok
	// bluez do not support muti a2dp devices
	// disconnect a2dp device before connect

	for _, uuid := range d.UUIDs {
		if uuid == A2DP_SINK_UUID {
			globalBluetooth.disconnectA2DPDeviceExcept(d)
		}
	}

	d.setConnectPhase(connectPhaseConnectProfilesStart)
	err = d.bluezDevice.Connect(0)
	d.setConnectPhase(connectPhaseConnectProfilesEnd)
	if err == nil {
		// connect succeeded
		logger.Infof("%s connect succeeded", d)
		globalBluetooth.config.setDeviceConfigConnected(d.connectAddress(), true)

		// auto trust device when connecting success
		trusted, _ := d.bluezDevice.Trusted().Get(0)
		if !trusted {
			err := d.bluezDevice.Trusted().Set(0, true)
			if err != nil {
				logger.Warning(err)
			}
		}
		notifyConnected(d.Alias)

	} else {
		// connect failed
		logger.Warningf("%s connect failed: %v", d, err)

		globalBluetooth.config.setDeviceConfigConnected(d.connectAddress(), false)

		errMsg := err.Error()
		if strings.Contains(errMsg, "Host is down") ||
			strings.Contains(errMsg, "Input/output error") {
			notifyConnectFailedHostDown(d.Alias)
		} else if strings.Contains(errMsg, "Resource temporarily unavailable") {
			notifyConnectFailedResourceUnavailable(d.Alias)
		} else if strings.Contains(errMsg, "Software caused connection abort") {
			notifyConnectFailedSoftwareCaused(d.Alias)
		} else {
			notifyConnectFailedOther(d.Alias)
		}
	}
}

func (d *device) Disconnect() {
	logger.Debugf("%s call Disconnect()", d)

	disconnectPhase := d.getDisconnectPhase()
	if disconnectPhase != disconnectPhaseNone {
		logger.Debugf("%s disconnect is in progress", d)
		return
	}

	d.setDisconnectPhase(disconnectPhaseStart)
	defer d.setDisconnectPhase(disconnectPhaseNone)

	connected, err := d.bluezDevice.Connected().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if !connected {
		logger.Debugf("%s not connected", d)
		return
	}

	globalBluetooth.config.setDeviceConfigConnected(d.connectAddress(), false)

	ch := d.goWaitDisconnect()

	d.setDisconnectPhase(disconnectPhaseDisconnectStart)
	d.bluezDevice.Disconnect(0)
	d.setDisconnectPhase(disconnectPhaseDisconnectEnd)

	if d.Icon == "phone" || d.Icon == "computer" {
		// do not block phone or computer
	} else {
		d.bluezDevice.Blocked().Set(0, true)
	}

	<-ch
	notifyDisconnected(d.Alias)
}

func (d *device) goWaitDisconnect() chan struct{} {
	ch := make(chan struct{})
	go func() {
		select {
		case <-d.disconnectChan:
			logger.Debugf("%s disconnectChan receive ok", d)
		case <-time.After(60 * time.Second):
			logger.Debugf("%s disconnectChan receive timed out", d)
		}
		ch <- struct{}{}
	}()
	return ch
}
