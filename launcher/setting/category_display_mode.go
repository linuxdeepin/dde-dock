/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package setting

// CategoryDisplayMode type for category display mode.
type CategoryDisplayMode int64

// category display
const (
	CategoryDisplayModeUnknown CategoryDisplayMode = iota - 1
	CategoryDisplayModeIcon
	CategoryDisplayModeText

	CategoryDisplayModeKey string = "category-display-mode"
)

func (c CategoryDisplayMode) String() string {
	switch c {
	case CategoryDisplayModeUnknown:
		return "unknown category display mode"
	case CategoryDisplayModeText:
		return "display text mode"
	case CategoryDisplayModeIcon:
		return "display icon mode"
	default:
		return "unknown mode"
	}
}
