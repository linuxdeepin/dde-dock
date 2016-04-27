/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package item

import (
	"fmt"
	"os"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	. "pkg.deepin.io/lib/gettext"
	"regexp"
	"strings"
)

var chromeShortcurtExecRegexp = regexp.MustCompile(`google-chrome.*--app-id=`)

// Check if the item passed is is a chrome shortcut or chrome app.
func itemIsChromeShortcut(i ItemInfo) bool {
	return strings.HasPrefix(string(i.ID()), "chrome-") &&
		chromeShortcurtExecRegexp.MatchString(i.ExecCmd())
}

// Uninstall chrome shortcut or chrome app. It just simply removes
// the desktop file chrome created.
func uninstallFromChromeShortcuts(item ItemInfo) error {
	icon := item.Icon()

	// Simply remove the desktop file will be OK.
	err := os.Remove(item.Path())
	if err == nil {
		Notify(icon, "Launcher",
			fmt.Sprintf(Tr("%q removed successfully."), item.Name()))
	}

	return err
}
