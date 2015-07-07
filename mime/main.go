package mime

import (
	"fmt"
	"os"
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	_DEFAULT_APPS_DEST = "com.deepin.daemon.Mime"
	_DEFAULT_APPS_PATH = "/com/deepin/daemon/DefaultApps"
	_DEFAULT_APPS_IFC  = "com.deepin.daemon.DefaultApps"
	_MEDIA_MOUNT_PATH  = "/com/deepin/daemon/MediaMount"
	_MEDIA_MOUNT_IFC   = "com.deepin.daemon.MediaMount"

	_HTTP_CONTENT_TYPE     = "x-scheme-handler/http"
	_HTTPS_CONTENT_TYPE    = "x-scheme-handler/https"
	_MAIL_CONTENT_TYPE     = "x-scheme-handler/mailto"
	_CALENDAR_CONTENT_TYPE = "text/calendar"
	_EDITOR_CONTENT_TYPE   = "text/plain"
	_AUDIO_CONTENT_TYPE    = "audio/mpeg"
	_VIDEO_CONTENT_TYPE    = "video/mp4"
)

var (
	logger = log.NewLogger("dde-daemon/mime")
)

func (dapp *DefaultApps) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       _DEFAULT_APPS_DEST,
		ObjectPath: _DEFAULT_APPS_PATH,
		Interface:  _DEFAULT_APPS_IFC,
	}
}

func (media *MediaMount) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       _DEFAULT_APPS_DEST,
		ObjectPath: _MEDIA_MOUNT_PATH,
		Interface:  _MEDIA_MOUNT_IFC,
	}
}

func (dapp *DefaultApps) Reset() bool {
	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		logger.Warning("Get homeDir failed")
		return false
	}
	if err := os.Remove(homeDir + "/" + MIME_CACHE_FILE); err != nil {
		logger.Warning("Delete '%s' failed: %v",
			homeDir+"/"+MIME_CACHE_FILE, err)
		return false
	}
	_TerminalGSettings().Reset("exec")

	return true
}

func (media *MediaMount) Reset() bool {
	media.settings.Reset(MEDIA_KEY_AUTOMOUNT)
	media.settings.Reset(MEDIA_KEY_AUTOOPEN)
	media.settings.Reset(MEDIA_KEY_AUTORUN_NEVER)
	media.settings.Reset(MEDIA_KEY_IGNORE)
	media.settings.Reset(MEDIA_KEY_OPEN_FOLDER)
	media.settings.Reset(MEDIA_KEY_START_SOFT)

	return true
}

var (
	_dapp  *DefaultApps
	_media *MediaMount
)

func startDefaultApps() error {
	_dapp = NewDefaultApps()
	if _dapp == nil {
		return fmt.Errorf("Create DefaultApps Failed")
	}

	err := dbus.InstallOnSession(_dapp)
	if err != nil {
		return err
	}

	return nil
}

func endDefaultApps() {
	if _dapp == nil {
		return
	}

	_dapp.destroy()
	_dapp = nil
}

func startMediaMount() error {
	_media = NewMediaMount()
	if _media == nil {
		return fmt.Errorf("Create MediaMount Failed")
	}
	err := dbus.InstallOnSession(_media)
	if err != nil {
		return err
	}

	return nil
}

func endMediaMount() {
	if _media == nil {
		return
	}

	_media.destroy()
	_media = nil
}

func finalize() {
	endDefaultApps()
	endMediaMount()
	logger.EndTracing()
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("mime", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Stop() error {
	if _dapp == nil {
		return nil
	}

	finalize()
	return nil
}

func (d *Daemon) Start() error {
	if _dapp != nil {
		return nil
	}

	logger.BeginTracing()
	err := startDefaultApps()
	if err != nil {
		logger.Error(err)
		logger.EndTracing()
		endDefaultApps()
		return err
	}

	err = startMediaMount()
	if err != nil {
		logger.Error(err)
		finalize()
		return err
	}
	return nil
}
