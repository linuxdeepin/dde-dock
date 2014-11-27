/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

package datetime

import (
	"fmt"
	"pkg.linuxdeepin.com/dde-daemon/datetime/ntp"
	"pkg.linuxdeepin.com/dde-daemon/datetime/timezone"
	. "pkg.linuxdeepin.com/dde-daemon/datetime/utils"
)

var (
	errInvalidDateArgs = fmt.Errorf("Invalid Date Argment")
)

func (date *DateTime) SetDate(year, month, day, hour, min, sec, nsec int32) error {
	if !IsYearValid(year) || !IsMonthValid(month) ||
		!IsDayValid(year, month, day) || !IsHourValid(hour) ||
		!IsMinuteValid(min) || !IsSecondValid(sec) {
		return errInvalidDateArgs
	}

	value := fmt.Sprintf("%v-%v-%v", year, month, day)
	err := SetDate(value)
	if err != nil {
		Warningf(date.logger, "Set Date '%s' Failed: %v", value, err)
		return err
	}

	value = fmt.Sprintf("%v:%v:%v", hour, min, sec)
	err = SetTime(value)
	if err != nil {
		Warningf(date.logger, "Set Date '%s' Failed: %v", value, err)
		return err
	}

	return nil
}

func (date *DateTime) SetTimezone(zone string) error {
	err := timezone.SetTimezone(zone)
	if err != nil {
		Warning(date.logger, err)
		return err
	}

	//date.settings.Reset(gsKeyDSTOffset)
	date.setPropString(&date.CurrentTimezone,
		"CurrentTimezone", zone)
	date.AddUserTimezone(zone)

	if date.NTPEnabled.Get() {
		ntp.Timezone = zone
		go func() {
			ok := ntp.SyncNetworkTime()
			if !ok {
				Warning(date.logger, "Sync Network Time Failed")
			}
		}()
	}

	return nil
}

func (date *DateTime) AddUserTimezone(zone string) {
	if !timezone.IsZoneValid(zone) {
		Warning(date.logger, "Invalid zone:", zone)
		return
	}

	list := date.settings.GetStrv(gsKeyTimezoneList)
	if IsStrInList(zone, list) {
		return
	}

	list = append(list, zone)
	date.settings.SetStrv(gsKeyTimezoneList, list)
}

func (date *DateTime) DeleteUserTimezone(zone string) {
	if !timezone.IsZoneValid(zone) {
		Warning(date.logger, "Invalid zone:", zone)
		return
	}

	list := date.settings.GetStrv(gsKeyTimezoneList)
	var tmp []string
	for _, s := range list {
		if s == zone {
			continue
		}

		tmp = append(tmp, s)
	}

	if len(tmp) == len(list) {
		return
	}

	date.settings.SetStrv(gsKeyTimezoneList, tmp)
}

func (date *DateTime) GetZoneInfo(zone string) (timezone.ZoneInfo, error) {
	info, err := timezone.GetZoneInfo(zone)
	if info == nil {
		return timezone.ZoneInfo{}, err
	}

	return *info, nil
}

func (date *DateTime) GetAllZoneSummary() []timezone.ZoneSummary {
	return timezone.GetZoneSummaryList()
}

func (date *DateTime) Destroy() {
	DestroySetDateTime()
	ntp.FiniNtpModule()
	date.settings.Unref()
}
