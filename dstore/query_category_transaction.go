/*
 * Copyright (C) 2015 ~ 2017 Deepin Technology Co., Ltd.
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
	"gir/gio-2.0"
)

type QueryCategoryTransaction struct {
	pkgQuery       *DQueryPkgNameTransaction
	deepinQuery    *DeepinQueryIDTransaction
	xCategoryQuery *XCategoryQueryIDTransaction
}

func NewQueryCategoryTransaction(desktopToPkgFile string, appInfoFile string, xcategoryFile string) (*QueryCategoryTransaction, error) {
	t := &QueryCategoryTransaction{}
	var err1 error
	var err2 error

	t.pkgQuery, err1 = NewDQueryPkgNameTransaction(desktopToPkgFile)
	t.deepinQuery, err2 = NewDeepinQueryIDTransaction(appInfoFile)
	t.xCategoryQuery, _ = NewXCategoryQueryIDTransaction(xcategoryFile, AllCategoryInfoFile)

	if err1 != nil {
		return t, err1
	} else if err2 != nil {
		return t, err2
	}

	return t, nil
}

func (t *QueryCategoryTransaction) Query(app *gio.DesktopAppInfo) (string, error) {
	if t.pkgQuery != nil && t.deepinQuery != nil {
		pkgName := t.pkgQuery.Query(app.GetId())
		cid, err := t.deepinQuery.Query(pkgName)
		if err == nil {
			return cid, nil
		}
	}

	cid, err := t.xCategoryQuery.Query(app.GetCategories())
	return cid, err
}
