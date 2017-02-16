package miracast

import (
	"dbus/org/freedesktop/miracle/wifi"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
)

type PeerInfo struct {
	Name      string
	P2PMac    string
	Interface string
	Connected bool
	Path      dbus.ObjectPath
	LinkPath  dbus.ObjectPath

	locker sync.Mutex
	core   *wifi.Peer
}
type PeerInfos []*PeerInfo

func newPeerInfo(dpath dbus.ObjectPath) (*PeerInfo, error) {
	peer, err := wifi.NewPeer(wifiDest, dpath)
	if err != nil {
		return nil, err
	}

	var info = &PeerInfo{
		Path: dpath,
		core: peer,
	}
	info.update()

	// TODO: handle signals 'GoNegRequest' and 'ProvisionDiscovery'
	return info, nil
}

func destroyPeerInfo(info *PeerInfo) {
	if info.core == nil {
		return
	}
	info.locker.Lock()
	wifi.DestroyPeer(info.core)
	info.core = nil
	info.locker.Unlock()
}

func (peer *PeerInfo) update() {
	peer.locker.Lock()
	defer peer.locker.Unlock()
	if peer.core == nil {
		return
	}
	peer.Name = peer.core.FriendlyName.Get()
	peer.P2PMac = peer.core.P2PMac.Get()
	peer.Interface = peer.core.Interface.Get()
	peer.Connected = peer.core.Connected.Get()
	peer.LinkPath = peer.core.Link.Get()
}

func (peers PeerInfos) Get(dpath dbus.ObjectPath) *PeerInfo {
	if !isPeerObjectPath(dpath) {
		return nil
	}

	for _, peer := range peers {
		if peer.Path == dpath {
			return peer
		}
	}
	return nil
}

func (peers PeerInfos) Remove(dpath dbus.ObjectPath) (PeerInfos, bool) {
	var (
		tmp    PeerInfos
		exists bool = false
	)
	for _, peer := range peers {
		if peer.Path == dpath {
			exists = true
			destroyPeerInfo(peer)
			continue
		}
		tmp = append(tmp, peer)
	}
	return tmp, exists
}

func isPeerObjectPath(dpath dbus.ObjectPath) bool {
	return strings.Contains(string(dpath), peerPath)
}
