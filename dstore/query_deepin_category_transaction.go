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
