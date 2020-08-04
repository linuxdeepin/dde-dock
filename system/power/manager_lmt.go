package power

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

// laptop mode tools config file
const lmtConfigFile = "/etc/laptop-mode/laptop-mode.conf"
const laptopModeBin = "/usr/sbin/laptop_mode"

const (
	lmtConfigAuto     = 1
	lmtConfigEnabled  = 2
	lmtConfigDisabled = 3
)

const lowBatteryThreshold = 20.0

func isLaptopModeBinOk() bool {
	_, err := os.Stat(laptopModeBin)
	return err == nil
}

func setLMTConfig(mode int) (changed bool, err error) {
	lines, err := loadLmtConfig()
	if err != nil {
		// ignore not exist error
		if os.IsNotExist(err) {
			return false, nil
		}
		logger.Warning(err)
		return false, err
	}

	dict := make(map[string]string)
	switch mode {
	case lmtConfigAuto:
		dict["ENABLE_LAPTOP_MODE_TOOLS"] = "1"
		dict["ENABLE_LAPTOP_MODE_ON_BATTERY"] = "1"
		dict["ENABLE_LAPTOP_MODE_ON_AC"] = "0"
	case lmtConfigEnabled:
		dict["ENABLE_LAPTOP_MODE_TOOLS"] = "1"
		dict["ENABLE_LAPTOP_MODE_ON_BATTERY"] = "1"
		dict["ENABLE_LAPTOP_MODE_ON_AC"] = "1"
	case lmtConfigDisabled:
		dict["ENABLE_LAPTOP_MODE_TOOLS"] = "1"
		dict["ENABLE_LAPTOP_MODE_ON_BATTERY"] = "0"
		dict["ENABLE_LAPTOP_MODE_ON_AC"] = "0"
	}
	lines, changed = modifyLMTConfig(lines, dict)
	if changed {
		logger.Debug("write LMT Config")
		err = writeLmtConfig(lines)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

func reloadLaptopModeService() error {
	if !isLaptopModeBinOk() {
		logger.Debug("laptop mode tools is not installed")
		return nil
	}

	systemBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	systemdObj := systemBus.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	return systemdObj.Call("org.freedesktop.systemd1.Manager.ReloadUnit",
		dbus.FlagNoAutoStart, "laptop-mode.service", "replace").Err
}

func modifyLMTConfig(lines []string, dict map[string]string) ([]string, bool) {
	var changed bool
	for idx := range lines {
		line := lines[idx]
		for key, value := range dict {
			if strings.HasPrefix(line, key) {
				newLine := key + "=" + value
				if line != newLine {
					changed = true
					lines[idx] = newLine
				}
				delete(dict, key)
			}
		}
		if len(dict) == 0 {
			break
		}
	}
	if len(dict) > 0 {
		for key, value := range dict {
			newLine := key + "=" + value
			lines = append(lines, newLine)
		}
		changed = true
	}
	return lines, changed
}

func loadLmtConfig() ([]string, error) {
	f, err := os.Open(lmtConfigFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

func writeLmtConfig(lines []string) error {
	tempFile, err := writeLmtConfigTemp(lines)
	if err != nil {
		if tempFile != "" {
			os.Remove(tempFile)
		}
		return err
	}
	return os.Rename(tempFile, lmtConfigFile)
}

func writeLmtConfigTemp(lines []string) (string, error) {
	dir := filepath.Dir(lmtConfigFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(dir, "laptop-mode.conf")
	logger.Debug("writeLmtConfig temp file", f.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = f.Chmod(0644)
	if err != nil {
		return f.Name(), err
	}

	bufWriter := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := bufWriter.WriteString(line)
		if err != nil {
			logger.Warning(err)
		}
		err = bufWriter.WriteByte('\n')
		if err != nil {
			logger.Warning(err)
		}
	}
	return f.Name(), bufWriter.Flush()
}

func (m *Manager) writePowerSavingModeEnabledCb(write *dbusutil.PropertyWrite) *dbus.Error {
	logger.Debug("set laptop mode enabled", write.Value)

	enabled := write.Value.(bool)
	var err error
	var lmtCfgChanged bool

	m.PropsMu.Lock()
	m.setPropPowerSavingModeAuto(false)
	m.setPropPowerSavingModeAutoWhenBatteryLow(false)
	m.PropsMu.Unlock()

	if enabled {
		lmtCfgChanged, err = setLMTConfig(lmtConfigEnabled)
	} else {
		lmtCfgChanged, err = setLMTConfig(lmtConfigDisabled)
	}

	if err != nil {
		logger.Warning("failed to set LMT config:", err)
	}

	if lmtCfgChanged {
		err := reloadLaptopModeService()
		if err != nil {
			logger.Warning(err)
		}
	}

	return nil
}

func (m *Manager) updatePowerSavingMode() { // 根据用户设置以及当前状态,修改节能模式
	if !m.initDone {
		// 初始化未完成时，暂不提供功能
		return
	}
	var enable bool
	var lmtCfgChanged bool
	var err error
	if m.PowerSavingModeAuto && m.PowerSavingModeAutoWhenBatteryLow {
		if m.OnBattery || m.batteryLow {
			enable = true
		} else {
			enable = false
		}
	} else if m.PowerSavingModeAuto && !m.PowerSavingModeAutoWhenBatteryLow {
		if m.OnBattery {
			enable = true
		} else {
			enable = false
		}
	} else if !m.PowerSavingModeAuto && m.PowerSavingModeAutoWhenBatteryLow {
		if m.batteryLow {
			enable = true
		} else {
			enable = false
		}
	} else {
		return // 未开启两个自动节能开关
	}
	logger.Info("updatePowerSavingMode PowerSavingModeEnabled: ", enable)
	m.PropsMu.Lock()
	changed := m.setPropPowerSavingModeEnabled(enable)
	m.PropsMu.Unlock()
	if changed {
		if enable {
			lmtCfgChanged, err = setLMTConfig(lmtConfigEnabled)
			if err != nil {
				logger.Warning("failed to set LMT config:", err)
			}
		} else {
			lmtCfgChanged, err = setLMTConfig(lmtConfigDisabled)
			if err != nil {
				logger.Warning("failed to set LMT config:", err)
			}
		}
		if lmtCfgChanged {
			err := reloadLaptopModeService()
			if err != nil {
				logger.Warning(err)
			}
		}
	}
}
