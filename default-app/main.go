package main

import (
	"dlib/dbus"
)

const (
	_DEFAULT_APPS_DEST = "com.deepin.daemon.DefaultApps"
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

func main() {
	dapp := DefaultApps{}
	dbus.InstallOnSession(&dapp)

	media := NewMediaMount()
	dbus.InstallOnSession (media)
	select {}
}
