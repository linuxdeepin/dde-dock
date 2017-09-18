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

#ifndef SHUTDOWNPLUGIN_H
#define SHUTDOWNPLUGIN_H

#include "pluginsiteminterface.h"
#include "pluginwidget.h"
#include "powerstatuswidget.h"
#include "dbus/dbuspower.h"

#include <QLabel>

#define BATTERY_DISCHARED   2
#define BATTERY_FULL        4

class ShutdownPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "shutdown.json")

public:
    explicit ShutdownPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);

    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);
    const QString itemCommand(const QString &itemKey);
    const QString itemContextMenu(const QString &itemKey);
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked);
    void displayModeChanged(const Dock::DisplayMode displayMode);

private:
    void updateBatteryVisible();
    void requestContextMenu(const QString &itemKey);
    void delayLoader();

private:
    PluginWidget *m_shutdownWidget;
    PowerStatusWidget *m_powerStatusWidget;
    QLabel *m_tipsLabel;

    DBusPower *m_powerInter;
};

#endif // SHUTDOWNPLUGIN_H
