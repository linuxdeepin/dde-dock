/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package bluetooth

import (
	"encoding/json"
	"os/exec"
	"strings"
)

func isStringInArray(str string, list []string) bool {
	for _, tmp := range list {
		if tmp == str {
			return true
		}
	}
	return false
}

func marshalJSON(v interface{}) (strJSON string) {
	byteJSON, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return
	}
	strJSON = string(byteJSON)
	return
}

// check if needed dialog exist
func isDialogExist(dialog string) bool {
	//cmd := `ps ux | awk '/` + dialog + `/ && !/awk/ {print $2}'`
	result, err := exec.Command("pgrep", "-la", "-f", dialog).Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(result)) != ""
}
