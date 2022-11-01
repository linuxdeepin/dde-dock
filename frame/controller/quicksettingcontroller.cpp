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
#include "pluginsiteminterface.h"
#include "proxyplugincontroller.h"
#include "pluginsitem.h"

QuickSettingController::QuickSettingController(QObject *parent)
    : AbstractPluginsController(parent)
{
    // 加载本地插件
    ProxyPluginController::instance(PluginType::QuickPlugin)->addProxyInterface(this);
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
    QPluginLoader *pluginLoader = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginLoader(itemInter);
    if (!pluginLoader)
        return PluginAttribute::Quick;

    if (pluginLoader->fileName().contains(TRAY_PATH)) {
        // 如果是从系统目录下加载的插件，则认为它是系统插件，此时需要放入到托盘中
        return PluginAttribute::Tray;
    }

    const QJsonObject &meta = pluginLoader->metaData().value("MetaData").toObject();
    if (meta.contains("tool") && meta.value("tool").toBool()) {
        // 如果有tool标记，则认为它是工具插件，例如回收站和窗管提供的相关插件
        return PluginAttribute::Tool;
    }

    if (meta.contains("system") && meta.value("system").toBool()) {
        // 如果有system标记，则认为它是右侧的关机按钮插件
        return PluginAttribute::System;
    }

    if (meta.contains("fixed") && meta.value("fixed").toBool()) {
        // 如果有fixed标记，则认为它是固定区域的插件，例如显示桌面和多任务视图
        return PluginAttribute::Fixed;
    }

    // 其他的都认为是快捷插件
    return PluginAttribute::Quick;
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
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &plugins = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginsMap();
    QList<PluginsItemInterface *> allPlugins = plugins.keys();
    for (PluginsItemInterface *plugin : allPlugins) {
        PluginAttribute pluginAttr = pluginAttribute(plugin);
        if (pluginAttr == QuickSettingController::PluginAttribute::Quick
                || pluginAttr == QuickSettingController::PluginAttribute::System
                || pluginAttr == QuickSettingController::PluginAttribute::Tool)
            settingPlugins << plugin;
    }

    return settingPlugins;
}
