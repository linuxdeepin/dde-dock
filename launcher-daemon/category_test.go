package main

import (
	"dlib/gio-2.0"
	"sort"
	"testing"
)

func TestInitCategoryTable(t *testing.T) {
	initCategory()
}

func TestGetCategoryInfos(t *testing.T) {
	infos := getCategoryInfos()
	if !(infos[len(infos)-1].Id == -2 &&
		sort.IsSorted(CategoryInfosResult(infos[1:len(infos)-1]))) {
		t.Error("Not Sorted Correctly", infos)
	}
}

func TestGetDeepinCategory(t *testing.T) {
	a := gio.NewDesktopAppInfo("deepin-game-center.desktop")
	if a == nil {
		t.Error("cannot create deepin-game-center.desktop")
	}
	defer a.Unref()

	id, err := getDeepinCategory(a)
	if err != nil {
		t.Error("get category id of deepin-game-center.desktop failed", err)
	}
	if id != GamesID {
		t.Error("deepin-game-center.desktop category id:", id)
	}
}
