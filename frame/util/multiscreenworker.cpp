// Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "multiscreenworker.h"
#include "mainwindow.h"
#include "utils.h"
#include "displaymanager.h"
#include "traymainwindow.h"
#include "mainwindow.h"
#include "menuworker.h"
#include "windowmanager.h"
#include "dockitemmanager.h"
#include "dockscreen.h"
#include "docksettings.h"

#include <QWidget>
#include <QScreen>
#include <QEvent>
#include <QRegion>
#include <QSequentialAnimationGroup>
#include <QVariantAnimation>
#include <QX11Info>
#include <QDBusConnection>
#include <QGuiApplication>
#include <QMenu>

#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>

const QString MonitorsSwitchTime = "monitorsSwitchTime";
const QString OnlyShowPrimary = "onlyShowPrimary";
const double DEFAULTOPACITY = 0.4;

#define DIS_INS DisplayManager::instance()
#define DOCK_SCREEN DockScreen::instance()

// 保证以下数据更新顺序(大环节顺序不要变，内部还有一些小的调整，比如任务栏显示区域更新的时候，里面内容的布局方向可能也要更新...)
// Monitor数据－＞屏幕是否可停靠更新－＞监视唤醒区域更新，任务栏显示区域更新－＞拖拽区域更新－＞通知后端接口，通知窗管

MultiScreenWorker::MultiScreenWorker(QObject *parent)
    : QObject(parent)
    , m_eventInter(new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus(), this))
    , m_extralEventInter(new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus(), this))
    , m_touchEventInter(new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus(), this))
    , m_launcherInter(new DBusLuncher(launcherService, launcherPath, QDBusConnection::sessionBus(), this))
    , m_appearanceInter(new Appearance("org.deepin.dde.Appearance1", "/org/deepin/dde/Appearance1", QDBusConnection::sessionBus(), this))
    , m_monitorUpdateTimer(new QTimer(this))
    , m_delayWakeTimer(new QTimer(this))
    , m_position(Dock::Position(-1))
    , m_hideMode(Dock::HideMode(-1))
    , m_hideState(Dock::HideState(-1))
    , m_displayMode(Dock::DisplayMode(-1))
    , m_state(AutoHide)
{
    initConnection();
    initMembers();
    initDockMode();
    QMetaObject::invokeMethod(this, &MultiScreenWorker::initDisplayData, Qt::QueuedConnection);
}

MultiScreenWorker::~MultiScreenWorker()
{
}

void MultiScreenWorker::updateDaemonDockSize(const int &dockSize)
{
    if (m_displayMode == Dock::DisplayMode::Fashion)
        DockSettings::instance()->setWindowSizeFashion(dockSize);
    else
        DockSettings::instance()->setWindowSizeEfficient(dockSize);
}

/**
 * @brief MultiScreenWorker::setStates 用于存储一些状态
 * @param state 标识是哪一种状态，后面有需求可以扩充
 * @param on 设置状态为true或false
 */
void MultiScreenWorker::setStates(RunStates state, bool on)
{
    RunState type = static_cast<RunState>(int(state & RunState::RunState_Mask));

    if (on)
        m_state |= type;
    else
        m_state &= ~(type);
}

void MultiScreenWorker::onAutoHideChanged(const bool autoHide)
{
    if (testState(AutoHide) != autoHide)
        setStates(AutoHide, autoHide);

    if (testState(AutoHide)) {
        QTimer::singleShot(500, this, &MultiScreenWorker::onDelayAutoHideChanged);
    }
}

void MultiScreenWorker::onRegionMonitorChanged(int x, int y, const QString &key)
{
    if (m_registerKey != key || testState(MousePress))
        return;

    tryToShowDock(x, y);
}

// 鼠标在任务栏之外移动时,任务栏该响应隐藏时需要隐藏
void MultiScreenWorker::onExtralRegionMonitorChanged(int x, int y, const QString &key)
{
    // TODO 后续可以考虑去掉这部分的处理，不用一直监听外部区域的移动，xeventmonitor有一个CursorInto和CursorOut的信号，使用这个也可以替代，但要做好测试工作
    Q_UNUSED(x);
    Q_UNUSED(y);
    Q_UNUSED(key);

    if (m_extralRegisterKey != key || testState(MousePress))
        return;

    // FIXME:每次都要重置一下，是因为qt中的QScreen类缺少nameChanged信号，后面会给上游提交patch修复
    DOCK_SCREEN->updateDockedScreen(getValidScreen(position()));

    // 鼠标移动到任务栏界面之外，停止计时器（延时2秒改变任务栏所在屏幕）
    m_delayWakeTimer->stop();

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Hide);
    }
}

