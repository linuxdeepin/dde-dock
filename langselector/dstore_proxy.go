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

package langselector

import (
	"pkg.deepin.io/dde/api/i18n_dependent"
	"pkg.deepin.io/dde/daemon/dstore"
	"strings"
	"time"
)

const (
	timeoutDuration = time.Second * 60 * 30
)

func installI18nDependent(locale string) error {
	installInfos, removeInfos, err := i18n_dependent.GetByLocale(locale)
	if err != nil {
		return err
	}
	logger.Debug("Install package infos:", installInfos)
	logger.Debug("Remove package infos:", removeInfos)

	installPkgs := getMissingPackages(installInfos, false)
	logger.Info("Need to install:", installPkgs)
	if err := installPackages(installPkgs); err != nil {
		return err
	}

	removePkgs := getMissingPackages(removeInfos, true)
	logger.Info("Need to remove:", removePkgs)
	if err := removePackages(removePkgs); err != nil {
		return err
	}
	return nil
}

func getMissingPackages(infos i18n_dependent.DependentInfos, removed bool) []string {
	var set = make(map[string]bool)
	for _, info := range infos {
		if len(info.Dependent) != 0 && !dstore.IsInstalled(info.Dependent) {
			continue
		}

		list := filterPackages(info.Packages, removed)
		for _, v := range list {
			set[v] = true
		}
	}

	var pkgs []string
	for k, _ := range set {
		pkgs = append(pkgs, k)
	}
	return pkgs
}

func filterPackages(pkgs []string, removed bool) []string {
	var list []string
	for _, pkg := range pkgs {
		if !dstore.IsExists(pkg) {
			continue
		}

		if removed && !dstore.IsInstalled(pkg) {
			continue
		}

		if !removed && dstore.IsInstalled(pkg) {
			continue
		}

		list = append(list, pkg)
	}
	return list
}

func installPackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}

	return dstore.NewDInstallTransaction(strings.Join(pkgs, " "),
		"", timeoutDuration).Exec()
}

func removePackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}

	return dstore.NewDUninstallTransaction(strings.Join(pkgs, " "),
		true, timeoutDuration).Exec()
}
