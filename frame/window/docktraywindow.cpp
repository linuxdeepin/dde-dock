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
#include "docktraywindow.h"
#include "datetimedisplayer.h"
#include "systempluginwindow.h"
#include "quickpluginwindow.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "tray_delegate.h"
#include "quicksettingcontroller.h"
#include "pluginsitem.h"

#include <DGuiApplicationHelper>

#include <QBoxLayout>
#include <QLabel>

#define FRONTSPACING 18
#define SPLITERSIZE 2
#define SPLITESPACE 10

DockTrayWindow::DockTrayWindow(DockInter *dockInter, QWidget *parent)
    : QWidget(parent)
    , m_dockInter(dockInter)
    , m_position(Dock::Position::Bottom)
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::RightToLeft, this))
    , m_toolWidget(new QWidget(this))
    , m_toolLayout(new QBoxLayout(QBoxLayout::RightToLeft, m_toolWidget))
    , m_toolLineLabel(new QLabel(this))
    , m_dateTimeWidget(new DateTimeDisplayer(true, this))
    , m_systemPuginWidget(new SystemPluginWindow(this))
    , m_quickIconWidget(new QuickPluginWindow(this))
    , m_trayView(new TrayGridView(this))
    , m_model(new TrayModel(m_trayView, false, true, this))
    , m_delegate(new TrayDelegate(m_trayView, this))
{
    initUi();
    initConnection();

    m_trayView->setModel(m_model);
    m_trayView->setItemDelegate(m_delegate);
    m_trayView->openPersistentEditor(m_model->index(0, 0));
}

void DockTrayWindow::setPositon(const Dock::Position &position)
{
    m_position = position;
    m_dateTimeWidget->setPositon(position);
    m_systemPuginWidget->setPositon(position);
    m_quickIconWidget->setPositon(position);
    m_trayView->setPosition(position);
    m_delegate->setPositon(position);
    QModelIndex index = m_model->index(0, 0);
    m_trayView->closePersistentEditor(index);
    m_trayView->openPersistentEditor(index);
    updateLayout(position);
    onResetLayout();
}

void DockTrayWindow::setDisplayMode(const Dock::DisplayMode &displayMode)
{
    m_displayMode = displayMode;
    moveToolPlugin();
}

QSize DockTrayWindow::suitableSize(const Dock::Position &position, const int &, const double &) const
{
    if (position == Dock::Position::Left || position == Dock::Position::Right) {
        // 左右的尺寸
        int height = FRONTSPACING
                + m_toolWidget->height()
                + (SPLITESPACE * 2)
                + SPLITERSIZE
                + m_dateTimeWidget->suitableSize(position).height()
                + m_systemPuginWidget->suitableSize(position).height()
                + m_quickIconWidget->suitableSize(position).height()
                + m_trayView->suitableSize(position).height();

        return QSize(-1, height);
    }
    // 上下的尺寸
    int width = FRONTSPACING
            + m_toolWidget->width()
            + (SPLITESPACE * 2)
            + SPLITERSIZE
            + m_dateTimeWidget->width()
            + m_systemPuginWidget->width()
            + m_quickIconWidget->width()
            + m_trayView->width();
    return QSize(width, -1);
}

QSize DockTrayWindow::suitableSize() const
{
    return suitableSize(m_position, 0, 0);
}

void DockTrayWindow::layoutWidget()
{
    resizeTool();
}

void DockTrayWindow::resizeEvent(QResizeEvent *event)
{
    Q_EMIT requestUpdate();
    // 当尺寸发生变化的时候，通知托盘区域刷新尺寸，让托盘图标始终保持居中显示
    Q_EMIT m_delegate->sizeHintChanged(m_model->index(0, 0));
    QWidget::resizeEvent(event);
    switch (m_position) {
    case Dock::Position::Left:
    case Dock::Position::Right:
        m_toolLineLabel->setFixedSize(width() * 0.6, SPLITERSIZE);
        break;
    case Dock::Position::Top:
    case Dock::Position::Bottom:
        m_toolLineLabel->setFixedSize(SPLITERSIZE, height() * 0.6);
        break;
    }
}

void DockTrayWindow::paintEvent(QPaintEvent *event)
{
    QWidget::paintEvent(event);
    QPainter painter(this);
    QColor color;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        color = Qt::black;
        painter.setOpacity(0.5);
    } else {
        color = Qt::white;
        painter.setOpacity(0.1);
    }

    painter.fillRect(m_toolLineLabel->geometry(), color);
}

/** 根据任务栏的位置来更新布局的方向
 * @brief DockTrayWindow::updateLayout
 * @param position
 */
void DockTrayWindow::updateLayout(const Dock::Position &position)
{
    switch (position) {
    case Dock::Position::Left:
    case Dock::Position::Right: {
        m_mainBoxLayout->setDirection(QBoxLayout::BottomToTop);
        m_toolLayout->setDirection(QBoxLayout::BottomToTop);
        break;
    }
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        m_mainBoxLayout->setDirection(QBoxLayout::RightToLeft);
        m_toolLayout->setDirection(QBoxLayout::RightToLeft);
        break;
    }
    }
}

