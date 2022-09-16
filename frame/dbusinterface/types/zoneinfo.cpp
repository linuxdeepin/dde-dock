/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
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

#include "zoneinfo.h"

ZoneInfo::ZoneInfo()
{

}

bool ZoneInfo::operator ==(const ZoneInfo &what) const
{
    // TODO: 这里只判断这两个成员应该就可以了
    return m_zoneName == what.m_zoneName &&
            m_utcOffset == what.m_utcOffset;
}

QDebug operator<<(QDebug argument, const ZoneInfo & info)
{
    argument << info.m_zoneName << ',' << info.m_zoneCity << ',' << info.m_utcOffset << ',';
    argument << info.i2 << ',' << info.i3 << ',' << info.i4 << endl;

    return argument;
}

QDBusArgument &operator<<(QDBusArgument & argument, const ZoneInfo & info)
{
    argument.beginStructure();
    argument << info.m_zoneName << info.m_zoneCity << info.m_utcOffset;
    argument.beginStructure();
    argument << info.i2 << info.i3 << info.i4;
    argument.endStructure();
    argument.endStructure();

    return argument;
}

QDataStream &operator<<(QDataStream & argument, const ZoneInfo & info)
{
    argument << info.m_zoneName << info.m_zoneCity << info.m_utcOffset;
    argument << info.i2 << info.i3 << info.i4;

    return argument;
}

const QDBusArgument &operator>>(const QDBusArgument & argument, ZoneInfo & info)
{
    argument.beginStructure();
    argument >> info.m_zoneName >> info.m_zoneCity >> info.m_utcOffset;
    argument.beginStructure();
    argument >> info.i2 >> info.i3 >> info.i4;
    argument.endStructure();
    argument.endStructure();

    return argument;
}

const QDataStream &operator>>(QDataStream & argument, ZoneInfo & info)
{
    argument >> info.m_zoneName >> info.m_zoneCity >> info.m_utcOffset;
    argument >> info.i2 >> info.i3 >> info.i4;

    return argument;
}

void registerZoneInfoMetaType()
{
    qRegisterMetaType<ZoneInfo>("ZoneInfo");
    qDBusRegisterMetaType<ZoneInfo>();
}
