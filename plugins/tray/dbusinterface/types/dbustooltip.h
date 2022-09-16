/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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
