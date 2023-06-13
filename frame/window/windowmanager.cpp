// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainwindow.h"
#include "taskmanager/entry.h"
#include "windowmanager.h"
#include "traymainwindow.h"
#include "multiscreenworker.h"
#include "menuworker.h"
#include "dockitemmanager.h"
#include "dockscreen.h"
#include "displaymanager.h"

#include <DWindowManagerHelper>
#include <DDBusSender>

#include <QScreen>
#include <QX11Info>
#include <QtConcurrent>
#include <QScreen>
#include <QtGlobal>

#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>

DGUI_USE_NAMESPACE

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"

#define DOCKSCREEN_INS DockScreen::instance()
#define DIS_INS DisplayManager::instance()

using org::kde::StatusNotifierWatcher;

WindowManager::WindowManager(MultiScreenWorker *multiScreenWorker, QObject *parent)
    : QObject(parent)
    , m_multiScreenWorker(multiScreenWorker)
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_position(Dock::Position::Bottom)
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    , m_sniWatcher(new StatusNotifierWatcher(SNI_WATCHER_SERVICE, SNI_WATCHER_PATH, QDBusConnection::sessionBus(), this))
{
    initSNIHost();
    initConnection();
    initMember();
}

WindowManager::~WindowManager()
{
}

void WindowManager::addWindow(MainWindowBase *window)
{
    connect(window, &MainWindowBase::requestUpdate, this, [ window, this ] {
        updateDockGeometry(window->geometry());
    });

    window->setPosition(m_multiScreenWorker->position());
    window->setDisplayMode(m_multiScreenWorker->displayMode());
    window->setOrder(m_topWindows.size());

    m_topWindows << window;
}

void WindowManager::launch()
{
    if (!qApp->property("CANSHOW").toBool())
        return;

    const QString &currentScreen = DOCKSCREEN_INS->current();
    if (m_multiScreenWorker->hideMode() == HideMode::KeepShowing) {
        onPlayAnimation(currentScreen, m_multiScreenWorker->position(), Dock::AniAction::Show);
    } else if (m_multiScreenWorker->hideMode() == HideMode::KeepHidden) {
        qApp->setProperty(PROP_HIDE_STATE, HideState::Hide);
        onUpdateDockGeometry(HideMode::KeepHidden);
    } else if (m_multiScreenWorker->hideMode() == HideMode::SmartHide) {
        switch(m_multiScreenWorker->hideState()) {
        case HideState::Show:
            onPlayAnimation(currentScreen, m_multiScreenWorker->position(), Dock::AniAction::Show);
            break;
        case HideState::Hide:
            onPlayAnimation(currentScreen, m_multiScreenWorker->position(), Dock::AniAction::Hide);
            break;
        default:
            break;
        }

        qApp->setProperty(PROP_HIDE_STATE, m_multiScreenWorker->hideState());
    }

    QMetaObject::invokeMethod(this, [ this ] {
        for (MainWindowBase *mainWindow : m_topWindows)
            mainWindow->setDisplayMode(m_multiScreenWorker->displayMode());
    }, Qt::QueuedConnection);
}

void WindowManager::sendNotifications()
{
    QStringList actionButton;
    actionButton << "reload" << tr("Exit Safe Mode");
    QVariantMap hints;
    hints["x-deepin-action-reload"] = QString("dbus-send,--session,--dest=org.deepin.dde.Dock1,--print-reply,/org/deepin/dde/Dock1,org.deepin.dde.Dock1.ReloadPlugins");
    // 在进入安全模式时，执行此DBUS耗时25S左右，导致任务栏显示阻塞，所以使用线程调用
    QtConcurrent::run(QThreadPool::globalInstance(), [ = ] {
        DDBusSender()
                .service(notificationService)
                .path(notificationPath)
                .interface(notificationInterface)
                .method(QString("Notify"))
                .arg(QString("dde-control-center"))                                            // appname
                .arg(static_cast<uint>(0))                                                     // id
                .arg(QString("preferences-system"))                                            // icon
                .arg(QString(tr("Dock - Safe Mode")))                                          // summary
                .arg(QString(tr("The Dock is in safe mode, please exit to show it properly"))) // content
                .arg(actionButton)                                                             // actions
                .arg(hints)                                                                    // hints
                .arg(15000)                                                                    // timeout
                .call();
    });
}

