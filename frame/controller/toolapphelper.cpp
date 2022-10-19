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
#include "quicksettingcontroller.h"

#include <QWidget>
#include <QBoxLayout>

ToolAppHelper::ToolAppHelper(QWidget *pluginAreaWidget, QWidget *toolAreaWidget, QObject *parent)
    : QObject(parent)
    , m_pluginAreaWidget(pluginAreaWidget)
    , m_toolAreaWidget(toolAreaWidget)
    , m_displayMode(DisplayMode::Efficient)
    , m_trashItem(nullptr)
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute &pluginClass) {
        if (pluginClass != QuickSettingController::PluginAttribute::Tool)
            return;

        pluginItemAdded(itemInter);
    });

    QList<PluginsItemInterface *> pluginItems = QuickSettingController::instance()->pluginItems(QuickSettingController::PluginAttribute::Tool);
    for (PluginsItemInterface *pluginItem : pluginItems)
        pluginItemAdded(pluginItem);
}

void ToolAppHelper::setDisplayMode(DisplayMode displayMode)
{
    m_displayMode = displayMode;
    updateWidgetStatus();
    moveToolWidget();
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

void ToolAppHelper::appendToToolArea(int index, DockItem *dockItem)
{
    dockItem->setParent(m_toolAreaWidget);
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_toolAreaWidget->layout());
    if (index >= 0)
        boxLayout->insertWidget(index, dockItem);
    else
        boxLayout->addWidget(dockItem);

    Q_EMIT requestUpdate();
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

void ToolAppHelper::moveToolWidget()
{
    for (int i = m_toolAreaWidget->layout()->count() - 1; i >= 0; i--) {
        QLayoutItem *layoutItem = m_toolAreaWidget->layout()->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginWidget = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginWidget)
            continue;

        m_toolAreaWidget->layout()->removeWidget(pluginWidget);
    }

    if (m_displayMode == Dock::DisplayMode::Fashion) {
        QuickSettingController *quickController = QuickSettingController::instance();
        QList<PluginsItemInterface *> plugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tool);
        for (PluginsItemInterface *plugin : plugins) {
            PluginsItem *pluginWidget = quickController->pluginItemWidget(plugin);
            m_toolAreaWidget->layout()->addWidget(pluginWidget);
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
        if (pluginLayout) {
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
    }

    return dockItems;
}

void ToolAppHelper::pluginItemAdded(PluginsItemInterface *itemInter)
{
    if (m_displayMode != Dock::DisplayMode::Fashion || pluginExists(itemInter))
        return;

    QuickSettingController *quickController = QuickSettingController::instance();
    PluginsItem *pluginItem = quickController->pluginItemWidget(itemInter);
    if (pluginInTool(pluginItem))
        appendToToolArea(0, pluginItem);
}

bool ToolAppHelper::pluginExists(PluginsItemInterface *itemInter) const
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_toolAreaWidget->layout());
    if (!boxLayout)
        return false;

    for (int i = 0; i < boxLayout->count() ; i++) {
        QLayoutItem *layoutItem = boxLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginItem = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginItem)
            continue;

        // 如果当前的插件的接口已经存在，则无需再次插入
        if (pluginItem->pluginItem() == itemInter)
            return true;
    }

    return false;
}
