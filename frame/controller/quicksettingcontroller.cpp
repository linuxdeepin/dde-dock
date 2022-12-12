/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "quicksettingcontroller.h"
#include "quicksettingitem.h"
#include "proxyplugincontroller.h"
#include "pluginsitem.h"

QuickSettingController::QuickSettingController(QObject *parent)
    : AbstractPluginsController(parent)
{
    // 加载本地插件
    ProxyPluginController *contoller = ProxyPluginController::instance(PluginType::QuickPlugin);
    contoller->addProxyInterface(this);
    connect(contoller, &ProxyPluginController::pluginLoaderFinished, this, &QuickSettingController::pluginLoaderFinished);
}

QuickSettingController::~QuickSettingController()
{
    ProxyPluginController::instance(PluginType::QuickPlugin)->removeProxyInterface(this);
}

void QuickSettingController::pluginItemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 根据读取到的metaData数据获取当前插件的类型，提供给外部
    PluginAttribute pluginAttr = pluginAttribute(itemInter);

    m_quickPlugins[pluginAttr] << itemInter;
    m_quickPluginsMap[itemInter] = itemKey;

    emit pluginInserted(itemInter, pluginAttr);
}

void QuickSettingController::pluginItemUpdate(PluginsItemInterface * const itemInter, const QString &)
{
    updateDockInfo(itemInter, DockPart::QuickPanel);
    updateDockInfo(itemInter, DockPart::QuickShow);
    updateDockInfo(itemInter, DockPart::SystemPanel);
}

void QuickSettingController::pluginItemRemoved(PluginsItemInterface * const itemInter, const QString &)
{
    for (auto it = m_quickPlugins.begin(); it != m_quickPlugins.end(); it++) {
        QList<PluginsItemInterface *> &plugins = m_quickPlugins[it.key()];
        if (!plugins.contains(itemInter))
            continue;

        plugins.removeOne(itemInter);
        if (plugins.isEmpty()) {
            QuickSettingController::PluginAttribute pluginclass = it.key();
            m_quickPlugins.remove(pluginclass);
        }

        break;
    }

    m_quickPluginsMap.remove(itemInter);
    Q_EMIT pluginRemoved(itemInter);
}

void QuickSettingController::requestSetPluginAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool show)
{
    // 设置插件列表可见事件
    if (show)
        Q_EMIT requestAppletShow(itemInter, itemKey);
}

void QuickSettingController::updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part)
{
    Q_EMIT pluginUpdated(itemInter, part);
}

QuickSettingController::PluginAttribute QuickSettingController::pluginAttribute(PluginsItemInterface * const itemInter) const
{
    // 工具插件，例如回收站
    if (hasFlag(itemInter, PluginFlag::Type_Tool))
        return PluginAttribute::Tool;

    // 系统插件，例如关机按钮
    if (hasFlag(itemInter, PluginFlag::Type_System))
        return PluginAttribute::System;

    // 托盘插件，例如磁盘图标
    if (hasFlag(itemInter, PluginFlag::Type_Tray))
        return PluginAttribute::Tray;

    // 固定插件，例如显示桌面和多任务试图
    if (hasFlag(itemInter, PluginFlag::Type_Fixed))
        return PluginAttribute::Fixed;

    // 通用插件，一般的插件都是通用插件，就是放在快捷插件区域的那些插件
    if (hasFlag(itemInter, PluginFlag::Type_Common))
        return PluginAttribute::Quick;

    // 基本插件，不在任务栏上显示的插件
    return PluginAttribute::None;
}

bool QuickSettingController::hasFlag(PluginsItemInterface *itemInter, PluginFlag flag) const
{
    return itemInter->flags() & flag;
}

QuickSettingController *QuickSettingController::instance()
{
    static QuickSettingController instance;
    return &instance;
}

QList<PluginsItemInterface *> QuickSettingController::pluginItems(const PluginAttribute &pluginClass) const
{
    return m_quickPlugins.value(pluginClass);
}

QString QuickSettingController::itemKey(PluginsItemInterface *pluginItem) const
{
    return m_quickPluginsMap.value(pluginItem);
}

QJsonObject QuickSettingController::metaData(PluginsItemInterface *pluginItem) const
{
    QPluginLoader *pluginLoader = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginLoader(pluginItem);
    if (!pluginLoader)
        return QJsonObject();

    return pluginLoader->metaData().value("MetaData").toObject();
}

PluginsItem *QuickSettingController::pluginItemWidget(PluginsItemInterface *pluginItem)
{
    if (m_pluginItemWidgetMap.contains(pluginItem))
        return m_pluginItemWidgetMap[pluginItem];

    PluginsItem *widget = new PluginsItem(pluginItem, itemKey(pluginItem), metaData(pluginItem));
    m_pluginItemWidgetMap[pluginItem] = widget;
    return widget;
}

QList<PluginsItemInterface *> QuickSettingController::pluginInSettings()
{
    QList<PluginsItemInterface *> settingPlugins;
    // 用于在控制中心显示可改变位置的插件，这里只提供
    QList<PluginsItemInterface *> allPlugins = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginCurrent();
    for (PluginsItemInterface *plugin : allPlugins) {
        if (plugin->pluginDisplayName().isEmpty())
            continue;

        if (hasFlag(plugin, PluginFlag::Attribute_CanSetting))
            settingPlugins << plugin;
    }

    return settingPlugins;
}