/**
 * @brief MainWindow::callShow
 * 此方法是被外部进程通过DBus调用的。
 * @note 当任务栏以-r参数启动时，其不会显示界面，需要在外部通过DBus调用此接口之后才会显示界面，
 * 这里是以前为了优化任务栏的启动速度做的处理，当任务栏启动时，此时窗管进程可能还未启动完全，
 * 部分设置未初始化完等，导致任务栏显示的界面异常，所以留下此接口，被startdde延后调用
 */

void WindowManager::callShow()
{
    static bool flag = false;
    if (flag) {
        return;
    }
    flag = true;

    qApp->setProperty("CANSHOW", true);

    launch();

    // 预留200ms提供给窗口初始化再通知startdde，不影响启动速度
    QTimer::singleShot(200, this, &WindowManager::RegisterDdeSession);
}

/** 调整任务栏的大小，这个接口提供给dbus使用，一般是控制中心来调用
 * @brief WindowManager::resizeDock
 * @param offset
 * @param dragging
 */
void WindowManager::resizeDock(int offset, bool dragging)
{
    QScreen *screen = DIS_INS->screen(DOCKSCREEN_INS->current());
    if (!screen)
        return;

    Utils::setIsDraging(dragging);

    int dockSize = qBound(DOCK_MIN_SIZE, offset, DOCK_MAX_SIZE);
    for (MainWindowBase *mainWindow : m_topWindows) {
        QRect windowRect = mainWindow->getDockGeometry(screen, m_multiScreenWorker->position(), m_multiScreenWorker->displayMode(), Dock::HideState::Hide);
        QRect newWindowRect;
        switch (m_multiScreenWorker->position()) {
        case Top: {
           newWindowRect.setX(windowRect.x());
           newWindowRect.setY(windowRect.y());
           newWindowRect.setWidth(windowRect.width());
           newWindowRect.setHeight(dockSize);
           break;
       }
        case Bottom: {
            newWindowRect.setX(windowRect.x());
            newWindowRect.setY(windowRect.y() + windowRect.height() - dockSize);
            newWindowRect.setWidth(windowRect.width());
            newWindowRect.setHeight(dockSize);
            break;
        }
        case Left: {
            newWindowRect.setX(windowRect.x());
            newWindowRect.setY(windowRect.y());
            newWindowRect.setWidth(dockSize);
            newWindowRect.setHeight(windowRect.height());
            break;
        }
        case Right: {
            newWindowRect.setX(windowRect.x() + windowRect.width() - dockSize);
            newWindowRect.setY(windowRect.y());
            newWindowRect.setWidth(dockSize);
            newWindowRect.setHeight(windowRect.height());
            break;
        }
        }

        // 更新界面大小
        mainWindow->blockSignals(true);
        mainWindow->setFixedSize(newWindowRect.size());
        mainWindow->resetPanelGeometry();
        mainWindow->move(newWindowRect.topLeft());
        mainWindow->blockSignals(false);
    }

    m_multiScreenWorker->updateDaemonDockSize(dockSize);
}

/** 获取任务栏的实际大小，这个接口用于获取任务栏的尺寸返回给dbus接口
 * @brief WindowManager::geometry
 * @return 任务栏的实际尺寸和位置信息
 */
