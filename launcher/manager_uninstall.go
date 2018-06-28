/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package launcher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"pkg.deepin.io/lib/appinfo/desktopappinfo"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
)

var chromeShortcurtExecRegexp = regexp.MustCompile(`google-chrome.*--app-id=`)

const flatpakBin = "flatpak"

func isChromeShortcut(item *Item) bool {
	logger.Debugf("isChromeShortcut item ID: %q, exec: %q", item.ID, item.exec)
	return strings.HasPrefix(item.ID, "chrome-") &&
		chromeShortcurtExecRegexp.MatchString(item.exec)
}

// Simply remove the desktop file
func (m *Manager) uninstallDesktopFile(item *Item) error {
	logger.Debugf("remove desktop file %q", item.Path)
	err := os.Remove(item.Path)
	go m.notifyUninstallDone(item, err == nil)
	return err
}

func (m *Manager) removeAutostart(id string) {
	name := filepath.Base(id) + desktopExt
	file := filepath.Join(basedir.GetUserConfigDir(), "autostart", name)
	logger.Debugf("removeAutostart id: %q, file: %q", id, file)
	os.Remove(file)
}

func (m *Manager) notifyUninstallDone(item *Item, succeed bool) {
	const icon = "deepin-appstore"
	var msg string
	if succeed {
		msg = fmt.Sprintf(Tr("%q removed successfully"), item.Name)
	} else {
		msg = fmt.Sprintf(Tr("Failed to uninstall %q"), item.Name)
	}
	m.sendNotification(msg, "", icon)
}

func isWineApp(appInfo *desktopappinfo.DesktopAppInfo) bool {
	createdBy, _ := appInfo.GetString(desktopappinfo.MainSection, "X-Created-By")
	return strings.HasPrefix(createdBy, "cxoffice-") ||
		strings.Contains(appInfo.GetCommandline(), "env WINEPREFIX=")
}

func (m *Manager) uninstallDeepinWineApp(item *Item) error {
	logger.Debug("uninstallDeepinWineApp", item.Path)
	cmd := exec.Command("/opt/deepinwine/tools/uninstall.sh", item.Path)
	err := cmd.Run()
	go m.notifyUninstallDone(item, err == nil)
	return err
}

type flatpakAppInfo struct {
	name, arch, branch string
}

func parseFlatpakAppCmdline(cmdline string) (*flatpakAppInfo, error) {
	parts := strings.Split(cmdline, " ")

	if !(len(parts) > 0 && filepath.Base(parts[0]) == flatpakBin) {
		return nil, errors.New("flatpak app bin is not " + flatpakBin)
	}
	var name, arch, branch string
	for idx, part := range parts {
		if branch == "" && strings.HasPrefix(part, "--branch=") {
			branch = strings.TrimPrefix(part, "--branch=")
		}

		if arch == "" && strings.HasPrefix(part, "--arch=") {
			arch = strings.TrimPrefix(part, "--arch=")
		}

		if name == "" && idx != 0 && !strings.HasPrefix(part, "--") && strings.Contains(part, ".") {
			name = part
		}

		if branch != "" && arch != "" && name != "" {
			break
		}
	}

	if branch == "" {
		return nil, errors.New("failed to get flatpak app branch")
	}

	if arch == "" {
		return nil, errors.New("failed to get flatpak app arch")
	}

	if name == "" {
		return nil, errors.New("failed to get flatpak app name")
	}

	return &flatpakAppInfo{
		name:   name,
		arch:   arch,
		branch: branch,
	}, nil
}

func isFlatpakApp(appInfo *desktopappinfo.DesktopAppInfo) (bool, *flatpakAppInfo, error) {
	// X-Flatpak maybe store true or flatpak app id
	xFlatpak, _ := appInfo.GetString(desktopappinfo.MainSection, "X-Flatpak")
	if xFlatpak == "" {
		return false, nil, nil
	}

	cmdline := appInfo.GetCommandline()
	fpAppInfo, err := parseFlatpakAppCmdline(cmdline)
	if err != nil {
		return false, nil, err
	}
	return true, fpAppInfo, nil
}

