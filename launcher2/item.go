/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher2

import (
	"dbus/com/deepin/api/pinyin"
	"errors"
	"fmt"
	"gir/gio-2.0"
	"regexp"
	"strings"
)

type SearchScore uint64

type Item struct {
	Path          string
	Name          string // display name
	enName        string
	ID            string
	Icon          string
	CategoryID    CategoryID
	TimeInstalled int64

	keywords      []string
	categories    []string
	exec          string
	genericName   string
	comment       string
	searchTargets map[string]SearchScore
}

func (item *Item) String() string {
	if item == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<item %v>", item.ID)
}

const (
	desktopExt           = ".desktop"
	desktopCategroyDelim = ";"
)

func getAppId(desktopId string) string {
	return strings.TrimSuffix(desktopId, desktopExt)
}

func NewItemWithDesktopAppInfo(app *gio.DesktopAppInfo) *Item {
	if app == nil {
		return nil
	}
	item := &Item{
		Path:          app.GetFilename(),
		Name:          app.GetDisplayName(),
		enName:        app.GetString("Name"),
		exec:          app.GetCommandline(),
		genericName:   app.GetString("GenericName"),
		comment:       app.GetString("Comment"),
		searchTargets: make(map[string]SearchScore),
	}
	icon := app.GetIcon()
	if icon != nil {
		item.Icon = icon.ToString()
	}
	for _, kw := range app.GetKeywords() {
		item.keywords = append(item.keywords, strings.ToLower(kw))
	}

	strCategories := app.GetCategories()
	categories := strings.Split(strings.TrimSuffix(strCategories, desktopCategroyDelim), desktopCategroyDelim)
	for _, c := range categories {
		item.categories = append(item.categories, strings.ToLower(c))
	}
	return item
}

var chromeShortcurtExecRegexp = regexp.MustCompile(`google-chrome.*--app-id=`)

func (item *Item) isChromeShortcut() bool {
	logger.Debugf("isChromeShortcut item %#v", item)
	return strings.HasPrefix(item.ID, "chrome-") &&
		chromeShortcurtExecRegexp.MatchString(item.exec)
}

func (item *Item) isWineApp() (bool, error) {
	appInfo := gio.NewDesktopAppInfoFromFilename(item.Path)
	if appInfo == nil {
		return false, errors.New("appInfo is nil")
	}
	defer appInfo.Unref()
	return strings.HasPrefix(appInfo.GetString("X-Created-By"), "cxoffice-") ||
		strings.Contains(appInfo.GetCommandline(), "env WINEPREFIX="), nil
}

func (item *Item) getXCategory() CategoryID {
	logger.Debug("getXCategory item.categories:", item.categories)
	categoriesCountMap := make(map[CategoryID]int)
	if len(item.categories) == 1 {
		return parseXCategoryString(item.categories[0])
	}

	for _, categoryStr := range item.categories {
		cid := parseXCategoryString(categoryStr)
		categoriesCountMap[cid] = categoriesCountMap[cid] + 1
	}

	// ignore CategoryOthers
	delete(categoriesCountMap, CategoryOthers)
	logger.Debug("getXCategory categoriesCountMap:", categoriesCountMap)

	if len(categoriesCountMap) > 0 {
		var categoryCountMax int
		categoryMax := CategoryOthers
		for cid, count := range categoriesCountMap {
			if count > categoryCountMax {
				categoryCountMax = count
				categoryMax = cid
			}
		}
		logger.Debugf("category max %v count %v", categoryMax, categoryCountMax)
		return categoryMax
	}
	return CategoryOthers
}

const (
	idScore          = 100
	nameScore        = 80
	genericNameScore = 70
	keywordScore     = 60
	categoryScore    = 60
	commentScore     = 50
)

func (item *Item) setSearchTargets(pinyin *pinyin.Pinyin) {
	item.addSearchTarget(idScore, item.ID)
	item.addSearchTarget(nameScore, item.Name)
	item.addSearchTarget(nameScore, item.enName)
	item.addSearchTarget(genericNameScore, item.genericName)
	item.addSearchTarget(categoryScore, item.CategoryID.String())
	item.addSearchTarget(commentScore, item.comment)
	for _, c := range item.categories {
		item.addSearchTarget(categoryScore, c)
	}
	if pinyin != nil {
		pinyinList, err := pinyin.Query(item.Name)
		if err == nil {
			for _, v := range pinyinList {
				item.addSearchTarget(nameScore, v)
			}
		} else {
			logger.Warning(err)
		}

		item.addSearchTarget(categoryScore, item.CategoryID.Pinyin())
	}

	// add keywords
	for _, kw := range item.keywords {
		item.addSearchTarget(keywordScore, kw)
	}
}

func (item *Item) addSearchTarget(score SearchScore, str string) {
	if str == "" {
		return
	}
	str = strings.ToLower(str)
	scoreInDict, ok := item.searchTargets[str]
	if !ok || (ok && scoreInDict < score) {
		item.searchTargets[str] = score
	}
}
