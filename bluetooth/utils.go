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
	"io/ioutil"
	"pkg.deepin.io/lib/procfs"
	"strconv"
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

// find process
func checkProcessExists(processName string) bool {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		logger.Warningf("read proc failed,err:%v", err)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		pid, err := strconv.ParseUint(f.Name(), 10, 32)
		if err != nil {
			continue
		}

		process := procfs.Process(pid)
		executablePath, err := process.Exe()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		//if !fullpath {
		//	executablePath = filepath.Base(executablePath)
		//}
		if executablePath == processName {
			return true
		}
	}

	return false
}
