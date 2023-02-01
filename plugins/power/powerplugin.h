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

#include "org_deepin_dde_systempower1.h"

#include <QLabel>

using SystemPowerInter = org::deepin::dde::Power1;
namespace Dock {
class TipsWidget;
}
class PowerPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "power.json")

public:
    explicit PowerPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    bool pluginIsAllowDisable() override { return true; }
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void refreshIcon(const QString &itemKey) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginFlags flags() const override;

private:
    void updateBatteryVisible();
    void loadPlugin();
    void onGSettingsChanged(const QString &key);
    void initUi();
    void initConnection();

private slots:
    void refreshTipsData();
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    bool m_pluginLoaded;
    bool m_showTimeToFull;

    QScopedPointer<PowerStatusWidget> m_powerStatusWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;

    SystemPowerInter *m_systemPowerInter;
    DBusPower *m_powerInter;
    QTimer *m_preChargeTimer;
    QWidget *m_quickPanel;
    QLabel *m_imageLabel;
    QLabel *m_labelText;
};

#endif // POWERPLUGIN_H