QRect WindowManager::geometry() const
{
    int x = 0;
    int y = 0;
    int width = 0;
    int height = 0;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        for (int i = 0; i < m_topWindows.size(); i++) {
            MainWindowBase *mainWindow = m_topWindows[i];
            if (!mainWindow->isVisible())
                continue;

            QRect windowRect = mainWindow->geometry();
            if (i == 0 || x > windowRect.x())
                x = windowRect.x();
            if (i == 0) {
                y = windowRect.y();
                height = windowRect.height();
            }
            width += windowRect.width() + mainWindow->dockSpace();
        }

        return QRect(x, y, width, height);
    }

    for (int i = 0; i < m_topWindows.size(); i++) {
        MainWindowBase *mainWindow = m_topWindows[i];
        if (!mainWindow->isVisible())
            continue;

        QRect windowRect = mainWindow->geometry();
        if (i == 0 || y > windowRect.y())
            y = windowRect.y();

        if (i == 0) {
            x = windowRect.x();
            width = windowRect.width();
        }

        height += windowRect.height() + mainWindow->dockSpace();
    }

    return QRect(x, y, width, height);
}

void WindowManager::onUpdateDockGeometry(const Dock::HideMode &hideMode)
{
    Dock::HideState hideState;
    if (hideMode == HideMode::KeepShowing || (hideMode == HideMode::SmartHide && m_multiScreenWorker->hideState() == HideState::Show))
        hideState = Dock::HideState::Show;
    else
        hideState = Dock::HideState::Hide;

    updateMainGeometry(hideState);
}

void WindowManager::onPositionChanged(const Dock::Position &position)
{
    Position lastPos = m_position;
    if (lastPos == position)
        return;

    m_position = position;
    // 调用设置位置，一会根据需要放到实际的位置
    for (MainWindowBase *mainWindow : m_topWindows)
        mainWindow->setPosition(position);

    // 在改变位置后，需要根据当前任务栏是隐藏还是显示，来调整左右两侧区域的大小
    onUpdateDockGeometry(HideMode::KeepHidden);
}

void WindowManager::onDisplayModeChanged(const Dock::DisplayMode &displayMode)
{
    m_displayMode = displayMode;

    DockItem::setDockDisplayMode(m_displayMode);
    qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(displayMode));

    for (MainWindowBase *mainWindow : m_topWindows)
        mainWindow->setDisplayMode(m_displayMode);

    onUpdateDockGeometry(m_multiScreenWorker->hideMode());
}

void WindowManager::onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (name == SNI_WATCHER_SERVICE && !newOwner.isEmpty()) {
        qDebug() << SNI_WATCHER_SERVICE << "SNI watcher daemon started, register dock to watcher as SNI Host";
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    }
}

void WindowManager::showAniFinish()
{
    qApp->setProperty(PROP_HIDE_STATE, HideState::Show);

    // 通知后端更新区域
    onRequestUpdateFrontendGeometry();
    onRequestNotifyWindowManager();
}

void WindowManager::animationFinish(bool showOrHide)
{
    for (MainWindowBase *mainWindow : m_topWindows) {
        if (!mainWindow->isVisible())
            continue;

        mainWindow->animationFinished(showOrHide);
    }
}

void WindowManager::hideAniFinish()
{
    DockItem::setDockPosition(m_position);
    qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));
    qApp->setProperty(PROP_HIDE_STATE, HideState::Hide);
    // 通知后端更新区域
    onRequestUpdateFrontendGeometry();
    onRequestNotifyWindowManager();
}

/**获取整个任务栏区域的位置和尺寸的信息，用于提供给后端设置位置等信息
 * @brief WindowManager::getDockGeometry
 * @param withoutScale
 * @return
 */
