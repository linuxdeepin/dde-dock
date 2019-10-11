/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "deviceitem.h"

#include <DDBusSender>
#include <QJsonDocument>

using namespace dde::network;

DeviceItem::DeviceItem(dde::network::NetworkDevice *device)
    : QWidget(nullptr),

      m_device(device),
      m_path(device->path())
{
}

QSize DeviceItem::sizeHint() const
{
    return QSize(26, 26);
}

const QString DeviceItem::itemCommand() const
{
    return QString();
}

const QString DeviceItem::itemContextMenu()
{
    if (m_device.isNull()) {
        return QString();
    }

    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> enable;
    enable["itemId"] = "enable";
    if (!m_device->enabled())
        enable["itemText"] = tr("Enable network");
    else
        enable["itemText"] = tr("Disable network");
    enable["isActive"] = true;
    items.push_back(enable);

    QMap<QString, QVariant> settings;
    settings["itemId"] = "settings";
    settings["itemText"] = tr("Network settings");
    settings["isActive"] = true;
    items.push_back(settings);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

QWidget *DeviceItem::itemTips()
{
    return nullptr;
}

void DeviceItem::invokeMenuItem(const QString &menuId)
{
    if (m_device.isNull()) {
        return;
    }

    if (menuId == "settings")
        //QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:network\"");
        DDBusSender()
                .service("com.deepin.dde.ControlCenter")
                .interface("com.deepin.dde.ControlCenter")
                .path("/com/deepin/dde/ControlCenter")
                .method("ShowPage")
                .arg(QString("network"))
                .arg(m_pageName)
                .call();

    else if (menuId == "enable")
        Q_EMIT requestSetDeviceEnable(m_path, !m_device->enabled());
}

QWidget *DeviceItem::itemApplet()
{
    return nullptr;
}
