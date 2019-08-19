/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     xuwenw <xuwenw@xuwenw.so>
 *
 * Maintainer: xuwenw <xuwenw@xuwenw.so>
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

#include "mainpanelcontrol.h"
#include "../item/dockitem.h"

#include <DAnchors>

#include <QTimer>

DWIDGET_USE_NAMESPACE

MainPanelControl::MainPanelControl(QWidget *parent)
    : QWidget(parent)
    , m_mainPanelLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_fixedAreaWidget(new QWidget(this))
    , m_appAreaWidget(new QWidget(this))
    , m_trayAreaWidget(new QWidget(this))
    , m_pluginAreaWidget(new QWidget(this))
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_trayAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_pluginLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_position(Qt::TopEdge)
{
    init();
    updateMainPanelLayout();
}

MainPanelControl::~MainPanelControl()
{
}

void MainPanelControl::init()
{
    m_mainPanelLayout->setMargin(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);
    m_fixedAreaLayout->setMargin(0);
    m_fixedAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_pluginLayout->setMargin(0);
    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_trayAreaLayout->setMargin(0);
    m_trayAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_appAreaSonLayout->setMargin(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_fixedAreaWidget);
    m_fixedAreaWidget->setLayout(m_fixedAreaLayout);
    m_mainPanelLayout->addWidget(m_appAreaWidget);
    m_mainPanelLayout->addWidget(m_trayAreaWidget);
    m_trayAreaWidget->setLayout(m_trayAreaLayout);
    m_mainPanelLayout->addWidget(m_pluginAreaWidget);
    m_pluginAreaWidget->setLayout(m_pluginLayout);
    m_appAreaSonWidget->setLayout(m_appAreaSonLayout);
    m_fixedAreaLayout->setSpacing(0);
    m_appAreaSonLayout->setSpacing(0);
    m_trayAreaLayout->setSpacing(0);
    m_pluginLayout->setSpacing(0);

    DAnchors<QWidget> anchors(m_appAreaSonWidget);
    anchors.setAnchor(Qt::AnchorHorizontalCenter, this, Qt::AnchorHorizontalCenter);
    anchors.setAnchor(Qt::AnchorVerticalCenter, this, Qt::AnchorVerticalCenter);
}

void MainPanelControl::updateMainPanelLayout()
{
    switch (m_position) {
    case Qt::TopEdge:
    case Qt::BottomEdge:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_mainPanelLayout->setDirection(QBoxLayout::LeftToRight);
        m_fixedAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_pluginLayout->setDirection(QBoxLayout::LeftToRight);
        m_trayAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_appAreaSonLayout->setDirection(QBoxLayout::LeftToRight);
        break;
    case Qt::RightEdge:
    case Qt::LeftEdge:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_mainPanelLayout->setDirection(QBoxLayout::TopToBottom);
        m_fixedAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_pluginLayout->setDirection(QBoxLayout::TopToBottom);
        m_trayAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_appAreaSonLayout->setDirection(QBoxLayout::TopToBottom);
        break;
    default:
        break;
    }

    QTimer::singleShot(0, this, &MainPanelControl::updateAppAreaSonWidgetSize);
}

void MainPanelControl::addFixedAreaItem(const int index, QWidget *wdg)
{
    m_fixedAreaLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addAppAreaItem(const int index, QWidget *wdg)
{
    m_appAreaSonLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addTrayAreaItem(const int index, QWidget *wdg)
{
    m_trayAreaLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addPluginAreaItem(const int index, QWidget *wdg)
{
    m_pluginLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::removeFixedAreaItem(QWidget *wdg)
{
    m_fixedAreaLayout->removeWidget(wdg);
}

void MainPanelControl::removeAppAreaItem(QWidget *wdg)
{
    m_appAreaSonLayout->removeWidget(wdg);
}

void MainPanelControl::removeTrayAreaItem(QWidget *wdg)
{
    m_trayAreaLayout->removeWidget(wdg);
}

void MainPanelControl::removePluginAreaItem(QWidget *wdg)
{
    m_pluginLayout->removeWidget(wdg);
}

void MainPanelControl::resizeEvent(QResizeEvent *event)
{
    updateAppAreaSonWidgetSize();

    return QWidget::resizeEvent(event);
}

void MainPanelControl::updateAppAreaSonWidgetSize()
{
    if ((m_position == Qt::TopEdge) || (m_position == Qt::BottomEdge)) {
        m_appAreaSonWidget->setMaximumHeight(QWIDGETSIZE_MAX);
        m_appAreaSonWidget->setMaximumWidth(qMin((m_appAreaWidget->geometry().right() - width() / 2) * 2, m_appAreaWidget->width()));
    } else {
        m_appAreaSonWidget->setMaximumWidth(QWIDGETSIZE_MAX);
        m_appAreaSonWidget->setMaximumHeight(qMin((m_appAreaWidget->geometry().bottom() - height() / 2) * 2, m_appAreaWidget->height()));
    }

    m_appAreaSonWidget->adjustSize();
}

void MainPanelControl::setPositonValue(const Qt::Edge val)
{
    m_position = val;
}

void MainPanelControl::insertItem(const int index, DockItem *item)
{
    switch (item->itemType()) {
    case DockItem::Launcher:
        addFixedAreaItem(index, item);
        break;
    case DockItem::App:
        addAppAreaItem(index, item);
        break;
    case DockItem::TrayPlugin:
        addTrayAreaItem(index, item);
        break;
    case DockItem::Plugins:
        addPluginAreaItem(index, item);
        break;
    default:
        break;
    }

    updateAppAreaSonWidgetSize();
}

void MainPanelControl::removeItem(DockItem *item)
{
    switch (item->itemType()) {
    case DockItem::Launcher:
        removeFixedAreaItem(item);
        break;
    case DockItem::App:
        removeAppAreaItem(item);
        break;
    case DockItem::TrayPlugin:
        removeTrayAreaItem(item);
        break;
    case DockItem::Plugins:
        removePluginAreaItem(item);
        break;
    default:
        break;
    }

    updateAppAreaSonWidgetSize();
}

 void MainPanelControl::movedItem(const int index, DockItem *item)
{
    // remove old item
    removeItem(item);
    // insert new position
    insertItem(index, item);
}