QRect WindowManager::getDockGeometry(bool withoutScale) const
{
    QScreen *screen = DIS_INS->screen(DOCKSCREEN_INS->current());
    if (!screen)
        return QRect();

    int x = 0;
    int y = 0;
    int width = 0;
    int height = 0;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        for (int i = 0; i < m_topWindows.size(); i++) {
            MainWindowBase *mainWindow = m_topWindows[i];
            if (!mainWindow->isVisible())
                continue;

            QRect windowRect = mainWindow->getDockGeometry(screen, m_position, m_displayMode, Dock::HideState::Show, withoutScale);
            if (i == 0 || x > windowRect.x())
                x = windowRect.x();

            if (i == 0) {
                y = windowRect.y();
                height = windowRect.height();
            }
            width += windowRect.width() + mainWindow->dockSpace();
        }
    } else {
        for (int i = 0; i < m_topWindows.size(); i++) {
            MainWindowBase *mainWindow = m_topWindows[i];
            if (!mainWindow->isVisible())
                continue;

            QRect windowRect = mainWindow->getDockGeometry(screen, m_position, m_displayMode, Dock::HideState::Show, withoutScale);
            if (i == 0 || y > windowRect.y())
                y = windowRect.y();

            if (i == 0) {
                x = windowRect.x();
                width = windowRect.width();
            }
            height += windowRect.height() + mainWindow->dockSpace();
        }
    }

    return QRect(x, y, width, height);
}

void WindowManager::RegisterDdeSession()
{
    QString envName("DDE_SESSION_PROCESS_COOKIE_ID");

    QByteArray cookie = qgetenv(envName.toUtf8().data());
    qunsetenv(envName.toUtf8().data());

    if (!cookie.isEmpty()) {
        QDBusPendingReply<bool> r = DDBusSender()
                .interface(sessionManagerService)
                .path(sessionManagerPath)
                .service(sessionManagerInterface)
                .method("Register")
                .arg(QString(cookie))
                .call();

        qDebug() << Q_FUNC_INFO << r.value();
    }
}

void WindowManager::updateDockGeometry(const QRect &rect)
{
    // 如果当前正在执行动画，则无需设置
    if (m_multiScreenWorker->testState(MultiScreenWorker::ChangePositionAnimationStart)
            || m_multiScreenWorker->testState(MultiScreenWorker::ShowAnimationStart)
            || m_multiScreenWorker->testState(MultiScreenWorker::HideAnimationStart))
        return;

    QScreen *screen = DIS_INS->screen(DOCKSCREEN_INS->current());
    if (!screen || m_position == Dock::Position(-1))
        return;

    for (MainWindowBase *mainWindow : m_topWindows) {
        if (!mainWindow->isVisible())
            continue;

        QRect windowShowSize = mainWindow->getDockGeometry(screen, m_multiScreenWorker->position(),
                                                            m_multiScreenWorker->displayMode(), Dock::HideState::Show);
        switch(m_position) {
        case Dock::Position::Top: {
            windowShowSize.setHeight(rect.height());
            break;
        }
        case Dock::Position::Bottom: {
            int bottomY = windowShowSize.y() + windowShowSize.height();
            windowShowSize.setY(bottomY - rect.height());
            windowShowSize.setHeight(rect.height());
            break;
        }
        case Dock::Position::Left: {
            windowShowSize.setWidth(rect.width());
            break;
        }
        case Dock::Position::Right: {
            int righyX = windowShowSize.x() + windowShowSize.width();
            windowShowSize.setX(righyX - rect.width());
            windowShowSize.setWidth(rect.width());
            break;
        }
        default: break;
        }

        mainWindow->blockSignals(true);
        mainWindow->raise();
        mainWindow->setFixedSize(windowShowSize.size());
        mainWindow->move(windowShowSize.topLeft());
        mainWindow->resetPanelGeometry();
        mainWindow->blockSignals(false);
    }

    // 抛出geometry变化的信号，通知控制中心调整尺寸
    Q_EMIT panelGeometryChanged();
}

