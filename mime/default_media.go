package mime

import (
	"pkg.linuxdeepin.com/lib/dbus/property"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	MEDIA_KEY_AUTOMOUNT     = "automount"
	MEDIA_KEY_AUTOOPEN      = "automount-open"
	MEDIA_KEY_AUTORUN_NEVER = "autorun-never"
	MEDIA_KEY_IGNORE        = "autorun-x-content-ignore"
	MEDIA_KEY_OPEN_FOLDER   = "autorun-x-content-open-folder"
	MEDIA_KEY_START_SOFT    = "autorun-x-content-start-app"
)

var (
	mediaGSettings = gio.NewSettings("org.gnome.desktop.media-handling")
)

type MediaMount struct {
	AutoMountOpen      *property.GSettingsBoolProperty `access:"readwrite"`
	MediaActionChanged func()
}

func NewMediaMount() *MediaMount {
	media := &MediaMount{}
	media.AutoMountOpen = property.NewGSettingsBoolProperty(
		media, "AutoMountOpen",
		mediaGSettings, MEDIA_KEY_AUTOMOUNT)
	media.listenGSettings()

	return media
}

func (op *MediaMount) SetMediaAppByMime(mime, appID string) {
	setActionByMime(mime, appID)
}

func (op *MediaMount) DefaultMediaAppByMime(mime string) AppInfo {
	return getActionByMime(mime)
}

func (op *MediaMount) MediaAppListByMime(mime string) []AppInfo {
	return getActionsByMime(mime)
}

func (op *MediaMount) listenGSettings() {
	mediaGSettings.Connect("changed::autorun-x-content-ignore",
		func(s *gio.Settings, key string) {
			op.MediaActionChanged()
		})

	mediaGSettings.Connect("changed::autorun-x-content-open-folder",
		func(s *gio.Settings, key string) {
			op.MediaActionChanged()
		})

	mediaGSettings.Connect("changed::autorun-x-content-start-app",
		func(s *gio.Settings, key string) {
			op.MediaActionChanged()
		})
}

func getActionByMime(mime string) AppInfo {
	ignoreList := mediaGSettings.GetStrv(MEDIA_KEY_IGNORE)
	openFolderList := mediaGSettings.GetStrv(MEDIA_KEY_OPEN_FOLDER)
	runSoftList := mediaGSettings.GetStrv(MEDIA_KEY_START_SOFT)

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

func getActionsByMime(mime string) []AppInfo {
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

func setActionByMime(mime, appID string) {
	ignoreList := mediaGSettings.GetStrv(MEDIA_KEY_IGNORE)
	openFolderList := mediaGSettings.GetStrv(MEDIA_KEY_OPEN_FOLDER)
	runSoftList := mediaGSettings.GetStrv(MEDIA_KEY_START_SOFT)

	switch appID {
	case "Nothing":
		if !isMimeExist(mime, ignoreList) {
			ignoreList = append(ignoreList, mime)
			mediaGSettings.SetStrv(MEDIA_KEY_IGNORE, ignoreList)
		}

		list, ok := delMimeFromList(mime, openFolderList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_OPEN_FOLDER, list)
		}

		list, ok = delMimeFromList(mime, runSoftList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}
	case "Open Folder":
		if !isMimeExist(mime, openFolderList) {
			openFolderList = append(openFolderList, mime)
			mediaGSettings.SetStrv(MEDIA_KEY_OPEN_FOLDER, openFolderList)
		}

		list, ok := delMimeFromList(mime, ignoreList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}

		list, ok = delMimeFromList(mime, runSoftList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}
	case "Run Soft":
		if !isMimeExist(mime, runSoftList) {
			runSoftList = append(runSoftList, mime)
			mediaGSettings.SetStrv(MEDIA_KEY_START_SOFT, runSoftList)
		}

		list, ok := delMimeFromList(mime, ignoreList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_START_SOFT, list)
		}

		list, ok = delMimeFromList(mime, openFolderList)
		if ok {
			mediaGSettings.SetStrv(MEDIA_KEY_OPEN_FOLDER, list)
		}
	default:
		m := DefaultApps{}
		m.SetDefaultAppViaType(mime, appID)
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
