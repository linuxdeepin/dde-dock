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
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"pkg.deepin.io/dde/api/inhibit_hint"
	"pkg.deepin.io/dde/daemon/grub_common"
	"pkg.deepin.io/lib/dbusutil"
)

var _g *Grub2

func RunAsDaemon() {
	allowNoCheckAuth()
	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal("failed to new system service", err)
	}
	_g = NewGrub2(service)
	ihObj := inhibit_hint.New("lastore-daemon")
	ihObj.SetName("Control Center")
	ihObj.SetIcon("preferences-system")

	err = service.Export(dbusPath, _g)
	if err != nil {
		logger.Fatal("failed to export grub2:", err)
	}

	err = service.Export(themeDBusPath, _g.theme)
	if err != nil {
		logger.Fatal("failed to export grub2 theme:", err)
	}

	err = ihObj.Export(service)
	if err != nil {
		logger.Warning("failed to export inhibit hint:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(5*time.Minute, _g.canSafelyExit)
	service.Wait()
}

func PrepareGfxmodeDetect() error {
	params, err := grub_common.LoadGrubParams()
	if err != nil {
		logger.Warning(err)
	}

	if grub_common.InGfxmodeDetectionMode(params) {
		return errors.New("already in detection mode")
	}

	gfxmodes, err := grub_common.GetGfxmodesFromXRandr()
	if err != nil {
		return err
	}
	gfxmodes.SortDesc()
	logger.Debug("gfxmodes:", gfxmodes)
	gfxmodesStr := joinGfxmodesForDetect(gfxmodes)
	getModifyFuncPrepareGfxmodeDetect(gfxmodesStr)(params)

	err = ioutil.WriteFile(grub_common.GfxmodeDetectReadyPath, nil, 0644)
	if err != nil {
		return err
	}

	err = writeGrubParams(params)
	if err != nil {
		return err
	}

	cmd := exec.Command(adjustThemeCmd, "-fallback-only")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Warning("failed to adjust theme:", err)
	}

	return nil
}
