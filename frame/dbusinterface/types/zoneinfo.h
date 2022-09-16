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
