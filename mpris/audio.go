// handle audio signals
package mpris

import (
	mpris2 "dbus/org/mpris/mediaplayer2"
	"strings"
)

const (
	actionTypePlay int = iota + 1
	actionTypePause
	actionTypeStop
	actionTypePrevious
	actionTypeNext
	actionTypeRewind
	actionTypeForward
	actionTypeRepeat
)

const (
	senderTypeMpris = "org.mpris.MediaPlayer2"
	mprisPath       = "/org/mpris/MediaPlayer2"

	playerDelta int64 = 5000 * 1000 // 5s
)

func (m *Manager) playerAction(action int, pressed bool) {
	if pressed {
		return
	}

	player := m.getActiveMpris()
	if player == nil {
		return
	}

	switch action {
	case actionTypePlay:
		player.PlayPause()
	case actionTypePause:
		player.Pause()
	case actionTypeStop:
		player.Stop()
	case actionTypePrevious:
		player.Previous()
		player.Play()
	case actionTypeNext:
		player.Next()
		player.Play()
	case actionTypeRewind:
		pos := player.Position.Get()
		if pos-playerDelta < 0 {
			pos = 0
		} else {
			pos = 0 - playerDelta
		}

		player.Seek(pos)
		if player.PlaybackStatus.Get() != "Playing" {
			player.PlayPause()
		}
	case actionTypeForward:
		player.Seek(playerDelta)
		if player.PlaybackStatus.Get() != "Playing" {
			player.PlayPause()
		}
	case actionTypeRepeat:
		player.Play()
	}
}

func (m *Manager) pauseAllPlayer(actived bool) {
	if !actived {
		return
	}

	for _, sender := range m.getMprisSender() {
		player, err := mpris2.NewPlayer(sender, mprisPath)
		if err != nil {
			continue
		}
		player.Pause()
	}
}

func (m *Manager) getMprisSender() []string {
	var senders []string
	names, _ := m.dbusDaemon.ListNames()
	for _, name := range names {
		if !strings.Contains(name, senderTypeMpris) {
			continue
		}

		senders = append(senders, name)
	}

	return senders
}

func (m *Manager) getActiveMpris() *mpris2.Player {
	var senders = m.getMprisSender()

	length := len(senders)
	if length == 0 {
		return nil
	}

	for _, sender := range senders {
		player, err := mpris2.NewPlayer(sender, mprisPath)
		if err != nil {
			continue
		}

		if length == 1 {
			return player
		}

		if length == 2 && strings.Contains(sender, "vlc") {
			return player
		}

		if player.PlaybackStatus.Get() == "Playing" {
			m.prevPlayer = sender
			return player
		}

		if m.prevPlayer == sender {
			return player
		}
	}

	player, _ := mpris2.NewPlayer(senders[0], mprisPath)
	return player
}
