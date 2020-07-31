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

	"github.com/nfnt/resize"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/imgutil"
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
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func updateModTime(file string) {
	now := time.Now()
	err := os.Chtimes(file, now, now)
	if err != nil {
		logger.Warning("failed to update cache file modify time:", err)
	}
}

func resizeImage(filename, cacheDir string) (outFilename, ext string, isResized bool) {
	img, err := imgutil.Load(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load image file %q: %v\n", filename, err)
		outFilename = filename
		return
	}

	const (
		stdWidth  = 3840
		stdHeight = 2400
	)

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	if imgWidth <= stdWidth && imgHeight <= stdHeight {
		// no need to resize
		outFilename = filename
		return
	}

	ext = "jpg"
	format := graphic.FormatJpeg
	_, err = os.Stat("/usr/share/wallpapers/deepin/desktop.bmp")
	if err == nil {
		ext = "bmp"
		format = graphic.FormatBmp
	}

	fh, err := ioutil.TempFile(cacheDir, "tmp-")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to create temp file:", err)
		outFilename = filename
		return
	}

	// tmp-###
	outFilename = fh.Name()
	err = fh.Close()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to close temp file:", err)
	}

	if float64(imgWidth)/float64(imgHeight) > float64(stdWidth)/float64(stdHeight) {
		// use std height
		imgWidth = 0
		imgHeight = stdHeight
	} else {
		// use std width
		imgWidth = stdWidth
		imgHeight = 0
	}

	img = resize.Resize(uint(imgWidth), uint(imgHeight), img, resize.Lanczos3)
	err = graphic.SaveImage(outFilename, img, format)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to save image file %q: %v\n",
			outFilename, err)
		outFilename = filename
		return
	}

	isResized = true
	return
}

func prepare(filename string) (string, error) {
	// image is not uri
	logger.Debug("prepare", filename)
	if strings.HasPrefix(filename, CustomWallpapersConfigDir) {
		updateModTime(filename)
		return filename, nil
	}

	filename, resizeExt, isResized := resizeImage(filename, CustomWallpapersConfigDir)

	md5sum, err := sumFileMd5(filename)
	if err != nil {
		return "", err
	}
	var ext string
	if isResized {
		ext = "." + resizeExt
	} else {
		ext = filepath.Ext(filename)
	}

	baseName := md5sum + ext
	dstFile := filepath.Join(CustomWallpapersConfigDir, baseName)
	_, err = os.Stat(dstFile)
	if err != nil {
		// copy image to cacheFile
		err = os.MkdirAll(CustomWallpapersConfigDir, 0755)
		if err != nil {
			return "", err
		}

		if isResized {
			err = os.Rename(filename, dstFile)
			if err != nil {
				return "", err
			}
		} else {
			err = dutils.CopyFile(filename, dstFile)
			if err != nil {
				return "", err
			}
		}

		time.AfterFunc(time.Second, func() {
			shrinkCache(baseName)
		})
	} else {
		updateModTime(dstFile)
		if isResized {
			// remove temp file
			err := os.Remove(filename)
			if err != nil && !os.IsNotExist(err) {
				_, _ = fmt.Fprintln(os.Stderr, "failed to remove temp file:", err)
			}
		}
	}

	return dstFile, nil
}

func shrinkCache(cacheFileBaseName string) {
	gs := gio.NewSettings("com.deepin.dde.appearance")
	defer gs.Unref()

	workspaceBackgrounds := gs.GetStrv("background-uris")
	var notDeleteFiles strv.Strv
	notDeleteFiles = append(notDeleteFiles, cacheFileBaseName)
	for _, uri := range workspaceBackgrounds {
		wbFile := dutils.DecodeURI(uri)
		if strings.HasPrefix(wbFile, CustomWallpapersConfigDir) {
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
	fileInfos, _ := ioutil.ReadDir(CustomWallpapersConfigDir)
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
			fullPath := filepath.Join(CustomWallpapersConfigDir, fileBaseName)
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
