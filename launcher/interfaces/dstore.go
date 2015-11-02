package interfaces

import (
	// "pkg.deepin.io/dde/daemon/dstore"
	"time"
)

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

// type UninstallTransaction interface {
// 	dstore.UninstallTransaction
// }

// type QueryPkgNameTransaction interface {
// 	dstore.QueryPkgNameTransaction
// }

// type QueryInstalledTimeTransaction interface {
// 	dstore.QueryTimeInstalledTransaction
// }

// DStore is interface for deepin store.
type DStore interface {
	NewUninstallTransaction(pkgName string, purge bool, timeout time.Duration) UninstallTransaction
	NewQueryTimeInstalledTransaction(file string) (QueryTimeInstalledTransaction, error)
	NewQueryPkgNameTransaction(path string) (QueryPkgNameTransaction, error)
}
