// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef BLUETOOTHPLUGIN_H
#define BLUETOOTHPLUGIN_H

#include "pluginsiteminterface.h"
#include "bluetoothitem.h"

#include <QScopedPointer>

class BluetoothMainWidget;

class AdaptersManager;

class BluetoothPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "bluetooth.json")

public:
    explicit BluetoothPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    bool pluginIsAllowDisable() override { return true; }
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void refreshIcon(const QString &itemKey) override;

    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginMode status() const override;
    QString description() const override;
    PluginFlags flags() const override;

private:
    AdaptersManager *m_adapterManager;
    QScopedPointer<BluetoothItem> m_bluetoothItem;
    QScopedPointer<BluetoothMainWidget> m_bluetoothWidget;
    bool m_enableState = true;
};

#endif // BLUETOOTHPLUGIN_H
