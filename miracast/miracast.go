package miracast

import (
	"dbus/org/freedesktop/miracle/wfd"
	"dbus/org/freedesktop/miracle/wifi"
	"dbus/org/freedesktop/networkmanager"
	"encoding/json"
	"fmt"
	"path"
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
	peers   PeerInfos
	sinks   SinkInfos
	devices WirelessInfos

	linkLocker   sync.Mutex
	peerLocker   sync.Mutex
	sinkLocker   sync.Mutex
	deviceLocker sync.Mutex

	managingLinks   map[dbus.ObjectPath]bool
	connectingPeers map[dbus.ObjectPath]bool

	Added   func(dbus.ObjectPath, string)
	Removed func(dbus.ObjectPath, string)
	Event   func(uint8, dbus.ObjectPath)
}

func newMiracast() (*Miracast, error) {
	devices, err := ListWirelessInfo()
	if err != nil {
		return nil, err
	}

	logger.Debugf("All devices: %#v", devices)
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
		devices:         devices.ListMiracastDevice(),
		connectingPeers: make(map[dbus.ObjectPath]bool),
		managingLinks:   make(map[dbus.ObjectPath]bool),
	}, nil
}

func (m *Miracast) init() {
	logger.Debug("Devices:", m.devices)
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
	logger.Debug("Peers:", m.peers)
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

	m.peerLocker.Lock()
	for _, peer := range m.peers {
		destroyPeerInfo(peer)
	}
	m.peers = nil
	m.peerLocker.Unlock()
}

func (m *Miracast) addObject(dpath dbus.ObjectPath) (interface{}, error) {
	if isLinkObjectPath(dpath) {
		return m.addLinkInfo(dpath)
	} else if isPeerObjectPath(dpath) {
		return m.addPeerInfo(dpath)
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

func (m *Miracast) addPeerInfo(dpath dbus.ObjectPath) (*PeerInfo, error) {
	m.peerLocker.Lock()
	defer m.peerLocker.Unlock()
	if peer := m.peers.Get(dpath); peer != nil {
		// exists, just update
		peer.update()
		return nil, fmt.Errorf("The peer '%v' has exists", dpath)
	}

	peer, err := newPeerInfo(dpath)
	if err != nil {
		logger.Warning("Failed to new peer:", err)
		return nil, err
	}
	m.peers = append(m.peers, peer)
	return peer, nil
}

func (m *Miracast) addSinkInfo(dpath dbus.ObjectPath) (*SinkInfo, error) {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	if sink := m.sinks.Get(dpath); sink != nil {
		// nothing to do
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
	} else if isPeerObjectPath(dpath) {
		m.peerLocker.Lock()
		defer m.peerLocker.Unlock()
		tmp := m.peers.Get(dpath)
		if tmp != nil {
			detail = toJSON(tmp)
			m.peers, removed = m.peers.Remove(dpath)
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
	m.network.ConnectDeviceAdded(func(dpath dbus.ObjectPath) {
		m.deviceLocker.Lock()
		defer m.deviceLocker.Unlock()
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
		devices, err := ListWirelessInfo()
		if err != nil {
			logger.Warning("[DeviceRemoved] Failed to list wireless devices:", err)
			return
		}
		m.devices = devices.ListMiracastDevice()
	})

	m.wifiObj.ConnectInterfacesAdded(func(dpath dbus.ObjectPath, detail map[string]map[string]dbus.Variant) {
		v, err := m.addObject(dpath)
		if err == nil {
			dbus.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wifiObj.ConnectInterfacesRemoved(func(dpath dbus.ObjectPath, details []string) {
		if ok, detail := m.removeObject(dpath); ok {
			dbus.Emit(m, "Removed", dpath, detail)
		}
	})

	m.wfdObj.ConnectInterfacesAdded(func(dpath dbus.ObjectPath, detail map[string]map[string]dbus.Variant) {
		v, err := m.addObject(dpath)
		if err == nil {
			dbus.Emit(m, "Added", dpath, toJSON(v))
		}
	})

	m.wfdObj.ConnectInterfacesRemoved(func(dpath dbus.ObjectPath, details []string) {
		if ok, detail := m.removeObject(dpath); ok {
			dbus.Emit(m, "Removed", dpath, detail)
		}
	})
}

func (m *Miracast) enableWirelessManaged(macAddress string, enabled bool) error {
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

func (m *Miracast) disconnectPeer(dpath dbus.ObjectPath) error {
	m.peerLocker.Lock()
	defer m.peerLocker.Unlock()
	peer := m.peers.Get(dpath)
	if peer == nil {
		logger.Warning("Not found the peer:", dpath)
		return fmt.Errorf("Not found the peer: %v", dpath)
	}

	m.sinkLocker.Lock()
	sink := m.sinks.Get(sinkPath + dbus.ObjectPath(path.Base(string(peer.Path))))
	if sink != nil {
		// Teardown session
		err := sink.TearDown()
		if err != nil {
			logger.Warning("[disconnectPeer] Teardown failed:", err)
		}
	}
	m.sinkLocker.Unlock()

	err := peer.core.Disconnect()
	if err != nil {
		logger.Warning("[DisconnectPeer] Failed to disconnect:", dpath, err)
		return err
	}
	delete(m.connectingPeers, dpath)
	return nil
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
