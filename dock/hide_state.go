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
	HideStateUnknown HideStateType = iota
	HideStateShow
	HideStateHide
)

func (s HideStateType) String() string {
	switch s {
	case HideStateShow:
		return "Show"
	case HideStateHide:
		return "Hide"
	default:
		return "Unknown"
	}
}
