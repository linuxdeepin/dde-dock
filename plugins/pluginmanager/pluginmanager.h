// Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINMANAGER_H
#define PLUGINMANAGER_H

#include "pluginsiteminterface.h"
#include "pluginmanagerinterface.h"

#include <QObject>

class DockPluginController;
class QuickSettingContainer;
class IconManager;

class PluginManager : public PluginManagerInterface, public PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "pluginmanager.json")

public:
    explicit PluginManager(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginFlags flags() const override;
    PluginSizePolicy pluginSizePolicy() const override;

protected:
    void positionChanged(const Dock::Position position) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;

protected:
    // 实现PluginManagerInterface接口，用于向dock提供所有已经加载的插件
    QList<PluginsItemInterface *> plugins() const override;
    QList<PluginsItemInterface *> pluginsInSetting() const override;
    QList<PluginsItemInterface *> currentPlugins() const override;
    QString itemKey(PluginsItemInterface *itemInter) const override;
    QJsonObject metaData(PluginsItemInterface *itemInter) const override;

private:
    QStringList getPluginPaths() const;

private:
    QSharedPointer<DockPluginController> m_dockController;
    QSharedPointer<QuickSettingContainer> m_quickContainer;
    QSharedPointer<IconManager> m_iconManager;
};

#endif // PLUGINMANAGER_H
