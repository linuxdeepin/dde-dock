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
