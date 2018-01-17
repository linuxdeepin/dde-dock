package accounts

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
)

type ImageBlur struct {
	mu    sync.Mutex
	tasks map[string]struct{}

	// signal:
	BlurDone func(imgFile, imgBlurFile string, ok bool)
}

func newImageBlur() *ImageBlur {
	return &ImageBlur{
		tasks: make(map[string]struct{}),
	}
}

func (ib *ImageBlur) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: "/com/deepin/daemon/ImageBlur",
		Interface:  "com.deepin.daemon.ImageBlur",
	}
}

func (ib *ImageBlur) Get(file string) (string, error) {
	ib.mu.Lock()
	_, ok := ib.tasks[file]
	ib.mu.Unlock()

	if ok {
		// generating
		return "", nil
	}

	blurFile := getImageBlurFile(file)

	fileInfo, err := os.Stat(file)
	if err != nil {
		logger.Warning(err)
		if os.IsNotExist(err) {
			// source file not exist
			os.Remove(blurFile)
		}
		return "", err
	}

	blurFileInfo, err := os.Stat(blurFile)
	if err == nil {
		fileChangeTime := getChangeTime(fileInfo)
		blurFileChangeTime := getChangeTime(blurFileInfo)

		if fileChangeTime.Before(blurFileChangeTime) {
			return blurFile, nil
		} else {
			// delete old, generate new
			os.Remove(blurFile)
		}
	}

	ib.gen(file)
	return "", nil
}

// getChangeTime get time when file status was last changed.
func getChangeTime(fileInfo os.FileInfo) time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
}

func (ib *ImageBlur) Delete(file string) error {
	ib.mu.Lock()
	_, ok := ib.tasks[file]
	ib.mu.Unlock()

	if ok {
		return errors.New("generation task is in progress")
	}

	blurFile := getImageBlurFile(file)
	err := os.Remove(blurFile)
	logger.Debugf("delete blur: %q, source: %q", blurFile, file)
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}

const imageBlurDir = "/var/cache/image-blur"

func getImageBlurFile(src string) string {
	md5sum, _ := dutils.SumStrMd5(src)
	return filepath.Join(imageBlurDir, md5sum+filepath.Ext(src))
}

func (ib *ImageBlur) gen(file string) {
	ib.mu.Lock()
	_, ok := ib.tasks[file]
	if ok {
		logger.Debug("ImageBlur.gen task exist:", file)
		ib.mu.Unlock()
		return
	}

	ib.tasks[file] = struct{}{}
	ib.mu.Unlock()

	go func() {
		logger.Debug("ImageBlur.gen will blur image:", file)
		output, err := exec.Command("/usr/lib/deepin-api/image-blur-helper", file).CombinedOutput()
		if len(output) > 0 {
			logger.Debugf("image-blur-helper output: %s", output)
		}
		if err != nil {
			logger.Warningf("failed to blur image %q: %v", file, err)
		}
		dbus.Emit(ib, "BlurDone", file, getImageBlurFile(file), err == nil)

		ib.mu.Lock()
		delete(ib.tasks, file)
		ib.mu.Unlock()
	}()
}

func genGaussianBlur(file string) {
	file = dutils.DecodeURI(file)
	if _imageBlur != nil {
		_imageBlur.gen(file)
	} else {
		logger.Warning("_imageBlur is nil")
	}
}
