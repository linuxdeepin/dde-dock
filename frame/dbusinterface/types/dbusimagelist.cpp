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

#include "dbusimagelist.h"

QDBusArgument &operator<<(QDBusArgument &argument, const DBusImage &image)
{
    argument.beginStructure();
    argument << image.width << image.height << image.pixels;
    argument.endStructure();
    return argument;
}

const QDBusArgument &operator>>(const QDBusArgument &argument, DBusImage &image)
{
    argument.beginStructure();
    argument >> image.width >> image.height >> image.pixels;
    argument.endStructure();
    return argument;
}

void registerDBusImageListMetaType()
{
    qRegisterMetaType<DBusImage>("DBusImage");
    qDBusRegisterMetaType<DBusImage>();

    qRegisterMetaType<DBusImageList>("DBusImageList");
    qDBusRegisterMetaType<DBusImageList>();
}

bool operator ==(const DBusImage &a, const DBusImage &b)
{
    return a.width == b.width
            && a.height == b.height
            && a.pixels == b.pixels;
}

bool operator !=(const DBusImage &a, const DBusImage &b)
{
    return !(a == b);
}
