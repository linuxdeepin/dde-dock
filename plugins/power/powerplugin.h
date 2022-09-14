// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef POWERPLUGIN_H
#define POWERPLUGIN_H

#include "pluginsiteminterface.h"
#include "powerstatuswidget.h"
#include "dbus/dbuspower.h"

#include <com_deepin_system_systempower.h>

#include <QLabel>

using SystemPowerInter = com::deepin::system::Power;
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
    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void refreshIcon(const QString &itemKey) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void pluginSettingsChanged() override;

private:
    void updateBatteryVisible();
    void loadPlugin();
    void refreshPluginItemsVisible();
    void onGSettingsChanged(const QString &key);
    void refreshTipsData();

private:
    bool m_pluginLoaded;
    bool m_showTimeToFull;

    QScopedPointer<PowerStatusWidget> m_powerStatusWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;

    SystemPowerInter *m_systemPowerInter;
    DBusPower *m_powerInter;
    QTimer *m_preChargeTimer;
};

#endif // POWERPLUGIN_H