void WindowManager::initConnection()
{
    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &WindowManager::onDbusNameOwnerChanged);

    connect(m_multiScreenWorker, &MultiScreenWorker::serviceRestart, this, &WindowManager::onServiceRestart);
    connect(m_multiScreenWorker, &MultiScreenWorker::positionChanged, this, &WindowManager::onPositionChanged);
    connect(m_multiScreenWorker, &MultiScreenWorker::displayModeChanged, this, &WindowManager::onDisplayModeChanged);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestPlayAnimation, this, &WindowManager::onPlayAnimation);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestChangeDockPosition, this, &WindowManager::onChangeDockPosition);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestUpdateDockGeometry, this, &WindowManager::onUpdateDockGeometry);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestUpdateFrontendGeometry, this, &WindowManager::onRequestUpdateFrontendGeometry);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestNotifyWindowManager, this, &WindowManager::onRequestNotifyWindowManager);
    connect(m_multiScreenWorker, &MultiScreenWorker::requestUpdateFrontendGeometry, DockItemManager::instance(), &DockItemManager::requestUpdateDockItem);
    connect(DockItemManager::instance(), &DockItemManager::requestWindowAutoHide, m_multiScreenWorker, &MultiScreenWorker::onAutoHideChanged);
}

void WindowManager::initSNIHost()
{
    QDBusConnection dbusConn = QDBusConnection::sessionBus();
    m_sniHostService = QString("org.kde.StatusNotifierHost-") + QString::number(qApp->applicationPid());
    dbusConn.registerService(m_sniHostService);
    dbusConn.registerObject("/StatusNotifierHost", this);

    if (m_sniWatcher->isValid()) {
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    } else {
        qDebug() << SNI_WATCHER_SERVICE << "SNI watcher daemon is not exist for now!";
    }
}

void WindowManager::initMember()
{
    m_displayMode = m_multiScreenWorker->displayMode();
    m_position = m_multiScreenWorker->position();
}

// 更新任务栏的位置和尺寸的信息
void WindowManager::updateMainGeometry(const Dock::HideState &hideState)
{
    QScreen *screen = DIS_INS->screen(DOCKSCREEN_INS->current());
    if (!screen)
        return;

    for (MainWindowBase *mainWindow : m_topWindows) {
        QRect windowRect = mainWindow->getDockGeometry(screen, m_multiScreenWorker->position(), m_multiScreenWorker->displayMode(), hideState);
        mainWindow->updateParentGeometry(m_position, windowRect);
    }

    // 在切换模式的时候，需要根据实际当前是隐藏还是显示来记录当前任务栏是隐藏还是显示，MainPanelWindow会根据这个状态来决定怎么获取图标的尺寸，
    // 如果不加上这一行，那么鼠标在唤醒任务栏的时候，左侧区域会有显示问题
    qApp->setProperty(PROP_HIDE_STATE, hideState);
}

