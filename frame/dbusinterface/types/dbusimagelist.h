// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSIMAGELIST_H
#define DBUSIMAGELIST_H

#include <QDBusMetaType>
#include <QRect>
#include <QList>

struct DBusImage
{
    int width;
    int height;
    QByteArray pixels;
};
Q_DECLARE_METATYPE(DBusImage)

typedef QList<DBusImage> DBusImageList;
Q_DECLARE_METATYPE(DBusImageList)

QDBusArgument &operator<<(QDBusArgument&, const DBusImage&);
const QDBusArgument &operator>>(const QDBusArgument&, DBusImage&);

bool operator ==(const DBusImage&, const DBusImage&);
bool operator !=(const DBusImage&, const DBusImage&);

void registerDBusImageListMetaType();

#endif // DBUSIMAGELIST_H
