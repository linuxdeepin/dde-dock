package miracast

import (
	"path"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func (m *Miracast) handleLinkManaged(link *LinkInfo) {
	link.core.Managed.ConnectChanged(func() {
		if link.core == nil {
			return
		}
		link.update()
		logger.Debugf("Link '%s' managed changed: %v", link.Path, link.Managed)
		v, ok := m.managingLinks[link.Path]
		if !ok {
			return
		}

		if link.Managed {
			link.core.WfdSubelements.Set("000600111c4400c8")
			// TODO: wait WfdSubelements changed
			dbus.Emit(m, "Event", EventLinkManaged, link.Path)
		} else {
			link.core.WfdSubelements.Set("")
			dbus.Emit(m, "Event", EventLinkUnmanaged, link.Path)
			err := m.enableWirelessManaged(link.MacAddress, true)
			if err != nil {
				logger.Error("Enable nm device failed:", link.MacAddress, err)
			}
		}

		logger.Debugf("[handleLinkChanged] link '%s' managed: %v, should be: %v",
			link.Path, link.Managed, v)
		if v == link.Managed {
			delete(m.managingLinks, link.Path)
		}
	})
}

func (m *Miracast) handlePeerConnected(peer *PeerInfo, x, y, w, h uint32) {
	var hasConnected = false
	peer.core.Connected.ConnectChanged(func() {
		if peer.core == nil {
			return
		}
		peer.update()
		logger.Debugf("Peer '%s' connected changed: %v", peer.Path, peer.Connected)
		_, ok := m.connectingPeers[peer.Path]
		if !ok {
			return
		}

		defer delete(m.connectingPeers, peer.Path)
		if !peer.Connected {
			logger.Debugf("Peer '%s' was disconnect", peer.Path)
			return
		}

		hasConnected = true
		dbus.Emit(m, "Event", EventPeerConnected, peer.Path)
		m.doConnect(peer, x, y, w, h)
	})

	// timeout
	time.AfterFunc(time.Second*60, func() {
		if peer.core == nil || hasConnected {
			return
		}
		m.Disconnect(peer.Path)
		dbus.Emit(m, "Event", EventPeerConnectedFailed, peer.Path)
	})
}

func (m *Miracast) doConnect(peer *PeerInfo, x, y, w, h uint32) {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	sink := m.sinks.Get(sinkPath + dbus.ObjectPath(path.Base(string(peer.Path))))
	if sink == nil {
		logger.Warning("Invalid peer path:", peer.Path)
		return
	}
	err := sink.StartSession(x, y, w, h)
	if err != nil {
		logger.Error("Failed to start session:", peer.Path, err)
	}
}
