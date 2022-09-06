// SPDX-FileCopyrightText: 2015 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dbusaccount.h"

/*
 * Implementation of interface class DBusAccount
 */

DBusAccount::DBusAccount(QObject *parent)
    : QDBusAbstractInterface(staticService(), staticInterfacePath(), staticInterfaceName(), QDBusConnection::systemBus(), parent)
{
    QDBusConnection::systemBus().connect(this->service(), this->path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged","sa{sv}as", this, SLOT(__propertyChanged__(QDBusMessage)));
}

DBusAccount::~DBusAccount()
{
    QDBusConnection::systemBus().disconnect(service(), path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged",  "sa{sv}as", this, SLOT(propertyChanged(QDBusMessage)));
}

