/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     lzm <lizhongming@uniontech.com>
 *
 * Maintainer: lzm <lizhongming@uniontech.com>
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
#include "datetimehelper.h"

QString DateTimeHelper::WeekDayFormatString(int weekdayFormat)
{
    switch (weekdayFormat) {
        case DefaultWeekdayFormat: return "dddd";
        case WeekdayFormat1:       return "ddd";
        default:                   return "dddd";
    }
}

QString DateTimeHelper::ShortDateFormatString(int shortDataFormat)
{
    switch (shortDataFormat) {
        case DefaultShortDateFormat: return "yyyy/M/d";
        case ShortDateFormat1:       return "yyyy-M-d";
        case ShortDateFormat2:       return "yyyy.M.d";
        case ShortDateFormat3:       return "yyyy/MM/dd";
        case ShortDateFormat4:       return "yyyy-MM-dd";
        case ShortDateFormat5:       return "yyyy.MM.dd";
        case ShortDateFormat6:       return "yy/M/d";
        case ShortDateFormat7:       return "yy-M-d";
        case ShortDateFormat8:       return "yy.M.d";
        default:                     return "yyyy/M/d";
    }
}

QString DateTimeHelper::LongDateFormatString(int longDateFormat)
{
    Q_UNUSED(longDateFormat);
    return "";
}

QString DateTimeHelper::ShortTimeFormatString(int shortTimeFormat)
{
    switch (shortTimeFormat) {
        case DefaultShortTimeFormat:  return "h:mm";
        case ShortTimeFormat1:        return "hh:mm";
        default:                      return "h:mm";
    }
}

QString DateTimeHelper::LongTimeFormatString(int longTimeFormat)
{
    Q_UNUSED(longTimeFormat);
    return "";
}
