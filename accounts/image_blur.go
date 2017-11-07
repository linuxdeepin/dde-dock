package accounts

import (
	"os/exec"
	"path/filepath"
	"sync"

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

func (ib *ImageBlur) Get(file string) string {
	ib.mu.Lock()
	_, ok := ib.tasks[file]
	ib.mu.Unlock()

	if ok {
		// generating
		return ""
	}

	blurFile := getImageBlurFile(file)
	if dutils.IsFileExist(blurFile) {
		return blurFile
	}

	ib.gen(file)
	return ""
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
		logger.Debugf("image-blur-helper output: %s", output)
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
