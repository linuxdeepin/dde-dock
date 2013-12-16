package main

import (
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

var (
	_mediaGSettings = gio.NewSettings("org.gnome.desktop.media-handling")
)

type MediaMount struct {
	AutoMountOpen      *property.GSettingsBoolProperty `access:"readwrite"`
	TypeIgnoreList     *property.GSettingsStrvProperty `access:"readwrite"`
	TypeOpenFolderList *property.GSettingsStrvProperty `access:"readwrite"`
	TypeExecList       *property.GSettingsStrvProperty `access:"readwrite"`
}

func NewMediaMount() *MediaMount {
	media := &MediaMount{}
	media.AutoMountOpen = property.NewGSettingsBoolProperty(
		media, "AutoMountOpen",
		_mediaGSettings, "automount-open")
	media.TypeIgnoreList = property.NewGSettingsStrvProperty(
		media, "TypeIgnoreList",
		_mediaGSettings, "autorun-x-content-ignore")
	media.TypeOpenFolderList = property.NewGSettingsStrvProperty(
		media, "TypeOpenFolderList",
		_mediaGSettings, "autorun-x-content-open-folder")
	media.TypeExecList = property.NewGSettingsStrvProperty(
		media, "TypeExecList",
		_mediaGSettings, "autorun-x-content-start-app")

	return media
}
