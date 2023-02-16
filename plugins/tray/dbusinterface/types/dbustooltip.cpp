// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
