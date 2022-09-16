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

#include "dbustooltip.h"

QDBusArgument &operator<<(QDBusArgument &argument, const DBusToolTip &tip)
{
    argument.beginStructure();
    argument << tip.iconName << tip.iconPixmap << tip.title << tip.description;
    argument.endStructure();
    return argument;
}

const QDBusArgument &operator>>(const QDBusArgument &argument, DBusToolTip &tip)
{
    argument.beginStructure();
    argument >> tip.iconName >> tip.iconPixmap >> tip.title >> tip.description;
    argument.endStructure();
    return argument;
}

bool operator ==(const DBusToolTip &a, const DBusToolTip &b)
{
    return a.iconName == b.iconName
            && a.iconPixmap == b.iconPixmap
            && a.title == b.title
            && a.description == b.description;
}

bool operator !=(const DBusToolTip &a, const DBusToolTip &b)
{
    return !(a == b);
}

void registerDBusToolTipMetaType()
{
    qRegisterMetaType<DBusToolTip>("DBusToolTip");
    qDBusRegisterMetaType<DBusToolTip>();
}
