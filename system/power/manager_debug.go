// +build debug

package power

import (
	"gir/gudev-1.0"
)

func (m *Manager) Debug(cmd string) {
	logger.Debug("Debug", cmd)
	switch cmd {
	case "init-batteries":
		devices := m.gudevClient.QueryBySubsystem("power_supply")
		for _, dev := range devices {
			m.addBattery(dev)
		}
		logger.Debug("initBatteries done")
		for _, dev := range devices {
			dev.Unref()
		}

	case "remove-all-batteries":
		var devices []*gudev.Device
		for _, bat := range m.batteries {
			devices = append(devices, bat.newDevice())
		}

		for _, dev := range devices {
			m.removeBattery(dev)
			dev.Unref()
		}

	case "destroy":
		m.destroy()

	default:
		logger.Warning("Command not support")
	}
}
