package launcher

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

func testGetDeepinCategory(t *testing.T, name string, id CategoryId) {
	a := gio.NewDesktopAppInfo(name)
	if a == nil {
		t.Error("cannot create", name)
	}
	defer a.Unref()

	_id, err := getDeepinCategory(a)
	if err != nil {
		t.Error("get category id of", name, err)
	}
	if _id != id {
		t.Error(name, "category id:", _id)
	}
}
func TestGetDeepinCategory(t *testing.T) {
	testGetDeepinCategory(t, "deepin-game-center.desktop", GamesID)
	testGetDeepinCategory(t, "firefox.desktop", NetworkID)
}
