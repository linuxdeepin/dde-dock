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
}
