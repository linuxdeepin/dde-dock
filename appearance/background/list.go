package background

import (
	"os"
	"path"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
	"sort"
	"strings"
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
		walls = append(walls, tmp)
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
