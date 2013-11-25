package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"fmt"
)

var (
	busConn        *dbus.Conn
	mediaGSettings = gio.NewSettings("org.gnome.desktop.media-handling")
)

type MediaMount struct {
	AllowAutoMount     bool `access:"read"`
	TypeIgnoreList     dbus.Property
	TypeOpenFolderList dbus.Property
	TypeExecList       dbus.Property
}

func (media *MediaMount) SetAllowAutoRun(allow bool) {
	mediaGSettings.SetBoolean("automount", allow)
	mediaGSettings.SetBoolean("automount-open", allow)
}

func NewMediaMount() *MediaMount {
	var err error

	media := MediaMount{}
	busConn, err = dbus.SessionBus()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	media.AllowAutoMount = GetAllowAutoMount()
	media.TypeIgnoreList = property.NewGSettingsPropertyFull(
		mediaGSettings, "autorun-x-content-ignore", []string{},
		busConn, _MEDIA_MOUNT_PATH, _MEDIA_MOUNT_IFC, "TypeIgnoreList")
	media.TypeOpenFolderList = property.NewGSettingsPropertyFull(
		mediaGSettings, "autorun-x-content-open-folder", []string{},
		busConn, _MEDIA_MOUNT_PATH, _MEDIA_MOUNT_IFC,
		"TypeOpenFolderList")
	media.TypeExecList = property.NewGSettingsPropertyFull(
		mediaGSettings, "autorun-x-content-start-app", []string{},
		busConn, _MEDIA_MOUNT_PATH, _MEDIA_MOUNT_IFC, "TypeExecList")

	ListenGSettingsChanged(&media)

	return &media
}

func GetAllowAutoMount() bool {
	if GetAutoMount() && GetAutoMountOpen() {
		return true
	}

	return false
}

func GetAutoMount() bool {
	return mediaGSettings.GetBoolean("automount")
}

func GetAutoMountOpen() bool {
	return mediaGSettings.GetBoolean("automount-open")
}

func ListenGSettingsChanged(media *MediaMount) {
	mediaGSettings.Connect("changed::automount",
		func(s *gio.Settings, name string) {
			media.AllowAutoMount = GetAllowAutoMount()
			dbus.NotifyChange(busConn, media, "AllowAutoMount")
		})

		mediaGSettings.Connect("changed::automount-open",
		func(s *gio.Settings, name string) {
			media.AllowAutoMount = GetAllowAutoMount()
			dbus.NotifyChange(busConn, media, "AllowAutoMount")
		})
}
