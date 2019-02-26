package background

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

func sumFileMd5(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func updateModTime(file string) {
	now := time.Now()
	err := os.Chtimes(file, now, now)
	if err != nil {
		logger.Warning("failed to update cache file modify time:", err)
	}
}

func prepare(image string) (string, error) {
	// image is not uri
	logger.Debug("prepare", image)
	if strings.HasPrefix(image, customWallpapersCacheDir) {
		updateModTime(image)
		return image, nil
	}

	md5sum, err := sumFileMd5(image)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(image)
	cacheFileBaseName := md5sum + ext
	cacheFile := filepath.Join(customWallpapersCacheDir, cacheFileBaseName)
	_, err = os.Stat(cacheFile)
	if err != nil {
		// copy image to cacheFile
		os.MkdirAll(customWallpapersCacheDir, 0755)
		err = dutils.CopyFile(image, cacheFile)
		if err != nil {
			return "", err
		}

		time.AfterFunc(time.Second, func() {
			shrinkCache(cacheFileBaseName)
		})
	} else {
		updateModTime(cacheFile)
	}

	return cacheFile, nil
}

func shrinkCache(cacheFileBaseName string) {
	gs := gio.NewSettings("com.deepin.dde.appearance")
	defer gs.Unref()

	workspaceBackgrounds := gs.GetStrv("background-uris")
	var notDeleteFiles strv.Strv
	notDeleteFiles = append(notDeleteFiles, cacheFileBaseName)
	for _, uri := range workspaceBackgrounds {
		wbFile := dutils.DecodeURI(uri)
		if strings.HasPrefix(wbFile, customWallpapersCacheDir) {
			// is custom wallpaper
			basename := filepath.Base(wbFile)
			if basename != cacheFileBaseName {
				notDeleteFiles = append(notDeleteFiles, basename)
			}
		}
	}
	deleteOld(notDeleteFiles)
}

func deleteOld(notDeleteFiles strv.Strv) {
	fileInfos, _ := ioutil.ReadDir(customWallpapersCacheDir)
	count := len(fileInfos) - customWallpapersLimit
	if count <= 0 {
		return
	}
	logger.Debugf("need delete %d file(s)", count)

	sort.Sort(byModTime(fileInfos))
	for _, fileInfo := range fileInfos {
		if count == 0 {
			break
		}

		// traverse from old to new
		fileBaseName := fileInfo.Name()
		if !notDeleteFiles.Contains(fileBaseName) {
			logger.Debug("delete", fileBaseName)
			fullPath := filepath.Join(customWallpapersCacheDir, fileBaseName)
			err := os.Remove(fullPath)
			if os.IsNotExist(err) {
				err = nil
			}

			if err == nil {
				count--

				if customWallpaperDeleteCallback != nil {
					customWallpaperDeleteCallback(fullPath)
				}
			} else {
				logger.Warning(err)
			}
		}
	}
}

type byModTime []os.FileInfo

func (a byModTime) Len() int      { return len(a) }
func (a byModTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byModTime) Less(i, j int) bool {
	return a[i].ModTime().Unix() < a[j].ModTime().Unix()
}
