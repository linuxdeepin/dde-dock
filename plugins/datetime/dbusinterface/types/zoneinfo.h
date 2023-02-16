// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ZONEINFO_H
#define ZONEINFO_H

#include <QDebug>
#include <QDataStream>
#include <QString>
#include <QDBusArgument>
#include <QDBusMetaType>

class ZoneInfo
{
public:
    ZoneInfo();

    friend QDebug operator<<(QDebug argument, const ZoneInfo &info);
    friend QDBusArgument &operator<<(QDBusArgument &argument, const ZoneInfo &info);
    friend QDataStream &operator<<(QDataStream &argument, const ZoneInfo &info);
    friend const QDBusArgument &operator>>(const QDBusArgument &argument, ZoneInfo &info);
    friend const QDataStream &operator>>(QDataStream &argument, ZoneInfo &info);

    bool operator==(const ZoneInfo &what) const;

public:
    inline QString getZoneName() const {return m_zoneName;}
    inline QString getZoneCity() const {return m_zoneCity;}
    inline int getUTCOffset() const {return m_utcOffset;}

private:
    QString m_zoneName;
    QString m_zoneCity;
    int m_utcOffset;
    qint64 i2;
    qint64 i3;
    int i4;
};

Q_DECLARE_METATYPE(ZoneInfo)

void registerZoneInfoMetaType();

#endif // ZONEINFO_H
