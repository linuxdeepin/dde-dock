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

void QuickSettingController::sortPlugins()
{
    QList<QuickSettingItem *> primarySettingItems;
    QList<QuickSettingItem *> quickItems;
    for (QuickSettingItem *item : m_quickSettingItems) {
        if (item->isPrimary())
            primarySettingItems << item;
        else
            quickItems << item;
    }

    static QStringList existKeys = {"network-item-key", "sound-item-key", "VPN", "PROJECTSCREEN"};
    qSort(primarySettingItems.begin(), primarySettingItems.end(), [ = ](QuickSettingItem *item1, QuickSettingItem *item2) {
        int index1 = existKeys.indexOf(item1->itemKey());
        int index2 = existKeys.indexOf(item2->itemKey());
        if (index1 >= 0 || index2 >= 0)
            return index1 < index2;

        return true;
    });

    m_quickSettingItems.clear();
    m_quickSettingItems << primarySettingItems << quickItems;
}

void QuickSettingController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<QuickSettingItem *>::iterator findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                 [ = ](QuickSettingItem *item) {
        return item->itemKey() == itemKey;
    });

    if (findItemIterator != m_quickSettingItems.end())
        return;

    QPluginLoader *pluginLoader = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginLoader(itemInter);
    QJsonObject metaData;
    if (pluginLoader)
        metaData = pluginLoader->metaData().value("MetaData").toObject();
    QuickSettingItem *quickItem = new QuickSettingItem(itemInter, itemKey, metaData);

    m_quickSettingItems << quickItem;
    sortPlugins();

    emit pluginInserted(quickItem);
}

void QuickSettingController::itemUpdate(PluginsItemInterface * const itemInter, const QString &)
{
    auto findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                                         [ = ](QuickSettingItem *item) {
        return item->pluginItem() == itemInter;
    });
    if (findItemIterator != m_quickSettingItems.end()) {
        QuickSettingItem *settingItem = *findItemIterator;
        settingItem->update();
    }
}

void QuickSettingController::itemRemoved(PluginsItemInterface * const itemInter, const QString &)
{
    // 删除本地记录的插件列表
    QList<QuickSettingItem *>::iterator findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                                         [ = ](QuickSettingItem *item) {
            return (item->pluginItem() == itemInter);
    });
    if (findItemIterator != m_quickSettingItems.end()) {
        QuickSettingItem *quickItem = *findItemIterator;
        m_quickSettingItems.removeOne(quickItem);
        Q_EMIT pluginRemoved(quickItem);
        quickItem->deleteLater();
    }
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