void WindowManager::onPlayAnimation(const QString &screenName, const Dock::Position &pos, Dock::AniAction act, bool containMouse, bool updatePos)
{
    // 如果containMouse= true,则需要计算鼠标是否包含在任务栏的位置， 如果鼠标在任务栏内部，则无需执行动画,这种情况一般用于执行隐藏模式发生变化的时候
    if (containMouse) {
        QRect dockGeometry = getDockGeometry(false);
        if (dockGeometry.contains(QCursor::pos()))
            return;
    }

    if (act == Dock::AniAction::Show) {
        for (MainWindowBase *mainWindow : m_topWindows) {
            if (mainWindow->windowType() != MainWindowBase::DockWindowType::MainWindow)
                continue;

            // 如果请求显示的动画，且当前任务栏已经显示，则不继续执行动画
            if (mainWindow->width() > 0 && mainWindow->height() > 0)
                return;
        }
    }

    QScreen *screen = DIS_INS->screen(screenName);
    if (!m_multiScreenWorker->testState(MultiScreenWorker::RunState::AutoHide) || qApp->property("DRAG_STATE").toBool()
            || m_multiScreenWorker->testState(MultiScreenWorker::RunState::ChangePositionAnimationStart)
            || m_multiScreenWorker->testState(MultiScreenWorker::RunState::HideAnimationStart)
            || m_multiScreenWorker->testState(MultiScreenWorker::RunState::ShowAnimationStart)
            || !screen)
        return;

    QParallelAnimationGroup *group = createAnimationGroup(act, screenName, pos);
    if (!group)
        return;

    switch (act) {
    case Dock::AniAction::Show:
        m_multiScreenWorker->setStates(MultiScreenWorker::ShowAnimationStart);
        break;
    case Dock::AniAction::Hide:
        m_multiScreenWorker->setStates(MultiScreenWorker::HideAnimationStart);
    }

    connect(group, &QParallelAnimationGroup::finished, this, [ = ] {
        switch (act) {
        case Dock::AniAction::Show:
            showAniFinish();
            if (updatePos)
                onPositionChanged(m_multiScreenWorker->position());
            m_multiScreenWorker->setStates(MultiScreenWorker::ShowAnimationStart, false);
            animationFinish(true);
            break;
        case Dock::AniAction::Hide:
            hideAniFinish();
            if (updatePos)
                onPositionChanged(m_multiScreenWorker->position());
            m_multiScreenWorker->setStates(MultiScreenWorker::HideAnimationStart, false);
            animationFinish(false);
            break;
        }
    });

    group->stop();
    group->start(QVariantAnimation::DeleteWhenStopped);
}

/**创建动画，在时尚模式先同时创建左区域和右区域的动画
 * @brief WindowManager::createAnimationGroup
 * @param aniAction  显示动画还是隐藏动画
 * @param screenName 执行动画的屏幕
 * @param position   执行动画的位置（上下左右）
 * @return           要执行的动画组（左右侧区域同时执行的动画组）
 */
QParallelAnimationGroup *WindowManager::createAnimationGroup(const Dock::AniAction &aniAction, const QString &screenName, const Dock::Position &position) const
{
    QScreen *screen = DIS_INS->screen(screenName);
    if (!screen)
        return nullptr;

    bool stopAnimation = false;
    QList<QVariantAnimation *> animations;
    for (MainWindowBase *mainWindow : m_topWindows) {
        if (!mainWindow->isVisible())
            continue;

        QVariantAnimation *ani = mainWindow->createAnimation(screen, position, aniAction);
        if (!ani) {
            stopAnimation = true;
            continue;
        }

        animations << ani;
    }

    if (stopAnimation) {
        qDeleteAll(animations.begin(), animations.end());
        return nullptr;
    }

    QParallelAnimationGroup *aniGroup = new QParallelAnimationGroup;
    for (QVariantAnimation *ani : animations) {
        ani->setParent(aniGroup);
        aniGroup->addAnimation(ani);
    }

    return aniGroup;
}

void WindowManager::onChangeDockPosition(QString fromScreen, QString toScreen, const Dock::Position &fromPos, const Dock::Position &toPos)
{
    QList<QParallelAnimationGroup *> animations;
    // 获取隐藏的动作
    QParallelAnimationGroup *hideGroup = createAnimationGroup(Dock::AniAction::Hide, fromScreen, fromPos);
    if (hideGroup) {
        connect(hideGroup, &QParallelAnimationGroup::finished, this, [ = ] {
            // 在隐藏动画结束的时候，开始设置位置信息
            onPositionChanged(m_multiScreenWorker->position());
            DockItem::setDockPosition(m_multiScreenWorker->position());
            qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_multiScreenWorker->position()));
        });
        animations << hideGroup;
    }
    // 获取显示的动作
    QParallelAnimationGroup *showGroup = createAnimationGroup(Dock::AniAction::Show, toScreen, toPos);
    if (showGroup)
        animations << showGroup;

    if (animations.size() == 0)
        return;

    m_multiScreenWorker->setStates(MultiScreenWorker::ChangePositionAnimationStart);

    QSequentialAnimationGroup *group = new QSequentialAnimationGroup;
    connect(group, &QVariantAnimation::finished, this, [ = ] {
        // 结束之后需要根据确定需要再隐藏
        showAniFinish();
        m_multiScreenWorker->setStates(MultiScreenWorker::ChangePositionAnimationStart, false);
        animationFinish(true);
        emit panelGeometryChanged();
    });

    for (QParallelAnimationGroup *ani : animations) {
        ani->setParent(group);
        group->addAnimation(ani);
    }

    group->start(QVariantAnimation::DeleteWhenStopped);
}

