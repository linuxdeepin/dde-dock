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
	"time"

	"pkg.deepin.io/dde/daemon/dstore"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type DStoreAdapter struct {
	store *dstore.DStore
}

func NewDStoreAdapter(store *dstore.DStore) *DStoreAdapter {
	return &DStoreAdapter{store: store}
}

func (s *DStoreAdapter) NewUninstallTransaction(pkgName string, purge bool, timeout time.Duration) UninstallTransaction {
	return s.store.NewUninstallTransaction(pkgName, purge, timeout)
}

// !!!!!
// NB: NewQueryTimeInstalledTransaction returns a pointer value and a error, event if thie pointer value is nil,
// the interface variable to which pointer value will be assigned WON'T be equal to nil on comparsion.
// To fix it, nil must be return explicitly as return value.
func (s *DStoreAdapter) NewQueryTimeInstalledTransaction(file string) (QueryTimeInstalledTransaction, error) {
	t, err := s.store.NewQueryTimeInstalledTransaction(file)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (s *DStoreAdapter) NewQueryPkgNameTransaction(path string) (QueryPkgNameTransaction, error) {
	t, err := s.store.NewQueryPkgNameTransaction(path)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (s *DStoreAdapter) NewQueryCategoryTransaction() (QueryCategoryTransaction, error) {
	t, err := s.store.NewQueryCategoryTransaction()
	if t == nil {
		return nil, err
	}
	return t, err
}
