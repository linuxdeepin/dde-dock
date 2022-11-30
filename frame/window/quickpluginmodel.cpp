/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
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
#include "quickpluginmodel.h"
#include "pluginsiteminterface.h"
#include "quicksettingcontroller.h"
#include "settingconfig.h"

#include <QWidget>

static QStringList fixedPluginNames { "network", "sound", "power" };
#define PLUGINNAMEKEY "Dock_Quick_Plugin_Name"

QuickPluginModel *QuickPluginModel::instance()
{
    static QuickPluginModel instance;
    return &instance;
}

void QuickPluginModel::addPlugin(PluginsItemInterface *itemInter, int index)
{
    // 这里只接受快捷面板的插件，因此，需要做一次判断
    if (QuickSettingController::instance()->pluginAttribute(itemInter) != QuickSettingController::PluginAttribute::Quick)
        return;

    if (index < 0) {
        // 如果索引值小于0，则认为它插在最后面
        index = m_dockedPluginIndex.size();
    }

    // 如果插入的插件在原来的插件列表中存在，并且位置相同，则不做任何的处理
    int oldIndex = m_dockedPluginIndex.contains(itemInter->pluginName());
    if (oldIndex == index && m_dockedPluginsItems.contains(itemInter))
        return;

    m_dockedPluginIndex[itemInter->pluginName()] = index;
    if (!m_dockedPluginsItems.contains(itemInter)) {
        m_dockedPluginsItems << itemInter;
        // 保存配置到dConfig中
        saveConfig();
    }
    // 向外发送更新列表的信号
    Q_EMIT requestUpdate();
}

void QuickPluginModel::removePlugin(PluginsItemInterface *itemInter)
{
    if (!m_dockedPluginsItems.contains(itemInter) && !m_dockedPluginIndex.contains(itemInter->pluginName()))
        return;

    if (m_dockedPluginIndex.contains(itemInter->pluginName())) {
        m_dockedPluginIndex.remove(itemInter->pluginName());
        // 保存配置到DConfig中
        saveConfig();
    }

    if (m_dockedPluginsItems.contains(itemInter)) {
        m_dockedPluginsItems.removeOne(itemInter);
        Q_EMIT requestUpdate();
    }
}

QList<PluginsItemInterface *> QuickPluginModel::dockedPluginItems() const
{
    // 先查找出固定插件，始终排列在最前面
    QList<PluginsItemInterface *> dockedItems;
    QList<PluginsItemInterface *> activedItems;
    for (PluginsItemInterface *itemInter : m_dockedPluginsItems) {
        if (fixedPluginNames.contains(itemInter->pluginName()))
            dockedItems << itemInter;
        else
            activedItems << itemInter;
    }
    std::sort(dockedItems.begin(), dockedItems.end(), [ this ](PluginsItemInterface *item1, PluginsItemInterface *item2) {
        return m_dockedPluginIndex.value(item1->pluginName()) < m_dockedPluginIndex.value(item2->pluginName());
    });
    std::sort(activedItems.begin(), activedItems.end(), [ this ](PluginsItemInterface *item1, PluginsItemInterface *item2) {
        return m_dockedPluginIndex.value(item1->pluginName()) < m_dockedPluginIndex.value(item2->pluginName());
    });
    return (QList<PluginsItemInterface *>() << dockedItems << activedItems);
}

bool QuickPluginModel::isDocked(PluginsItemInterface *itemInter) const
{
    return (m_dockedPluginsItems.contains(itemInter));
}

bool QuickPluginModel::isFixed(PluginsItemInterface *itemInter) const
{
    return fixedPluginNames.contains(itemInter->pluginName());
}

QuickPluginModel::QuickPluginModel(QObject *parent)
    : QObject(parent)
{
    initConnection();
    initConfig();
}