void MultiScreenWorker::updateDisplay()
{
    //1、屏幕停靠信息，
    //2、任务栏当前显示在哪个屏幕也需要更新
    //3、任务栏高度或宽度调整的拖拽区域，
    //4、通知窗管的任务栏显示区域信息，
    //5、通知后端的任务栏显示区域信息
    if (DIS_INS->screens().size() == 0) {
        qWarning() << "No Screen Can Display.";
        return;
    }

    // 更新所在屏幕
    resetDockScreen();
    // 通知后端
    Q_EMIT requestUpdateFrontendGeometry();
    // 通知窗管
    Q_EMIT requestNotifyWindowManager();
}

void MultiScreenWorker::onWindowSizeChanged(uint value)
{
    Q_UNUSED(value);

    m_monitorUpdateTimer->start();
}

void MultiScreenWorker::onPrimaryScreenChanged()
{
    // 先更新主屏信息
    DOCK_SCREEN->updatePrimary(DIS_INS->primary());

    // 无效值
    if (DIS_INS->screenRawHeight() == 0 || DIS_INS->screenRawWidth() == 0) {
        qWarning() << "screen raw data is not valid:"
                   << DIS_INS->screenRawHeight() << DIS_INS->screenRawWidth();
        return;
    }

    m_monitorUpdateTimer->start();
}

void MultiScreenWorker::onPositionChanged(int position)
{
    Position lastPos = m_position;
    if (lastPos == position)
        return;
#ifdef QT_DEBUG
    qDebug() << "position change from: " << lastPos << " to: " << position;
#endif
    m_position = static_cast<Position>(position);

    if (m_hideMode == HideMode::KeepHidden || (m_hideMode == HideMode::SmartHide && m_hideState == HideState::Hide)) {
        // 这种情况切换位置,任务栏不需要显示
        // 参数说明 1 当前屏幕名称 2 改变位置之前的位置，因为需要从之前的位置完成隐藏的动画
        // 3 隐藏动画 4 无需考虑当前鼠标是否在任务栏上，这个参数是通过其他方式隐藏唤醒任务栏的时候考虑鼠标是否在任务栏的位置来决定是否做隐藏动画
        // 默认是false，也就是无需考虑 5 当前动画是否为执行位置改变的动画，如果该值为true，那么在动画执行完成后，WindowManager需要给其管理的
        // 子窗口来更新当前的位置的信息
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), lastPos, Dock::AniAction::Hide, false, true);
        // 更新当前屏幕信息,下次显示从目标屏幕显示
        DOCK_SCREEN->updateDockedScreen(getValidScreen(m_position));
        // 需要更新frontendWindowRect接口数据，否则会造成HideState属性值不变
        emit requestUpdateFrontendGeometry();
    } else {
        // 一直显示的模式才需要显示
        emit requestUpdatePosition(lastPos, m_position);
    }
}

void MultiScreenWorker::onDisplayModeChanged(int displayMode)
{
    if (displayMode == m_displayMode)
        return;

    qInfo() << "display mode change:" << displayMode;

    m_displayMode = static_cast<DisplayMode>(displayMode);

    emit displayModeChanged(m_displayMode);
    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::onHideModeChanged(int hideMode)
{
    if (m_hideMode == hideMode)
        return;

    qInfo() << "hidemode change:" << hideMode;

    m_hideMode = static_cast<HideMode>(hideMode);

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Hide);
    }

    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::onHideStateChanged(int state)
{
    if (state == m_hideState)
        return;

    m_hideState = static_cast<HideState>(state);

    // 检查当前屏幕的当前位置是否允许显示,不允许需要更新显示信息(这里应该在函数外部就处理好,不应该走到这里)

    //TODO 这里是否存在屏幕找不到的问题，m_ds的当前屏幕是否可以做成实时同步的，公用一个指针？
    //TODO 这里真的有必要加以下代码吗，只是隐藏模式的切换，理论上不需要检查屏幕是否允许任务栏停靠
    const QString currentScreen = DOCK_SCREEN->current();
    QScreen *curScreen = DIS_INS->screen(currentScreen);
    if (!DIS_INS->canDock(curScreen, m_position)) {
        DOCK_SCREEN->updateDockedScreen(getValidScreen(m_position));
    }

    qInfo() << "hidestate change:" << m_hideMode << m_hideState;

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        Q_EMIT requestPlayAnimation(currentScreen, m_position, Dock::AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        // 最后一个参数，当任务栏的隐藏状态发生变化的时候（从一直显示变成一直隐藏或者智能隐藏），需要考虑鼠标是否在任务栏上，如果在任务栏上，此时无需执行隐藏动画
        Q_EMIT requestPlayAnimation(currentScreen, m_position, Dock::AniAction::Hide);
    }
}

