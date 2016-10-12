/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

type ItemInfo struct {
	Path          string
	Name          string // display name
	ID            string
	Icon          string
	CategoryID    CategoryID
	TimeInstalled int64
}

func (item *Item) newItemInfo() ItemInfo {
	iInfo := ItemInfo{
		Path:          item.Path,
		Name:          item.Name,
		ID:            item.ID,
		Icon:          item.Icon,
		CategoryID:    item.CategoryID,
		TimeInstalled: item.TimeInstalled,
	}
	return iInfo
}
