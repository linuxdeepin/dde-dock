/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

DeviceItem::DeviceItem(const QString &path)
    : QWidget(nullptr),
      m_devicePath(path),

      m_networkManager(NetworkManager::instance(this))
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
    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> enable;
    enable["itemId"] = "enable";
    if (!enabled())
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

QWidget *DeviceItem::itemPopup()
{
    return nullptr;
}

void DeviceItem::invokeMenuItem(const QString &menuId)
{
    if (menuId == "settings")
        QProcess::startDetached("dde-control-center", QStringList() << "network");
    else if (menuId == "enable")
        setEnabled(!enabled());
}

bool DeviceItem::enabled() const
{
    return m_networkManager->deviceEnabled(m_devicePath);
}

void DeviceItem::setEnabled(const bool enable)
{
    m_networkManager->setDeviceEnabled(m_devicePath, enable);
}

QWidget *DeviceItem::itemApplet()
{
    return nullptr;
}
