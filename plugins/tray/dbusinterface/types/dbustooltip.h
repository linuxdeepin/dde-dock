// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSTOOLTIP_H
#define DBUSTOOLTIP_H

#include "dbusimagelist.h"

#include <QDBusMetaType>
#include <QRect>
#include <QList>

struct DBusToolTip
{
    QString iconName;
    DBusImageList iconPixmap;
    QString title;
    QString description;
};
Q_DECLARE_METATYPE(DBusToolTip)

QDBusArgument &operator<<(QDBusArgument&, const DBusToolTip&);
const QDBusArgument &operator>>(const QDBusArgument&, DBusToolTip&);

bool operator ==(const DBusToolTip&, const DBusToolTip&);
bool operator !=(const DBusToolTip&, const DBusToolTip&);

void registerDBusToolTipMetaType();

#endif // DBUSTOOLTIP_H
