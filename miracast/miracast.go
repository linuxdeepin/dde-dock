/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package miracast

import (
	"dbus/org/freedesktop/miracle/wfd"
	"dbus/org/freedesktop/miracle/wifi"
	"dbus/org/freedesktop/networkmanager"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/dde/daemon/iw"
	oldDBusLib "pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.Miracast"
	dbusPath        = "/com/deepin/daemon/Miracast"
	dbusInterface   = dbusServiceName
)

const (
	wifiDBusServiceName = "org.freedesktop.miracle.wifi"
	wifiDBusPath        = "/org/freedesktop/miracle/wifi"
	linkDBusPath        = "/org/freedesktop/miracle/wifi/link/"
	linkDBusInterface   = "org.freedesktop.miracle.wifi.Link"
	peerDBusPath        = "/org/freedesktop/miracle/wifi/peer/"
	peerDBusInterface   = "org.freedesktop.miracle.wifi.Peer"

	wfdDBusServiceName = "org.freedesktop.miracle.wfd"
	wfdDBusPath        = "/org/freedesktop/miracle/wfd"
	sinkDBusPath       = "/org/freedesktop/miracle/wfd/sink/"
	sinkDBusInterface  = "org.freedesktop.miracle.wfd.Sink"

	nmDBusServiceName = "org.freedesktop.NetworkManager"
	nmDBusPath        = "/org/freedesktop/NetworkManager"
)

const (
	defaultTimeout = time.Second * 10
)

type Miracast struct {
	wifiObj *wifi.ObjectManager
	wfdObj  *wfd.ObjectManager
	network *networkmanager.Manager
	links   LinkInfos
	sinks   SinkInfos
	devices iw.WirelessInfos

	linkLocker   sync.Mutex
	sinkLocker   sync.Mutex
	deviceLocker sync.Mutex

	managingLinks   map[dbus.ObjectPath]bool
	connectingSinks map[dbus.ObjectPath]bool

	service *dbusutil.Service
	signals *struct {
		Added, Removed struct {
			path       dbus.ObjectPath
			detailJSON string
		}
		Event struct {
			eventType uint8
			path      dbus.ObjectPath
		}
	}

	methods *struct {
		ListLinks   func() `out:"links"`
		ListSinks   func() `out:"sinks"`
		Enable      func() `in:"link,enabled"`
		SetLinkName func() `in:"link,name"`
		Scanning    func() `in:"link,enabled"`
		Connect     func() `in:"sink,x,y,w,h"`
		Disconnect  func() `in:"sink"`
	}
}

func newMiracast(service *dbusutil.Service) (*Miracast, error) {
	network, err := networkmanager.NewManager(nmDBusServiceName, nmDBusPath)
	if err != nil {
		return nil, err
	}

	wifiObj, err := wifi.NewObjectManager(wifiDBusServiceName, wifiDBusPath)
	if err != nil {
		networkmanager.DestroyManager(network)
		return nil, err
	}

	wfdObj, err := wfd.NewObjectManager(wfdDBusServiceName, wfdDBusPath)
	if err != nil {
		networkmanager.DestroyManager(network)
		wifi.DestroyObjectManager(wifiObj)
		return nil, err
	}

	return &Miracast{
		service:         service,
		network:         network,
		wifiObj:         wifiObj,
		wfdObj:          wfdObj,
		connectingSinks: make(map[dbus.ObjectPath]bool),
		managingLinks:   make(map[dbus.ObjectPath]bool),
	}, nil
}

func (m *Miracast) init() {
	devices, err := iw.ListWirelessInfo()
	if err != nil {
		logger.Error("Failed to list wireless info:", err)
	}
	logger.Debugf("All devices: %#v", devices)

	m.devices = devices.ListMiracastDevice()
	if len(m.devices) == 0 {
		m.handleEvent()
		return
	}
	objs, err := m.wifiObj.GetManagedObjects()
	if err != nil {
		logger.Error("Failed to get wifi objects:", err)
	}
	for dpath, _ := range objs {
		_, err := m.addObject(dbus.ObjectPath(dpath))
		if err != nil {
			logger.Warning("Failed to add path:", dpath, err)
		}
	}
	objs, err = m.wfdObj.GetManagedObjects()
	if err != nil {
		logger.Error("Failed to get wfd objects:", err)
	}
	for dpath, _ := range objs {
		_, err := m.addObject(dbus.ObjectPath(dpath))
		if err != nil {
			logger.Warning("Failed to add path:", dpath, err)
		}
	}
	m.handleEvent()
	logger.Debug("Links:", m.links)
	logger.Debug("Sinks:", m.sinks)
}

