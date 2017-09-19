/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package appearance

import (
	"io/ioutil"
	"os"
	"strings"
)

var pamEnvFile = os.Getenv("HOME") + "/.pam_environment"

func writeKeyToEnvFile(key, value, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var lines = strings.Split(string(content), "\n")
	var idx = -1
	for i, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		if !strings.Contains(line, key+"=") {
			continue
		}

		if line == key+"="+value {
			return nil
		}
		idx = i
		break
	}

	if idx != -1 {
		lines[idx] = key + "=" + value
	} else {
		lines[len(lines)-1] = key + "=" + value
	}
	return ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
