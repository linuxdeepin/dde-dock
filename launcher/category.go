/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
