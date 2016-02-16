/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package setting

// SortMethod type for sort method.
type SortMethod int64

// sort method.
const (
	SortMethodUnknown SortMethod = iota - 1
	SortMethodByName
	SortMethodByCategory
	SortMethodByTimeInstalled
	SortMethodByFrequency

	SortMethodkey string = "sort-method"
)

func (s SortMethod) String() string {
	switch s {
	case SortMethodUnknown:
		return "unknown sort method"
	case SortMethodByName:
		return "sort by name"
	case SortMethodByCategory:
		return "sort by category"
	case SortMethodByTimeInstalled:
		return "sort by time installed"
	case SortMethodByFrequency:
		return "sort by frequency"
	default:
		return "unknown sort method"
	}
}
