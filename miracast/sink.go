package miracast

import (
	"dbus/com/deepin/daemon/audio"
	"dbus/org/freedesktop/miracle/wfd"
	"fmt"
	"os"
	"pkg.deepin.io/lib/dbus"
	"strings"
)

type SinkInfo struct {
	Path dbus.ObjectPath

	core *wfd.Sink
}
type SinkInfos []*SinkInfo

func newSinkInfo(dpath dbus.ObjectPath) (*SinkInfo, error) {
	core, err := wfd.NewSink(wfdDest, dpath)
	if err != nil {
		return nil, err
	}
	return &SinkInfo{
		Path: dpath,
		core: core,
	}, nil
}

func (sink *SinkInfo) StartSession(x, y, w, h uint32) error {
	var (
		// format: 'x://:0.0'
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

func (sink *SinkInfo) TearDown() error {
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
