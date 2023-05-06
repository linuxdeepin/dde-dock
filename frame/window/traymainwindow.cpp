// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "traymainwindow.h"
#include "traymanagerwindow.h"
#include "dragwidget.h"
#include "dockscreen.h"
#include "displaymanager.h"

#include <DSysInfo>
#include <DPlatformTheme>
#include <DStyleHelper>

#include <QBitmap>
#include <QBoxLayout>
#include <QX11Info>
#include <qpa/qplatformwindow.h>
#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>

#define DOCK_SCREEN DockScreen::instance()
#define DIS_INS DisplayManager::instance()

TrayMainWindow::TrayMainWindow(MultiScreenWorker *multiScreenWorker, QWidget *parent)
    : MainWindowBase(multiScreenWorker, parent)
    , m_trayManager(new TrayManagerWindow(this))
    , m_multiScreenWorker(multiScreenWorker)
{
    initUI();
    initConnection();
}

void TrayMainWindow::setPosition(const Dock::Position &position)
{
    MainWindowBase::setPosition(position);
    m_trayManager->setPositon(position);
}

TrayManagerWindow *TrayMainWindow::trayManagerWindow() const
{
    return m_trayManager;
}

void TrayMainWindow::setDisplayMode(const Dock::DisplayMode &displayMode)
{
    // 只有在时尚模式下才显示
    setVisible(displayMode == Dock::DisplayMode::Fashion);
    MainWindowBase::setDisplayMode(displayMode);
    m_trayManager->setDisplayMode(displayMode);
}

MainWindowBase::DockWindowType TrayMainWindow::windowType() const
{
    return DockWindowType::TrayWindow;
}

void TrayMainWindow::updateParentGeometry(const Dock::Position &position, const QRect &rect)
{
    QSize trayPanelSize = m_trayManager->suitableSize(position);
    // 设置trayManagerWindow的大小和位置
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        setFixedSize(trayPanelSize.width(), rect.height());
        move(rect.topLeft());
    } else {
        setFixedSize(rect.width(), trayPanelSize.height());
        move(rect.topLeft());
    }

    int panelSize = windowSize();
    QRect panelRect = rect;
    switch(position) {
    case Dock::Position::Left:
        m_trayManager->move(width() - panelSize, 0);
        panelRect.setWidth(panelSize);
        panelRect.setHeight(trayPanelSize.height());
        break;
    case Dock::Position::Top:
        m_trayManager->move(0, height() - panelSize);
        panelRect.setWidth(trayPanelSize.width());
        panelRect.setHeight(panelSize);
        break;
    case Dock::Position::Right: {
        m_trayManager->move(0, 0);
        panelRect.setWidth(panelSize);
        panelRect.setHeight(trayPanelSize.height());
        break;
    }
    case Dock::Position::Bottom: {
        m_trayManager->move(0, 0);
        panelRect.setWidth(trayPanelSize.width());
        panelRect.setHeight(panelSize);
        break;
    }
    }

    // 在从高效模式切换到时尚模式的时候，需要调用该函数来设置托盘区域的尺寸，在设置尺寸的时候会触发
    // 托盘区域的requestUpdate信号，WindowManager接收到requestUpdate会依次对每个顶层界面设置尺寸，此时又会触发该函数
    // 引起无限循环，因此，在设置尺寸的时候阻塞信号，防止进入死循环
    m_trayManager->blockSignals(true);
    m_trayManager->setFixedSize(panelRect.size());
    m_trayManager->updateLayout();
    m_trayManager->blockSignals(false);
}

QSize TrayMainWindow::suitableSize() const
{
    return m_trayManager->suitableSize();
}

void TrayMainWindow::resetPanelGeometry()
{
    m_trayManager->setFixedSize(size());
    m_trayManager->move(0, 0);
    m_trayManager->updateLayout();
}

int TrayMainWindow::dockSpace() const
{
    return 0;
}

QSize TrayMainWindow::suitableSize(const Dock::Position &pos, const int &, const double &) const
{
    return m_trayManager->suitableSize(pos);
}

void TrayMainWindow::initUI()
{
    m_trayManager->move(0, 0);
    m_trayManager->updateBorderRadius(MainWindowBase::m_platformWindowHandle.windowRadius());
}

void TrayMainWindow::initConnection()
{
    connect(m_trayManager, &TrayManagerWindow::requestUpdate, this, &TrayMainWindow::onRequestUpdate);
    connect(&(MainWindowBase::m_platformWindowHandle), &DTK_NAMESPACE::Widget::DPlatformWindowHandle::windowRadiusChanged, m_trayManager, [=]{
        m_trayManager->updateBorderRadius(MainWindowBase::m_platformWindowHandle.windowRadius());
    });
}

void TrayMainWindow::onRequestUpdate()
{
    // 如果当前是高效模式，则无需发送信号
    if (displayMode() == Dock::DisplayMode::Efficient)
        return;

    Q_EMIT requestUpdate();
}
