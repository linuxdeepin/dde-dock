// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRASHPLUGIN_H
#define TRASHPLUGIN_H

#include "pluginsiteminterface.h"
#include "trashwidget.h"

#include <QLabel>
#include <QSettings>
namespace Dock{
class TipsWidget;
}
class TrashPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "trash.json")

public:
    explicit TrashPlugin(QObject *parent = 0);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void refreshIcon(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;
    void pluginStateSwitched() override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;
    void pluginSettingsChanged() override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginFlags flags() const override;

private:
    void refreshPluginItemsVisible();

    QScopedPointer<TrashWidget> m_trashWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // TRASHPLUGIN_H
