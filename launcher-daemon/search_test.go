package main

import (
	"os"
	"testing"

	pinyin "dbus/com/deepin/api/search"
	"dlib/gio-2.0"
)

func _TestGetMatchers(t *testing.T) {
	for k, v := range getMatchers("chrome") {
		t.Log(k.String(), v)
	}
}

func TestSearchContent(t *testing.T) {
	r := search("fi")
	// r := search("chome")
	for _, id := range r {
		item := itemTable[id]
		t.Logf(`id: %s
		Name: %s
		Path: %s
		keywords: %v
		GenericName: %s
		Description: %s
		Exec: %s
		`, id, item.Name, item.Path, item.xinfo.keywords,
			item.xinfo.genericName, item.xinfo.description,
			item.xinfo.exec)
	}
}

func _TestSearch(t *testing.T) {
	r := search("chrome")
	j := search("chrome")
	if len(r) != len(j) {
		t.Error("not equal: get different length.")
	}
	title := false
	for i := 0; i < len(r); i++ {
		if itemTable[r[i]].Id != itemTable[j[i]].Id {
			if !title {
				t.Error("not equal: get different sequence.")
				title = true
			}
			item := itemTable[r[i]]
			t.Errorf("\tindex: %d, the first search: Id: %s, Name: %s", i, r[i], item.Name)
			item = itemTable[j[i]]
			t.Errorf("\tindex: %d, the second search: Id: %s, Name: %s", i, j[i], item.Name)
		}
	}
}

func _TestPinYin(t *testing.T) {
	tree, err := pinyin.NewSearch("/com/deepin/dde/api/Search")
	if err != nil {
		return
	}
	names := make(map[string]string, 0)
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")
	addName := func(m map[string]string, n string) {
		app := gio.NewDesktopAppInfo(n)
		defer app.Unref()
		name := app.GetDisplayName()
		t.Log(name)
		m[name] = name
	}
	addName(names, "deepin-software-center.desktop")
	addName(names, "firefox.desktop")
	t.Log(names)
	treeId, _ := tree.NewTrieWithString(names, "DDELauncherDaemonTest")
	var keys []string
	keys, _ = tree.SearchKeys("ruan", treeId)
	t.Log(keys)
	keys, _ = tree.SearchKeys("firefox", treeId)
	t.Log(keys)
	keys, _ = tree.SearchKeys("wang", treeId)
	t.Log(keys)
	keys, _ = tree.SearchKeys("网络", treeId)
	t.Log(keys)
	tree.DestroyTrie(treeId)
}
