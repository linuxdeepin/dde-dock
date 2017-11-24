/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package mime

import (
	"encoding/json"
	"io/ioutil"
)

type defaultAppTable struct {
	Apps defaultAppInfos `json:"DefaultApps"`
}

type defaultAppInfo struct {
	AppId   []string `json:"AppId"`
	AppType string   `json:"AppType"`
	Types   []string `json:"SupportedType"`
}
type defaultAppInfos []*defaultAppInfo

func unmarshal(file string) (*defaultAppTable, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var table defaultAppTable
	err = json.Unmarshal(content, &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func marshal(v interface{}) (string, error) {
	content, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func genMimeAppsFile(data string) error {
	table, err := unmarshal(data)
	if err != nil {
		logger.Warning("[genMimeAppsFile] unmarshal failed:", err)
		return err
	}

	for _, info := range table.Apps {
		var validId = ""
		for _, ty := range info.Types {
			if validId != "" {
				SetAppInfo(ty, validId)
				continue
			}

			for _, id := range info.AppId {
				err := SetAppInfo(ty, id)
				if err != nil {
					logger.Warningf("[genMimeAppsFile] set '%s' to parse '%s' failed: %v\n",
						info.AppId, ty, err)
					continue
				}
				validId = id
				break
			}
		}
	}

	return nil
}
