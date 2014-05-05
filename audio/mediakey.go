package main

import libsound "dbus/com/deepin/api/sound"

var _player interface{}

func (audio *Audio) listenMediaKey() {
	player, err := libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	if err != nil {
		Logger.Error("Can't create com.deepin.api.Sound! Sound feedback support will be disabled", err)
	}
	_player = player

	mediakeyObj.ConnectAudioMute(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Info("^^ Sink is nil")
				return
			}
			sink.SetSinkMute(!sink.Mute)
			player.PlaySystemSound("audio-volume-change")
		}
	})
	mediakeyObj.ConnectAudioUp(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Info("^^ Sink is nil")
				return
			}
			volume := int32(sink.Volume + _VOLUME_STEP)
			if volume < 0 {
				volume = 0
			} else if volume > 100 {
				volume = 100
			}
			if sink.Volume < 100 {
				sink.setSinkVolume(uint32(volume))
				sink.setSinkMute(false)
			}
			player.PlaySystemSound("audio-volume-change")
		}
	})
	mediakeyObj.ConnectAudioDown(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				Logger.Info("^^ Sink is nil")
				return
			}
			volume := int32(sink.Volume - _VOLUME_STEP)
			if volume < 0 {
				volume = 0
			}
			sink.setSinkVolume(uint32(volume))
			sink.setSinkMute(false)
			player.PlaySystemSound("audio-volume-change")
		}
	})
}
