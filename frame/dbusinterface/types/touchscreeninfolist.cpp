/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     quezhiyong <quezhiyong@uniontech.com>
 *
 * Maintainer: quezhiyong <quezhiyong@uniontech.com>
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

#include "touchscreeninfolist.h"

QDBusArgument &operator<<(QDBusArgument &arg, const TouchscreenInfo &info)
{
    arg.beginStructure();
    arg << info.id << info.name << info.deviceNode << info.serialNumber;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, TouchscreenInfo &info)
{
    arg.beginStructure();
    arg >> info.id >> info.name >> info.deviceNode >> info.serialNumber;
    arg.endStructure();

    return arg;
}

bool TouchscreenInfo::operator==(const TouchscreenInfo &info)
{
    return id == info.id && name == info.name && deviceNode == info.deviceNode && serialNumber == info.serialNumber;
}

void registerTouchscreenInfoMetaType()
{
    qRegisterMetaType<TouchscreenInfo>("TouchscreenInfo");
    qDBusRegisterMetaType<TouchscreenInfo>();
}

void registerTouchscreenInfoListMetaType()
{
    registerTouchscreenInfoMetaType();

    qRegisterMetaType<TouchscreenInfoList>("TouchscreenInfoList");
    qDBusRegisterMetaType<TouchscreenInfoList>();
}