void QuickPluginModel::onPluginRemoved(PluginsItemInterface *itemInter)
{
    // 如果插件移除，无需移除下方的排序设置，因为下次插件插入的时候还会插入到下方任务栏
    // 因此，此处只需要从列表中移除当前插件
    if (m_dockedPluginsItems.contains(itemInter))
        m_dockedPluginsItems.removeOne(itemInter);
    // 向外发送更新列表的信号
    Q_EMIT requestUpdate();
}

void QuickPluginModel::onSettingChanged(const QString &key, const QVariant &value)
{
    if (key != PLUGINNAMEKEY)
        return;
    QStringList localOrder = m_dockedPluginIndex.keys();
    std::sort(localOrder.begin(), localOrder.end(), [ = ](const QString &key1, const QString &key2) {
        return m_dockedPluginIndex.value(key1) < m_dockedPluginIndex.value(key2);
    });
    if (localOrder == value.toStringList())
        return;

    // 当配置发生变化的时候，更新任务栏的插件显示
    // 1、将当前现有的插件列表中不在配置中的插件移除
    localOrder = value.toStringList();
    for (PluginsItemInterface *itemInter : m_dockedPluginsItems) {
        if (localOrder.contains(itemInter->pluginName()))
            continue;

        m_dockedPluginsItems.removeOne(itemInter);
        m_dockedPluginIndex.remove(itemInter->pluginName());
    }
    // 2、将配置中已有的但是插件列表中没有的插件移动到任务栏上
    QList<PluginsItemInterface *> plugins = QuickSettingController::instance()->pluginItems(QuickSettingController::PluginAttribute::Quick);
    for (PluginsItemInterface *plugin : plugins) {
        if (m_dockedPluginsItems.contains(plugin) || !localOrder.contains(plugin->pluginName()))
            continue;

        m_dockedPluginsItems << plugin;
    }

    m_dockedPluginIndex.clear();
    for (int i = 0; i < localOrder.size(); i++)
        m_dockedPluginIndex[localOrder[i]] = i;

    Q_EMIT requestUpdate();
}

void QuickPluginModel::initConnection()
{
    QuickSettingController *quickController = QuickSettingController::instance();
    connect(quickController, &QuickSettingController::pluginInserted, this, [ this ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute plugAttr) {
        if (plugAttr != QuickSettingController::PluginAttribute::Quick)
            return;

        QWidget *quickWidget = itemInter->itemWidget(QUICK_ITEM_KEY);
        if (quickWidget && !quickWidget->parentWidget())
            quickWidget->setVisible(false);

        // 用来读取已经固定在下方的插件
        if (!m_dockedPluginIndex.contains(itemInter->pluginName()))
            return;

        m_dockedPluginsItems << itemInter;

        // 向外发送更新列表的信号
        Q_EMIT requestUpdate();
    });

    connect(quickController, &QuickSettingController::pluginRemoved, this, &QuickPluginModel::onPluginRemoved);
    connect(quickController, &QuickSettingController::pluginUpdated, this, &QuickPluginModel::requestUpdatePlugin);
    connect(SETTINGCONFIG, &SettingConfig::valueChanged, this, &QuickPluginModel::onSettingChanged);
}

void QuickPluginModel::initConfig()
{
    // 此处用于读取dConfig配置，记录哪些插件是固定在任务栏上面的
    QStringList dockPluginsName = SETTINGCONFIG->value(PLUGINNAMEKEY).toStringList();
    for (int i = 0; i < dockPluginsName.size(); i++)
        m_dockedPluginIndex[dockPluginsName[i]]  = i;
}

void QuickPluginModel::saveConfig()
{
    QStringList plugins;
    for (auto it = m_dockedPluginIndex.begin(); it != m_dockedPluginIndex.end(); it++) {
        plugins << it.key();
    }
    std::sort(plugins.begin(), plugins.end(), [ this ](const QString &p1, const QString &p2) {
        return m_dockedPluginIndex.value(p1) < m_dockedPluginIndex.value(p2);
    });
    SETTINGCONFIG->setValue(PLUGINNAMEKEY, plugins);
}
