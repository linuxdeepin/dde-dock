package item

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"
import (
	"os"
	p "path"

	. "pkg.linuxdeepin.com/dde-daemon/launcher/utils"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"pkg.linuxdeepin.com/lib/utils"
)

func getDesktopPath(name string) string {
	C.g_reload_user_special_dirs_cache()
	return p.Join(glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop), p.Base(name))
}

func isOnDesktop(name string) bool {
	path := getDesktopPath(name)
	return utils.IsFileExist(path)
}

func sendToDesktop(itemPath string) error {
	path := getDesktopPath(itemPath)
	err := CopyFile(itemPath, path,
		CopyFileNotKeepSymlink|CopyFileOverWrite)
	if err != nil {
		return err
	}
	s, err := os.Stat(path)
	if err != nil {
		removeFromDesktop(itemPath)
		return err
	}
	var execPerm os.FileMode = 0100
	os.Chmod(path, s.Mode().Perm()|execPerm)
	return nil
}

func removeFromDesktop(itemPath string) error {
	path := getDesktopPath(itemPath)
	return os.Remove(path)
}
