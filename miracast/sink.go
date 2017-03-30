package miracast

import (
	"dbus/com/deepin/daemon/audio"
	"dbus/org/freedesktop/miracle/wfd"
	"dbus/org/freedesktop/miracle/wifi"
	"fmt"
	"os"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
)

type SinkInfo struct {
	Name      string
	P2PMac    string
	Interface string
	Connected bool
	Path      dbus.ObjectPath
	LinkPath  dbus.ObjectPath

	core   *wfd.Sink
	peer   *wifi.Peer
	locker sync.Mutex
}
type SinkInfos []*SinkInfo

func newSinkInfo(dpath dbus.ObjectPath) (*SinkInfo, error) {
	core, err := wfd.NewSink(wfdDest, dpath)
	if err != nil {
		return nil, err
	}
	peer, err := wifi.NewPeer(wifiDest, core.Peer.Get())
	if err != nil {
		wfd.DestroySink(core)
		return nil, err
	}
	var sink = &SinkInfo{
		Path: dpath,
		core: core,
		peer: peer,
	}
	sink.update()
	return sink, nil
}

func destroySinkInfo(info *SinkInfo) {
	info.locker.Lock()
	defer info.locker.Unlock()
	if info.core != nil {
		wfd.DestroySink(info.core)
		info.core = nil
	}
	if info.peer != nil {
		wifi.DestroyPeer(info.peer)
		info.peer = nil
	}
}

func (sink *SinkInfo) update() {
	sink.locker.Lock()
	defer sink.locker.Unlock()
	if sink.core == nil || sink.peer == nil {
		return
	}
	sink.Name = sink.peer.FriendlyName.Get()
	sink.P2PMac = sink.peer.P2PMac.Get()
	sink.Interface = sink.peer.Interface.Get()
	sink.Connected = sink.peer.Connected.Get()
	sink.LinkPath = sink.peer.Link.Get()
}

func (sink *SinkInfo) StartSession(x, y, w, h uint32) error {
	var (
		// format: 'x://:0'
		dpy       = "x://" + os.Getenv("DISPLAY")
		xauth     = os.Getenv("XAUTHORITY")
		audioSink = getAudioSink()
	)
	logger.Debug("[StartSession] args:", xauth, dpy, x, y, w, h, audioSink)
	stateId, err := sink.core.StartSession(xauth, dpy, x, y, w, h, audioSink)
	if err != nil {
		return err
	}
	// TODO: handle state
	logger.Debug("[StartSession] state id:", stateId)
	return nil
}

func (sink *SinkInfo) Teardown() error {
	var p = sink.core.Session.Get()
	if p == "/" {
		return fmt.Errorf("No session found")
	}

	session, err := wfd.NewSession(wfdDest, p)
	if err != nil {
		return err
	}
	defer wfd.DestroySession(session)
	return session.Teardown()
}

func (sinks SinkInfos) Get(dpath dbus.ObjectPath) *SinkInfo {
	if !isSinkObjectPath(dpath) {
		return nil
	}
	for _, sink := range sinks {
		if sink.Path == dpath {
			return sink
		}
	}
	return nil
}

func (sinks SinkInfos) Remove(dpath dbus.ObjectPath) (SinkInfos, bool) {
	var (
		tmp    SinkInfos
		exists bool
	)
	for _, sink := range sinks {
		if sink.Path == dpath {
			exists = true
			continue
		}
		tmp = append(tmp, sink)
	}
	return tmp, exists
}

func isSinkObjectPath(dpath dbus.ObjectPath) bool {
	return strings.Contains(string(dpath), sinkPath)
}

func getAudioSink() string {
	return ""
	obj, err := audio.NewAudio("com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio")
	if err != nil {
		return ""
	}
	defer audio.DestroyAudio(obj)

	sink, err := audio.NewAudioSink("com.deepin.daemon.Audio",
		obj.DefaultSink.Get())
	if err != nil {
		return ""
	}
	defer audio.DestroyAudioSink(sink)
	return sink.Name.Get() + ".monitor"
}

func isPeerObjectPath(dpath dbus.ObjectPath) bool {
	return strings.Contains(string(dpath), peerPath)
}
