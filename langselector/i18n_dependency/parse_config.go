/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package i18n_dependency

import "encoding/json"
import "io/ioutil"

type dependentPkgInfo struct {
	LangCode   string `json:"LangCode"`
	FormatType int32  `json:"FormatType"`
	DependPkg  string `json:"DependentPkg"`
	PkgPull    string `json:"PkgPull"`
}

type dependentPkgGroup struct {
	Category string             `json:"Category"`
	PkgInfos []dependentPkgInfo `json:"PkgInfos"`
}

type dependentPkgList struct {
	PkgDepends []dependentPkgGroup `json:"PkgDepends"`
}

func getPkgDependList(filename string) (*dependentPkgList, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var list dependentPkgList
	err = json.Unmarshal(contents, &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}
