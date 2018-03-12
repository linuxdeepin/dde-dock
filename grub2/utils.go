/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package grub2

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"pkg.deepin.io/lib/polkit"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

type InvalidResolutionError struct {
	Resolution string
}

func (err InvalidResolutionError) Error() string {
	return fmt.Sprintf("invalid resolution %q", err.Resolution)
}

func parseResolution(v string) (w, h uint16, err error) {
	if v == "auto" {
		err = errors.New("unknown auto")
		return
	}

	arr := strings.Split(v, "x")
	if len(arr) != 2 {
		err = InvalidResolutionError{v}
		return
	}
	// parse width
	tmpw, err := strconv.ParseUint(arr[0], 10, 16)
	if err != nil {
		err = InvalidResolutionError{v}
		return
	}

	// parse height
	tmph, err := strconv.ParseUint(arr[1], 10, 16)
	if err != nil {
		err = InvalidResolutionError{v}
		return
	}

	w = uint16(tmpw)
	h = uint16(tmph)

	if w == 0 || h == 0 {
		err = InvalidResolutionError{v}
		return
	}

	return
}

func checkResolution(v string) error {
	if v == "auto" {
		return nil
	}

	_, _, err := parseResolution(v)
	return err
}

func getStringIndexInArray(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}

var noCheckAuth bool

func initPolkit() {
	if os.Getenv("NO_CHECK_AUTH") == "1" {
		noCheckAuth = true
		return
	}

	polkit.Init()
}

func checkAuthWithPid(pid uint32) (bool, error) {
	subject := polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
	subject.SetDetail("start-time", uint64(0))
	const actionId = dbusServiceName
	details := make(map[string]string)
	result, err := polkit.CheckAuthorization(subject, actionId, details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}

var errAuthFailed = errors.New("authentication failed")

func getBytesMD5Sum(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func getFileMD5sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}