void MultiScreenWorker::onOpacityChanged(const double value)
{
    if (int(m_opacity * 100) == int(value * 100))
        return;

    m_opacity = value;

    emit opacityChanged(quint8(value * 255));
}

/**
 * @brief onRequestUpdateRegionMonitor  更新监听区域信息
 * 触发时机:屏幕大小,屏幕坐标,屏幕数量,发生变化
 *          任务栏位置发生变化
 *          任务栏'模式'发生变化
 */
void MultiScreenWorker::onRequestUpdateRegionMonitor()
{
    if (!m_registerKey.isEmpty()) {
#ifdef QT_DEBUG
        bool ret = m_eventInter->UnregisterArea(m_registerKey);
        qDebug() << "取消唤起区域监听:" << ret;
#else
        m_eventInter->UnregisterArea(m_registerKey);
#endif
        m_registerKey.clear();
    }

    if (!m_extralRegisterKey.isEmpty()) {
#ifdef QT_DEBUG
        bool ret = m_extralEventInter->UnregisterArea(m_extralRegisterKey);
        qDebug() << "取消任务栏外部区域监听:" << ret;
#else
        m_extralEventInter->UnregisterArea(m_extralRegisterKey);
#endif
        m_extralRegisterKey.clear();
    }

    if (!m_touchRegisterKey.isEmpty()) {
        m_touchEventInter->UnregisterArea(m_touchRegisterKey);
        m_touchRegisterKey.clear();
    }

    const static int flags = Motion | Button | Key;
    const static int monitorHeight = static_cast<int>(15 * qApp->devicePixelRatio());
    // 后端认为的任务栏大小(无缩放因素影响)
    const int realDockSize = int((m_displayMode == DisplayMode::Fashion ? m_windowFashionSize + 20 : m_windowEfficientSize) * qApp->devicePixelRatio());

    // 任务栏唤起区域
    m_monitorRectList.clear();
    for (auto s : DIS_INS->screens()) {
        // 屏幕此位置不可停靠时,不用监听这块区域
        if (!DIS_INS->canDock(s, m_position))
            continue;

        MonitRect monitorRect;
        QRect screenRect = s->geometry();
        screenRect.setSize(screenRect.size() * s->devicePixelRatio());

        switch (m_position) {
        case Top: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + monitorHeight;
        }
            break;
        case Bottom: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y() + screenRect.height() - monitorHeight;
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Left: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + monitorHeight;
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Right: {
            monitorRect.x1 = screenRect.x() + screenRect.width() - monitorHeight;
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        }

        if (!m_monitorRectList.contains(monitorRect)) {
            m_monitorRectList << monitorRect;
#ifdef QT_DEBUG
            qDebug() << "监听区域：" << monitorRect.x1 << monitorRect.y1 << monitorRect.x2 << monitorRect.y2;
#endif
        }
    }

    m_extralRectList.clear();
    for (auto s : DIS_INS->screens()) {
        // 屏幕此位置不可停靠时,不用监听这块区域
        if (!DIS_INS->canDock(s, m_position))
            continue;

        MonitRect monitorRect;
        QRect screenRect = s->geometry();
        screenRect.setSize(screenRect.size() * s->devicePixelRatio());

        switch (m_position) {
        case Top: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + realDockSize;
        }
            break;
        case Bottom: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y() + screenRect.height() - realDockSize;
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Left: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + realDockSize;
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Right: {
            monitorRect.x1 = screenRect.x() + screenRect.width() - realDockSize;
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        }

        if (!m_extralRectList.contains(monitorRect)) {
            m_extralRectList << monitorRect;
#ifdef QT_DEBUG
            qDebug() << "任务栏内部区域：" << monitorRect.x1 << monitorRect.y1 << monitorRect.x2 << monitorRect.y2;
#endif
        }
    }

    // 触屏监控高度固定调整为最大任务栏高度100+任务栏与屏幕边缘间距
    const int monitHeight = 100 + WINDOWMARGIN;

    // 任务栏触屏唤起区域
    m_touchRectList.clear();
    for (auto s : DIS_INS->screens()) {
        // 屏幕此位置不可停靠时,不用监听这块区域
        if (!DIS_INS->canDock(s, m_position))
            continue;

        MonitRect monitorRect;
        QRect screenRect = s->geometry();
        screenRect.setSize(screenRect.size() * s->devicePixelRatio());

        switch (m_position) {
        case Top: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + monitHeight;
        }
            break;
        case Bottom: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y() + screenRect.height() - monitHeight;
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Left: {
            monitorRect.x1 = screenRect.x();
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + monitHeight;
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        case Right: {
            monitorRect.x1 = screenRect.x() + screenRect.width() - monitHeight;
            monitorRect.y1 = screenRect.y();
            monitorRect.x2 = screenRect.x() + screenRect.width();
            monitorRect.y2 = screenRect.y() + screenRect.height();
        }
            break;
        }

        if (!m_touchRectList.contains(monitorRect)) {
            m_touchRectList << monitorRect;
        }

    }

    m_registerKey = m_eventInter->RegisterAreas(m_monitorRectList, flags);
    m_extralRegisterKey = m_extralEventInter->RegisterAreas(m_extralRectList, flags);
    m_touchRegisterKey = m_touchEventInter->RegisterAreas(m_touchRectList, flags);
}

