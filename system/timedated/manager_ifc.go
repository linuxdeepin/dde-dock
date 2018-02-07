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

package timedated

import (
	"pkg.deepin.io/lib/dbus"
)

// SetTime set the current time and date,
// pass a value of microseconds since 1 Jan 1970 UTC
func (m *Manager) SetTime(dmsg dbus.DMessage, usec int64, relative bool) error {
	err := m.checkAuthorization("SetTime",
		Tr("Authentication is required to set the system time."),
		dmsg.GetSenderPID())
	if err != nil {
		return err
	}

	// TODO: check usec validity
	return m.core.SetTime(usec, relative, false)
}

// SetTimezone set the system time zone, the value from /usr/share/zoneinfo/zone.tab
func (m *Manager) SetTimezone(dmsg dbus.DMessage, timezone string) error {
	err := m.checkAuthorization("SetTimezone",
		Tr("Authentication is required to set the system timezone."),
		dmsg.GetSenderPID())
	if err != nil {
		return err
	}

	// TODO: check timezone validity
	if m.core.Timezone.Get() == timezone {
		return nil
	}
	return m.core.SetTimezone(timezone, false)
}

// SetLocalRTC to control whether the RTC is the local time or UTC.
func (m *Manager) SetLocalRTC(dmsg dbus.DMessage, enabled bool, fixSystem bool) error {
	err := m.checkAuthorization("SetLocalRTC",
		Tr("Authentication is required to control whether the RTC stores the local or UTC time."),
		dmsg.GetSenderPID())
	if err != nil {
		return err
	}

	if m.core.LocalRTC.Get() == enabled {
		return nil
	}
	return m.core.SetLocalRTC(enabled, fixSystem, false)
}

// SetNTP to control whether the system clock is synchronized with the network
func (m *Manager) SetNTP(dmsg dbus.DMessage, enabled bool) error {
	err := m.checkAuthorization("SetNTP",
		Tr("Authentication is required to control whether network time synchronization shall be enabled."),
		dmsg.GetSenderPID())
	if err != nil {
		return err
	}

	if m.core.NTP.Get() == enabled {
		return nil
	}
	return m.core.SetNTP(enabled, false)
}

func Tr(str string) string {
	return str
}
