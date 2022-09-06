// SPDX-FileCopyrightText: 2015 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dbusmenu.h"

/*
 * Implementation of interface class DBusMenu
 */

DBusMenu::DBusMenu(const QString &path, QObject *parent)
    : QDBusAbstractInterface(staticServerPath(), path, staticInterfaceName(), QDBusConnection::sessionBus(), parent)
{
}

DBusMenu::~DBusMenu()
{
}

