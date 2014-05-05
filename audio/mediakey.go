package main

import libsound "dbus/com/deepin/api/sound"
import "dbus/com/deepin/daemon/keybinding"
import "fmt"

func (audio *Audio) listenMediaKey() {
	mediakey, err := keybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	if err != nil {
		Logger.Error("Can't create keybinding.MediaKey! Mediakey support will be disabled", err)
		return
	}

	player, err := libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	if err != nil {
		Logger.Error("Can't create com.deepin.api.Sound! Sound feedback support will be disabled", err)
	}

	mediakey.ConnectAudioMute(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			if sink == nil {
				fmt.Println("^^ Sink is nil")
				//Logger.Info("^^ Sink is nil")
				return
			}
			sink.SetSinkMute(!sink.Mute)
			player.PlaySystemSound("audio-volume-change")
		}
	})
	mediakey.ConnectAudioUp(func(pressed bool) {
		if !pressed {
			sink := audio.GetDefaultSink()
			fmt.Println("^^^ set volume up")
			if sink == nil {
				fmt.Println("^^ Sink is nil")
				//Logger.Info("^^ Sink is nil")
				return
			}
			volume := int32(sink.Volume + _VOLUME_STEP)
			fmt.Println("Sink Volume:", volume)
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
	mediakey.ConnectAudioDown(func(pressed bool) {
		if !pressed {
			fmt.Println("^^^ set volume down")
			sink := audio.GetDefaultSink()
			if sink == nil {
				fmt.Println("^^ Sink is nil")
				//Logger.Info("^^ Sink is nil")
				return
			}
			volume := int32(sink.Volume - _VOLUME_STEP)
			fmt.Println("Sink Volume:", volume)
			if volume < 0 {
				volume = 0
			}
			sink.setSinkVolume(uint32(volume))
			sink.setSinkMute(false)
			player.PlaySystemSound("audio-volume-change")
		}
	})
}