void DockTrayWindow::resizeTool() const
{
    int toolSize = 0;
    int size = 0;
    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right)
        size = width();
    else
        size = height();

    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *toolWidget = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!toolWidget)
            continue;

        toolWidget->setFixedSize(size, size);
        toolSize += size;
    }

    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right)
        m_toolWidget->setFixedSize(QWIDGETSIZE_MAX, toolSize);
    else
        m_toolWidget->setFixedSize(toolSize, QWIDGETSIZE_MAX);
}

bool DockTrayWindow::pluginExists(PluginsItemInterface *itemInter) const
{
    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginItem = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginItem)
            continue;

        if (pluginItem->pluginItem() == itemInter)
            return true;
    }

    return false;
}

void DockTrayWindow::moveToolPlugin()
{
    for (int i = m_toolLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginWidget = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginWidget)
            continue;

        m_toolLayout->removeWidget(pluginWidget);
    }
    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 如果当前是高效模式，则将所有的工具插件移动到当前的工具区域
        QuickSettingController *quickController = QuickSettingController::instance();
        QList<PluginsItemInterface *> plugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tool);
        for (PluginsItemInterface *pluginInter : plugins) {
            PluginsItem *pluginWidget = quickController->pluginItemWidget(pluginInter);
            m_toolLayout->addWidget(pluginWidget);
        }
    }
}

void DockTrayWindow::initUi()
{
    m_toolLayout->setContentsMargins(0, 0, 0, 0);
    m_toolLayout->setSpacing(0);

    m_systemPuginWidget->setDisplayMode(Dock::DisplayMode::Efficient);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(0);
    m_mainBoxLayout->addSpacing(FRONTSPACING);
    m_mainBoxLayout->addWidget(m_toolWidget);
    m_mainBoxLayout->addSpacing(SPLITESPACE);
    m_mainBoxLayout->addWidget(m_toolLineLabel);
    m_mainBoxLayout->addSpacing(SPLITESPACE);
    m_mainBoxLayout->addWidget(m_dateTimeWidget);
    m_mainBoxLayout->addWidget(m_systemPuginWidget);
    m_mainBoxLayout->addWidget(m_quickIconWidget);
    m_mainBoxLayout->addWidget(m_trayView);
    m_mainBoxLayout->setAlignment(m_toolLineLabel, Qt::AlignCenter);

    WinInfo info;
    info.type = TrayIconType::ExpandIcon;
    m_model->addRow(info);
    m_trayView->openPersistentEditor(m_model->index(0, 0));

    m_toolLineLabel->setFixedSize(0, 0);

    m_mainBoxLayout->addStretch();
}

void DockTrayWindow::initConnection()
{
    connect(m_systemPuginWidget, &SystemPluginWindow::itemChanged, this, &DockTrayWindow::onResetLayout);
    connect(m_dateTimeWidget, &DateTimeDisplayer::requestUpdate, this, &DockTrayWindow::onResetLayout);
    connect(m_quickIconWidget, &QuickPluginWindow::itemCountChanged, this, &DockTrayWindow::onResetLayout);
    connect(m_trayView, &TrayGridView::requestRemove, this, &DockTrayWindow::onResetLayout);

    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ this ] (PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute &pluginAttr) {
        switch (pluginAttr) {
        case QuickSettingController::PluginAttribute::Tool:
            // 下方只处理回收站等插件
            onItemAdded(itemInter);
            break;
        default:
            break;
        }
    });

    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, &DockTrayWindow::onItemRemove);
}

void DockTrayWindow::onResetLayout()
{
    switch(m_position) {
    case Dock::Position::Left:
    case Dock::Position::Right: {
        m_dateTimeWidget->setFixedSize(QWIDGETSIZE_MAX, m_dateTimeWidget->suitableSize().height());
        m_systemPuginWidget->setFixedSize(QWIDGETSIZE_MAX, m_systemPuginWidget->suitableSize().height());
        m_quickIconWidget->setFixedSize(QWIDGETSIZE_MAX, m_quickIconWidget->suitableSize().height());
        m_trayView->setFixedSize(QWIDGETSIZE_MAX, m_trayView->suitableSize().height());
        break;
    }
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        m_dateTimeWidget->setFixedSize(m_dateTimeWidget->suitableSize().width(), QWIDGETSIZE_MAX);
        m_systemPuginWidget->setFixedSize(m_systemPuginWidget->suitableSize().width(), QWIDGETSIZE_MAX);
        m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), QWIDGETSIZE_MAX);
        m_trayView->setFixedSize(m_trayView->suitableSize().width(), QWIDGETSIZE_MAX);
        break;
    }
    }
    Q_EMIT requestUpdate();
}

void DockTrayWindow::onItemAdded(PluginsItemInterface *itemInter)
{
    if (m_displayMode != Dock::DisplayMode::Efficient || pluginExists(itemInter))
        return;

    QuickSettingController *quickController = QuickSettingController::instance();
    PluginsItem *pluginItem = quickController->pluginItemWidget(itemInter);
    pluginItem->setVisible(true);

    m_toolLayout->addWidget(pluginItem);

    Q_EMIT requestUpdate();
}

void DockTrayWindow::onItemRemove(PluginsItemInterface *itemInter)
{
    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginItem = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginItem || pluginItem->pluginItem() != itemInter)
            continue;

        m_toolLayout->removeWidget(pluginItem);

        Q_EMIT requestUpdate();
        break;
    }
}
