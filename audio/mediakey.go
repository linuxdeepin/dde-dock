package main

import libsound "dbus/com/deepin/api/sound"
import "dbus/com/deepin/daemon/keybinding"

var __keepMediakeyManagerAlive interface{}
var __keepPlayerAlive interface{}

func (audio *Audio) listenMediaKey() {
	player, err := libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	mediaKeyManager, err := keybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	__keepMediakeyManagerAlive = mediaKeyManager
	__keepPlayerAlive = player

	if err != nil {
		Logger.Error("Can't create com.deepin.api.Sound! Sound feedback support will be disabled", err)
	}

	mediaKeyManager.ConnectAudioMute(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Error("Default Sink is nil", audio.DefaultSink)
				return
			}
			sink.SetMute(!sink.Mute)
			player.PlaySystemSound("audio-volume-change")
		}
	})
	mediaKeyManager.ConnectAudioUp(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Error("Default Sink is nil", audio.DefaultSink)
				return
			}
			if sink.Volume > 1 {
				Logger.Warning("ignore add volume bigger than 100% when use MediaKey")
				return
			}
			player.PlaySystemSound("audio-volume-change")

			if sink.Mute {
				sink.SetMute(false)
			}

			nv := sink.Volume + 0.1
			if nv > 1 {
				nv = 1
			}
			sink.SetVolume(nv)
		}
	})
	mediaKeyManager.ConnectAudioDown(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Info("Default Sink is nil", audio.DefaultSink)
				return
			}
			if sink.Mute {
				sink.SetMute(false)
			}
			nv := sink.Volume - 0.1
			if nv < 0 {
				nv = 0
			}
			sink.SetVolume(nv)
			player.PlaySystemSound("audio-volume-change")
		}
	})
}
