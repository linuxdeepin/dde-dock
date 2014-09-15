package audio

import "dbus/com/deepin/daemon/keybinding"

var __keepMediakeyManagerAlive interface{}

func (audio *Audio) setupMediaKeyMonitor() {
	mediaKeyManager, err := keybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	__keepMediakeyManagerAlive = mediaKeyManager
	if err != nil {
		logger.Error("Can't create com.deepin.daemon.Keybinding! mediakey support will be disabled", err)
	}

	mediaKeyManager.ConnectAudioMute(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				logger.Error("Default Sink is nil", audio.DefaultSink)
				return
			}
			sink.SetMute(!sink.Mute)
		}
	})
	mediaKeyManager.ConnectAudioUp(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				logger.Error("Default Sink is nil", audio.DefaultSink)
				return
			}
			if sink.Volume > 1 {
				logger.Warning("ignore add volume bigger than 100% when use MediaKey")
				return
			}
			playFeedbackWithDevice(sink.Name)

			if sink.Mute {
				sink.SetMute(false)
			}

			nv := sink.Volume + 0.1
			if nv > 1 {
				nv = 1
			}
			sink.SetVolume(nv, true)
		}
	})
	mediaKeyManager.ConnectAudioDown(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				logger.Info("Default Sink is nil", audio.DefaultSink)
				return
			}
			if sink.Mute {
				sink.SetMute(false)
			}
			nv := sink.Volume - 0.1
			if nv < 0 {
				nv = 0
			}
			sink.SetVolume(nv, true)
			playFeedbackWithDevice(sink.Name)
		}
	})
}
