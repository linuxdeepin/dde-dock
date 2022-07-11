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

#include "toolapphelper.h"
#include "dockitem.h"
#include "pluginsitem.h"

#include <QWidget>
#include <QBoxLayout>

ToolAppHelper::ToolAppHelper(QWidget *pluginAreaWidget, QWidget *toolAreaWidget, QObject *parent)
    : QObject(parent)
    , m_pluginAreaWidget(pluginAreaWidget)
    , m_toolAreaWidget(toolAreaWidget)
    , m_displayMode(DisplayMode::Efficient)
    , m_trashItem(nullptr)
{
}

void ToolAppHelper::setDisplayMode(DisplayMode displayMode)
{
    m_displayMode = displayMode;
    resetPluginItems();
    updateWidgetStatus();
}

void ToolAppHelper::addPluginItem(int index, DockItem *dockItem)
{
    if (pluginInTool(dockItem))
        appendToToolArea(index, dockItem);
    else
        appendToPluginArea(index, dockItem);

    // 将插件指针顺序保存到列表中
    if (index >= 0 && index < m_sequentPluginItems.size())
        m_sequentPluginItems.insert(index, dockItem);
    else
        m_sequentPluginItems << dockItem;

    // 保存垃圾箱插件指针
    PluginsItem *pluginsItem = qobject_cast<PluginsItem *>(dockItem);
    if (pluginsItem && pluginsItem->pluginName() == "trash")
        m_trashItem = pluginsItem;

    if (!toolIsVisible())
        updateWidgetStatus();

    Q_EMIT requestUpdate();
}

void ToolAppHelper::removePluginItem(DockItem *dockItem)
{
    if (dockItem == m_trashItem)
        m_trashItem = nullptr;

    if (!removePluginArea(dockItem))
        removeToolArea(dockItem);

    if (m_toolAreaWidget->layout()->count() == 0 && toolIsVisible())
        updateWidgetStatus();

    Q_EMIT requestUpdate();
}

PluginsItem *ToolAppHelper::trashPlugin() const
{
    return m_trashItem;
}

bool ToolAppHelper::toolIsVisible() const
{
    return m_toolAreaWidget->isVisible();
}

void ToolAppHelper::appendToPluginArea(int index, DockItem *dockItem)
{
    // 因为日期时间插件和其他插件的大小有异，为了方便设置边距，在插件区域布局再添加一层布局设置边距
    // 因此在处理插件图标时，需要通过两层布局判断是否为需要的插件，例如拖动插件位置等判断
    QBoxLayout *boxLayout = new QBoxLayout(QBoxLayout::LeftToRight, m_pluginAreaWidget);
    boxLayout->addWidget(dockItem, 0, Qt::AlignCenter);
    QBoxLayout *pluginLayout = static_cast<QBoxLayout *>(m_pluginAreaWidget->layout());
    pluginLayout->insertLayout(index, boxLayout, 0);
}

void ToolAppHelper::appendToToolArea(int index, DockItem *dockItem)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_toolAreaWidget->layout());
    if (index >= 0)
        boxLayout->insertWidget(index, dockItem);
    else
        boxLayout->addWidget(dockItem);
}

bool ToolAppHelper::removePluginArea(DockItem *dockItem)
{
    bool removeResult = false;
    QBoxLayout *pluginLayout = static_cast<QBoxLayout *>(m_pluginAreaWidget->layout());
    for (int i = 0; i < pluginLayout->count(); ++i) {
        QLayoutItem *layoutItem = pluginLayout->itemAt(i);
        QLayout *boxLayout = layoutItem->layout();
        if (boxLayout && boxLayout->itemAt(0)->widget() == dockItem) {
            boxLayout->removeWidget(dockItem);
            pluginLayout->removeItem(layoutItem);
            delete layoutItem;
            layoutItem = nullptr;
            removeResult = true;
        }
    }

    return removeResult;
}

bool ToolAppHelper::removeToolArea(DockItem *dockItem)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_toolAreaWidget->layout());
    for (int i = 0; i < boxLayout->count(); i++) {
        if (boxLayout->itemAt(i)->widget() == dockItem) {
            boxLayout->removeWidget(dockItem);
            return true;
        }
    }

    return false;
}

