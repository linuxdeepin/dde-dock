/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package accounts

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/encoding/kv"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/utils"
)

const (
	polkitActionUserAdministration     = "com.deepin.daemon.accounts.user-administration"
	polkitActionChangeOwnData          = "com.deepin.daemon.accounts.change-own-user-data"
	polkitActionEnableAutoLogin        = "com.deepin.daemon.accounts.enable-auto-login"
	polkitActionDisableAutoLogin       = "com.deepin.daemon.accounts.disable-auto-login"
	polkitActionEnableNoPasswordLogin  = "com.deepin.daemon.accounts.enable-nopass-login"
	polkitActionDisableNoPasswordLogin = "com.deepin.daemon.accounts.disable-nopass-login"
	polkitActionSetKeyboardLayout      = "com.deepin.daemon.accounts.set-keyboard-layout"

	systemLocaleFile  = "/etc/default/locale"
	systemdLocaleFile = "/etc/locale.conf"
	defaultLocale     = "en_US.UTF-8"

	layoutDelimiter   = ";"
	defaultLayout     = "us" + layoutDelimiter
	defaultLayoutFile = "/etc/default/keyboard"
)

type ErrCodeType int32

const (
	// 未知错误
	ErrCodeUnkown ErrCodeType = iota
	// 权限认证失败
	ErrCodeAuthFailed
	// 执行命令失败
	ErrCodeExecFailed
	// 传入的参数不合法
	ErrCodeParamInvalid
)

func (code ErrCodeType) String() string {
	switch code {
	case ErrCodeUnkown:
		return "Unkown error"
	case ErrCodeAuthFailed:
		return "Policykit authentication failed"
	case ErrCodeExecFailed:
		return "Exec command failed"
	case ErrCodeParamInvalid:
		return "Invalid parameters"
	}

	return "Unkown error"
}

// return icons uris
func getUserStandardIcons() []string {
	imgs, err := graphic.GetImagesInDir(userIconsDir)
	if err != nil {
		return nil
	}

	var icons []string
	for _, img := range imgs {
		img = utils.EncodeURI(img, utils.SCHEME_FILE)
		if strings.Contains(img, "guest") || img == defaultUserIcon {
			continue
		}

		icons = append(icons, img)
	}

	return icons
}

func getNewUserCustomIconDest(username string) string {
	ns := time.Now().UnixNano()
	base := username + "-" + strconv.FormatInt(ns, 36)
	return filepath.Join(userCustomIconsDir, base)
}

func isStrInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}

func isStrvEqual(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	sort.Strings(l1)
	sort.Strings(l2)
	for i, v := range l1 {
		if v != l2[i] {
			return false
		}
	}
	return true
}

func checkAuth(actionId string, sysBusName string) error {
	success, err := checkAuthByPolkit(actionId, sysBusName)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf(ErrCodeAuthFailed.String())
	}

	return nil
}

func checkAuthByPolkit(actionId string, sysBusName string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", sysBusName)

	ret, err := authority.CheckAuthorization(0, subject,
		actionId, nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}

func getDefaultLocale() (locale string) {
	files := [...]string{
		systemLocaleFile,
		systemdLocaleFile,
	}
	for _, file := range files {
		locale = getLocaleFromFile(file)
		if locale != "" {
			// get locale success
			break
		}
	}
	if locale == "" {
		return defaultLocale
	}
	return locale
}

func getLocaleFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	r := kv.NewReader(f)
	r.Delim = '='
	r.Comment = '#'
	r.TrimSpace = kv.TrimLeadingTailingSpace
	for {
		pair, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return ""
		}

		if pair.Key == "LANG" {
			return pair.Value
		}
	}
	return ""
}

func getDefaultLayout() string {
	layout, err := getSystemLayout(defaultLayoutFile)
	if err != nil {
		logger.Warning("failed to get system default layout:", err)
		return defaultLayout
	}
	return layout
}

func getSystemLayout(file string) (string, error) {
	fr, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	var (
		found   int
		layout  string
		variant string

		regLayout  = regexp.MustCompile(`^XKBLAYOUT=`)
		regVariant = regexp.MustCompile(`^XKBVARIANT=`)

		scanner = bufio.NewScanner(fr)
	)
	for scanner.Scan() {
		if found == 2 {
			break
		}

		var line = scanner.Text()
		if regLayout.MatchString(line) {
			layout = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
			continue
		}

		if regVariant.MatchString(line) {
			variant = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
		}
	}

	if len(layout) == 0 {
		return "", fmt.Errorf("not found default layout")
	}

	return layout + layoutDelimiter + variant, nil
}

func getValueFromLine(line, delim string) string {
	array := strings.Split(line, delim)
	if len(array) != 2 {
		return ""
	}

	return strings.TrimSpace(array[1])
}

// Get available shells from '/etc/shells'
func getAvailableShells(file string) []string {
	contents, err := ioutil.ReadFile(file)
	if err != nil || len(contents) == 0 {
		return nil
	}
	var shells []string
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		shells = append(shells, line)
	}
	return shells
}
