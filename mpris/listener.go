/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mpris

func (m *Manager) listenMediakey() {
	m.mediakey.ConnectAudioPlay(func(pressed bool) {
		m.playerAction(actionTypePlay, pressed)
	})

	m.mediakey.ConnectAudioPause(func(pressed bool) {
		m.playerAction(actionTypePause, pressed)
	})

	m.mediakey.ConnectAudioStop(func(pressed bool) {
		m.playerAction(actionTypeStop, pressed)
	})

	m.mediakey.ConnectAudioPrevious(func(pressed bool) {
		m.playerAction(actionTypePrevious, pressed)
	})

	m.mediakey.ConnectAudioNext(func(pressed bool) {
		m.playerAction(actionTypeNext, pressed)
	})

	m.mediakey.ConnectAudioRewind(func(pressed bool) {
		m.playerAction(actionTypeRewind, pressed)
	})

	m.mediakey.ConnectAudioForward(func(pressed bool) {
		m.playerAction(actionTypeForward, pressed)
	})

	m.mediakey.ConnectAudioRepeat(func(pressed bool) {
		m.playerAction(actionTypeRepeat, pressed)
	})

	m.mediakey.ConnectLaunchBrowser(func(pressed bool) {
		go execByMime(mimeTypeBrowser, pressed)
	})

	m.mediakey.ConnectLaunchEmail(func(pressed bool) {
		go execByMime(mimeTypeEmail, pressed)
	})

	m.mediakey.ConnectLaunchCalculator(func(pressed bool) {
		go execByMime(mimeTypeCalc, pressed)
	})

	m.mediakey.ConnectBrightnessUp(func(pressed bool) {
		m.changeBrightness(true, pressed)
	})

	m.mediakey.ConnectBrightnessDown(func(pressed bool) {
		m.changeBrightness(false, pressed)
	})

	m.mediakey.ConnectAudioUp(func(pressed bool) {
		m.changeVolume(true, pressed)
	})

	m.mediakey.ConnectAudioDown(func(pressed bool) {
		m.changeVolume(false, pressed)
	})

	m.mediakey.ConnectAudioMute(func(pressed bool) {
		m.setMute(pressed)
	})

	m.login.ConnectPrepareForSleep(func(actived bool) {
		m.pauseAllPlayer(actived)
	})

	m.mediakey.ConnectPowerSleep(func(pressed bool) {
		m.suspend(pressed)
	})

	m.mediakey.ConnectPowerSuspend(func(pressed bool) {
		m.suspend(pressed)
	})

	m.mediakey.ConnectEject(func(pressed bool) {
		go m.eject(pressed)
	})

	m.mediakey.ConnectAudioMedia(func(pressed bool) {
		go execByMime(mimeTypeAudioMedia, pressed)
	})

}
