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
#include "item/deviceitem.h"

#include <NetworkWorker>
#include <NetworkModel>

#define NETWORK_KEY "network-item-key"

class NetworkItem;
class NetworkPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "network.json")

public:
    explicit NetworkPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void refreshIcon(const QString &itemKey) override;
    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    void pluginSettingsChanged() override;

//    static bool isConnectivity();

private slots:
    void onDeviceListChanged(const QList<dde::network::NetworkDevice *> devices);
    void onActivateConnection(const QString &path, const QString &uuid);

private:
    void loadPlugin();
    void refreshPluginItemsVisible();

private:
    QScopedPointer<dde::network::NetworkModel> m_networkModel;
    QScopedPointer<dde::network::NetworkWorker> m_networkWorker;

    QScopedPointer<NetworkItem> m_networkItem;

    bool m_hasDevice;
};

#endif // NETWORKPLUGIN_H
