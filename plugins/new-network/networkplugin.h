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

#ifndef NETWORKPLUGIN_H
#define NETWORKPLUGIN_H

#include "pluginsiteminterface.h"
#include <item/deviceitem.h>

#include <QSettings>
#include <NetworkWorker>
#include <NetworkModel>

class NetworkPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "network.json")

public:
    explicit NetworkPlugin(QObject *parent = 0);

    const QString pluginName() const;
    const QString pluginDisplayName() const;
    void init(PluginProxyInterface *proxyInter);
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked);
    void refershIcon(const QString &itemKey);
    void pluginStateSwitched();
    bool pluginIsAllowDisable() { return true; }
    bool pluginIsDisable();
    const QString itemCommand(const QString &itemKey);
    const QString itemContextMenu(const QString &itemKey);
    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);
    QWidget *itemPopupApplet(const QString &itemKey);

    int itemSortKey(const QString &itemKey);
    void setSortKey(const QString &itemKey, const int order);

private slots:
    void onDeviceListChanged(const QList<dde::network::NetworkDevice *> devices);
    void contextMenuRequested();

private:
    DeviceItem *itemByPath(const QString &path);

private:
    dde::network::NetworkModel *m_networkModel;
    dde::network::NetworkWorker *m_networkWorker;

    QMap<QString, DeviceItem *> m_itemsMap;
    QSettings m_settings;
};

#endif // NETWORKPLUGIN_H
