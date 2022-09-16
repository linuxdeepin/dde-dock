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

#include "touchscreeninfolist_v2.h"

QDBusArgument &operator<<(QDBusArgument &arg, const TouchscreenInfo_V2 &info)
{
    arg.beginStructure();
    arg << info.id << info.name << info.deviceNode << info.serialNumber << info.UUID;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, TouchscreenInfo_V2 &info)
{
    arg.beginStructure();
    arg >> info.id >> info.name >> info.deviceNode >> info.serialNumber >> info.UUID;
    arg.endStructure();

    return arg;
}

bool TouchscreenInfo_V2::operator==(const TouchscreenInfo_V2 &info)
{
    return id == info.id && name == info.name && deviceNode == info.deviceNode && serialNumber == info.serialNumber && UUID == info.UUID;
}

void registerTouchscreenInfoV2MetaType()
{
    qRegisterMetaType<TouchscreenInfo_V2>("TouchscreenInfo_V2");
    qDBusRegisterMetaType<TouchscreenInfo_V2>();
}

void registerTouchscreenInfoList_V2MetaType()
{
    registerTouchscreenInfoV2MetaType();

    qRegisterMetaType<TouchscreenInfoList_V2>("TouchscreenInfoList_V2");
    qDBusRegisterMetaType<TouchscreenInfoList_V2>();
}
