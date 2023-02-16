// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    argument << info.i2 << ',' << info.i3 << ',' << info.i4 << Qt::endl;

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
