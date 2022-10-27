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
#include "screenspliter_wayland.h"
#include "appitem.h"

#include <QWindow>
#include <QApplication>
#include <QX11Info>
#include <QtWaylandClient>
#define private public
#include <private/qwaylandintegration_p.h>
#include <private/qwaylandshellsurface_p.h>
#include <private/qwaylandwindow_p.h>
#include <private/qwaylandcursor_p.h>
#undef private

#include <registry.h>
#include <ddeshell.h>
#include <event_queue.h>
#include <plasmashell.h>
#include <compositor.h>
#include <clientmanagement.h>
#include <connection_thread.h>

SplitWindowManager *ScreenSpliter_Wayland::m_splitManager = nullptr;

/** wayland下的分屏功能
 * @brief ScreenSpliter_Wayland::ScreenSpliter_Wayland
 * @param parent
 */
ScreenSpliter_Wayland::ScreenSpliter_Wayland(AppItem *appItem, DockEntryInter *entryInter, QObject *parent)
    : ScreenSpliter(appItem, entryInter, parent)
    , m_checkedNotSupport(false)
{
    if (!m_splitManager)
        m_splitManager = new SplitWindowManager;

    connect(m_splitManager, &SplitWindowManager::splitStateChange, this, &ScreenSpliter_Wayland::onSplitStateChange);
}

ScreenSpliter_Wayland::~ScreenSpliter_Wayland()
{
}

void ScreenSpliter_Wayland::startSplit(const QRect &rect)
{
    if (entryInter()->windowInfos().size() == 0) {
        // 如果默认打开的子窗口的数量为0，则无需操作，同时记录标记，在打开新的窗口的时候，设置遮罩
        m_splitRect = rect;
        entryInter()->Activate(QX11Info::getTimestamp());
        return;
    }

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
        palette.setBrush(QPalette::ColorRole::Background, backColor);
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
   const QString windowUuid = splitUuid();
   if (windowUuid.isEmpty())
      return false;

   std::string sUuid = windowUuid.toStdString();
   const char *uuid = sUuid.c_str();
   m_splitManager->requestSplitWindow(uuid, direction);
   return true;
}

QString ScreenSpliter_Wayland::splitUuid() const
{
#ifdef USE_AM
    WindowInfoMap windowsInfo = entryInter()->windowInfos();
    if (windowsInfo.isEmpty())
        return QString();

    const QString uuid = windowsInfo.values()[0].uuid;
    if (windowSupportSplit(uuid))
        return uuid;

#endif
    return QString();
}

bool ScreenSpliter_Wayland::windowSupportSplit(const QString &uuid) const
{
    return m_splitManager->canSplit(uuid);
}

QString ScreenSpliter_Wayland::firstWindowUuid() const
{
#ifdef USE_AM
    WindowInfoMap winInfos = entryInter()->windowInfos();
    if (winInfos.size() == 0)
        return QString();

    return winInfos.begin().value().uuid;
#else
    return QString();
#endif
}

void ScreenSpliter_Wayland::onSplitStateChange(const char *uuid, int splitable)
{
#ifdef USE_AM
    const QString windowUuid = firstWindowUuid();
    qDebug() << "Split State Changed, window uuid:" << windowUuid << "split uuid:" << uuid << "split value:" << splitable;
    if (QString(uuid) != windowUuid)
        return;

    if (m_splitRect.isEmpty())
        return;

    if (splitable > 0) {
        setMaskVisible(m_splitRect, true);
    } else {
        // 如果不支持二分屏，则退出当前的窗体，且标记当前不支持二分屏，下次打开的时候不再进行打开窗口来检测
        entryInter()->ForceQuit();
        m_checkedNotSupport = true;
    }
    m_splitRect = QRect(0, 0, 0, 0);
#endif
}

bool ScreenSpliter_Wayland::suportSplitScreen()
{
    // 如果之前检测过是否不支持分屏(m_checkedNotSupport默认为false，如果不支持分屏，m_checkedNotSupport就会变为true),则直接返回不支持分屏
    if (m_checkedNotSupport)
        return false;

    // 如果存在未打开的窗口，就默认让其认为支持，后续会根据这个来打开一个新的窗口
    if (entryInter()->windowInfos().size() == 0)
        return true;

    // 如果存在已经打开的窗口
    m_checkedNotSupport = splitUuid().isEmpty();
    return (!m_checkedNotSupport);
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
    if (!m_clientManagement)
        return false;

    const QVector <ClientManagement::WindowState> &clientWindowStates = m_clientManagement->getWindowStates();
    qInfo() << "client window states count:" << clientWindowStates.size();
    for (ClientManagement::WindowState windowState : clientWindowStates) {
        qDebug() << "window uuid:" << uuid << "window state uuid:" << windowState.uuid
                 << "active:" << windowState.isActive << "resource name:" << windowState.resourceName;
        if (windowState.splitable > 0 && QString(windowState.uuid) == uuid)
            return true;
    };

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
        connect(m_clientManagement, &ClientManagement::splitStateChange, this, &SplitWindowManager::splitStateChange);
    });
    registry->setEventQueue(eventQueue);
    registry->create(m_connectionThreadObject);
    registry->setup();
}
