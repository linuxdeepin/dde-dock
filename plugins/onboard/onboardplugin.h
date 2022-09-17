// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ONBOARDPLUGIN_H
#define ONBOARDPLUGIN_H

#include "pluginsiteminterface.h"
#include "onboarditem.h"

#include <QLabel>
#include <QScopedPointer>

#include <com_deepin_dde_daemon_dock.h>
#include <com_deepin_dde_daemon_dock_entry.h>

namespace Dock {
class TipsWidget;
}
class OnboardPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "onboard.json")

public:
    explicit OnboardPlugin(QObject *parent = nullptr);

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
    void displayModeChanged(const Dock::DisplayMode displayMode) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    void pluginSettingsChanged() override;

private:
    void loadPlugin();
    void refreshPluginItemsVisible();

private:
    bool m_pluginLoaded;
    bool m_startupState;

    QScopedPointer<OnboardItem> m_onboardItem;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // ONBOARDPLUGIN_H
