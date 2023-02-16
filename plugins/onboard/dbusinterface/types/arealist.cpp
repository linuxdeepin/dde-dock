// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "arealist.h"

bool MonitRect::operator ==(const MonitRect &rect)
{
    return x1 == rect.x1 && y1 == rect.y1 && x2 == rect.x2 && y2 == rect.y2;
}

QDBusArgument &operator<<(QDBusArgument &arg, const MonitRect &rect)
{
    arg.beginStructure();
    arg << rect.x1 << rect.y1 << rect.x2 << rect.y2;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, MonitRect &rect)
{
    arg.beginStructure();
    arg >> rect.x1 >> rect.y1 >> rect.x2 >> rect.y2;
    arg.endStructure();

    return arg;
}

void registerAreaListMetaType()
{
    qRegisterMetaType<MonitRect>("MonitRect");
    qDBusRegisterMetaType<MonitRect>();

    qRegisterMetaType<AreaList>("AreaList");
    qDBusRegisterMetaType<AreaList>();
}