func (m *Miracast) destroy() {
	if m.network != nil {
		networkmanager.DestroyManager(m.network)
		m.network = nil
	}
	if m.wifiObj != nil {
		wifi.DestroyObjectManager(m.wifiObj)
	}

	m.linkLocker.Lock()
	for _, link := range m.links {
		destroyLinkInfo(link)
	}
	m.links = nil
	m.linkLocker.Unlock()

	m.sinkLocker.Lock()
	for _, sink := range m.sinks {
		destroySinkInfo(sink)
	}
	m.sinks = nil
	m.sinkLocker.Unlock()
}

func (m *Miracast) addObject(dpath dbus.ObjectPath) (interface{}, error) {
	if isLinkObjectPath(dpath) {
		return m.addLinkInfo(dpath)
	} else if isSinkObjectPath(dpath) {
		return m.addSinkInfo(dpath)
	}
	return nil, fmt.Errorf("Unknow object dpath: %v", dpath)
}

func (m *Miracast) addLinkInfo(dpath dbus.ObjectPath) (*LinkInfo, error) {
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	if link := m.links.Get(dpath); link != nil {
		// exists, just update
		link.update()
		return nil, fmt.Errorf("The link '%v' has exists", dpath)
	}

	link, err := newLinkInfo(dpath)
	if err != nil {
		logger.Warning("Failed to new link:", err)
		return nil, err
	}
	m.deviceLocker.Lock()
	defer m.deviceLocker.Unlock()
	if m.devices.Get(link.MacAddress) == nil {
		logger.Warningf("The link '%v' unsupported p2p", dpath)
		return nil, fmt.Errorf("Unsupported p2p: %v", dpath)
	}
	m.links = append(m.links, link)
	return link, nil
}

func (m *Miracast) addSinkInfo(dpath dbus.ObjectPath) (*SinkInfo, error) {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	if sink := m.sinks.Get(dpath); sink != nil {
		sink.update()
		return nil, fmt.Errorf("The sink '%v' has exists", dpath)
	}
	sink, err := newSinkInfo(dpath)
	if err != nil {
		logger.Warning("Failed to new sink:", dpath)
		return nil, err
	}
	m.sinks = append(m.sinks, sink)
	return sink, nil
}

func (m *Miracast) removeObject(dpath dbus.ObjectPath) (bool, string) {
	var (
		removed bool = false
		detail  string
	)
	if isLinkObjectPath(dpath) {
		m.linkLocker.Lock()
		defer m.linkLocker.Unlock()
		tmp := m.links.Get(dpath)
		if tmp != nil {
			detail = toJSON(tmp)
			m.links, removed = m.links.Remove(dpath)
		}
	} else if isSinkObjectPath(dpath) {
		m.sinkLocker.Lock()
		defer m.sinkLocker.Unlock()
		tmp := m.sinks.Get(dpath)
		if tmp != nil {
			detail = toJSON(tmp)
			m.sinks, removed = m.sinks.Remove(dpath)
		}
	}
	return removed, detail
}

