package accounts

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	imageBlurDBusPath      = "/com/deepin/daemon/ImageBlur"
	imageBlurDBusInterface = "com.deepin.daemon.ImageBlur"
)

type ImageBlur struct {
	service *dbusutil.Service
	mu      sync.Mutex
	tasks   map[string]struct{}

	signals *struct {
		BlurDone struct {
			imgFile     string
			imgBlurFile string
			ok          bool
		}
	}

	methods *struct {
		Get    func() `in:"source" out:"blurred"`
		Delete func() `in:"file"`
	}
}

func newImageBlur(service *dbusutil.Service) *ImageBlur {
	return &ImageBlur{
		service: service,
		tasks:   make(map[string]struct{}),
	}
}

func (ib *ImageBlur) GetInterfaceName() string {
	return imageBlurDBusInterface
}

func (ib *ImageBlur) Get(file string) (string, *dbus.Error) {
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
		return "", dbusutil.ToError(err)
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

func (ib *ImageBlur) Delete(file string) *dbus.Error {
	ib.mu.Lock()
	_, ok := ib.tasks[file]
	ib.mu.Unlock()

	if ok {
		return dbusutil.ToError(errors.New("generation task is in progress"))
	}

	blurFile := getImageBlurFile(file)
	err := os.Remove(blurFile)
	logger.Debugf("delete blur: %q, source: %q", blurFile, file)
	if os.IsNotExist(err) {
		err = nil
	}
	return dbusutil.ToError(err)
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
		ib.emitBlurDone(file, err == nil)

		ib.mu.Lock()
		delete(ib.tasks, file)
		ib.mu.Unlock()
	}()
}

func (ib *ImageBlur) emitBlurDone(file string, ok bool) {
	err := ib.service.Emit(ib, "BlurDone", file, getImageBlurFile(file), ok)
	if err != nil {
		logger.Warning(err)
	}
}

func genGaussianBlur(file string) {
	file = dutils.DecodeURI(file)
	if _imageBlur != nil {
		_imageBlur.Get(file)
	} else {
		logger.Warning("_imageBlur is nil")
	}
}
