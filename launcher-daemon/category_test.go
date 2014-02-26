package main

import (
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
