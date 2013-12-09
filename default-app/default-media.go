package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

var (
	mediaGSettings = gio.NewSettings("org.gnome.desktop.media-handling")
)

type MediaMount struct {
	AllowAutoMount     bool
	TypeIgnoreList     *property.GSettingsStrvProperty `access:"readwrite"`
	TypeOpenFolderList *property.GSettingsStrvProperty `access:"readwrite"`
	TypeExecList       *property.GSettingsStrvProperty `access:"readwrite"`
}

func (media *MediaMount) SetAllowAutoRun(allow bool) {
	mediaGSettings.SetBoolean("automount", allow)
	mediaGSettings.SetBoolean("automount-open", allow)
}

func NewMediaMount() *MediaMount {
	media := &MediaMount{}
	media.AllowAutoMount = GetAllowAutoMount()
	media.TypeIgnoreList = property.NewGSettingsStrvProperty(media,
		"TypeIgnoreList", mediaGSettings,
		"autorun-x-content-ignore")
	media.TypeOpenFolderList = property.NewGSettingsStrvProperty(media,
		"TypeOpenFolderList", mediaGSettings,
		"autorun-x-content-open-folder")
	media.TypeExecList = property.NewGSettingsStrvProperty(media,
		"TypeExecList", mediaGSettings,
		"autorun-x-content-start-app")

	ListenGSettingsChanged(media)

	return media
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
			dbus.NotifyChange(media, "AllowAutoMount")
		})

	mediaGSettings.Connect("changed::automount-open",
		func(s *gio.Settings, name string) {
			media.AllowAutoMount = GetAllowAutoMount()
			dbus.NotifyChange(media, "AllowAutoMount")
		})
}
