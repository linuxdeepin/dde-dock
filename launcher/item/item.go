/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package item

import (
	"os"
	"path"
	"strings"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/appinfo"
	"pkg.deepin.io/dde/daemon/launcher/category"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/utils"
)

const (
	_DesktopSuffix    = ".desktop"
	_DesktopSuffixLen = len(_DesktopSuffix)
)

// #define FILENAME_WEIGHT 0.3
// #define GENERIC_NAME_WEIGHT 0.01
// #define KEYWORD_WEIGHT 0.1
// #define CATEGORY_WEIGHT 0.01
// #define NAME_WEIGHT 0.01
// #define DISPLAY_NAME_WEIGHT 0.1
// #define DESCRIPTION_WEIGHT 0.01
// #define EXECUTABLE_WEIGHT 0.05

// Xinfo stores some information in desktop.
type Xinfo struct {
	keywords    []string
	exec        string
	genericName string
	description string
}

// Info stores some information for app.
type Info struct {
	path          string
	name          string
	enName        string
	id            ItemID
	icon          string
	categoryID    CategoryID
	timeInstalled int64
	xinfo         Xinfo
}

// Path returns desktop's path.
func (i *Info) Path() string {
	return i.path
}

// Name returns app's english name.
func (i *Info) Name() string {
	return i.enName
}

// LocaleName returns app's locale name.
func (i *Info) LocaleName() string {
	return i.name
}

// ID returns appid.
func (i *Info) ID() ItemID {
	return i.id
}

// Keywords returns keywords for searching.
func (i *Info) Keywords() []string {
	return i.xinfo.keywords
}

// GenericName returns generic name in desktop file.
func (i *Info) GenericName() string {
	return i.xinfo.genericName
}

// New creates a new Info object from GDesktopAppInfo.
func New(app *gio.DesktopAppInfo) *Info {
	if app == nil {
		return nil
	}
	item := &Info{
		id: getID(app),
	}
	item.init(app)
	return item
}

func (i *Info) Refresh() {
	app := gio.NewDesktopAppInfo(i.path)
	if app != nil {
		i.init(app)
		app.Unref()
	}
}

func (i *Info) init(app *gio.DesktopAppInfo) {
	i.path = app.GetFilename()
	i.name = app.GetDisplayName()
	i.enName = app.GetString("Name")
	icon := app.GetIcon()
	if icon != nil {
		i.icon = icon.ToString()
		if path.IsAbs(i.icon) && !utils.IsFileExist(i.icon) {
			i.icon = ""
		}
	}

	i.xinfo.keywords = make([]string, 0)
	keywords := app.GetKeywords()
	for _, keyword := range keywords {
		i.xinfo.keywords = append(i.xinfo.keywords, strings.ToLower(keyword))
	}
	i.xinfo.exec = app.GetCommandline()
	i.xinfo.genericName = app.GetGenericName()
	i.xinfo.description = app.GetDescription()
	i.categoryID = category.OthersID
}

// Description returns the description storing in desktop file.
func (i *Info) Description() string {
	return i.xinfo.description
}

// ExecCmd returns the exec stroing in desktop file.
func (i *Info) ExecCmd() string {
	return i.xinfo.exec
}

// Icon returns the app's icon.
func (i *Info) Icon() string {
	return i.icon
}

// CategoryID returns category id in deepin store.
func (i *Info) CategoryID() CategoryID {
	return i.categoryID
}

// SetCategoryID sets the category id in deepin store.
func (i *Info) SetCategoryID(id CategoryID) {
	i.categoryID = id
}

// TimeInstalled returns the time installed.
func (i *Info) TimeInstalled() int64 {
	return i.timeInstalled
}

// SetTimeInstalled sets the time installed.
func (i *Info) SetTimeInstalled(timeInstalled int64) {
	i.timeInstalled = timeInstalled
}

func (i *Info) LastModifiedTime() int64 {
	stat, e := os.Stat(i.path)
	if e != nil {
		return 0
	}
	return stat.ModTime().UnixNano()
}

// GenID returns a valid item id.
func GenID(filename string) ItemID {
	if len(filename) <= _DesktopSuffixLen {
		return ItemID("")
	}

	basename := path.Base(filename)
	return ItemID(appinfo.NormalizeAppID(strings.TrimSuffix(basename, _DesktopSuffix)))
}

func getID(app *gio.DesktopAppInfo) ItemID {
	id := app.GetId()
	return ItemID(appinfo.NormalizeAppID(strings.TrimSuffix(id, _DesktopSuffix)))
}
