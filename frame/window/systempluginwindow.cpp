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
#include "systempluginwindow.h"
#include "systemplugincontroller.h"
#include "systempluginitem.h"
#include "dockpluginscontroller.h"

#include <DListView>
#include <QBoxLayout>
#include <QMetaObject>

#define MAXICONSIZE 48
#define MINICONSIZE 24
#define ICONMARGIN 8

SystemPluginWindow::SystemPluginWindow(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_pluginController(new FixedPluginController(this))
    , m_listView(new DListView(this))
    , m_position(Dock::Position::Bottom)
    , m_mainLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
{
    initUi();
    connect(m_pluginController, &DockPluginsController::pluginItemInserted, this, &SystemPluginWindow::onPluginItemAdded);
    connect(m_pluginController, &DockPluginsController::pluginItemRemoved, this, &SystemPluginWindow::onPluginItemRemoved);
    connect(m_pluginController, &DockPluginsController::pluginItemUpdated, this, &SystemPluginWindow::onPluginItemUpdated);
    QMetaObject::invokeMethod(m_pluginController, &DockPluginsController::startLoader, Qt::QueuedConnection);
}

SystemPluginWindow::~SystemPluginWindow()
{
}

void SystemPluginWindow::setPositon(Position position)
{
    if (m_position == position)
        return;

    m_position = position;

    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    else
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
}

QSize SystemPluginWindow::suitableSize()
{
    QMargins m = m_mainLayout->contentsMargins();
    if (m_mainLayout->direction() == QBoxLayout::Direction::LeftToRight) {
        int itemSize = height() - m_mainLayout->contentsMargins().top() - m_mainLayout->contentsMargins().bottom();
        int itemWidth = m.left() + m.right();
        for (int i = 0; i < m_mainLayout->count(); i++) {
            QWidget *widget = m_mainLayout->itemAt(i)->widget();
            if (!widget)
                continue;

            PluginsItem *item = qobject_cast<PluginsItem *>(widget);
            if (!item)
                continue;

            // 如果是横向的，则高度是固定，高宽一致，因此读取高度作为它的尺寸值
            itemWidth += itemSize;
            if (i < m_mainLayout->count() - 1)
                itemWidth += m_mainLayout->spacing();
        }

        itemWidth += m.right();
        return QSize(itemWidth, height());
    }

    int itemSize = width() - m_mainLayout->contentsMargins().left() - m_mainLayout->contentsMargins().right();
    int itemHeight = m.top();
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QWidget *widget = m_mainLayout->itemAt(i)->widget();
        if (!widget)
            continue;

        PluginsItem *item = qobject_cast<PluginsItem *>(widget);
        if (!item)
            continue;

        itemHeight += itemSize;
        if (i < m_mainLayout->count() - 1)
            itemHeight += m_mainLayout->spacing();
    }

    itemHeight += m.bottom();

    return QSize(width(), itemHeight);
}

void SystemPluginWindow::resizeEvent(QResizeEvent *event)
{
    DBlurEffectWidget::resizeEvent(event);
    Q_EMIT pluginSizeChanged();
}

void SystemPluginWindow::initUi()
{
    m_mainLayout->setContentsMargins(8, 8, 8, 8);
    m_mainLayout->setSpacing(5);
}

int SystemPluginWindow::calcIconSize() const
{
    switch (m_position) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        if (height() >= 56)
            return MAXICONSIZE;
        if (height() <= 40)
            return MINICONSIZE;
        return height() - ICONMARGIN * 2;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        if (width() >= 56)
            return MAXICONSIZE;
        if (width() <= 40)
            return MINICONSIZE;
        return width() - ICONMARGIN * 2;
    }
    }
    return -1;
}

void SystemPluginWindow::onPluginItemAdded(PluginsItem *pluginItem)
{
    if (m_mainLayout->children().contains(pluginItem))
        return;

    m_mainLayout->addWidget(pluginItem);
    Q_EMIT pluginSizeChanged();
}

void SystemPluginWindow::onPluginItemRemoved(PluginsItem *pluginItem)
{
    if (!m_mainLayout->children().contains(pluginItem))
        return;

    m_mainLayout->removeWidget(pluginItem);
    Q_EMIT pluginSizeChanged();
}

void SystemPluginWindow::onPluginItemUpdated(PluginsItem *pluginItem)
{
    pluginItem->refreshIcon();
}

// can loader plugins
FixedPluginController::FixedPluginController(QObject *parent)
    : DockPluginsController(parent)
{
}

PluginsItem *FixedPluginController::createPluginsItem(PluginsItemInterface * const itemInter, const QString &itemKey, const QString &pluginApi)
{
    return new StretchPluginsItem(itemInter, itemKey, pluginApi);
}

bool FixedPluginController::needLoad(PluginsItemInterface *itemInter)
{
    return (itemInter->pluginName().compare("shutdown") == 0);
}
