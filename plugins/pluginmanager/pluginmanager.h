/*
 * Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
    bool eventHandler(QEvent *event) override;
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
  