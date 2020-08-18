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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/godbus/dbus"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.miracle.wfd"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.miracle.wifi"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	"pkg.deepin.io/dde/daemon/iw"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	dbusServiceName = "com.deepin.daemon.Miracast"
	dbusPath        = "/com/deepin/daemon/Miracast"
	dbusInterface   = dbusServiceName
)

const (
	linkDBusPathPrefix    = "/org/freedesktop/miracle/wifi/link/"
	peerDBusPathPrefix    = "/org/freedesktop/miracle/wifi/peer/"
	sinkDBusPathPrefix    = "/org/freedesktop/miracle/wfd/sink/"
	sessionDBusPathPrefix = "/org/freedesktop/miracle/wfd/session/"
)

const (
	defaultTimeout  = time.Second * 30
	defaultInterval = 600 * time.Millisecond
)

type Miracast struct {
	sysSigLoop *dbusutil.SignalLoop
	wifiObj    *wifi.Wifi
	wfdObj     *wfd.Wfd
	network    *networkmanager.Manager
	sysBusObj  *ofdbus.DBus

	links      LinkInfos
	linkLocker sync.Mutex

	sinks      SinkInfos
	sinkLocker sync.Mutex

	devices      iw.WirelessInfos
	deviceLocker sync.Mutex

	inited bool
	locker sync.Mutex

	service *dbusutil.Service
	//nolint
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
	//nolint
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
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	network := networkmanager.NewManager(sysBus)
	wifiObj := wifi.NewWifi(sysBus)
	wfdObj := wfd.NewWfd(sysBus)
	sysBusObj := ofdbus.NewDBus(sysBus)

	return &Miracast{
		inited:     false,
		service:    service,
		sysSigLoop: dbusutil.NewSignalLoop(sysBus, 10),
		network:    network,
		wifiObj:    wifiObj,
		wfdObj:     wfdObj,
		sysBusObj:  sysBusObj,
	}, nil
}

func (m *Miracast) init() {
	m.locker.Lock()
	if m.inited {
		m.locker.Unlock()
		return
	}
	m.inited = true
	m.locker.Unlock()

	devices, err := iw.ListWirelessInfo()
	if err != nil {
		logger.Error("failed to list wireless info:", err)
	}
	logger.Debugf("all devices: %#v", devices)

	m.devices = devices.ListMiracastDevice()
	if len(m.devices) == 0 {
		m.handleEvent()
		return
	}
	objs, err := m.wifiObj.GetManagedObjects(0)
	if err != nil {
		logger.Error("failed to get wifi objects:", err)
	}
	for objPath := range objs {
		_, err := m.addObject(objPath)
		if err != nil {
			logger.Warning("failed to add path:", objPath, err)
		}
	}
	objs, err = m.wfdObj.GetManagedObjects(0)
	if err != nil {
		logger.Error("failed to get wfd objects:", err)
	}
	for objPath := range objs {
		_, err := m.addObject(objPath)
		if err != nil {
			logger.Warning("failed to add path:", objPath, err)
		}
	}
	m.handleEvent()
	logger.Debug("Links:", m.links)
	logger.Debug("Sinks:", m.sinks)
}

