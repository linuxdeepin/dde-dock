package background

import (
	"os"
	"path"
	"sort"
	"strings"

	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	dthemeDir = "personalization/themes"
)

// ListDirs list all background dirs
func ListDirs() []string {
	var dirs = []string{
		"/usr/share/backgrounds",
		path.Join(glib.GetUserSpecialDir(
			glib.UserDirectoryDirectoryPictures),
			"Wallpapers"),
	}

	dirs = append(dirs, getDirsFromDTheme(path.Join("/usr/share",
		dthemeDir))...)
	dirs = append(dirs, getDirsFromDTheme(path.Join(os.Getenv("HOME"),
		dthemeDir))...)
	return dirs
}

func getBgFiles() []string {
	var walls []string
	for _, dir := range ListDirs() {
		walls = append(walls, scanner(dir)...)
	}
	return walls
}

func isDeletable(file string) bool {
	if strings.Contains(file, os.Getenv("HOME")) {
		return true
	}
	return false
}

func scanner(dir string) []string {
	fr, err := os.Open(dir)
	if err != nil {
		return []string{}
	}
	defer fr.Close()

	names, err := fr.Readdirnames(0)
	if err != nil {
		return []string{}
	}

	var walls []string
	for _, name := range names {
		tmp := path.Join(dir, name)
		if !IsBackgroundFile(tmp) {
			continue
		}
		// TODO: 1. if a file link to two files in current dir; 2. link file is relative path
		if dutils.IsSymlink(tmp) {
			walls = delBgFromList(readLink(tmp), walls)
		}
		walls = addBgToList(tmp, walls)
	}
	return walls
}

// dir: ex '/usr/share/personalization/themes
func getDirsFromDTheme(dir string) []string {
	fr, err := os.Open(dir)
	if err != nil {
		return []string{}
	}
	defer fr.Close()

	names, err := fr.Readdirnames(0)
	if err != nil {
		return []string{}
	}

	var dirs []string
	for _, name := range names {
		tmp := path.Join(dir, name)
		if !dutils.IsDir(tmp) {
			continue
		}

		wall := path.Join(tmp, "wallpapers")
		if !dutils.IsDir(wall) {
			continue
		}
		dirs = append(dirs, wall)
	}

	sort.Strings(dirs)
	return dirs
}

func addBgToList(bg string, list []string) []string {
	for _, v := range list {
		if (v == bg) || (dutils.IsSymlink(v) && readLink(v) == bg) {
			return list
		}
	}
	list = append(list, bg)
	return list
}

func delBgFromList(bg string, list []string) []string {
	var ret []string
	for _, v := range list {
		if (v == bg) || (dutils.IsSymlink(v) && readLink(v) == bg) {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func readLink(file string) string {
	for {
		file, _ = os.Readlink(file)
		if !dutils.IsSymlink(file) {
			break
		}
	}
	return file
}