/**
 * @brief 判断屏幕是否为复制模式的依据，第一个屏幕的X和Y值是否和其他的屏幕的X和Y值相等
 * 对于复制模式，这两个值肯定是相等的，如果不是复制模式，这两个值肯定不等，目前支持双屏
 */
bool MultiScreenWorker::isCopyMode()
{
    QList<QScreen *> screens = DIS_INS->screens();
    if (screens.size() < 2)
        return false;

    // 在多个屏幕的情况下，如果所有屏幕的位置的X和Y值都相等，则认为是复制模式
    QRect rect0 = screens[0]->availableGeometry();
    for (int i = 1; i < screens.size(); i++) {
        QRect rect = screens[i]->availableGeometry();
        if (rect0.x() != rect.x() || rect0.y() != rect.y())
            return false;
    }

    return true;
}

void MultiScreenWorker::onRequestUpdatePosition(const Position &fromPos, const Position &toPos)
{
    qInfo() << "request change pos from: " << fromPos << " to: " << toPos;
    // 更新要切换到的屏幕
    if (!DIS_INS->canDock(DIS_INS->screen(DOCK_SCREEN->current()), m_position))
        DOCK_SCREEN->updateDockedScreen(getValidScreen(m_position));

    qInfo() << "update allow screen: " << DOCK_SCREEN->current();

    // 无论什么模式,都先显示
    changeDockPosition(DOCK_SCREEN->last(), DOCK_SCREEN->current(), fromPos, toPos);
}

void MultiScreenWorker::onRequestUpdateMonitorInfo()
{
    resetDockScreen();

    // 只需要在屏幕信息变化的时候更新，其他时间不需要更新
    onRequestUpdateRegionMonitor();

    m_monitorUpdateTimer->start();
}

