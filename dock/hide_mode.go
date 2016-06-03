/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

type HideModeType int32

const (
	HideModeKeepShowing HideModeType = iota
	HideModeKeepHidden
	HideModeAutoHide // invalid
	HideModeSmartHide
)

func validHideModeNum(mode int32) bool {
	if mode != 2 && (0 <= mode && mode <= 3) {
		return true
	}
	return false
}

func (t HideModeType) String() string {
	switch t {
	case HideModeKeepShowing:
		return "Keep showing mode"
	case HideModeKeepHidden:
		return "Keep hidden mode"
	case HideModeAutoHide:
		return "Auto hide mode"
	case HideModeSmartHide:
		return "Smart hide mode"
	default:
		return "Unknown mode"
	}
}
