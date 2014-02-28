/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"os"
	"strconv"
	"strings"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

func unquoteString(str string) string {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		s, _ := strconv.Unquote(str)
		return s
	} else if strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`) {
		return str[1 : len(str)-1]
	}
	return str
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isFileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// Get all screen's best resolution and choose a smaller one for there
// is no screen is primary.
func getPrimaryScreenBestResolution() (w uint16, h uint16) {
	w, h = 1024, 768 // default value

	X, err := xgb.NewConn()
	if err != nil {
		return
	}
	err = randr.Init(X)
	if err != nil {
		return
	}
	_, err = randr.QueryVersion(X, 1, 4).Reply()
	if err != nil {
		return
	}
	Root := xproto.Setup(X).DefaultScreen(X).Root
	resources, err := randr.GetScreenResources(X, Root).Reply()
	if err != nil {
		return
	}

	bestModes := make([]uint32, 0)
	for _, output := range resources.Outputs {
		reply, err := randr.GetOutputInfo(X, output, 0).Reply()
		if err == nil && reply.NumModes > 1 {
			bestModes = append(bestModes, uint32(reply.Modes[0]))
		}
	}

	w, h = 0, 0
	for _, m := range resources.Modes {
		for _, id := range bestModes {
			if id == m.Id {
				bw, bh := m.Width, m.Height
				if w*h == 0 {
					w, h = bw, bh
				} else if bw*bh < w*h {
					w, h = bw, bh
				}
			}
		}
	}

	logger.Info("primary screen's best resolution is %dx%d", w, h)
	return
}