void MultiScreenWorker::onRequestDelayShowDock()
{
    // 移动Dock至相应屏相应位置
    if (testState(LauncherDisplay))//启动器显示,则dock不显示
        return;

    // 复制模式．不需要响应切换屏幕
    if (DIS_INS->screens().size() == 2 && DIS_INS->screens().first()->geometry() == DIS_INS->screens().last()->geometry()) {
        qInfo() << "copy mode　or merge mode";
        Q_EMIT requestUpdateDockGeometry(m_hideMode);
        return;
    }

    DOCK_SCREEN->updateDockedScreen(m_delayScreen);

    // 检查当前屏幕的当前位置是否允许显示,不允许需要更新显示信息(这里应该在函数外部就处理好,不应该走到这里)
    // 检查边缘是否允许停靠
    QScreen *curScreen = DIS_INS->screen(m_delayScreen);
    if (curScreen && DIS_INS->canDock(curScreen, m_position)) {
        if (m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) {
            Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Hide);
        } else if (m_hideMode == HideMode::KeepShowing) {
            changeDockPosition(DOCK_SCREEN->last(), DOCK_SCREEN->current(), m_position, m_position);
        }
    }
}

void MultiScreenWorker::initMembers()
{
    m_monitorUpdateTimer->setInterval(100);
    m_monitorUpdateTimer->setSingleShot(true);

    m_delayWakeTimer->setSingleShot(true);

    m_windowFashionSize = int(DockSettings::instance()->getWindowSizeFashion() * qApp->devicePixelRatio());
    m_windowEfficientSize = int(DockSettings::instance()->getWindowSizeEfficient() * qApp->devicePixelRatio());

    setStates(LauncherDisplay, m_launcherInter->isValid() ? m_launcherInter->visible() : false);

    // init check
    checkXEventMonitorService();
}

void MultiScreenWorker::initConnection()
{
    connect(DIS_INS, &DisplayManager::primaryScreenChanged, this, &MultiScreenWorker::onPrimaryScreenChanged);
    connect(DIS_INS, &DisplayManager::screenInfoChanged, this, &MultiScreenWorker::requestUpdateMonitorInfo);

    connect(m_launcherInter, static_cast<void (DBusLuncher::*)(bool)>(&DBusLuncher::VisibleChanged), this, [ = ](bool value) { setStates(LauncherDisplay, value); });
    connect(m_appearanceInter, &Appearance::OpacityChanged, this, &MultiScreenWorker::onOpacityChanged);

    connect(this, &MultiScreenWorker::requestUpdatePosition, this, &MultiScreenWorker::onRequestUpdatePosition);
    connect(this, &MultiScreenWorker::requestUpdateMonitorInfo, this, &MultiScreenWorker::onRequestUpdateMonitorInfo);

    connect(m_delayWakeTimer, &QTimer::timeout, this, &MultiScreenWorker::onRequestDelayShowDock);

    // 刷新所有显示的内容，布局，方向，大小，位置等
    connect(m_monitorUpdateTimer, &QTimer::timeout, this, &MultiScreenWorker::updateDisplay);

    connect(DockSettings::instance(), &DockSettings::windowSizeEfficientChanged, this, [=]( uint size){ m_windowEfficientSize = size; });
    connect(DockSettings::instance(), &DockSettings::windowSizeFashionChanged, this, [=]( uint size){ m_windowFashionSize = size; });

    connect(DockSettings::instance(), &DockSettings::windowSizeEfficientChanged, this, &MultiScreenWorker::onWindowSizeChanged);
    connect(DockSettings::instance(), &DockSettings::windowSizeFashionChanged, this, &MultiScreenWorker::onWindowSizeChanged);
    connect(DockSettings::instance(), &DockSettings::positionModeChanged, this, &MultiScreenWorker::onPositionChanged);
    connect(DockSettings::instance(), &DockSettings::displayModeChanged, this, &MultiScreenWorker::onDisplayModeChanged);
    connect(DockSettings::instance(), &DockSettings::hideModeChanged, this, &MultiScreenWorker::onHideModeChanged);
    connect(TaskManager::instance(), &TaskManager::hideStateChanged, this, &MultiScreenWorker::onHideStateChanged);
}

void MultiScreenWorker::initDockMode()
{
        onPositionChanged(DockSettings::instance()->getPositionMode());
        onDisplayModeChanged(DockSettings::instance()->getDisplayMode());
        onHideModeChanged(DockSettings::instance()->getHideMode());
        onHideStateChanged(TaskManager::instance()->getHideState());
        onOpacityChanged(m_appearanceInter? m_appearanceInter->opacity(): DEFAULTOPACITY);

        DockItem::setDockPosition(m_position);
        qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));
        DockItem::setDockDisplayMode(m_displayMode);
        qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(m_displayMode));
}

