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
	"fmt"
)

type DeepinQueryIDTransaction struct {
	pkgToCategory map[string]string
}

func NewDeepinQueryIDTransaction(pkgToCategoryFile string) (*DeepinQueryIDTransaction, error) {
	t := &DeepinQueryIDTransaction{
		pkgToCategory: map[string]string{},
	}

	var err error
	t.pkgToCategory, err = getCategoryInfo(pkgToCategoryFile)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *DeepinQueryIDTransaction) Query(pkgName string) (string, error) {
	cidName, ok := t.pkgToCategory[pkgName]
	if !ok {
		return OthersName, fmt.Errorf("No such a category for package %q", pkgName)
	}
	return cidName, nil
}
