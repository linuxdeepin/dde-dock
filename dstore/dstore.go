/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