/**
 * @brief initDisplayData
 * 初始化任务栏的所有必要信息,并更新其位置
 */
void MultiScreenWorker::initDisplayData()
{
    //3\初始化监视区域
    onRequestUpdateRegionMonitor();

    //4\初始化任务栏停靠屏幕
    resetDockScreen();
}

/**
 * @brief reInitDisplayData
 * 重新初始化任务栏的所有必要信息,并更新其位置
 */
void MultiScreenWorker::reInitDisplayData()
{
    initDockMode();
    initDisplayData();
}

/**
 * @brief changeDockPosition    做一个动画操作
 * @param fromScreen            上次任务栏所在的屏幕
 * @param toScreen              任务栏要移动到的屏幕
 * @param fromPos               任务栏上次的方向
 * @param toPos                 任务栏打算移动到的方向
 */
void MultiScreenWorker::changeDockPosition(QString fromScreen, QString toScreen, const Position &fromPos, const Position &toPos)
{
    if (fromScreen == toScreen && fromPos == toPos) {
        qWarning() << "shouldn't be here,nothing happend!";
        return;
    }

    // 该动画放到WindowManager中来实现
    // 更新屏幕信息
    DOCK_SCREEN->updateDockedScreen(toScreen);

    // TODO: 考虑切换过快的情况,这里需要停止上一次的动画,可增加信号控制,暂时无需要
    qInfo() << "from: " << fromScreen << "  to: " << toScreen;
    Q_EMIT requestChangeDockPosition(fromScreen, toScreen, fromPos, toPos);
}

/**
 * @brief getValidScreen        获取一个当前任务栏可以停靠的屏幕，优先使用主屏
 * @return
 */
QString MultiScreenWorker::getValidScreen(const Position &pos)
{
    //TODO 考虑在主屏幕名变化时自动更新，是不是就不需要手动处理了
    DOCK_SCREEN->updatePrimary(DIS_INS->primary());

    if (DIS_INS->canDock(DIS_INS->screen(DOCK_SCREEN->current()), pos))
        return DOCK_SCREEN->current();

    if (DIS_INS->canDock(DIS_INS->screen(DIS_INS->primary()), pos))
        return DIS_INS->primary();

    for (auto s : DIS_INS->screens()) {
        if (DIS_INS->canDock(s, pos))
            return s->name();
    }

    return QString();
}

/**
 * @brief resetDockScreen     检查一下当前屏幕所在边缘是够允许任务栏停靠，不允许的情况需要更换下一块屏幕
 */
void MultiScreenWorker::resetDockScreen()
{
    if (testState(ChangePositionAnimationStart)
            || testState(HideAnimationStart)
            || testState(ShowAnimationStart)
            || Utils::isDraging())
        return;

    DOCK_SCREEN->updateDockedScreen(getValidScreen(position()));
    // 更新任务栏自身信息
    Q_EMIT requestUpdateDockGeometry(m_hideMode);
}

bool MultiScreenWorker::isCursorOut(int x, int y)
{
    const int realDockSize = int((m_displayMode == DisplayMode::Fashion ? m_windowFashionSize : m_windowEfficientSize) * qApp->devicePixelRatio());
    for (auto s : DIS_INS->screens()) {
        // 屏幕此位置不可停靠时,不用监听这块区域
        if (!DIS_INS->canDock(s, m_position))
            continue;

        QRect screenRect = s->geometry();
        screenRect.setSize(screenRect.size() * s->devicePixelRatio());

        if (m_position == Top) {
            // 如果任务栏在顶部
            if (x < screenRect.x() || x > (screenRect.x() + screenRect.width()))
                continue;

            return (y > (screenRect.y() + realDockSize) || y < screenRect.y());
        }

        if (m_position == Bottom) {
            // 如果任务栏在底部
            if (x < screenRect.x() || x > (screenRect.x() + screenRect.width()))
                continue;

            return (y < (screenRect.y() + screenRect.height() - realDockSize) || y > (screenRect.y() + screenRect.height()));
        }

        if (m_position == Left) {
            // 如果任务栏在左侧
            if (y < screenRect.y() || y > (screenRect.y() + screenRect.height()))
                continue;

            return (x > (screenRect.x() + realDockSize) || x < screenRect.x());
        }

        if (m_position == Right) {
            // 如果在任务栏右侧
            if (y < screenRect.y() || y > (screenRect.y() + screenRect.height()))
                continue;

            return (x < (screenRect.x() + screenRect.width() - realDockSize) || x > (screenRect.x() + screenRect.width()));
        }
    }

    return false;
}

