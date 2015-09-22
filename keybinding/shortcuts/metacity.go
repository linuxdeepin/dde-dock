package shortcuts

import (
	"pkg.deepin.io/lib/gettext"
)

func ListMetacityShortcut() Shortcuts {
	s := newMetacityGSetting()
	defer s.Unref()
	return doListShortcut(s, metacityIdNameMap(), KeyTypeMetacity)
}

func resetMetacityAccels() {
	meta := newMetacityGSetting()
	defer meta.Unref()
	doResetAccels(meta)

	gala := newGalaGSetting()
	defer gala.Unref()
	doResetAccels(gala)
}

func disableMetacityAccels(key string) {
	meta := newMetacityGSetting()
	defer meta.Unref()
	doDisableAccles(meta, key)

	gala := newGalaGSetting()
	defer gala.Unref()
	doDisableAccles(gala, key)
}

func addMetacityAccel(key, accel string) {
	meta := newMetacityGSetting()
	defer meta.Unref()
	doAddAccel(meta, key, accel)

	gala := newGalaGSetting()
	defer gala.Unref()
	doAddAccel(gala, key, accel)
}

func delMetacityAccel(key, accel string) {
	meta := newMetacityGSetting()
	defer meta.Unref()
	doDelAccel(meta, key, accel)

	gala := newGalaGSetting()
	defer gala.Unref()
	doDelAccel(gala, key, accel)
}

func metacityIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"expose-all-windows": gettext.Tr("Display windows of all workspaces"),
		"expose-windows":     gettext.Tr("Display windows of current workspace"),
		"preview-workspace":  gettext.Tr("Display workspace"),
	}

	return idNameMap
}
