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

#ifndef POWERPLUGIN_H
#define POWERPLUGIN_H

#include "pluginsiteminterface.h"
#include "powerstatuswidget.h"
#include "dbus/dbuspower.h"
#include "../widgets/tipswidget.h"

#include <com_deepin_system_systempower.h>

#include <QLabel>

using SystemPowerInter = com::deepin::system::Power;

// from https://upower.freedesktop.org/docs/Device.html#Device:State
enum BatteryState
{
    UNKNOWN = 0,
    CHARGING = 1,
    DISCHARGING = 2,
    EMPTY = 3,
    FULLY_CHARGED = 4,
    PENDING_CHARGE = 5,
    PENDING_DISCHARGE = 6
};

class PowerPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)

public:
    explicit PowerPlugin(QObject *parent = 0);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void refreshIcon(const QString &itemKey) Q_DECL_OVERRIDE;
    int itemSortKey(const QString &itemKey) Q_DECL_OVERRIDE;
    void setSortKey(const QString &itemKey, const int order) Q_DECL_OVERRIDE;
    void pluginSettingsChanged() override;

private:
    void updateBatteryVisible();
    void loadPlugin();
    void refreshPluginItemsVisible();
    void refreshTipsData();

private:
    bool m_pluginLoaded;

    PowerStatusWidget *m_powerStatusWidget;
    TipsWidget *m_tipsLabel;

    SystemPowerInter *m_systemPowerInter;
    DBusPower *m_powerInter;
};

#endif // POWERPLUGIN_H
