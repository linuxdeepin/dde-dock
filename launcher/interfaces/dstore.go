/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package interfaces

import (
	// "pkg.deepin.io/dde/daemon/dstore"
	"gir/gio-2.0"
	"time"
)

type QueryCategoryTransaction interface {
	Query(*gio.DesktopAppInfo) (string, error)
}

type UninstallTransaction interface {
	Exec() error
}

type InstallTransaction interface {
	Exec() error
}

type QueryPkgNameTransaction interface {
	Query(string) string
}

type QueryTimeInstalledTransaction interface {
	Query(string) int64
}

// DStore is interface for deepin store.
type DStore interface {
	NewUninstallTransaction(pkgName string, purge bool, timeout time.Duration) UninstallTransaction
	NewQueryTimeInstalledTransaction(file string) (QueryTimeInstalledTransaction, error)
	NewQueryPkgNameTransaction(path string) (QueryPkgNameTransaction, error)
	NewQueryCategoryTransaction() (QueryCategoryTransaction, error)
}
