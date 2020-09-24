package audio

import (
	"regexp"
	"strings"

	dbus "github.com/godbus/dbus"
	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/pulse"
)

func isBluezAudio(name string) bool {
	return strings.Contains(strings.ToLower(name), "bluez")
}

func isBluezDeviceValid(bluezPath string) bool {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		logger.Warning("[isDeviceValid] dbus connect failed:", err)
		return false
	}
	bluezDevice, err := bluez.NewDevice(systemBus, dbus.ObjectPath(bluezPath))
	if err != nil {
		logger.Warning("[isDeviceValid] new device failed:", err)
		return false
	}
	icon, err := bluezDevice.Icon().Get(0)
	if err != nil {
		logger.Warning("[isDeviceValid] get icon failed:", err)
		return false
	}
	if icon == "computer" {
		return false
	}
	return true
}

func createBluezVirtualCardPorts(ports pulse.CardPortInfos) pulse.CardPortInfos {
	var virtualPorts = make(pulse.CardPortInfos, 0)
	for _, port := range ports {
		if port.Direction == pulse.DirectionSource {
			if port.Profiles.Exists("headset_head_unit") {
				headsetPort := port
				headsetPort.Name += "(headset_head_unit)"
				headsetPort.Description += "(Headset)"
				if headsetPort.Available == pulse.AvailableTypeNo {
					headsetPort.Available = pulse.AvailableTypeUnknow
				}
				virtualPorts = append(virtualPorts, headsetPort)
				logger.Debug("create virtual bluez port headset")
			}
		} else {
			// 这里的顺序不能改，默认a2dp优先
			// 在优先级模块中，默认后接入的端口优先
			// 因此a2dp放在后面
			if port.Profiles.Exists("headset_head_unit") {
				headsetPort := port
				headsetPort.Name += "(headset_head_unit)"
				headsetPort.Description += "(Headset)"
				virtualPorts = append(virtualPorts, headsetPort)
				logger.Debug("create virtual bluez port headset")
			}

			if port.Profiles.Exists("a2dp_sink") {
				a2dpPort := port
				a2dpPort.Name += "(a2dp_sink)"
				a2dpPort.Description += "(A2DP)"
				virtualPorts = append(virtualPorts, a2dpPort)
				logger.Debug("create virtual bluez port a2dp")
			}

		}
	}

	return virtualPorts
}

func createBluezVirtualSinkPorts(ports []Port) []Port {
	var virtualPorts = make([]Port, 0)
	for _, port := range ports {
		headsetPort := port
		headsetPort.Name += "(headset_head_unit)"
		headsetPort.Description += "(Headset)"
		virtualPorts = append(virtualPorts, headsetPort)
		a2dpPort := port
		a2dpPort.Name += "(a2dp_sink)"
		a2dpPort.Description += "(A2DP)"
		virtualPorts = append(virtualPorts, a2dpPort)
	}
	return virtualPorts
}

func createBluezVirtualSourcePorts(ports []Port) []Port {
	var virtualPorts = make([]Port, 0)
	for _, port := range ports {
		headsetPort := port
		headsetPort.Name += "(headset_head_unit)"
		headsetPort.Description += "(Headset)"
		virtualPorts = append(virtualPorts, headsetPort)
	}
	return virtualPorts
}

func bluezAudioParseVirtualPort(virtualPortName string) (string, string) {
	r, err := regexp.Compile(`(.*?)\((.*?)\)`)
	if err != nil {
		logger.Warning(err)
		return "", ""
	}

	result := r.FindStringSubmatch(virtualPortName)
	if len(result) != 3 {
		logger.Warningf("cannot understand bluez virtual port %s", virtualPortName)
		return "", ""
	}

	port := result[1]
	profile := result[2]
	logger.Debugf("bluez port %s with profile %s", port, profile)

	return port, profile
}

func bluezAudioGetSinkProfile(s *Sink) string {
	a := s.audio
	card, err := a.cards.get(s.Card)
	if err != nil {
		logger.Warning(err)
		return ""
	}
	return card.ActiveProfile.Name
}
