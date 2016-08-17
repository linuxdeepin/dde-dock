/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package interfaces

// SearchID is type for pinyin search.
type SearchID string

// Search is interface for search transaction.
type Search interface {
	Search(string, []ItemInfo)
	Cancel()
}

// PinYin is interface for pinyin search transaction.
type PinYin interface {
	Search(string) ([]string, error)
	IsValid() bool
	Update([]string) error
}
