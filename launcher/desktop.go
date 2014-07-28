package launcher

import (
	// "fmt"
	"os"
	p "path"

	"pkg.linuxdeepin.com/lib/glib-2.0"
)

func getDesktopPath(name string) string {
	return p.Join(glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop), p.Base(name))
}

func isOnDesktop(name string) bool {
	path := getDesktopPath(name)
	// logger.Info(path)
	return exist(path)
}

func sendToDesktop(name string) error {
	path := getDesktopPath(name)
	// logger.Info(path)
	err := copyFile(name, path,
		CopyFileNotKeepSymlink|CopyFileOverWrite)
	if err != nil {
		return err
	}
	s, err := os.Stat(path)
	if err != nil {
		removeFromDesktop(name)
		return err
	}
	var execPerm os.FileMode = 0100
	os.Chmod(path, s.Mode().Perm()|execPerm)
	return nil
}

func removeFromDesktop(name string) error {
	path := getDesktopPath(name)
	return os.Remove(path)
}
