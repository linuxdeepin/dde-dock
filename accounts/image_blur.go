package accounts

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus"
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
	//nolint
	signals *struct {
		BlurDone struct {
			imgFile     string
			imgBlurFile string
			ok          bool
		}
	}
	//nolint
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
	logger.Debugf("Get %q", file)

	blurFile := getImageBlurFile(file)
	fileInfo, err := os.Lstat(file)
	if err != nil {
		logger.Warning(err)
		if os.IsNotExist(err) {
			// source file not exist
			_ = os.Remove(blurFile)
		}
		return "", dbusutil.ToError(err)
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// file is a symbolic link
		logger.Debug("the file is a symbolic link")
		realFile, err := filepath.EvalSymlinks(file)
		if err != nil {
			logger.Warning(err)
			return "", dbusutil.ToError(err)
		}

		logger.Debugf("real file: %q", realFile)
		file = realFile
		blurFile = getImageBlurFile(file)
	}

	ib.mu.Lock()
	_, ok := ib.tasks[file]
	ib.mu.Unlock()

	if ok {
		// generating
		return "", nil
	}

	_, err = os.Stat(blurFile)
	if err == nil {
		return blurFile, nil
	}

	ib.gen(file, blurFile)
	return "", nil
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

func (ib *ImageBlur) gen(file, blurFile string) {
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

		var blurOk bool
		if err != nil {
			logger.Warningf("failed to blur image %q: %v", file, err)
		} else {
			_, err = os.Stat(blurFile)
			if err != nil {
				logger.Warning("failed to stat blurFile:", err)
			} else {
				blurOk = true
			}
		}
		ib.emitBlurDone(file, blurOk)

		ib.mu.Lock()
		delete(ib.tasks, file)
		ib.mu.Unlock()
	}()
}

func (ib *ImageBlur) emitBlurDone(file string, ok bool) {
	blurFile := getImageBlurFile(file)
	logger.Debugf("emit signal BlurDone (%q,%q,%v)", file, blurFile, ok)
	err := ib.service.Emit(ib, "BlurDone", file, blurFile, ok)
	if err != nil {
		logger.Warning(err)
	}
}

func genGaussianBlur(file string) {
	file = dutils.DecodeURI(file)
	if _imageBlur != nil {
		_, _ = _imageBlur.Get(file)
	} else {
		logger.Warning("_imageBlur is nil")
	}
}
