package launcher

import (
	// "fmt"
	"os"
	p "path"

	"pkg.linuxdeepin.com/lib/glib-2.0"
)

func isOnDesktop(name string) bool {
	path := p.Join(glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop), p.Base(name))
	// logger.Info(path)
	return exist(path)
}

func sendToDesktop(name string) {
	path := p.Join(glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop), p.Base(name))
	// logger.Info(path)
	copyFile(name, path,
		CopyFileNotKeepSymlink|CopyFileOverWrite)
	s, _ := os.Stat(path)
	var execPerm os.FileMode = 0100
	os.Chmod(path, s.Mode().Perm()|execPerm)
}
