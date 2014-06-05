package mime

import (
	"dlib/dbus"
	"dlib/logger"
	libutils "dlib/utils"
	"os"
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
	logObject = logger.NewLogger("daemon/mime")
	objUtils  = libutils.NewUtils()
)

func (dapp *DefaultApps) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_DEFAULT_APPS_DEST,
		_DEFAULT_APPS_PATH,
		_DEFAULT_APPS_IFC,
	}
}

func (media *MediaMount) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_DEFAULT_APPS_DEST,
		_MEDIA_MOUNT_PATH,
		_MEDIA_MOUNT_IFC,
	}
}

func (dapp *DefaultApps) Reset() bool {
	homeDir, ok := objUtils.GetHomeDir()
	if !ok {
		logObject.Warning("Get homeDir failed")
		return false
	}
	if err := os.Remove(homeDir + "/" + MIME_CACHE_FILE); err != nil {
		logObject.Warningf("Delete '%s' failed: %v",
			homeDir+"/"+MIME_CACHE_FILE, err)
		return false
	}
	_TerminalGSettings.Reset("exec")

	return true
}

func (media *MediaMount) Reset() bool {
	mediaGSettings.Reset(MEDIA_KEY_AUTOMOUNT)
	mediaGSettings.Reset(MEDIA_KEY_AUTOOPEN)
	mediaGSettings.Reset(MEDIA_KEY_AUTORUN_NEVER)
	mediaGSettings.Reset(MEDIA_KEY_IGNORE)
	mediaGSettings.Reset(MEDIA_KEY_OPEN_FOLDER)
	mediaGSettings.Reset(MEDIA_KEY_START_SOFT)

	return true
}

func Start() {
	logObject.BeginTracing()
	defer logObject.EndTracing()

	var err error

	dapp := NewDefaultApps()
	if dapp == nil {
		return
	}
	err = dbus.InstallOnSession(dapp)
	if err != nil {
		logObject.Infof("Install Session Failed: %v", err)
		panic(err)
	}

	media := NewMediaMount()
	err = dbus.InstallOnSession(media)
	if err != nil {
		logObject.Infof("Install Session Failed: %v", err)
		panic(err)
	}
	dbus.DealWithUnhandledMessage()
}
