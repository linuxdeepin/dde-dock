// Copyright (C) 2015 Digia Plc and/or its subsidiary(-ies).
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dbustraymanager.h"

/*
 * Implementation of interface class DBusTrayManager
 */

DBusTrayManager::DBusTrayManager(QObject *parent)
    : QDBusAbstractInterface("org.deepin.dde.TrayManager1", "/org/deepin/dde/TrayManager1", staticInterfaceName(), QDBusConnection::sessionBus(), parent)
{
    qRegisterMetaType<TrayList>("TrayList");
    qDBusRegisterMetaType<TrayList>();

    QDBusConnection::sessionBus().connect(this->service(), this->path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged","sa{sv}as", this, SLOT(__propertyChanged__(QDBusMessage)));
}

DBusTrayManager::~DBusTrayManager()
{
    QDBusConnection::sessionBus().disconnect(service(), path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged",  "sa{sv}as", this, SLOT(propertyChanged(QDBusMessage)));
}

