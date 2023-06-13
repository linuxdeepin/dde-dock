// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "screenspliter_wayland.h"
#include "appitem.h"

#include <QWindow>
#include <QApplication>
#include <QX11Info>
#include <QtWaylandClient>

#include <DWayland/Client/registry.h>
#include <DWayland/Client/ddeshell.h>
#include <DWayland/Client/event_queue.h>
#include <DWayland/Client/plasmashell.h>
#include <DWayland/Client/compositor.h>
#include <DWayland/Client/clientmanagement.h>
#include <DWayland/Client/connection_thread.h>

SplitWindowManager *ScreenSpliter_Wayland::m_splitManager = nullptr;

/** wayland下的分屏功能
 * @brief ScreenSpliter_Wayland::ScreenSpliter_Wayland
 * @param parent
 */
ScreenSpliter_Wayland::ScreenSpliter_Wayland(AppItem *appItem, QObject *parent)
    : ScreenSpliter(appItem, parent)
{
    if (!m_splitManager)
        m_splitManager = new SplitWindowManager;
}

ScreenSpliter_Wayland::~ScreenSpliter_Wayland()
{
}

void ScreenSpliter_Wayland::startSplit(const QRect &rect)
{
    if (!suportSplitScreen())
        return;

    setMaskVisible(rect, true);
}

void ScreenSpliter_Wayland::setMaskVisible(const QRect &rect, bool visible)
{
    static QWidget *desktopWidget = nullptr;
    if (!desktopWidget) {
        desktopWidget = new QWidget;
        DPalette palette = DGuiApplicationHelper::instance()->applicationPalette();
        QColor backColor = palette.color(QPalette::Highlight);
        backColor.setAlpha(255 * 0.3);
        palette.setBrush(QPalette::ColorRole::Window, backColor);
        desktopWidget->setPalette(palette);
        desktopWidget->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool);
    }
    desktopWidget->setVisible(visible);
    desktopWidget->setGeometry(rect);
    desktopWidget->raise();
}

bool ScreenSpliter_Wayland::split(SplitDirection direction)
{
   setMaskVisible(QRect(), false);
   // 如果当前不支持分屏，则返回false
   if (!suportSplitScreen())
       return false;

   WindowInfoMap windowInfos = appItem()->windowsInfos();
   m_splitManager->requestSplitWindow(windowInfos.first().uuid.toStdString().c_str(), direction);

   return true;
}

bool ScreenSpliter_Wayland::windowSupportSplit(const QString &uuid) const
{
    return m_splitManager->canSplit(uuid);
}

bool ScreenSpliter_Wayland::suportSplitScreen()
{
    // 判断所有打开的窗口列表，只要有一个窗口支持分屏，就认为它支持分屏
    const WindowInfoMap &windowsInfo = appItem()->windowsInfos();
    for (const WindowInfo &windowInfo : windowsInfo) {
        if (windowSupportSplit(windowInfo.uuid))
            return true;
    }

    // 如果所有的窗口都不支持分屏，就认为它不支持分屏，包括没有打开窗口的情况
    return false;
}

bool ScreenSpliter_Wayland::releaseSplit()
{
    setMaskVisible(QRect(), false);
    return true;
}

/**
 * @brief SplitWindowManager::SplitWindowManager
 * @param wayland下的分屏的管理
 */
SplitWindowManager::SplitWindowManager(QObject *parent)
    : QObject(parent)
    , m_clientManagement(nullptr)
    , m_connectionThread(new QThread(nullptr))
    , m_connectionThreadObject(new ConnectionThread)
{
    connect(m_connectionThreadObject, &ConnectionThread::connected, this, &SplitWindowManager::onConnectionFinished, Qt::QueuedConnection);

    m_connectionThreadObject->moveToThread(m_connectionThread);
    m_connectionThread->start();
    m_connectionThreadObject->initConnection();
}

SplitWindowManager::~SplitWindowManager()
{
}

bool SplitWindowManager::canSplit(const QString &uuid) const
{
    const QVector <ClientManagement::WindowState> &windowStates = m_clientManagement->getWindowStates();
    for (const ClientManagement::WindowState &winState : windowStates)
        if (winState.uuid == uuid && winState.splitable > 0)
            return true;

    return false;
}

static ClientManagement::SplitType convertSplitType(ScreenSpliter::SplitDirection direction)
{
    static QMap<ScreenSpliter::SplitDirection, ClientManagement::SplitType> direcionMapping = {
        { ScreenSpliter::Left, ClientManagement::SplitType::Left },
        { ScreenSpliter::Right, ClientManagement::SplitType::Right},
        { ScreenSpliter::Top, ClientManagement::SplitType::Top },
        { ScreenSpliter::Bottom, ClientManagement::SplitType::Bottom },
        { ScreenSpliter::LeftTop, ClientManagement::SplitType::LeftTop },
        { ScreenSpliter::RightTop, ClientManagement::SplitType::RightTop },
        { ScreenSpliter::LeftBottom, ClientManagement::SplitType::LeftBottom },
        { ScreenSpliter::RightBottom, ClientManagement::SplitType::RightBottom }
    };

    return direcionMapping.value(direction, ClientManagement::SplitType::None);
}

void SplitWindowManager::requestSplitWindow(const char *uuid, const ScreenSpliter::SplitDirection &direction)
{
    m_clientManagement->requestSplitWindow(uuid, convertSplitType(direction));
}

void SplitWindowManager::onConnectionFinished()
{
    EventQueue *eventQueue = new EventQueue(this);
    eventQueue->setup(m_connectionThreadObject);

    Registry *registry = new Registry(this);
    connect(registry, &Registry::clientManagementAnnounced, this, [ this, registry ](quint32 name, quint32 version) {
        m_clientManagement = registry->createClientManagement(name, version, this);
    });
    registry->setEventQueue(eventQueue);
    registry->create(m_connectionThreadObject);
    registry->setup();
}
