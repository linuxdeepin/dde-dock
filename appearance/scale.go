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

package appearance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"pkg.deepin.io/dde/api/userenv"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
)

func (m *Manager) getScaleFactor() float64 {
	scaleFactor, err := m.xSettings.GetScaleFactor(dbus.FlagNoAutoStart)
	if err != nil {
		logger.Warning("failed to get scale factor:", err)
		scaleFactor = 1.0
	}
	return scaleFactor
}

func (m *Manager) setScaleFactor(scale float64) error {
	err := m.xSettings.SetScaleFactor(dbus.FlagNoAutoStart, scale)
	if err == nil {
		sendNotify(gettext.Tr("Display scaling"),
			gettext.Tr("Setting display scaling"), "dialog-window-scale")
	} else {
		logger.Warning("failed to set scale factor:", err)
	}
	return err
}

var _notifier *notify.Notification

func sendNotify(summary, body, icon string) {
	if _notifier == nil {
		notify.Init("dde-daemon")
		_notifier = notify.NewNotification(summary, body, icon)
	} else {
		_notifier.Update(summary, body, icon)
	}
	err := _notifier.Show()
	if err != nil {
		logger.Warning("Failed to send notify:", summary, body)
	}
}

func (m *Manager) handleSetScaleFactorDone() {
	sendNotify(gettext.Tr("Set successfully"),
		gettext.Tr("Please log back in to view the changes"), "dialog-window-scale")
}

const (
	envQtScaleFactor           = "QT_SCALE_FACTOR"
	envQtScreenScaleFactors    = "QT_SCREEN_SCALE_FACTORS"
	envQtAutoScreenScaleFactor = "QT_AUTO_SCREEN_SCALE_FACTOR"
	envQtFontDPI               = "QT_FONT_DPI"
)

func (m *Manager) setScreenScaleFactors(factors map[string]float64) error {
	individualScalingJoined := joinIndividualScaling(factors)
	ok := m.xSettingsGs.SetString(gsKeyIndividualScaling, individualScalingJoined)
	if !ok {
		return errors.New("failed to save individual scaling to gsettings")
	}

	if len(factors) == 0 {
		scaleFactor := m.xSettingsGs.GetDouble("scale-factor")
		err := userenv.Modify(func(v map[string]string) {
			if scaleFactor == 1 {
				delete(v, envQtScaleFactor)
			} else {
				v[envQtScaleFactor] = fmt.Sprintf("%.2f", scaleFactor)
			}

			delete(v, envQtScreenScaleFactors)
			delete(v, envQtAutoScreenScaleFactor)
			delete(v, envQtFontDPI)
		})
		if err != nil {
			return err
		}
	} else {
		primaryScreenName, err := getPrimaryScreenName(m.xConn)
		if err != nil {
			return err
		}
		primaryScreenFactor, ok := factors[primaryScreenName]
		if !ok {
			primaryScreenFactor = 1
		}

		err = userenv.Modify(func(v map[string]string) {
			delete(v, envQtScaleFactor)
			v[envQtScreenScaleFactors] = individualScalingJoined
			v[envQtAutoScreenScaleFactor] = "0"
			if primaryScreenFactor != 1 {
				fontDPI := int(96 * primaryScreenFactor)
				v[envQtFontDPI] = strconv.Itoa(fontDPI)
			} else {
				delete(v, envQtFontDPI)
			}
		})
		if err != nil {
			return err
		}
	}

	err := m.greeter.SetScreenScaleFactors(0, factors)
	if err != nil {
		logger.Warning(err)
	}

	return nil
}

func (m *Manager) getScreenScaleFactors() map[string]float64 {
	str := m.xSettingsGs.GetString(gsKeyIndividualScaling)
	return parseIndividualScaling(str)
}

func parseIndividualScaling(str string) map[string]float64 {
	pairs := strings.Split(str, ";")
	result := make(map[string]float64)
	for _, value := range pairs {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) != 2 {
			continue
		}

		value, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			logger.Warning(err)
			continue
		}

		result[kv[0]] = value
	}

	return result
}

func joinIndividualScaling(v map[string]float64) string {
	pairs := make([]string, len(v))
	idx := 0
	for key, value := range v {
		pairs[idx] = fmt.Sprintf("%s=%.2f", key, value)
		idx++
	}
	return strings.Join(pairs, ";")
}

func getPrimaryScreenName(xConn *x.Conn) (string, error) {
	rootWin := xConn.GetDefaultScreen().Root
	getPrimaryReply, err := randr.GetOutputPrimary(xConn, rootWin).Reply(xConn)
	if err != nil {
		return "", err
	}
	outputInfo, err := randr.GetOutputInfo(xConn, getPrimaryReply.Output,
		x.CurrentTime).Reply(xConn)
	if err != nil {
		return "", err
	}
	return outputInfo.Name, nil
}
