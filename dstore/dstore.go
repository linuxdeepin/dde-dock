/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