void ToolAppHelper::resetPluginItems()
{
    if (m_displayMode == DisplayMode::Efficient) {
        // 高效模式下, 让工具区域的插件移动到插件区域显示
        QList<DockItem *> dockItems = dockItemOnWidget(true);
        for (DockItem *dockItem : dockItems) {
            // 从工具列表中移除插件, 将这些插件放入到插件区域
            removeToolArea(dockItem);
            int index = itemIndex(dockItem, false);
            appendToPluginArea(index, dockItem);
        }
    } else {
        // 时尚模式下，将插件区域对应的插件移动到工具区域
        QList<DockItem *> dockItems = dockItemOnWidget(false);
        for (DockItem *dockItem : dockItems) {
            if (!pluginInTool(dockItem))
                continue;

            // 从插件区域中移除相关插件，并将其插入到工具区域中
            removePluginArea(dockItem);
            int index = itemIndex(dockItem, true);
            appendToToolArea(index, dockItem);
        }
    }
}

void ToolAppHelper::updateWidgetStatus()
{
    bool oldVisible = toolIsVisible();
    if (m_displayMode == DisplayMode::Efficient) {
        // 高效模式
        m_pluginAreaWidget->setVisible(true);
        m_toolAreaWidget->setVisible(false);
    } else {
        // 时尚模式
        m_pluginAreaWidget->setVisible(false);
        m_toolAreaWidget->setVisible(m_toolAreaWidget->layout()->count() > 0);
    }
    bool visible = toolIsVisible();
    if (oldVisible != visible)
        Q_EMIT toolVisibleChanged(visible);
}

bool ToolAppHelper::pluginInTool(DockItem *dockItem) const
{
    if (m_displayMode != DisplayMode::Fashion)
        return false;

    PluginsItem *pluginItem = qobject_cast<PluginsItem *>(dockItem);
    if (!pluginItem)
        return false;

    QJsonObject metaData = pluginItem->metaData();
    if (metaData.contains("tool"))
        return metaData.value("tool").toBool();

    return false;
}

/**
 * @brief ToolAppHelper::itemIndex 返回该插件在工具区域（isTool == true）或插件区域（isTool == false）的正确位置
 * @param dockItem
 * @param isTool
 * @return
 */
int ToolAppHelper::itemIndex(DockItem *dockItem, bool isTool) const
{
    int index = m_sequentPluginItems.indexOf(dockItem);
    if (index < 0 || index >= m_sequentPluginItems.size() - 1)
        return -1;

    QList<DockItem *> dockItems = dockItemOnWidget(isTool);
    for (int i = index + 1; i < m_sequentPluginItems.size(); i++) {
        DockItem *nextItem = m_sequentPluginItems[i];
        if (dockItems.contains(nextItem)) {
            // 如果当前包含当前插入的下一个item，则直接返回下一个item的插入位置
            return dockItems.indexOf(nextItem);
        }
    }

    return -1;
}

QList<DockItem *> ToolAppHelper::dockItemOnWidget(bool isTool) const
{
    QList<DockItem *> dockItems;
    if (isTool) {
        QLayout *layout = m_toolAreaWidget->layout();
        for (int i = 0; i < layout->count(); i++) {
            DockItem *dockItem = qobject_cast<DockItem *>(layout->itemAt(i)->widget());
            if (!dockItem)
                continue;

            dockItems << dockItem;
        }
    } else {
        QBoxLayout *pluginLayout = static_cast<QBoxLayout *>(m_pluginAreaWidget->layout());
        for (int i = 0; i < pluginLayout->count(); ++i) {
            QLayoutItem *layoutItem = pluginLayout->itemAt(i);
            QLayout *boxLayout = layoutItem->layout();
            if (!boxLayout)
                continue;

            DockItem *dockItem = qobject_cast<DockItem *>(boxLayout->itemAt(0)->widget());
            if (!dockItem)
                continue;

            dockItems << dockItem;
        }
    }

    return dockItems;
}
