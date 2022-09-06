// SPDX-FileCopyrightText: 2015 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dbusmenumanager.h"

/*
 * Implementation of interface class DBusMenuManager
 */

DBusMenuManager::DBusMenuManager(QObject *parent)
    : QDBusAbstractInterface(staticServerPath(), staticInterfacePath(), staticInterfaceName(), QDBusConnection::sessionBus(), parent)
{
}

DBusMenuManager::~DBusMenuManager()
{
}