func (m *Miracast) destroy() {
	m.wifiObj.RemoveHandler(proxy.RemoveAllHandlers)
	m.wfdObj.RemoveHandler(proxy.RemoveAllHandlers)
	m.network.RemoveHandler(proxy.RemoveAllHandlers)

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

func (m *Miracast) addObject(objPath dbus.ObjectPath) (interface{}, error) {
	if isLinkObjectPath(objPath) {
		logger.Debug("add link", objPath)
		return m.addLinkInfo(objPath)
	} else if isSinkObjectPath(objPath) {
		logger.Debug("add sink", objPath)
		return m.addSinkInfo(objPath)
	} else if isPeerObjectPath(objPath) {
		logger.Debug("add peer", objPath)
	} else if isSessionObjectPath(objPath) {
		logger.Debug("add session", objPath)
	} else {
		logger.Debug("add", objPath)
	}
	return nil, fmt.Errorf("unknow object objPath: %v", objPath)
}

func (m *Miracast) addLinkInfo(objPath dbus.ObjectPath) (*LinkInfo, error) {
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	if link := m.links.Get(objPath); link != nil {
		// exists, just update
		link.update()
		return nil, fmt.Errorf("link '%v' has exists", objPath)
	}

	link, err := newLinkInfo(objPath)
	if err != nil {
		logger.Warning("failed to new link:", err)
		return nil, err
	}
	link.connectSignal(m)

	m.deviceLocker.Lock()
	defer m.deviceLocker.Unlock()
	if m.devices.Get(link.MacAddress) == nil {
		logger.Warningf("link '%v' unsupported p2p", objPath)
		return nil, fmt.Errorf("unsupported p2p: %v", objPath)
	}
	m.links = append(m.links, link)
	return link, nil
}

func (m *Miracast) addSinkInfo(objPath dbus.ObjectPath) (*SinkInfo, error) {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	if sink := m.sinks.Get(objPath); sink != nil {
		sink.update()
		return nil, fmt.Errorf("sink '%v' has exists", objPath)
	}
	sink, err := newSinkInfo(objPath)
	if err != nil {
		logger.Warning("failed to new sink:", objPath)
		return nil, err
	}

	sink.connectSignal(m)
	m.sinks = append(m.sinks, sink)
	return sink, nil
}

func (m *Miracast) removeObject(objPath dbus.ObjectPath) (bool, interface{}) {
	var (
		removed    bool
		removedObj interface{}
	)
	if isLinkObjectPath(objPath) {
		logger.Debug("remove link", objPath)
		m.linkLocker.Lock()

		link := m.links.Get(objPath)
		if link != nil {
			destroyLinkInfo(link)
			m.links, removed = m.links.Remove(objPath)
			removedObj = link
		}

		m.linkLocker.Unlock()
	} else if isSinkObjectPath(objPath) {
		logger.Debug("remove sink", objPath)
		m.sinkLocker.Lock()

		sink := m.sinks.Get(objPath)
		if sink != nil {
			destroySinkInfo(sink)
			m.sinks, removed = m.sinks.Remove(objPath)
			removedObj = sink
		}

		m.sinkLocker.Unlock()
	} else if isSessionObjectPath(objPath) {
		logger.Debug("remove session", objPath)
	} else if isPeerObjectPath(objPath) {
		logger.Debug("remove peer", objPath)
	} else {
		logger.Debug("remove", objPath)
	}
	return removed, removedObj
}

func (m *Miracast) emitSignalAdded(path dbus.ObjectPath, detailJSON string) {
	err := m.service.Emit(m, "Added", path, detailJSON)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Miracast) emitSignalRemoved(path dbus.ObjectPath, detailJSON string) {
	err := m.service.Emit(m, "Removed", path, detailJSON)
	if err != nil {
		logger.Warning(err)
	}
}

func eventTypeToStr(eventType uint8) string {
	switch eventType {
	case EventLinkManaged:
		return "LinkManaged"
	case EventLinkUnmanaged:
		return "LinkUnmanaged"
	case EventSinkConnected:
		return "SinkConnected"
	case EventSinkConnectedFailed:
		return "SinkConnectedFailed"
	case EventSinkDisconnected:
		return "SinkDisconnected"
	default:
		panic(fmt.Errorf("unknown event type %d", eventType))
	}
}

func (m *Miracast) emitSignalEvent(eventType uint8, path dbus.ObjectPath) {
	logger.Debug("emit signal event", eventTypeToStr(eventType), path)
	err := m.service.Emit(m, "Event", eventType, path)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Miracast) startSession(sink *SinkInfo, x, y, w, h uint32) {
	session, err := sink.core.Session().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if session != "/" {
		logger.Debug("sink session had connected:", sink.Path, session)
		return
	}
	err = sink.StartSession(x, y, w, h)
	if err != nil {
		logger.Error("failed to start session:", sink.Path, err)
	}
}

func (m *Miracast) handleEvent() {
	m.sysSigLoop.Start()
	m.wifiObj.InitSignalExt(m.sysSigLoop, true)
	_, err := m.wifiObj.ConnectInterfacesAdded(func(objectPath dbus.ObjectPath,
		interfacesAndProperties map[string]map[string]dbus.Variant) {
		v, err := m.addObject(objectPath)
		if err != nil {
			m.emitSignalAdded(objectPath, toJSON(v))
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	_, err = m.wifiObj.ConnectInterfacesRemoved(func(objectPath dbus.ObjectPath, interfaces []string) {
		if ok, v := m.removeObject(objectPath); ok {
			m.emitSignalRemoved(objectPath, toJSON(v))
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	m.wfdObj.InitSignalExt(m.sysSigLoop, true)
	_, err = m.wfdObj.ConnectInterfacesAdded(func(objectPath dbus.ObjectPath,
		interfacesAndProperties map[string]map[string]dbus.Variant) {
		v, err := m.addObject(objectPath)
		if err == nil {
			m.emitSignalAdded(objectPath, toJSON(v))
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	_, err = m.wfdObj.ConnectInterfacesRemoved(func(objectPath dbus.ObjectPath, interfaces []string) {
		if ok, v := m.removeObject(objectPath); ok {
			m.emitSignalRemoved(dbus.ObjectPath(objectPath), toJSON(v))
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	m.network.InitSignalExt(m.sysSigLoop, true)
	_, err = m.network.ConnectDeviceAdded(func(devicePath dbus.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()

		logger.Debug("device added", devicePath)
		devices, err := iw.ListWirelessInfo()
		if err != nil {
			logger.Warning(err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = m.network.ConnectDeviceRemoved(func(devicePath dbus.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()

		logger.Debug("device removed", devicePath)
		devices, err := iw.ListWirelessInfo()
		if err != nil {
			logger.Warning(err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Miracast) enableWirelessManaged(interfaceName string, enabled bool) error {
	logger.Debug("call enableWirelessManaged", interfaceName, enabled)
	has, err := m.sysBusObj.NameHasOwner(0, m.network.ServiceName_())
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	devPaths, err := m.network.GetAllDevices(0)
	if err != nil {
		return err
	}

	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	const nmDeviceTypeWifi = 2

	for _, devPath := range devPaths {
		d, _ := networkmanager.NewDevice(sysBus, devPath)
		devType, err := d.DeviceType().Get(0)
		if err != nil {
			logger.Warning(err)
			continue
		}

		if devType != nmDeviceTypeWifi {
			continue
		}

		ifcName, err := d.Interface().Get(0)
		if err != nil {
			logger.Warning(err)
			continue
		}
		if ifcName != interfaceName {
			continue
		}

		managed, err := d.Managed().Get(0)
		if err != nil {
			return err
		}

		if managed != enabled {
			err = d.Managed().Set(0, enabled)
			if err != nil {
				return err
			}
		}

		err = waitNmDeviceManaged(d, enabled)
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func waitNmDeviceManaged(device *networkmanager.Device, wantManged bool) error {
	name := fmt.Sprintf("device %s manged", device.Path_())
	return waitChange(name, wantManged, func() (b bool, err error) {
		return device.Managed().Get(0)
	})
}

func (m *Miracast) disconnectSink(objPath dbus.ObjectPath) error {
	if !isSinkObjectPath(objPath) {
		return fmt.Errorf("invalid sink objPath: %v", objPath)
	}

	m.sinkLocker.Lock()
	sink := m.sinks.Get(objPath)
	m.sinkLocker.Unlock()
	if sink == nil {
		logger.Warning("not found the sink:", objPath)
		return fmt.Errorf("not found the sink: %v", objPath)
	}

	sink.locker.Lock()
	defer sink.locker.Unlock()
	err := sink.TeardownSession()
	if err != nil {
		logger.Warning("[disconnectSink] Failed to teardown:", err)
	}

	if sink.peer == nil {
		logger.Warning("no peer found in sink:", objPath)
		return fmt.Errorf("not found the peer in sink: %v", objPath)
	}

	err = sink.peer.Disconnect(0)
	if err != nil {
		logger.Warning("[DisconnectSink] Failed to disconnect:", objPath, err)
		return err
	}
	return nil
}

func (*Miracast) GetInterfaceName() string {
	return dbusInterface
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
