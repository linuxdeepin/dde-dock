package dstore

import (
	"time"
)

type DStore struct {
}

func New() (*DStore, error) {
	return &DStore{}, nil
}

func IsInstalled(pkgName string) bool {
	proxy, err := newDStoreManager()
	if err != nil {
		return false
	}
	defer destroyDStoreManager(proxy)

	installed, _ := proxy.PackageExists(pkgName)
	return installed
}

func IsExists(pkgName string) bool {
	proxy, err := newDStoreManager()
	if err != nil {
		return false
	}
	defer destroyDStoreManager(proxy)

	exists, _ := proxy.PackageInstallable(pkgName)
	return exists
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
