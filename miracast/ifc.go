package miracast

import (
	"fmt"
	"path"
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	EventLinkManaged uint8 = iota + 1
	EventLinkUnmanaged
	EventPeerConnected
	EventPeerConnectedFailed
)

func (m *Miracast) ListLinks() LinkInfos {
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	return m.links
}

func (m *Miracast) ListPeers() PeerInfos {
	m.peerLocker.Lock()
	defer m.peerLocker.Unlock()
	var tmp PeerInfos
	for _, peer := range m.peers {
		sink := m.sinks.Get(sinkPath + dbus.ObjectPath(path.Base(string(peer.Path))))
		if sink == nil {
			logger.Debugf("No sink found for peer: %#v", peer)
			continue
		}
		tmp = append(tmp, peer)
	}
	return tmp
}

func (m *Miracast) ListSinks() SinkInfos {
	m.peerLocker.Lock()
	defer m.peerLocker.Unlock()
	return m.sinks
}

func (m *Miracast) Enable(dpath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(dpath) {
		return fmt.Errorf("Invalid link dpath: %v", dpath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(dpath)
	if link == nil {
		logger.Warning("Not found the link:", dpath)
		return fmt.Errorf("Not found the link: %v", dpath)
	}

	if link.core.Managed.Get() == enabled {
		logger.Debug("Work right, nothing to do:", enabled)
		return nil
	}

	if v, ok := m.managingLinks[dpath]; ok && v == enabled {
		logger.Debug("Link's managed '%s' has been setting to ", enabled)
		return nil
	}
	m.managingLinks[dpath] = enabled

	if enabled {
		err := m.enableWirelessManaged(link.MacAddress, false)
		if err != nil {
			delete(m.managingLinks, dpath)
			logger.Error("Failed to disable manage wireless device:", err)
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}

	m.handleLinkManaged(link)
	return link.EnableManaged(enabled)
}

func (m *Miracast) Scanning(dpath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(dpath) {
		return fmt.Errorf("Invalid link dpath: %v", dpath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(dpath)
	if link == nil {
		logger.Warning("Not found the link:", dpath)
		return fmt.Errorf("Not found the link: %v", dpath)
	}

	link.EnableP2PScanning(enabled)
	// TODO: wait P2PScanning changed
	link.update()
	return nil
}

func (m *Miracast) Connect(dpath dbus.ObjectPath, x, y, w, h uint32) error {
	if !isPeerObjectPath(dpath) {
		return fmt.Errorf("Invalid peer dpath: %v", dpath)
	}

	m.peerLocker.Lock()
	defer m.peerLocker.Unlock()
	peer := m.peers.Get(dpath)
	if peer == nil {
		logger.Warning("Not found the peer:", dpath)
		return fmt.Errorf("Not found the peer: %v", dpath)
	}

	if peer.core.Connected.Get() {
		logger.Debug("Has connected, start session")
		m.doConnect(peer, x, y, w, h)
		return nil
	}

	if v, ok := m.connectingPeers[dpath]; ok && v {
		logger.Debug("[ConnectPeer] peer has connecting:", dpath)
		return nil
	}
	m.connectingPeers[dpath] = true

	m.handlePeerConnected(peer, x, y, w, h)
	err := peer.core.Connect("auto", "")
	if err != nil {
		delete(m.connectingPeers, dpath)
		logger.Error("Failed to connect peer:", err)
		return err
	}

	return nil
}

func (m *Miracast) Disconnect(dpath dbus.ObjectPath) error {
	if !isPeerObjectPath(dpath) {
		return fmt.Errorf("Invalid peer dpath: %v", dpath)
	}
	return m.disconnectPeer(dpath)
}
