package mpris

import (
	"dbus/com/deepin/daemon/audio"
	"dbus/com/deepin/daemon/display"
	"dbus/com/deepin/daemon/keybinding"
	"dbus/org/freedesktop/dbus"
	"dbus/org/freedesktop/login1"
	"fmt"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/mpris")

type Manager struct {
	mediakey    *keybinding.Mediakey
	login       *login1.Manager
	disp        *display.Display
	dbusDaemon  *dbus.DBusDaemon
	audioDaemon *audio.Audio

	prevPlayer string
}

func NewManager() (*Manager, error) {
	var m = new(Manager)

	var err error
	m.mediakey, err = keybinding.NewMediakey("com.deepin.daemon.Keybinding",
		"/com/deepin/daemon/Keybinding/Mediakey")
	if err != nil {
		return nil, err
	}

	m.login, err = login1.NewManager("org.freedesktop.login1",
		"/org/freedesktop/login1")
	if err != nil {
		return nil, err
	}

	m.dbusDaemon, err = dbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return nil, err
	}

	m.disp, err = display.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		logger.Warning("Create display connection failed:", err)
	}

	m.audioDaemon, err = audio.NewAudio("com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio")
	if err != nil {
		logger.Warning("Create audio connection failed:", err)
	}

	return m, nil
}

func (m *Manager) destroy() {
	keybinding.DestroyMediakey(m.mediakey)
	login1.DestroyManager(m.login)
}

func (m *Manager) changeBrightness(raised, pressed bool) {
	if m.disp == nil || !pressed {
		return
	}

	var delta float64 = 0.05
	if !raised {
		delta = -0.05
	}
	for output, v := range m.disp.Brightness.Get() {
		v = v + delta
		err := m.disp.SetBrightness(output, v)
		if err != nil {
			logger.Warning("[SetBrightness] failed:", output, v)
		}
	}
}

func (m *Manager) setMute(pressed bool) {
	if !pressed {
		return
	}

	sink, err := m.getDefaultSink()
	if err != nil {
		logger.Warning("[GetDefaultSink] failed:", err)
		return
	}
	sink.SetMute(!sink.Mute.Get())
}

func (m *Manager) changeVolume(raised, pressed bool) {
	if m.audioDaemon == nil || !pressed {
		return
	}

	sink, err := m.getDefaultSink()
	if err != nil {
		logger.Warning("[GetDefaultSink] failed:", err)
		return
	}

	if sink.Mute.Get() {
		sink.SetMute(false)
	}

	var delta float64 = 0.1
	if !raised {
		delta = -0.1
	}

	v := sink.Volume.Get() + delta
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1.0
	}
	sink.SetVolume(v, true)
}

func (m *Manager) getDefaultSink() (*audio.AudioSink, error) {
	if m.audioDaemon == nil {
		return nil, fmt.Errorf("Can not connect audio daemon")
	}

	sinkPath, err := m.audioDaemon.GetDefaultSink()
	if err != nil {
		return nil, err
	}

	sink, err := audio.NewAudioSink("com.deepin.daemon.Audio", sinkPath)
	if err != nil {
		return nil, err
	}

	return sink, nil
}