/**
 * @brief checkDaemonXEventMonitorService
 * org.deepin.dde.XEventMonitor1服务比dock晚启动，导致dock启动后的状态错误
 */
void MultiScreenWorker::checkXEventMonitorService()
{
    auto connectionInit = [ = ](XEventMonitor * eventInter, XEventMonitor * extralEventInter, XEventMonitor * touchEventInter) {
        connect(eventInter, &XEventMonitor::CursorMove, this, &MultiScreenWorker::onRegionMonitorChanged);
        connect(eventInter, &XEventMonitor::ButtonPress, this, [ = ] { setStates(MousePress, true); });
        connect(eventInter, &XEventMonitor::ButtonRelease, this, [ = ] { setStates(MousePress, false); });

        connect(extralEventInter, &XEventMonitor::CursorOut, this, [ = ](int x, int y, const QString &key) {
            if (isCursorOut(x, y)) {
                if (testState(ShowAnimationStart)) {
                    // 在OUT后如果检测到当前的动画正在进行，在out后延迟500毫秒等动画结束再执行移出动画
                    QTimer::singleShot(500, this, [ = ] {
                        onExtralRegionMonitorChanged(x, y, key);
                    });
                } else {
                    onExtralRegionMonitorChanged(x, y, key);
                }
            }
        });

        // 触屏时，后端只发送press、release消息，有move消息则为鼠标，press置false
        connect(touchEventInter, &XEventMonitor::CursorMove, this, [ = ] { setStates(TouchPress, false); });
        connect(touchEventInter, &XEventMonitor::ButtonPress, this, &MultiScreenWorker::onTouchPress);
        connect(touchEventInter, &XEventMonitor::ButtonRelease, this, &MultiScreenWorker::onTouchRelease);
    };

    QDBusConnectionInterface *ifc = QDBusConnection::sessionBus().interface();

    if (!ifc->isServiceRegistered(xEventMonitorService)) {
        connect(ifc, &QDBusConnectionInterface::serviceOwnerChanged, this, [ = ](const QString & name, const QString & oldOwner, const QString & newOwner) {
            Q_UNUSED(oldOwner)
            if (name == xEventMonitorService && !newOwner.isEmpty()) {
                FREE_POINT(m_eventInter);
                FREE_POINT(m_extralEventInter);
                FREE_POINT(m_touchEventInter);

                m_eventInter = new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus());
                m_extralEventInter = new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus());
                m_touchEventInter = new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus());
                // connect
                connectionInit(m_eventInter, m_extralEventInter, m_touchEventInter);

                disconnect(ifc);
            }
        });
    } else {
        connectionInit(m_eventInter, m_extralEventInter, m_touchEventInter);
    }
}

bool MultiScreenWorker::onScreenEdge(const QString &screenName, const QPoint &point)
{
    QScreen *screen = DIS_INS->screen(screenName);
    if (screen) {
        const QRect r { screen->geometry() };
        const QRect rect { r.topLeft(), r.size() *screen->devicePixelRatio() };

        // 除了要判断鼠标的x坐标和当前区域的位置外，还需要判断当前的坐标的y坐标是否在任务栏的区域内
        // 因为有如下场景：任务栏在左侧，双屏幕屏幕上下拼接，此时鼠标沿着最左侧x=0的位置移动到另外一个屏幕
        // 如果不判断y坐标的话，此时就认为鼠标在当前任务栏的边缘，导致任务栏在这种状况下没有跟随鼠标
        if ((rect.x() == point.x() || rect.x() + rect.width() == point.x())
                && point.y() >= rect.top() && point.y() <= rect.bottom()) {
            return true;
        }

        // 同上，不过此时屏幕是左右拼接，任务栏位于上方或者下方
        if ((rect.y() == point.y() || rect.y() + rect.height() == point.y())
                && point.x() >= rect.left() && point.x() <= rect.right()) {
            return true;
        }
    }

    return false;
}

