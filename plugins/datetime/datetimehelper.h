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
#ifndef DATETIMEHELPER_H
#define DATETIMEHELPER_H

#include <QString>

class DateTimeHelper
{
public:
    enum WeekdayFormat {
         DefaultWeekdayFormat = 0,    //    星期一 (默认项）
         WeekdayFormat1,              //    周一
    };

    enum ShortDateFormat {
         DefaultShortDateFormat = 0,  //    2020/4/5（默认项）
         ShortDateFormat1,            //    2020-4-5
         ShortDateFormat2,            //    2020.4.5
         ShortDateFormat3,            //    2020/04/05
         ShortDateFormat4,            //    2020-04-05
         ShortDateFormat5,            //    2020.04.05
         ShortDateFormat6,            //    20/4/5
         ShortDateFormat7,            //    20-4-5
         ShortDateFormat8,            //    20.4.5
    };

    enum LongDateFormat {
         DefaultLongDateFormat = 0,   //    2020年4月5日（默认项）
         LongDateFormat1,             //    2020年4月5日 星期三
         LongDateFormat2,             //    星期三 020年4月5日
    };

    enum ShortTimeFormat {
         DefaultShortTimeFormat = 0,  //    9:40（默认项)
         ShortTimeFormat1,            //    09:40
    };

    enum LongTimeFormat {
         DefaultLongTimeFormat = 0,   //    9:40:07（默认项）
         LongTimeFormat1,             //    09:40:07
    };

    static QString WeekDayFormatString(int weekdayFormat);
    static QString ShortDateFormatString(int shortDataFormat);
    static QString LongDateFormatString(int longDateFormat);
    static QString ShortTimeFormatString(int shortTimeFormat);
    static QString LongTimeFormatString(int longTImeFormat);
};

#endif
