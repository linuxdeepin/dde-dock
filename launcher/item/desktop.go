/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package item

import (
	"os"
	p "path"

	"gir/glib-2.0"
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	"pkg.deepin.io/lib/utils"
)

func getDesktopPath(name string) string {
	GReloadUserSpecialDirsCache()
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
