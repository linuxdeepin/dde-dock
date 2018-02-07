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

package langselector

import (
	"strings"
	"time"

	"pkg.deepin.io/dde/api/language_support"
	"pkg.deepin.io/dde/daemon/dstore"
)

const (
	timeoutDuration = time.Second * 60 * 30
)

func installLangSupportPackages(locale string) error {
	logger.Debug("install language support packages for locale", locale)
	ls, err := language_support.NewLanguageSupport()
	if err != nil {
		return err
	}

	pkgs := ls.ByLocale(locale, false)
	ls.Destroy()
	logger.Info("need to install:", pkgs)
	return installPackages(pkgs)
}

func installPackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}

	return dstore.NewDInstallTransaction(strings.Join(pkgs, " "),
		"", timeoutDuration).Exec()
}
