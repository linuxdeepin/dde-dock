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

QuickSettingController::QuickSettingController(QObject *parent)
    : AbstractPluginsController(parent)
{
    // 异步加载本地插件
    QMetaObject::invokeMethod(this, &QuickSettingController::startLoader, Qt::QueuedConnection);
}

QuickSettingController::~QuickSettingController()
{
}

void QuickSettingController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<QuickSettingItem *>::iterator findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                 [ = ](QuickSettingItem *item) {
        return item->itemKey() == itemKey;
    });

    if (findItemIterator != m_quickSettingItems.end())
        return;

    QuickSettingItem *quickItem = new QuickSettingItem(itemInter, itemKey);

    m_quickSettingItems << quickItem;

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

QuickSettingController *QuickSettingController::instance()
{
    static QuickSettingController instance;
    return &instance;
}

void QuickSettingController::startLoader()
{
    QString pluginsDir("../plugins/quick-trays");
    if (!QDir(pluginsDir).exists())
        pluginsDir = "/usr/lib/dde-dock/plugins/quick-trays";

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}