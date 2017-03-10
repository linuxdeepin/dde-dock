/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

type WarnLevel uint32

const (
	WarnLevelNone WarnLevel = iota
	WarnLevelLow
	WarnLevelCritical
	WarnLevelAction
)

func (lv WarnLevel) String() string {
	switch lv {
	case WarnLevelNone:
		return "None"
	case WarnLevelLow:
		return "Low"
	case WarnLevelCritical:
		return "Critical"
	case WarnLevelAction:
		return "Action"
	default:
		return "Unknown"
	}
}

func _getWarnLevel(config *WarnLevelConfig, onBattery bool,
	percentage float64, timeToEmpty uint64) WarnLevel {

	if !onBattery {
		return WarnLevelNone
	}

	usePercentageForPolicy := config.UsePercentageForPolicy
	logger.Debugf("_getWarnLevel onBattery %v, percentage %v, timeToEmpty %v, usePercentage %v",
		onBattery, percentage, timeToEmpty, usePercentageForPolicy)
	if usePercentageForPolicy {
		if percentage > config.LowPercentage || percentage == 0.0 {
			return WarnLevelNone
		}
		if percentage > config.CriticalPercentage {
			return WarnLevelLow
		}
		if percentage > config.ActionPercentage {
			return WarnLevelCritical
		}
		return WarnLevelAction
	} else {
		if timeToEmpty > config.LowTime || timeToEmpty == 0 {
			return WarnLevelNone
		}
		if timeToEmpty > config.CriticalTime {
			return WarnLevelLow
		}
		if timeToEmpty > config.ActionTime {
			return WarnLevelCritical
		}
		return WarnLevelAction
	}
}

func (m *Manager) getWarnLevel(percentage float64, timeToEmpty uint64) WarnLevel {
	return _getWarnLevel(m.warnLevelConfig, m.OnBattery, percentage, timeToEmpty)
}
