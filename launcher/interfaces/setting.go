/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package interfaces

// SettingCore is interface for setting.
type SettingCore interface {
	GetEnum(string) int32
	SetEnum(string, int32) bool
	Connect(string, interface{})
	Unref()
}

// Setting is the interface for setting.
type Setting interface {
	GetCategoryDisplayMode() int64
	SetCategoryDisplayMode(newMode int64)
	GetSortMethod() int64
	SetSortMethod(newMethod int64)
	Destroy()
}
