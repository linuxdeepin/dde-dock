/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package apps

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func intSliceContains(slice []int, a int) bool {
	for _, elem := range slice {
		if elem == a {
			return true
		}
	}
	return false
}

func intSliceRemove(slice []int, a int) (ret []int) {
	for _, elem := range slice {
		if elem != a {
			ret = append(ret, elem)
		}
	}
	return
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return names, nil
}

func Walk(root string, walkFn WalkFunc) {
	info, err := os.Lstat(root)
	if err != nil {
		return
	}
	walk(root, ".", info, walkFn)
}

func walk(root, name0 string, info os.FileInfo, walkFn WalkFunc) {
	walkFn(name0, info)
	if !info.IsDir() {
		return
	}
	path := filepath.Join(root, name0)
	names, err := readDirNames(path)
	if err != nil {
		return
	}
	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			continue
		}
		walk(root, filepath.Join(name0, name), fileInfo, walkFn)
	}
}

type WalkFunc func(name string, info os.FileInfo)

func getDirsAndApps(root string) (dirs, apps []string) {
	Walk(root, func(name string, info os.FileInfo) {
		if info.IsDir() {
			dirs = append(dirs, name)
		} else if filepath.Ext(name) == desktopExt {
			apps = append(apps, strings.TrimSuffix(name, desktopExt))
		}
	})
	return
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

const desktopExt = ".desktop"

func isDesktopFile(path string) bool {
	return filepath.Ext(path) == desktopExt
}

func removeDesktopExt(name string) string {
	return strings.TrimSuffix(name, desktopExt)
}

func getSystemDataDirs() []string {
	return []string{"/usr/share", "/usr/local/share"}
}

// get user home
func getHomeByUid(uid int) (string, error) {
	user, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

func getDirPerm(uid int) os.FileMode {
	if uid == 0 {
		// rwx r-x r-x
		return 0755
	}
	// rwx --- ---
	return 0700
}

// copy from go source src/os/path.go
func MkdirAll(path string, uid int, perm os.FileMode) error {
	logger.Debug("MkdirAll", path, uid, perm)
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = MkdirAll(path[0:j-1], uid, perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = os.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}

	err = os.Chown(path, uid, uid)
	return err
}
