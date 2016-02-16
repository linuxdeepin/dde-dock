/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package timedate

import "pkg.deepin.io/dde/daemon/loader"

func init() {
	loader.Register(NewDaemon(logger))
}
