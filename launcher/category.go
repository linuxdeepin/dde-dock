package launcher

import (
	"gir/gio-2.0"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

func queryCategoryID(cm CategoryManager, app *gio.DesktopAppInfo) (CategoryID, error) {
	cm.LoadCategoryInfo()
	defer cm.FreeAppCategoryInfo()

	return cm.QueryID(app)
}
