/*
 * Copyright (C) 2018 ~ 2018 Deepin Technology Co., Ltd.
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

package debug

import (
	"bufio"
	"fmt"
	"os"
	"pkg.deepin.io/lib/strv"
	"strconv"
	"strings"
)

func getMemoryUsage() (int64, error) {
	filename := fmt.Sprintf("/proc/%d/status", os.Getpid())
	return sumMemByFile(filename)
}

func sumMemByFile(filename string) (int64, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fr.Close()

	var count = 0
	var memSize int64
	var scanner = bufio.NewScanner(fr)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if !strings.Contains(line, "RssAnon:") &&
			!strings.Contains(line, "VmPTE:") &&
			!strings.Contains(line, "VmPMD:") {
			continue
		}

		v, err := getInterge(line)
		if err != nil {
			return 0, err
		}
		memSize += v

		count++
		if count == 3 {
			break
		}
	}

	return memSize, nil
}

func getInterge(line string) (int64, error) {
	list := strings.Split(line, " ")
	list = strv.Strv(list).FilterEmpty()
	if len(list) != 3 {
		return 0, fmt.Errorf("Bad format: %s", line)
	}
	return strconv.ParseInt(list[1], 10, 64)
}
