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

ToolAppHelper::ToolAppHelper(QWidget *toolAreaWidget, QObject *parent)
    : QObject(parent)
    , m_toolAreaWidget(toolAreaWidget)
    , m_displayMode(DisplayMode::Efficient)
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::Tool)
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
    removeToolArea(dockItem);

    if (m_toolAreaWidget->layout()->count() == 0 && toolIsVisible())
        updateWidgetStatus();

    Q_EMIT requestUpdate();
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
        m_toolAreaWidget->setVisible(false);
    } else {
        // 时尚模式
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
