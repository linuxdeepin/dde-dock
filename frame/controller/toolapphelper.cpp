// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    , m_position(Dock::Position::Bottom)
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::Tool)
            return;

        pluginItemAdded(itemInter);
    });
    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, [ = ](PluginsItemInterface *itemInter) {
        pluginItemRemoved(itemInter);
    });

    QList<PluginsItemInterface *> pluginItems = QuickSettingController::instance()->pluginItems(QuickSettingController::PluginAttribute::Tool);
    for (PluginsItemInterface *pluginItem : pluginItems)
        pluginItemAdded(pluginItem);

    updateToolArea();
}

void ToolAppHelper::setDisplayMode(DisplayMode displayMode)
{
    m_displayMode = displayMode;
    moveToolWidget();
    updateWidgetStatus();
}

void ToolAppHelper::setPosition(Position position)
{
    m_toolAreaWidget->setFixedSize(QWIDGETSIZE_MAX, QWIDGETSIZE_MAX);
    m_position = position;
    updateWidgetStatus();
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
}

bool ToolAppHelper::removeToolArea(PluginsItemInterface *itemInter)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_toolAreaWidget->layout());
    for (int i = 0; i < boxLayout->count(); i++) {
        PluginsItem *dockItem = qobject_cast<PluginsItem *>(boxLayout->itemAt(i)->widget());
        if (dockItem && dockItem->pluginItem() == itemInter) {
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

void ToolAppHelper::updateToolArea()
{
    bool oldVisible = m_toolAreaWidget->isVisible();
    QLayout *layout = m_toolAreaWidget->layout();
    if (m_position == Dock::Position::Bottom || m_position == Dock::Position::Top) {
        int size = 0;
        for (int i = 0; i < layout->count(); i++) {
            PluginsItem *dockItem = qobject_cast<PluginsItem *>(layout->itemAt(i)->widget());
            if (!dockItem)
                continue;

            size += dockItem->width();
        }
        m_toolAreaWidget->setFixedWidth(size);
        m_toolAreaWidget->setVisible(size > 0);
    } else {
        int size = 0;
        for (int i = 0; i < layout->count(); i++) {
            PluginsItem *dockItem = qobject_cast<PluginsItem *>(layout->itemAt(i)->widget());
            if (!dockItem)
                continue;

            size += dockItem->height();
        }
        m_toolAreaWidget->setFixedHeight(size);
        m_toolAreaWidget->setVisible(size > 0);
    }
    bool isVisible = m_toolAreaWidget->isVisible();
    if (oldVisible != isVisible)
        Q_EMIT toolVisibleChanged(isVisible);
}

void ToolAppHelper::updateWidgetStatus()
{
    bool oldVisible = toolIsVisible();
    if (m_displayMode == DisplayMode::Efficient) {
        // 高效模式
        m_toolAreaWidget->setVisible(false);
    } else {
        // 时尚模式
        updateToolArea();
    }
    bool visible = toolIsVisible();
    if (oldVisible != visible)
        Q_EMIT toolVisibleChanged(visible);
}

bool ToolAppHelper::pluginInTool(PluginsItemInterface *itemInter) const
{
    return (QuickSettingController::instance()->pluginAttribute(itemInter) == QuickSettingController::PluginAttribute::Tool);
}

void ToolAppHelper::pluginItemAdded(PluginsItemInterface *itemInter)
{
    if (m_displayMode != Dock::DisplayMode::Fashion || pluginExists(itemInter))
        return;

    QuickSettingController *quickController = QuickSettingController::instance();
    if (pluginInTool(itemInter)) {
        PluginsItem *pluginItem = quickController->pluginItemWidget(itemInter);
        appendToToolArea(0, pluginItem);
        updateToolArea();
        Q_EMIT requestUpdate();
    }
}

void ToolAppHelper::pluginItemRemoved(PluginsItemInterface *itemInter)
{
    QuickSettingController *quickController = QuickSettingController::instance();
    if (pluginInTool(itemInter)) {
        PluginsItem *pluginItem = quickController->pluginItemWidget(itemInter);
        removeToolArea(pluginItem->pluginItem());
        updateToolArea();
        Q_EMIT requestUpdate();
    }
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
