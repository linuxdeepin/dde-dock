/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

type HideStateType int32

const (
	HideStateShowing HideStateType = iota
	HideStateShown
	HideStateHidding
	HideStateHidden
)

func (s HideStateType) String() string {
	switch s {
	case HideStateShowing:
		return "HideStateShowing"
	case HideStateShown:
		return "HideStateShown"
	case HideStateHidding:
		return "HideStateHidding"
	case HideStateHidden:
		return "HideStateHidden"
	default:
		return "Unknown state"
	}
}
