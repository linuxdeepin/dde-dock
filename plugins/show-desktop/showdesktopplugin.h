// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SHOWDESKTOPPLUGIN_H
#define SHOWDESKTOPPLUGIN_H

#include "pluginsiteminterface.h"
#include "showdesktopwidget.h"

#include <QLabel>
namespace Dock{
class TipsWidget;
}
class ShowDesktopPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "show-desktop.json")

public:
    explicit ShowDesktopPlugin(QObject *parent = nullptr);

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
    PluginType type() override;

private:
    void updateVisible();
    void loadPlugin();
    void refreshPluginItemsVisible();

private:
    bool m_pluginLoaded;

    QScopedPointer<ShowDesktopWidget> m_showDesktopWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // SHOWDESKTOPPLUGIN_H
