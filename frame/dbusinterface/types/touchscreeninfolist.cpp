// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