func (m *Miracast) handleEvent() {
	m.wifiObj.ConnectInterfacesAdded(func(dpath oldDBusLib.ObjectPath, detail map[string]map[string]oldDBusLib.Variant) {
		logger.Debug("[WIFI Added]:", dpath)
		v, err := m.addObject(dbus.ObjectPath(dpath))
		if err == nil {
			m.service.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wifiObj.ConnectInterfacesRemoved(func(dpath oldDBusLib.ObjectPath, details []string) {
		logger.Debug("[WIFI Removed]:", dpath)
		if ok, detail := m.removeObject(dbus.ObjectPath(dpath)); ok {
			m.service.Emit(m, "Removed", dpath, detail)
		}
	})

	m.wfdObj.ConnectInterfacesAdded(func(dpath oldDBusLib.ObjectPath, detail map[string]map[string]oldDBusLib.Variant) {
		logger.Debug("[WFD Added]:", dpath)
		v, err := m.addObject(dbus.ObjectPath(dpath))
		if err == nil {
			m.service.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wfdObj.ConnectInterfacesRemoved(func(dpath oldDBusLib.ObjectPath, details []string) {
		logger.Debug("[WFD Added]:", dpath)
		if ok, detail := m.removeObject(dbus.ObjectPath(dpath)); ok {
			m.service.Emit(m, "Removed", dpath, detail)
		}
	})

	if !ddbus.IsSystemBusActivated(m.network.DestName) {
		logger.Warning("Network service no activation")
		return
	}
	m.network.ConnectDeviceAdded(func(dpath oldDBusLib.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()
		logger.Debug("[Device Added]:", dpath)
		devices, err := iw.ListWirelessInfo()
		if err != nil {
			logger.Warning("[DeviceAdded] Failed to list wireless devices:", err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})
	m.network.ConnectDeviceRemoved(func(dpath oldDBusLib.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()
		logger.Debug("[Device Removed]:", dpath)
		devices, err := iw.ListWirelessInfo()
		if err != nil {
			logger.Warning("[DeviceRemoved] Failed to list wireless devices:", err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})
}

func (m *Miracast) enableWirelessManaged(macAddress string, enabled bool) error {
	if !ddbus.IsSystemBusActivated(m.network.DestName) {
		return fmt.Errorf("Network service no activation")
	}

	devPaths, err := m.network.GetAllDevices()
	if err != nil {
		return err
	}

	for _, devPath := range devPaths {
		wireless, err := networkmanager.NewDeviceWireless(nmDBusServiceName, devPath)
		if err != nil {
			logger.Warning("Failed to create device:", err)
			continue
		}
		if strings.ToLower(wireless.HwAddress.Get()) != macAddress {
			networkmanager.DestroyDeviceWireless(wireless)
			continue
		}

		networkmanager.DestroyDeviceWireless(wireless)
		dev, err := networkmanager.NewDevice(nmDBusServiceName, devPath)
		if err != nil {
			return err
		}

		if dev.Managed.Get() != enabled {
			dev.Managed.Set(enabled)
		}

		// wait 'Managed' value changed
		for {
			time.Sleep(time.Millisecond * 10)
			if dev.Managed.Get() == enabled {
				logger.Info("[enableWirelessManaged] Device managed has changed:", macAddress, enabled)
				break
			}
		}
		networkmanager.DestroyDevice(dev)
		break
	}
	return nil
}

func (m *Miracast) disconnectSink(dpath dbus.ObjectPath) error {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	sink := m.sinks.Get(dpath)
	if sink == nil {
		logger.Warning("Not found the sink:", dpath)
		return fmt.Errorf("Not found the sink: %v", dpath)
	}

	sink.locker.Lock()
	defer sink.locker.Unlock()
	err := sink.Teardown()
	if err != nil {
		logger.Warning("[disconnectSink] Failed to teardown:", err)
	}

	if sink.peer == nil {
		logger.Warning("No peer found in sink:", dpath)
		return fmt.Errorf("Not found the peer in sink: %v", dpath)
	}

	err = sink.peer.Disconnect()
	if err != nil {
		logger.Warning("[DisconnectSink] Failed to disconnect:", dpath, err)
		return err
	}
	delete(m.connectingSinks, dpath)
	return nil
}

// Remove it if miracle-dispd work fine
func (m *Miracast) ensureMiracleActive() {
	var failedCount = 0
	for {
		if failedCount > 10 {
			logger.Warning("Miracle failure too many, break")
			break
		}
		time.Sleep(time.Second * 10)
		if len(m.devices) == 0 {
			continue
		}

		_, err := m.wfdObj.GetManagedObjects()
		_, err = m.wifiObj.GetManagedObjects()
		if err != nil {
			logger.Debug("Failed to connect miracle:", err)
			failedCount += 1
		}
	}
}

func (*Miracast) GetInterfaceName() string {
	return dbusInterface
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
