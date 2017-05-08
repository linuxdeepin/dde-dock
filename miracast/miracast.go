package miracast

import (
	"dbus/org/freedesktop/miracle/wfd"
	"dbus/org/freedesktop/miracle/wifi"
	"dbus/org/freedesktop/networkmanager"
	"encoding/json"
	"fmt"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
	"time"
)

const (
	dbusDest = "com.deepin.daemon.Miracast"
	dbusPath = "/com/deepin/daemon/Miracast"
	dbusIFC  = dbusDest
)

const (
	wifiDest = "org.freedesktop.miracle.wifi"
	wifiPath = "/org/freedesktop/miracle/wifi"
	linkPath = "/org/freedesktop/miracle/wifi/link/"
	linkIFC  = "org.freedesktop.miracle.wifi.Link"
	peerPath = "/org/freedesktop/miracle/wifi/peer/"
	peerIFC  = "org.freedesktop.miracle.wifi.Peer"

	wfdDest  = "org.freedesktop.miracle.wfd"
	wfdPath  = "/org/freedesktop/miracle/wfd"
	sinkPath = "/org/freedesktop/miracle/wfd/sink/"
	sinkIFC  = "org.freedesktop.miracle.wfd.Sink"

	nmDest = "org.freedesktop.NetworkManager"
	nmPath = "/org/freedesktop/NetworkManager"
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
	devices WirelessInfos

	linkLocker   sync.Mutex
	sinkLocker   sync.Mutex
	deviceLocker sync.Mutex

	managingLinks   map[dbus.ObjectPath]bool
	connectingSinks map[dbus.ObjectPath]bool

	Added   func(dbus.ObjectPath, string)
	Removed func(dbus.ObjectPath, string)
	Event   func(uint8, dbus.ObjectPath)
}

func newMiracast() (*Miracast, error) {
	network, err := networkmanager.NewManager(nmDest, nmPath)
	if err != nil {
		return nil, err
	}

	wifiObj, err := wifi.NewObjectManager(wifiDest, wifiPath)
	if err != nil {
		networkmanager.DestroyManager(network)
		return nil, err
	}

	wfdObj, err := wfd.NewObjectManager(wfdDest, wfdPath)
	if err != nil {
		networkmanager.DestroyManager(network)
		wifi.DestroyObjectManager(wifiObj)
		return nil, err
	}

	return &Miracast{
		network:         network,
		wifiObj:         wifiObj,
		wfdObj:          wfdObj,
		connectingSinks: make(map[dbus.ObjectPath]bool),
		managingLinks:   make(map[dbus.ObjectPath]bool),
	}, nil
}

func (m *Miracast) init() {
	devices, err := ListWirelessInfo()
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
		_, err := m.addObject(dpath)
		if err != nil {
			logger.Warning("Failed to add path:", dpath, err)
		}
	}
	objs, err = m.wfdObj.GetManagedObjects()
	if err != nil {
		logger.Error("Failed to get wfd objects:", err)
	}
	for dpath, _ := range objs {
		_, err := m.addObject(dpath)
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
	m.wifiObj.ConnectInterfacesAdded(func(dpath dbus.ObjectPath, detail map[string]map[string]dbus.Variant) {
		logger.Debug("[WIFI Added]:", dpath)
		v, err := m.addObject(dpath)
		if err == nil {
			dbus.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wifiObj.ConnectInterfacesRemoved(func(dpath dbus.ObjectPath, details []string) {
		logger.Debug("[WIFI Removed]:", dpath)
		if ok, detail := m.removeObject(dpath); ok {
			dbus.Emit(m, "Removed", dpath, detail)
		}
	})

	m.wfdObj.ConnectInterfacesAdded(func(dpath dbus.ObjectPath, detail map[string]map[string]dbus.Variant) {
		logger.Debug("[WFD Added]:", dpath)
		v, err := m.addObject(dpath)
		if err == nil {
			dbus.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wfdObj.ConnectInterfacesRemoved(func(dpath dbus.ObjectPath, details []string) {
		logger.Debug("[WFD Added]:", dpath)
		if ok, detail := m.removeObject(dpath); ok {
			dbus.Emit(m, "Removed", dpath, detail)
		}
	})

	if !ddbus.IsSystemBusActivated(m.network.DestName) {
		logger.Warning("Network service no activation")
		return
	}
	m.network.ConnectDeviceAdded(func(dpath dbus.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()
		logger.Debug("[Device Added]:", dpath)
		devices, err := ListWirelessInfo()
		if err != nil {
			logger.Warning("[DeviceAdded] Failed to list wireless devices:", err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})
	m.network.ConnectDeviceRemoved(func(dpath dbus.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()
		logger.Debug("[Device Removed]:", dpath)
		devices, err := ListWirelessInfo()
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
		wireless, err := networkmanager.NewDeviceWireless(nmDest, devPath)
		if err != nil {
			logger.Warning("Failed to create device:", err)
			continue
		}
		if strings.ToLower(wireless.HwAddress.Get()) != macAddress {
			networkmanager.DestroyDeviceWireless(wireless)
			continue
		}

		networkmanager.DestroyDeviceWireless(wireless)
		dev, err := networkmanager.NewDevice(nmDest, devPath)
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

func (*Miracast) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
