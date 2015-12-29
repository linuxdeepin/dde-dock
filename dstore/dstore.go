package dstore

import (
	"time"
)

type DStore struct {
}

func New() (*DStore, error) {
	return &DStore{}, nil
}

func (*DStore) NewUninstallTransaction(pkgName string, purge bool, timeout time.Duration) *DUninstallTransaction {
	return NewDUninstallTransaction(pkgName, purge, timeout)
}

func (*DStore) NewQueryTimeInstalledTransaction(file string) (*DQueryTimeInstalledTransaction, error) {
	return NewDQueryTimeInstalledTransaction(file)
}
func (*DStore) NewQueryPkgNameTransaction(path string) (*DQueryPkgNameTransaction, error) {
	return NewDQueryPkgNameTransaction(path)
}

func (*DStore) NewInstallTransaction(pkgs string, desc string, timeout time.Duration) *DInstallTransaction {
	return NewDInstallTransaction(pkgs, desc, timeout)
}

func (*DStore) NewQueryCategoryTransaction() (*QueryCategoryTransaction, error) {
	return NewQueryCategoryTransaction(DesktopPkgMapFile, AppInfoFile, XCategoryAppInfoFile)
}
