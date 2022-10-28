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
    PluginAttribute pluginClass = PluginAttribute::Quick;
    QPluginLoader *pluginLoader = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginLoader(itemInter);
    if (pluginLoader) {
        if (pluginLoader->fileName().contains("/plugins/system-trays")) {
            // 如果是从系统托盘目录下加载的插件，则认为它是托盘插件，此时需要放入到托盘中
            pluginClass = PluginAttribute::System;
        } else {
            QJsonObject meta = pluginLoader->metaData().value("MetaData").toObject();
            if (meta.contains("tool") && meta.value("tool").toBool())
                pluginClass = PluginAttribute::Tool;
            else if (meta.contains("fixed") && meta.value("fixed").toBool())
                pluginClass = PluginAttribute::Fixed;
        }
    }

    m_quickPlugins[pluginClass] << itemInter;
    m_quickPluginsMap[itemInter] = itemKey;

    emit pluginInserted(itemInter, pluginClass);
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

void QuickSettingController::updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part)
{
    Q_EMIT pluginUpdated(itemInter, part);
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
