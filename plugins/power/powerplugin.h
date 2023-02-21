// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef POWERPLUGIN_H
#define POWERPLUGIN_H

#include "pluginsiteminterface.h"
#include "powerstatuswidget.h"
#include "dbus/dbuspower.h"

#include "org_deepin_dde_systempower1.h"

#include <QLabel>

using SystemPowerInter = org::deepin::dde::Power1;

DCORE_BEGIN_NAMESPACE
class DConfig;
DCORE_END_NAMESPACE

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
    Dtk::Core::DConfig *m_dconfig; // 配置
    QTimer *m_preChargeTimer;
    QWidget *m_quickPanel;
    QLabel *m_imageLabel;
    QLabel *m_labelText;
};

#endif // POWERPLUGIN_H
