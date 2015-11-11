package soundeffect

import (
	"pkg.deepin.io/dde/daemon/soundplayer"
)

func (m *Manager) PlaySystemSound(event string) error {
	return soundplayer.PlaySystemSound(event, "", false)
}

func (m *Manager) PlaySystemSoundSync(event string) error {
	return soundplayer.PlaySystemSound(event, "", true)
}

func (m *Manager) PlayThemeSound(theme, event string) error {
	return soundplayer.PlayThemeSound(theme, event, "", false)
}

func (m *Manager) PlayThemeSoundSync(theme, event string) error {
	return soundplayer.PlayThemeSound(theme, event, "", true)
}

func (m *Manager) PlaySystemSoundWithDevice(event, device string) error {
	return soundplayer.PlaySystemSound(event, device, false)
}

func (m *Manager) PlaySystemSoundWithDeviceSync(event, device string) error {
	return soundplayer.PlaySystemSound(event, device, true)
}

func (m *Manager) PlayThemeSoundWithDevice(theme, event, device string) error {
	return soundplayer.PlayThemeSound(theme, event, device, false)
}

func (m *Manager) PlayThemeSoundWithDeviceSync(theme, event,
	device string) error {
	return soundplayer.PlayThemeSound(theme, event, device, true)
}

func (m *Manager) PlaySoundFile(file string) error {
	return soundplayer.PlayThemeSoundFile(file, "", false)
}

func (m *Manager) PlaySoundFileSync(file string) error {
	return soundplayer.PlayThemeSoundFile(file, "", true)
}

func (m *Manager) PlaySoundFileWithDevice(file, device string) error {
	return soundplayer.PlayThemeSoundFile(file, device, false)
}

func (m *Manager) PlaySoundFileWithDeviceSync(file, device string) error {
	return soundplayer.PlayThemeSoundFile(file, device, true)
}