void WindowManager::onRequestUpdateFrontendGeometry()
{
    QRect rect = getDockGeometry(true);
    // org.deepin.dde.daemon.Dock1的SetFrontendWindowRect接口设置区域时,此区域的高度或宽度不能为0,否则会导致其HideState属性循环切换,造成任务栏循环显示或隐藏
    if (rect.width() == 0 || rect.height() == 0)
        return;

    int x = rect.x();
    int y = rect.y();
    if (m_displayMode == Dock::DisplayMode::Fashion) {
        QScreen *screen = DIS_INS->screen(DOCKSCREEN_INS->current());
        if (screen) {
            if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
                x = screen->handle()->geometry().x() + qMax(0, (int)((screen->handle()->geometry().width() - (rect.width() * qApp->devicePixelRatio())) / 2));
            else
                y = screen->handle()->geometry().y() + qMax(0, (int)((screen->handle()->geometry().height() - (rect.height() * qApp->devicePixelRatio())) / 2));
        }
    }

    TaskManager::instance()->setFrontendWindowRect(x, y, uint(rect.width()), uint(rect.height()));
}

void WindowManager::onRequestNotifyWindowManager()
{
    static QRect lastRect = QRect();
    static int lastScreenWidth = 0;
    static int lastScreenHeight = 0;

    // 从列表中查找主窗口
    MainWindowBase *mainWindow = nullptr;
    for (MainWindowBase *window : m_topWindows) {
        if (window->windowType() != MainWindowBase::DockWindowType::MainWindow)
            continue;

        mainWindow = window;
        break;
    }

    if (!mainWindow)
        return;

    /* 在非主屏或非一直显示状态时，清除任务栏区域，不挤占应用 */
    if ((!DIS_INS->isCopyMode() && DOCKSCREEN_INS->current() != DOCKSCREEN_INS->primary()) || m_multiScreenWorker->hideMode() != HideMode::KeepShowing) {
        lastRect = QRect();

        if (Utils::IS_WAYLAND_DISPLAY) {
            QList<QVariant> varList;
            varList.append(0);//dock位置
            varList.append(0);//dock高度/宽度
            varList.append(0);//start值
            varList.append(0);//end值

            // 此处只需获取左侧主窗口部分即可
            QPlatformWindow *windowHandle = mainWindow->windowHandle()->handle();
            if (windowHandle)
                QGuiApplication::platformNativeInterface()->setWindowProperty(windowHandle, "_d_dwayland_dockstrut", varList);

        } else {
            const auto display = QX11Info::display();
            if (!display) {
                qWarning() << "QX11Info::display() is " << display;
                return;
            }

            XcbMisc::instance()->clear_strut_partial(xcb_window_t(mainWindow->winId()));
        }

        return;
    }

    QRect dockGeometry = getDockGeometry(true);
    if (lastRect == dockGeometry
            && lastScreenWidth == DIS_INS->screenRawWidth()
            && lastScreenHeight == DIS_INS->screenRawHeight()) {
        return;
    }

    lastRect = dockGeometry;
    lastScreenWidth = DIS_INS->screenRawWidth();
    lastScreenHeight = DIS_INS->screenRawHeight();
    qDebug() << "dock real geometry:" << dockGeometry;
    qDebug() << "screen width:" << DIS_INS->screenRawWidth() << ", height:" << DIS_INS->screenRawHeight();

    const qreal &ratio = qApp->devicePixelRatio();
    if (Utils::IS_WAYLAND_DISPLAY) {
        QList<QVariant> varList = {0, 0, 0, 0};
        switch (m_position) {
        case Position::Top:
            varList[0] = 1;
            varList[1] = dockGeometry.y() + dockGeometry.height() + WINDOWMARGIN * ratio;
            varList[2] = dockGeometry.x();
            varList[3] = dockGeometry.x() + dockGeometry.width();
            break;
        case Position::Bottom:
            varList[0] = 3;
            varList[1] = DIS_INS->screenRawHeight() - dockGeometry.y() + WINDOWMARGIN * ratio;
            varList[2] = dockGeometry.x();
            varList[3] = dockGeometry.x() + dockGeometry.width();
            break;
        case Position::Left:
            varList[0] = 0;
            varList[1] = dockGeometry.x() + dockGeometry.width() + WINDOWMARGIN * ratio;
            varList[2] = dockGeometry.y();
            varList[3] = dockGeometry.y() + dockGeometry.height();
            break;
        case Position::Right:
            varList[0] = 2;
            varList[1] = DIS_INS->screenRawWidth() - dockGeometry.x() + WINDOWMARGIN * ratio;
            varList[2] = dockGeometry.y();
            varList[3] = dockGeometry.y() + dockGeometry.height();
            break;
        }

        QPlatformWindow *windowHandle = mainWindow->windowHandle()->handle();
        if (windowHandle) {
            QGuiApplication::platformNativeInterface()->setWindowProperty(windowHandle,"_d_dwayland_dockstrut", varList);
        }
    } else {
        XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
        double strut = 0;
        double strutStart = 0;
        double strutEnd = 0;

        switch (m_position) {
        case Position::Top:
            orientation = XcbMisc::OrientationTop;
            strut = dockGeometry.y() + dockGeometry.height();
            strutStart = dockGeometry.x();
            strutEnd = qMin(dockGeometry.x() + dockGeometry.width(), dockGeometry.right());
            break;
        case Position::Bottom:
            orientation = XcbMisc::OrientationBottom;
            strut = DIS_INS->screenRawHeight() - dockGeometry.y();
            strutStart = dockGeometry.x();
            strutEnd = qMin(dockGeometry.x() + dockGeometry.width(), dockGeometry.right());
            break;
        case Position::Left:
            orientation = XcbMisc::OrientationLeft;
            strut = dockGeometry.x() + dockGeometry.width();
            strutStart = dockGeometry.y();
            strutEnd = qMin(dockGeometry.y() + dockGeometry.height(), dockGeometry.bottom());
            break;
        case Position::Right:
            orientation = XcbMisc::OrientationRight;
            strut = DIS_INS->screenRawWidth() - dockGeometry.x();
            strutStart = dockGeometry.y();
            strutEnd = qMin(dockGeometry.y() + dockGeometry.height(), dockGeometry.bottom());
            break;
        }

        qDebug() << "set reserved area to xcb:" << strut << strutStart << strutEnd;

        const auto display = QX11Info::display();
        if (!display) {
            qWarning() << "QX11Info::display() is " << display;
            return;
        }

        XcbMisc::instance()->set_strut_partial(static_cast<xcb_window_t>(mainWindow->winId()), orientation,
                                               static_cast<uint>(strut + WINDOWMARGIN * ratio),                 // 设置窗口与屏幕边缘距离，需要乘缩放
                                               static_cast<uint>(strutStart),                                   // 设置任务栏起点坐标（上下为x，左右为y）
                                               static_cast<uint>(strutEnd));                                    // 设置任务栏终点坐标（上下为x，左右为y）
    }
}

void WindowManager::onServiceRestart()
{
    for (MainWindowBase *mainWindow : m_topWindows)
        mainWindow->serviceRestart();
}
