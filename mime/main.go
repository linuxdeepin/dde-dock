package main

import (
        "dlib"
        "dlib/dbus"
        "dlib/logger"
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
        defer func() {
                if err := recover(); err != nil {
                        logger.Println("Recover Error:", err)
                }
        }()

        dapp := NewDefaultApps()
        if dapp == nil {
                return
        }
        err := dbus.InstallOnSession(dapp)
        if err != nil {
                logger.Println("Install Session Failed:", err)
                panic(err)
        }

        media := NewMediaMount()
        err = dbus.InstallOnSession(media)
        if err != nil {
                logger.Println("Install Session Failed:", err)
                panic(err)
        }
        dbus.DealWithUnhandledMessage()

        go dlib.StartLoop()
        if err = dbus.Wait(); err != nil {
                logger.Println("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