const QPoint MultiScreenWorker::rawXPosition(const QPoint &scaledPos)
{
    QScreen const *screen = Utils::screenAtByScaled(scaledPos);

    return screen ? screen->geometry().topLeft() +
                    (scaledPos - screen->geometry().topLeft()) *
                    screen->devicePixelRatio()
                  : scaledPos;
}

void MultiScreenWorker::onTouchPress(int type, int x, int y, const QString &key)
{
    Q_UNUSED(type);
    if (key != m_touchRegisterKey) {
        return;
    }

    setStates(TouchPress);
    m_touchPos = QPoint(x, y);
}

void MultiScreenWorker::onTouchRelease(int type, int x, int y, const QString &key)
{
    Q_UNUSED(type);
    if (key != m_touchRegisterKey) {
        return;
    }

    if (!testState(TouchPress)) {
        return;
    }
    setStates(TouchPress, false);

    // 不从指定方向划入，不进行任务栏唤醒；如当任务栏在下，需从下往上划
    switch (m_position) {
    case Top:
        if (m_touchPos.y() >= y) {
            return;
        }
        break;
    case Bottom:
        if (m_touchPos.y() <= y) {
            return;
        }
        break;
    case Left:
        if (m_touchPos.x() >= x) {
            return;
        }
        break;
    case Right:
        if (m_touchPos.x() <= x) {
            return;
        }
        break;
    }

    tryToShowDock(x, y);
}

void MultiScreenWorker::onDelayAutoHideChanged()
{
    switch (m_hideMode) {
    case HideMode::KeepHidden: {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Hide, true);
        break;
    }
    case HideMode::SmartHide: {
        if (m_hideState == HideState::Show)
            Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Show);
        else if (m_hideState == HideState::Hide)
            Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Hide);
        break;
    }
    case HideMode::KeepShowing: {
        Q_EMIT requestPlayAnimation(DOCK_SCREEN->current(), m_position, Dock::AniAction::Show);
        break;
    }
    }
}

/**
 * @brief tryToShowDock 根据xEvent监控区域信号的x，y坐标处理任务栏唤醒显示
 * @param eventX        监控信号x坐标
 * @param eventY        监控信号y坐标
 */
void MultiScreenWorker::tryToShowDock(int eventX, int eventY)
{
    DockItem::setDockPosition(m_position);
    if (qApp->property("DRAG_STATE").toBool() || testState(ChangePositionAnimationStart)) {
        qWarning() << "dock is draging or animation is running";
        return;
    }

    QString toScreen;
    QScreen *screen = Utils::screenAtByScaled(QPoint(eventX, eventY));
    if (!screen) {
        qWarning() << "cannot find the screen" << QPoint(eventX, eventY);
        return;
    }

    toScreen = screen->name();

    /**
     * 坐标位于当前屏幕边缘时,当做屏幕内移动处理(防止鼠标移动到边缘时不唤醒任务栏)
     * 使用screenAtByScaled获取屏幕名时,实际上获取的不一定是当前屏幕
     * 举例:点(100,100)不在(0,0,100,100)的屏幕上
     */
    const QString &currentScreen = DOCK_SCREEN->current();
    if (onScreenEdge(currentScreen, QPoint(eventX, eventY))) {
        toScreen = currentScreen;
    }

    // 过滤重复坐标
    static QPoint lastPos(0, 0);
    if (lastPos == QPoint(eventX, eventY)) {
        return;
    }
    lastPos = QPoint(eventX, eventY);

    // 任务栏显示状态，但需要切换屏幕
    if (toScreen != currentScreen) {
        if (!m_delayWakeTimer->isActive()) {
            m_delayScreen = toScreen;
            m_delayWakeTimer->start(Utils::SettingValue("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", MonitorsSwitchTime, 2000).toInt());
        }
    } else {
        // 任务栏隐藏状态，但需要显示
        if (hideMode() == HideMode::KeepShowing) {
            Q_EMIT requestUpdateDockGeometry(m_hideMode);
            return;
        }

        if (testState(ShowAnimationStart)) {
            qDebug() << "animation is running";
            return;
        }

        if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide)) {
            Q_EMIT requestPlayAnimation(currentScreen, m_position, Dock::AniAction::Show);
        }
    }
}