func (m *Manager) uninstallFlatpakApp(item *Item, fpAppInfo *flatpakAppInfo) error {
	logger.Debug("uninstallFlatpakApp", item.Path, fpAppInfo.name)

	homeDir := basedir.GetUserHomeDir()
	if homeDir == "" {
		return errors.New("get home dir failed")
	}

	userInstallation := strings.HasPrefix(item.Path, homeDir)

	sysOrUser := "--user"
	if !userInstallation {
		// system wide installation
		sysOrUser = "--system"
		pkgFile := filepath.Join("/usr/share/deepin-flatpak/app/",
			fpAppInfo.name, fpAppInfo.arch, fpAppInfo.branch, "pkg")
		logger.Debug("pkg file:", pkgFile)
		content, err := ioutil.ReadFile(pkgFile)
		if err == nil {
			pkgName := string(bytes.TrimSpace(content))
			return m.uninstallSystemPackage(item.Name, pkgName)
		}
	}

	ref := fmt.Sprintf("app/%s/%s/%s", fpAppInfo.name, fpAppInfo.arch, fpAppInfo.branch)
	logger.Debug("ref:", ref)
	cmd := exec.Command("flatpak", sysOrUser, "uninstall", ref)
	err := cmd.Run()
	go m.notifyUninstallDone(item, err == nil)
	return err
}

func (m *Manager) uninstall(id string) error {
	item := m.getItemById(id)
	if item == nil {
		logger.Warning("RequestUninstall failed", errorInvalidID)
		return errorInvalidID
	}

	appInfo, err := desktopappinfo.NewDesktopAppInfoFromFile(item.Path)
	if err != nil {
		return err
	}

	// uninstall system package
	if pkg := m.queryPkgName(item.ID); pkg != "" {
		// is pkg installed?
		installed, err := m.lastore.PackageExists(0, pkg)
		if err != nil {
			return err
		}
		if installed {
			return m.uninstallSystemPackage(item.Name, pkg)
		}
	}

	// uninstall flatpak app
	isFlatpakApp, fpAppInfo, err := isFlatpakApp(appInfo)
	if err != nil {
		return err
	}
	if isFlatpakApp {
		logger.Debugf("fpAppInfo: %#v", fpAppInfo)
		return m.uninstallFlatpakApp(item, fpAppInfo)
	}

	// uninstall chrome shortcut
	if isChromeShortcut(item) {
		logger.Debug("item is chrome shortcut")
		return m.uninstallDesktopFile(item)
	}

	// uninstall wine app
	if isWineApp(appInfo) {
		logger.Debug("item is wine app")
		return m.uninstallDeepinWineApp(item)
	}

	return m.uninstallDesktopFile(item)
}

const (
	JobStatusSucceed = "succeed"
	JobStatusFailed  = "failed"
	JobStatusEnd     = "end"
)

func (m *Manager) uninstallSystemPackage(jobName, pkg string) error {
	jobPath, err := m.lastore.RemovePackage(0, jobName, pkg)
	logger.Debugf("uninstallSystemPackage pkg: %q jobPath: %q", pkg, jobPath)
	if err != nil {
		return err
	}
	return m.waitJobDone(string(jobPath))
}

func (m *Manager) waitJobDone(jobPath string) error {
	logger.Debug("waitJobDone", jobPath)
	defer logger.Debug("waitJobDone end")
	return m.monitorJobStatusChange(jobPath)
}

func (m *Manager) monitorJobStatusChange(jobPath string) error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	rule := dbusutil.NewMatchRuleBuilder().ExtPropertiesChanged(jobPath,
		"com.deepin.lastore.Job").Sender(lastoreDBusDest).Build()
	err = rule.AddTo(sysBus)
	if err != nil {
		return err
	}

	ch := make(chan *dbus.Signal, 10)
	sysBus.Signal(ch)

	defer func() {
		sysBus.RemoveSignal(ch)
		err := rule.RemoveFrom(sysBus)
		if err != nil {
			logger.Warning("RemoveMatch failed:", err)
		}
		logger.Debug("monitorJobStatusChange return")
	}()

	for v := range ch {
		if len(v.Body) != 3 {
			continue
		}

		props, _ := v.Body[1].(map[string]dbus.Variant)
		status, ok := props["Status"]
		if !ok {
			continue
		}
		statusStr, _ := status.Value().(string)
		logger.Debug("job status changed", statusStr)
		switch statusStr {
		case JobStatusSucceed:
			return nil
		case JobStatusFailed:
			return errors.New("job Failed")
		case JobStatusEnd:
			return nil
		}
	}
	return nil
}
