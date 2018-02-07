/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef DISKINFO_H
#define DISKINFO_H

#include <QString>
#include <QDataStream>
#include <QDebug>
#include <QtDBus>

class DiskInfo
{
public:
    DiskInfo();
    static void registerMetaType();

    friend QDebug operator<<(QDebug debug, const DiskInfo &info);
    friend QDBusArgument &operator<<(QDBusArgument &args, const DiskInfo &info);
    friend QDataStream &operator<<(QDataStream &args, const DiskInfo &info);
    friend const QDBusArgument &operator>>(const QDBusArgument &args, DiskInfo &info);
    friend const QDataStream &operator>>(QDataStream &args, DiskInfo &info);

public:
    QString m_id;
    QString m_name;
    QString m_type;
    QString m_path;
    QString m_mountPoint;
    QString m_icon;

    bool m_unmountable;
    bool m_ejectable;

    quint64 m_usedSize;
    quint64 m_totalSize;
};

typedef QList<DiskInfo> DiskInfoList;

Q_DECLARE_METATYPE(DiskInfo)
Q_DECLARE_METATYPE(DiskInfoList)

#endif // DISKINFO_H
