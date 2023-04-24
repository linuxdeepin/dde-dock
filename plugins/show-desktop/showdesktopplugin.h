// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
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
    bool pluginIsAllowDisable() override { return true; }
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void refreshIcon(const QString &itemKey) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    PluginType type() override;
    PluginFlags flags() const override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;

private:
    QScopedPointer<ShowDesktopWidget> m_showDesktopWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // SHOWDESKTOPPLUGIN_H
