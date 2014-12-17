package mime

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/dbus/property"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	MEDIA_KEY_AUTOMOUNT     = "automount"
	MEDIA_KEY_AUTOOPEN      = "automount-open"
	MEDIA_KEY_AUTORUN_NEVER = "autorun-never"
	MEDIA_KEY_IGNORE        = "autorun-x-content-ignore"
	MEDIA_KEY_OPEN_FOLDER   = "autorun-x-content-open-folder"
	MEDIA_KEY_START_SOFT    = "autorun-x-content-start-app"
)

const (
	mediaSchema = "org.gnome.desktop.media-handling"
)

type MediaMount struct {
	AutoMountOpen      *property.GSettingsBoolProperty `access:"readwrite"`
	MediaActionChanged func()

	settings *gio.Settings
}

func NewMediaMount() *MediaMount {
	if !dutils.IsGSchemaExist(mediaSchema) {
		return nil
	}

	media := &MediaMount{}

	media.settings = gio.NewSettings(mediaSchema)
	media.AutoMountOpen = property.NewGSettingsBoolProperty(
		media, "AutoMountOpen",
		media.settings, MEDIA_KEY_AUTOMOUNT)
	media.listenGSettings()

	return media
}

func (media *MediaMount) destroy() {
	dbus.UnInstallObject(media)
	media.settings.Unref()
}

func (media *MediaMount) SetMediaAppByMime(mime, appID string) {
	media.setActionByMime(mime, appID)
}

func (media *MediaMount) DefaultMediaAppByMime(mime string) AppInfo {
	return media.getActionByMime(mime)
}

func (media *MediaMount) MediaAppListByMime(mime string) []AppInfo {
	return media.getActionsByMime(mime)
}

func (media *MediaMount) listenGSettings() {
	media.settings.Connect("changed::autorun-x-content-ignore",
		func(s *gio.Settings, key string) {
			dbus.Emit(media, "MediaActionChanged")
		})

	media.settings.Connect("changed::autorun-x-content-open-folder",
		func(s *gio.Settings, key string) {
			dbus.Emit(media, "MediaActionChanged")
		})

	media.settings.Connect("changed::autorun-x-content-start-app",
		func(s *gio.Settings, key string) {
			dbus.Emit(media, "MediaActionChanged")
		})
}

func (media *MediaMount) getActionByMime(mime string) AppInfo {
	ignoreList := media.settings.GetStrv(MEDIA_KEY_IGNORE)
	openFolderList := media.settings.GetStrv(MEDIA_KEY_OPEN_FOLDER)
	runSoftList := media.settings.GetStrv(MEDIA_KEY_START_SOFT)

	if isMimeExist(mime, ignoreList) {
		return AppInfo{ID: "Nothing", Name: Tr("Nothing"), Exec: ""}
	}
	if isMimeExist(mime, openFolderList) {
		return AppInfo{ID: "Open Folder", Name: Tr("Open Folder"), Exec: ""}
	}

	if isMimeExist(mime, runSoftList) {
		return AppInfo{ID: "Run Soft", Name: Tr("Run Software"), Exec: ""}
	}

	m := DefaultApps{}
	return m.DefaultAppViaType(mime)
}

func (media *MediaMount) getActionsByMime(mime string) []AppInfo {
	apps := []AppInfo{}
	defaultApps := []AppInfo{
		AppInfo{ID: "Nothing", Name: Tr("Nothing"), Exec: ""},
		AppInfo{ID: "Open Folder", Name: Tr("Open Folder"), Exec: ""},
	}
	m := DefaultApps{}
	apps = m.AppsListViaType(mime)
	apps = append(apps, defaultApps...)
	//if mime == "x-content/unix-software" {
	//apps = append(apps,
	//AppInfo{ID: "Run Soft", Name: Tr("Run Soft", Exec: ""})
	//}

	return apps
}

func (media *MediaMount) setActionByMime(mime, appID string) {
	ignoreList := media.settings.GetStrv(MEDIA_KEY_IGNORE)
	openFolderList := media.settings.GetStrv(MEDIA_KEY_OPEN_FOLDER)
	runSoftList := media.settings.GetStrv(MEDIA_KEY_START_SOFT)

	switch appID {
	case "Nothing":
		if !isMimeExist(mime, ignoreList) {
			ignoreList = append(ignoreList, mime)
			media.settings.SetStrv(MEDIA_KEY_IGNORE, ignoreList)
		}

		list, ok := delMimeFromList(mime, openFolderList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_OPEN_FOLDER, list)
		}

		list, ok = delMimeFromList(mime, runSoftList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}
	case "Open Folder":
		if !isMimeExist(mime, openFolderList) {
			openFolderList = append(openFolderList, mime)
			media.settings.SetStrv(MEDIA_KEY_OPEN_FOLDER, openFolderList)
		}

		list, ok := delMimeFromList(mime, ignoreList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_IGNORE, list)
		}

		list, ok = delMimeFromList(mime, runSoftList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}
	case "Run Soft", "nautilus-autorun-software.desktop":
		if !isMimeExist(mime, runSoftList) {
			runSoftList = append(runSoftList, mime)
			media.settings.SetStrv(MEDIA_KEY_START_SOFT, runSoftList)
		}

		list, ok := delMimeFromList(mime, ignoreList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_IGNORE, list)
		}

		list, ok = delMimeFromList(mime, openFolderList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_OPEN_FOLDER, list)
		}
	default:
		m := DefaultApps{}
		m.SetDefaultAppViaType(mime, appID)

		list, ok := delMimeFromList(mime, ignoreList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_IGNORE, list)
		}

		list, ok = delMimeFromList(mime, openFolderList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_OPEN_FOLDER, list)
		}

		list, ok = delMimeFromList(mime, runSoftList)
		if ok {
			media.settings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}
	}
}

func isMimeExist(mime string, mimes []string) bool {
	for _, v := range mimes {
		if mime == v {
			return true
		}
	}

	return false
}

func delMimeFromList(mime string, mimes []string) ([]string, bool) {
	if !isMimeExist(mime, mimes) {
		return mimes, false
	}

	rets := []string{}
	for _, v := range mimes {
		if v == mime {
			continue
		}
		rets = append(rets, v)
	}

	return rets, true
}
