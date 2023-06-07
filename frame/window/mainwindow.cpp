// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainwindow.h"
#include "mainpanelcontrol.h"
#include "dockitemmanager.h"
#include "menuworker.h"
#include "windowmanager.h"
#include "dockscreen.h"
#include "dragwidget.h"
#include "multiscreenworker.h"
#include "constants.h"
#include "displaymanager.h"

#include <DStyle>
#include <DPlatformWindowHandle>
#include <DSysInfo>
#include <DPlatformTheme>
#include <DDBusSender>

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <QtConcurrent>
#include <qpa/qplatformwindow.h>
#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>
#include <QMenu>

#include <X11/X.h>
#include <X11/Xutil.h>

#define DOCK_SCREEN DockScreen::instance()
#define DIS_INS DisplayManager::instance()

MainWindow::MainWindow(MultiScreenWorker *multiScreenWorker, QWidget *parent)
    : MainWindowBase(multiScreenWorker, parent)
    , m_mainPanel(new MainPanelControl(this))
    , m_multiScreenWorker(multiScreenWorker)
    , m_needUpdateUi(false)
{
    m_mainPanel->setDisplayMode(m_multiScreenWorker->displayMode());

    initConnections();

    for (auto item : DockItemManager::instance()->itemList())
        m_mainPanel->insertItem(-1, item);
}

void MainWindow::initConnections()
{
    connect(DockItemManager::instance(), &DockItemManager::itemInserted, m_mainPanel, &MainPanelControl::insertItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemRemoved, m_mainPanel, &MainPanelControl::removeItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemUpdated, m_mainPanel, &MainPanelControl::itemUpdated, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::trayVisableCountChanged, this, &MainWindow::resizeDockIcon, Qt::QueuedConnection);
    connect(m_mainPanel, &MainPanelControl::itemMoved, DockItemManager::instance(), &DockItemManager::itemMoved, Qt::DirectConnection);
    connect(m_mainPanel, &MainPanelControl::itemAdded, DockItemManager::instance(), &DockItemManager::itemAdded, Qt::DirectConnection);
    connect(m_mainPanel, &MainPanelControl::requestUpdate, this, &MainWindow::requestUpdate);
}

/**
 * @brief MainWindow::getTrayVisableItemCount
 * 重新获取以下当前托盘区域有多少个可见的图标，并更新图标的大小
 */
void MainWindow::resizeDockIcon()
{
    m_mainPanel->resizeDockIcon();
}

/**
 * @brief MainWindow::setGeometry
 * @param rect 设置任务栏的位置和大小，重写此函数时为了及时发出panelGeometryChanged信号，最终供外部DBus调用方使用
 */
void MainWindow::setGeometry(const QRect &rect)
{
    if (rect == this->geometry())
        return;

    DBlurEffectWidget::setGeometry(rect);
}

MainWindowBase::DockWindowType MainWindow::windowType() const
{
    return DockWindowType::MainWindow;
}

void MainWindow::setPosition(const Position &position)
{
    MainWindowBase::setPosition(position);
    m_mainPanel->setPositonValue(position);

    // 更新鼠标拖拽样式，在类内部设置到qApp单例上去
    if ((Top == position) || (Bottom == position))
        m_mainPanel->setCursor(Qt::SizeVerCursor);
    else
        m_mainPanel->setCursor(Qt::SizeHorCursor);
}

void MainWindow::setDisplayMode(const Dock::DisplayMode &displayMode)
{
    m_mainPanel->setDisplayMode(displayMode);
    MainWindowBase::setDisplayMode(displayMode);
}

void MainWindow::updateParentGeometry(const Position &pos, const QRect &rect)
{
    setFixedSize(rect.size());
    setGeometry(rect);

    int panelSize = windowSize();
    QRect panelRect = rect;
    switch (pos) {
    case Position::Top:
        m_mainPanel->move(0, rect.height() - panelSize);
        panelRect.setHeight(panelSize);
        break;
    case Position::Left:
        m_mainPanel->move(width() - panelSize, 0);
        panelRect.setWidth(panelSize);
        break;
    case Position::Bottom:
        m_mainPanel->move(0, 0);
        panelRect.setHeight(panelSize);
        break;
    case Position::Right:
        m_mainPanel->move(0, 0);
        panelRect.setWidth(panelSize);
        break;
    }
    m_mainPanel->setFixedSize(panelRect.size());
}

QSize MainWindow::suitableSize(const Position &pos, const int &screenSize, const double &deviceRatio) const
{
    return m_mainPanel->suitableSize(pos, screenSize, deviceRatio);
}

void MainWindow::resetPanelGeometry()
{
    m_mainPanel->setFixedSize(size());
    m_mainPanel->move(0, 0);
}

void MainWindow::serviceRestart()
{
    // 在重启服务后，MultiScreenWorker会通知WindowManager类执行PositionChanged动画，在执行此动作过程中
    // 会执行动画，动画需要消耗时间，因此， 在重启服务后，需要标记更新UI,在稍后动画执行结束后，需要重新刷新界面的显示，否则任务栏显示错误
    m_needUpdateUi = true;
}

void MainWindow::animationFinished(bool showOrHide)
{
    if (m_needUpdateUi) {
        // 在动画执行结束后，如果收到需要更新UI的标记，那么则需要重新请求更新界面，在更新结束后，再将更新UI标记为false,那么在下次进来的时候，无需再次更新UI
        Q_EMIT requestUpdate();
        m_needUpdateUi = false;
    }
}
