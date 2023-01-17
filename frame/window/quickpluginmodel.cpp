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

#define PLUGINNAMEKEY "Dock_Quick_Plugins"

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

    // 获取当前插件在插件区的位置索引(所有在任务栏上显示的插件)
    int oldIndex = getCurrentIndex(itemInter);
    // 计算插入之前的顺序
    if (oldIndex == index && m_dockedPluginsItems.contains(itemInter))
        return;

    // 根据插件区域的位置计算新的索引值
    int newIndex = generaIndex(index, oldIndex);
    m_dockedPluginIndex[itemInter->pluginName()] = newIndex;
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
        if (isFixed(itemInter))
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
    return !(itemInter->flags() & PluginFlag::Attribute_CanInsert);
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
    for (PluginsItemInterface *plugin : plugins)
        m_dockedPluginsItems << plugin;

    m_dockedPluginIndex.clear();
    for (int i = 0; i < localOrder.size(); i++)
        m_dockedPluginIndex[localOrder[i]] = i;

    Q_EMIT requestUpdate();
}

void QuickPluginModel::initConnection()
{
    QuickSettingController *quickController = QuickSettingController::instance();
    connect(quickController, &QuickSettingController::pluginInserted, this, [ this, quickController ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute plugAttr) {
        if (plugAttr != QuickSettingController::PluginAttribute::Quick)
            return;

        QWidget *quickWidget = itemInter->itemWidget(QUICK_ITEM_KEY);
        if (quickWidget && !quickWidget->parentWidget())
            quickWidget->setVisible(false);

        if (!m_dockedPluginIndex.contains(itemInter->pluginName())) {
            QJsonObject json = quickController->metaData(itemInter);
            if (json.contains("order"))
                m_dockedPluginIndex[itemInter->pluginName()] = json.value("order").toInt();
        }

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

int QuickPluginModel::getCurrentIndex(PluginsItemInterface *itemInter)
{
    QList<PluginsItemInterface *> dockedPluginsItems = m_dockedPluginsItems;
    std::sort(dockedPluginsItems.begin(), dockedPluginsItems.end(), [ this ](PluginsItemInterface *plugin1, PluginsItemInterface *plugin2) {
        return m_dockedPluginIndex.value(plugin1->pluginName()) < m_dockedPluginIndex.value(plugin2->pluginName());
    });
    return dockedPluginItems().indexOf(itemInter);
}

int QuickPluginModel::generaIndex(int insertIndex, int oldIndex)
{
    int newIndex = insertIndex;
    if (oldIndex < 0) {
        newIndex = insertIndex + 1;
        // 如果该插件在列表中存在，则需要将原来的索引值加一
        if (insertIndex < 0) {
            // 如果新插入的索引值为-1,则表示需要插入到末尾的位置，此时需要从索引值中找到最大值
            int lastIndex = -1;
            for (PluginsItemInterface *itemInter : m_dockedPluginsItems) {
                int index = m_dockedPluginIndex.value(itemInter->pluginName());
                if (lastIndex < index)
                    lastIndex = index;
            }
            newIndex = lastIndex + 1;
        }
        if (m_dockedPluginIndex.values().contains(newIndex)) {
            // 遍历map列表，检查列表中是否存在等于新索引的插件，如果存在，将其后面的索引值向后加一
            for (auto it = m_dockedPluginIndex.begin(); it != m_dockedPluginIndex.end(); it++) {
                if (it.value() < newIndex)
                    continue;

                m_dockedPluginIndex[it.key()] = it.value() + 1;
            }
        }
    } else {
        newIndex = insertIndex;
        // 如果该插件已经存在于下面的列表中，则分两种情况
        if (insertIndex < 0) {
            // 如果插入在末尾，则计算最大值
            if (m_dockedPluginIndex.size() > 0) {
                int maxIndex = m_dockedPluginIndex.first();
                for (auto it = m_dockedPluginIndex.begin(); it != m_dockedPluginIndex.end(); it++) {
                    if (maxIndex < it.value())
                        maxIndex = it.value();
                }
                return maxIndex;
            }
            return 0;
        }
        if (insertIndex > oldIndex) {
            int minIndex = NGROUPS_MAX;
            // 新的位置的索引值大于原来位置的索引值，则认为插入在原来的任务栏的后面，将前面的插件的索引值减去1
            for (PluginsItemInterface *itemInter : m_dockedPluginsItems) {
                int pluginDockIndex = getCurrentIndex(itemInter);
                qInfo() << itemInter->pluginDisplayName() << m_dockedPluginIndex[itemInter->pluginName()] << pluginDockIndex;
                if (pluginDockIndex > oldIndex) {
                    if (pluginDockIndex <= insertIndex) {
                        int tmpIndex = m_dockedPluginIndex[itemInter->pluginName()];
                        if (tmpIndex < minIndex)
                            minIndex = tmpIndex;
                    }
                    m_dockedPluginIndex[itemInter->pluginName()]--;
                }
                qInfo() << itemInter->pluginDisplayName() << m_dockedPluginIndex[itemInter->pluginName()];
            }

            if (minIndex != NGROUPS_MAX)
                newIndex = minIndex;
        } else {
            // 新的位置索引小于原来的索引值，则认为是插在任务栏的前面，将任务栏后面的插件的索引值加一
            for (PluginsItemInterface *itemInter : m_dockedPluginsItems) {
                int pluginDockIndex = getCurrentIndex(itemInter);
                if (pluginDockIndex >= insertIndex) {
                    m_dockedPluginIndex[itemInter->pluginName()]++;
                }
            }
        }
    }

    return newIndex;
}
