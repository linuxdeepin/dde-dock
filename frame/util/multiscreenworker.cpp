﻿// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "multiscreenworker.h"
#include "mainwindow.h"
#include "utils.h"
#include "displaymanager.h"

#include <DWindowManagerHelper>

#include <QWidget>
#include <QScreen>
#include <QEvent>
#include <QRegion>
#include <QSequentialAnimationGroup>
#include <QVariantAnimation>
#include <QX11Info>
#include <QDBusConnection>
#include <QGuiApplication>

#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>

const QString MonitorsSwitchTime = "monitorsSwitchTime";
const QString OnlyShowPrimary = "onlyShowPrimary";

#define WINDOW_MARGIN ((m_displayMode == Dock::Efficient) ? 0 : 10)
#define DIS_INS DisplayManager::instance()

// 保证以下数据更新顺序(大环节顺序不要变，内部还有一些小的调整，比如任务栏显示区域更新的时候，里面内容的布局方向可能也要更新...)
// Monitor数据－＞屏幕是否可停靠更新－＞监视唤醒区域更新，任务栏显示区域更新－＞拖拽区域更新－＞通知后端接口，通知窗管

MultiScreenWorker::MultiScreenWorker(QWidget *parent)
    : QObject(parent)
    , m_parent(parent)
    , m_eventInter(new XEventMonitor("com.deepin.api.XEventMonitor", "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus(), this))
    , m_extralEventInter(new XEventMonitor("com.deepin.api.XEventMonitor", "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus(), this))
    , m_touchEventInter(new XEventMonitor("com.deepin.api.XEventMonitor", "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus(), this))
    , m_dockInter(new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
    , m_launcherInter(new DBusLuncher("com.deepin.dde.Launcher", "/com/deepin/dde/Launcher", QDBusConnection::sessionBus(), this))
    , m_monitorUpdateTimer(new QTimer(this))
    , m_delayWakeTimer(new QTimer(this))
    , m_delayDisplay(new QTimer(this))
    , m_ds(DIS_INS->primary())
    , m_screenMonitor(new ScreenChangeMonitor(&m_ds, this))
    , m_state(AutoHide)
{
    qInfo() << "init dock screen: " << m_ds.current();

    initConnection();
    initMembers();
    initDBus();
    initDisplayData();
    initUI();
}

void MultiScreenWorker::initShow()
{
    // 仅在初始化时调用一次
    static bool first = true;
    if (!first)
        return;
    first = false;

    //　这里是为了在调用时让MainWindow更新界面布局方向
    emit requestUpdateLayout();

    if (m_hideMode == HideMode::KeepShowing)
        displayAnimation(m_ds.current(), AniAction::Show);
    else if (m_hideMode == HideMode::KeepHidden) {
        QRect rect = getDockShowGeometry(m_ds.current(), m_position, m_displayMode);
        parent()->panel()->setFixedSize(rect.size());
        parent()->panel()->move(0, 0);

        rect = getDockHideGeometry(m_ds.current(), m_position, m_displayMode);
        parent()->setFixedSize(rect.size());
        parent()->setGeometry(rect);
    } else if (m_hideMode == HideMode::SmartHide) {
        if (m_hideState == HideState::Show)
            displayAnimation(m_ds.current(), AniAction::Show);
        else if (m_hideState == HideState::Hide)
            displayAnimation(m_ds.current(), AniAction::Hide);
    }
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

/**
 * @brief dockRect
 * @param screenName            屏幕名
 * @param pos                   任务栏位置
 * @param hideMode              模式
 * @param displayMode           状态
 * @return                      按照给定的数据计算出任务栏所在位置
 */
QRect MultiScreenWorker::dockRect(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode)
{
    if (hideMode == HideMode::KeepShowing || (hideMode == HideMode::SmartHide && m_hideState == HideState::Show))
        return getDockShowGeometry(screenName, pos, displayMode);
    else
        return getDockHideGeometry(screenName, pos, displayMode);
}

/**
 * @brief dockRect
 * @param screenName        屏幕名
 * @return                  按照当前屏幕的当前属性给出任务栏所在区域
 */
QRect MultiScreenWorker::dockRect(const QString &screenName)
{
    return dockRect(screenName, m_position, m_hideMode, m_displayMode);
}

/**
 * @brief realDockRect      给出不计算缩放情况的区域信息(和后端接口保持一致)
 * @param screenName        屏幕名
 * @param pos               任务栏位置
 * @param hideMode          模式
 * @param displayMode       状态
 * @return
 */
QRect MultiScreenWorker::dockRectWithoutScale(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode)
{
    if (hideMode == HideMode::KeepShowing || m_currentHideState == HideState::Show)
        return getDockShowGeometry(screenName, pos, displayMode, true);
    else
        return getDockHideGeometry(screenName, pos, displayMode, true);
}

void MultiScreenWorker::onAutoHideChanged(bool autoHide)
{
    if (testState(AutoHide) != autoHide)
        setStates(AutoHide, autoHide);

    if (testState(AutoHide)) {
        QTimer::singleShot(500, this, &MultiScreenWorker::onDelayAutoHideChanged);
    }
}

/**
 * @brief updateDaemonDockSize
 * @param dockSize              这里的高度是通过qt获取的，不能使用后端的接口数据
 */
void MultiScreenWorker::updateDaemonDockSize(int dockSize)
{
    m_dockInter->setWindowSize(uint(dockSize));
    if (m_displayMode == DisplayMode::Fashion)
        m_dockInter->setWindowSizeFashion(uint(dockSize));
    else
        m_dockInter->setWindowSizeEfficient(uint(dockSize));
}

void MultiScreenWorker::handleDBusSignal(QDBusMessage msg)
{
    QList<QVariant> arguments = msg.arguments();
    // 参数固定长度
    if (3 != arguments.count())
        return;
    // 返回的数据中,这一部分对应的是数据发送方的interfacename,可判断是否是自己需要的服务
    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName == "com.deepin.dde.daemon.Dock") {
        QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
        QStringList keys = changedProps.keys();
        foreach (const QString &prop, keys) {
            if (prop == "Position") {
                onPositionChanged(static_cast<Position>(changedProps.value(prop).toInt()));
            } else if (prop == "DisplayMode") {
                onDisplayModeChanged(static_cast<DisplayMode>(changedProps.value(prop).toInt()));
            } else if (prop == "HideMode") {
                onHideModeChanged(static_cast<HideMode>(changedProps.value(prop).toInt()));
            } else if (prop == "HideState") {
                onHideStateChanged(static_cast<HideState>(changedProps.value(prop).toInt()));
            }
        }
    }
}

void MultiScreenWorker::onRegionMonitorChanged(int x, int y, const QString &key)
{
    if (m_registerKey != key || testState(MousePress))
        return;

    CheckShouldDisplay(x, y, key);

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
    m_ds.updateDockedScreen(getValidScreen(position()));

    // 鼠标移动到任务栏界面之外，停止计时器（延时2秒改变任务栏所在屏幕）
    m_delayWakeTimer->stop();

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        displayAnimation(m_ds.current(), AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        displayAnimation(m_ds.current(), AniAction::Hide);
    }
}

void MultiScreenWorker::CheckShouldDisplay(int x, int y, const QString &key)
{
    if(m_registerKey != key) {
        if(m_delayDisplay->isActive()){
            qDebug()<<"not in check area....";
            m_delayDisplay->stop();
            setStates(CheckDockShouldDisplay, false);
        }
        return;
    }

    if (testState(MousePress)) {
        if(m_delayDisplay->isActive()) {
            qDebug()<<"MousePress stop check....";
            m_delayDisplay->stop();
            setStates(CheckDockShouldDisplay, false);
        }
        return;
    }

    m_delayDisplayPos = QPoint(x, y);

    if(!m_delayDisplay->isActive() && m_currentHideState == HideState::Hide) {
        qDebug()<<"start check....";
        m_delayDisplay->start();
        setStates(CheckDockShouldDisplay);
    }
}

void MultiScreenWorker::showAniFinished()
{
    const QRect rect = dockRect(m_ds.current(), m_position, HideMode::KeepShowing, m_displayMode);

    parent()->setFixedSize(rect.size());
    parent()->setGeometry(rect);

    parent()->panel()->setFixedSize(rect.size());
    parent()->panel()->move(0, 0);

    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::hideAniFinished()
{
    // 提前更新界面布局
    emit requestUpdateLayout();

    const QRect rect = dockRect(m_ds.current(), m_position, HideMode::KeepHidden, m_displayMode);
    parent()->setFixedSize(rect.size());
    parent()->setGeometry(rect);
    parent()->panel()->setFixedSize(rect.size());
    parent()->panel()->move(0, 0);

    DockItem::setDockPosition(m_position);
    qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));

    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::updateDisplay()
{
    //1、屏幕停靠信息，
    //2、任务栏当前显示在哪个屏幕也需要更新
    //2、监视任务栏唤醒区域信息，
    //3、任务栏高度或宽度调整的拖拽区域，
    //4、通知窗管的任务栏显示区域信息，
    //5、通知后端的任务栏显示区域信息

    // wayland环境下，QScreen销毁也就是拔掉显示器的时候，会隐藏当前窗口，这里手动显示一次
    parent()->setVisible(true);

    if (DIS_INS->screens().size() == 0) {
        qWarning() << "No Screen Can Display.";
        return;
    }
    // 更新所在屏幕
    resetDockScreen();
    // 通知后端
    onRequestUpdateFrontendGeometry();
    // 通知窗管
    onRequestNotifyWindowManager();
}

void MultiScreenWorker::onWindowSizeChanged(uint value)
{
    Q_UNUSED(value);

    m_monitorUpdateTimer->start();
}

void MultiScreenWorker::primaryScreenChanged(QScreen *screen)
{
    Q_ASSERT(screen);

    // 先更新主屏信息
    m_ds.updatePrimary(screen->name());

    // 无效值
    if (DIS_INS->screenRawHeight() == 0 || DIS_INS->screenRawWidth() == 0) {
        qWarning() << "screen raw data is not valid:"
                   << DIS_INS->screenRawHeight() << DIS_INS->screenRawWidth();
        return;
    }

    m_monitorUpdateTimer->start();
}

void MultiScreenWorker::updateParentGeometry(const QVariant &value, const Position &pos)
{
    Q_ASSERT_X(value.type() == QVariant::Rect, "", "ERROR OCCURED!");

    const QRect &rect = value.toRect();
    parent()->setFixedSize(rect.size());
    parent()->setGeometry(rect);

    const int panelSize = ((pos == Position::Top || pos == Position::Bottom) ? parent()->panel()->height() : parent()->panel()->width());

    switch (pos) {
    case Position::Top: {
        parent()->panel()->move(0, rect.height() - panelSize);
    }
        break;
    case Position::Left: {
        parent()->panel()->move(rect.width() - panelSize, 0);
    }
        break;
    case Position::Bottom:
    case Position::Right: {
        parent()->panel()->move(0, 0);
    }
        break;
    }
}

void MultiScreenWorker::updateParentGeometry(const QVariant &value)
{
    if (!testState(ShowAnimationStart) && !testState(HideAnimationStart))
        return;

    updateParentGeometry(value, m_position);
}

void MultiScreenWorker::onPositionChanged(const Position &position)
{
    Position lastPos = m_position;
    if (lastPos == position)
        return;

    qInfo() << "position changed from: " << lastPos << " to: " << position;
    m_position = position;

    // 更新鼠标拖拽样式，在类内部设置到qApp单例上去
    if ((Top == m_position) || (Bottom == m_position)) {
        parent()->panel()->setCursor(Qt::SizeVerCursor);
    } else {
        parent()->panel()->setCursor(Qt::SizeHorCursor);
    }

    if (m_hideMode == HideMode::KeepHidden || (m_hideMode == HideMode::SmartHide && m_hideState == HideState::Hide)) {
        // 这种情况切换位置,任务栏不需要显示
        displayAnimation(m_ds.current(), lastPos, AniAction::Hide);
        // 更新当前屏幕信息,下次显示从目标屏幕显示
        m_ds.updateDockedScreen(getValidScreen(m_position));
        // 需要更新frontendWindowRect接口数据，否则会造成HideState属性值不变
        emit requestUpdateFrontendGeometry();
    } else {
        // 一直显示的模式才需要显示
        emit requestUpdatePosition(lastPos, position);
    }
}

void MultiScreenWorker::onDisplayModeChanged(const DisplayMode &displayMode)
{
    if (displayMode == m_displayMode)
        return;

    qInfo() << "display mode changed:" << displayMode;

    m_displayMode = displayMode;

    DockItem::setDockDisplayMode(displayMode);
    qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(displayMode));

    QRect rect;
    if (m_hideMode == HideMode::KeepShowing || (m_hideMode == HideMode::SmartHide && m_hideState == HideState::Show)) {
        rect = dockRect(m_ds.current(), m_position, HideMode::KeepShowing, m_displayMode);
    } else {
        rect = dockRect(m_ds.current(), m_position, HideMode::KeepHidden, m_displayMode);
    }

    parent()->setFixedSize(rect.size());
    parent()->move(rect.topLeft());

    parent()->panel()->setFixedSize(rect.size());
    parent()->panel()->move(0, 0);
    parent()->panel()->setDisplayMode(m_displayMode);

    emit displayModeChanegd();
    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::onHideModeChanged(const HideMode &hideMode)
{
    if (m_hideMode == hideMode)
        return;

    qInfo() << "hidemode changed:" << hideMode;

    m_hideMode = hideMode;

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        displayAnimation(m_ds.current(), AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        displayAnimation(m_ds.current(), AniAction::Hide);
    }

    emit requestUpdateFrontendGeometry();
    emit requestNotifyWindowManager();
}

void MultiScreenWorker::onHideStateChanged(const Dock::HideState &state)
{
    if (state == Dock::Unknown)
        return;

    m_hideState = state;

    // 检查当前屏幕的当前位置是否允许显示,不允许需要更新显示信息(这里应该在函数外部就处理好,不应该走到这里)

    //TODO 这里是否存在屏幕找不到的问题，m_ds的当前屏幕是否可以做成实时同步的，公用一个指针？
    //TODO 这里真的有必要加以下代码吗，只是隐藏模式的切换，理论上不需要检查屏幕是否允许任务栏停靠
    QScreen *curScreen = DIS_INS->screen(m_ds.current());
    if (!DIS_INS->canDock(curScreen, m_position)) {
        m_ds.updateDockedScreen(getValidScreen(m_position));
    }

    qInfo() << "hidestate changed:" << m_hideMode << m_hideState;

    if (m_hideMode == HideMode::KeepShowing
            || ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Show)) {
        displayAnimation(m_ds.current(), AniAction::Show);
    } else if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) && m_hideState == HideState::Hide) {
        // 如果鼠标正在任务栏要显示的区域,就可以不用隐藏(相当于智能隐藏被唤醒一样)
        if (getDockShowGeometry(m_ds.current(), m_position, m_displayMode).contains(QCursor::pos()))
            return;

        displayAnimation(m_ds.current(), AniAction::Hide);
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
    const int realDockSize = int((m_displayMode == DisplayMode::Fashion ? m_dockInter->windowSizeFashion() + 2 * 10 /*上下的边距各10像素*/ : m_dockInter->windowSizeEfficient()) * qApp->devicePixelRatio());

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
    const int monitHeight = 100 + WINDOW_MARGIN;

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

void MultiScreenWorker::onRequestUpdateFrontendGeometry()
{
    HideMode hideMode = HideMode::KeepShowing;
    // 当设置为一直隐藏模式时，按照当前隐藏的状态去设置任务栏区域信息，否则，按照当前显示模式去设置
    if (m_hideMode == HideMode::KeepHidden) {
        hideMode = m_hideMode;
    }

    const QRect rect = dockRectWithoutScale(m_ds.current(), m_position, hideMode, m_displayMode);

#ifdef QT_DEBUG
    qDebug() << rect;
#endif

    m_dockInter->SetFrontendWindowRect(int(rect.x()), int(rect.y()), uint(rect.width()), uint(rect.height()));
    emit requestUpdateDockEntry();
}

void MultiScreenWorker::onRequestUpdateLayout()
{
    parent()->panel()->setFixedSize(dockRect(m_ds.current(), position(), HideMode::KeepShowing, displayMode()).size());
    parent()->panel()->move(0, 0);
    parent()->panel()->setDisplayMode(displayMode());
    parent()->panel()->setPositonValue(position());
    parent()->panel()->update();
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

/**
 * @brief MultiScreenWorker::getScreenRect
 * @param s
 * @param ratio
 * @return 返回屏幕 \s 的原始尺寸，\ratio 为屏幕 \s 的当前缩放值
 */
QRect MultiScreenWorker::getScreenRect(QScreen *s)
{
    if (!s)
        return QRect();

    auto geo = s->geometry();
    QRect screenRect;
    screenRect.setX(geo.x());
    screenRect.setY(geo.y());
    screenRect.setWidth(static_cast<int>(geo.width() * s->devicePixelRatio()));
    screenRect.setHeight(static_cast<int>(geo.height() * s->devicePixelRatio()));

    /* QScreen的handle会返回显示器缩放之前的分辨率
    * 比如1920*1080,调整缩放为1.25，那么QScreen::geometry返回的为QSize(1536, 864),QScreen::handle()->geometry()返回的为QSize(1920, 1080);
    *
    * @note 但在特殊情况下，比如频繁插拔显示器时，小概率出现handle返回的值并不等于缩放前的分辨率的情况，从而导致任务栏位置设置出错，
    * 这里仍然直接使用geometry*ratio的方式,去获取缩放前的分辨率
    * 但考虑到缩放值类型为double，可能在高分屏下计算得到的尺寸和通过handle函数的原始尺寸有1像素的误差，
    * 所以计算完毕后和handle比对一下，如果存在1像素的误差，则仍然使用handle函数，否则使用QScreen::geometry() * ratio
    */
    auto size = s->handle()->geometry().size() - screenRect.size();
    if (Q_LIKELY(size.width() <= 1 && size.height() <= 1)) {
        screenRect = s->handle()->geometry();
    }

    return screenRect;
}

/**
 * @brief 这里用到xcb去设置任务栏的高度，比较特殊，参考_NET_WM_STRUT_PARTIAL属性
 * 在屏幕旋转后，所有参数以控制中心自定义设置里主屏显示的图示为准（旋转不用特殊处理）
 */
void MultiScreenWorker::onRequestNotifyWindowManager()
{
    static QRect lastRect = QRect();
    static int lastScreenWidth = 0;
    static int lastScreenHeight = 0;

    /* 在非主屏或非一直显示状态时，清除任务栏区域，不挤占应用 */
    if ((!DIS_INS->isCopyMode() && m_ds.current() != m_ds.primary()) || m_hideMode != HideMode::KeepShowing) {
        lastRect = QRect();

        qInfo() << "clear wm struct";

        if (Utils::IS_WAYLAND_DISPLAY) {
            QList<QVariant> varList;
            varList.append(0);//dock位置
            varList.append(0);//dock高度/宽度
            varList.append(0);//start值
            varList.append(0);//end值
            if (parent()->windowHandle()->handle()) {
                QGuiApplication::platformNativeInterface()->setWindowProperty(parent()->windowHandle()->handle(),"_d_dwayland_dockstrut", varList);
            }
        } else {
            const auto display = QX11Info::display();
            if (!display) {
                qWarning() << "QX11Info::display() is " << display;
                return;
            }

            XcbMisc::instance()->clear_strut_partial(xcb_window_t(parent()->winId()));
        }

        return;
    }

    QRect dockGeometry = getDockShowGeometry(m_ds.current(), m_position, m_displayMode, true);
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
            varList[1] = dockGeometry.height();
            varList[2] = dockGeometry.x();
            varList[3] = dockGeometry.x() + dockGeometry.width();
            break;
        case Position::Bottom:
            varList[0] = 3;
            varList[1] = dockGeometry.height();
            varList[2] = dockGeometry.x();
            varList[3] = dockGeometry.x() + dockGeometry.width();
            break;
        case Position::Left:
            varList[0] = 0;
            varList[1] = dockGeometry.width();
            varList[2] = dockGeometry.y();
            varList[3] = dockGeometry.y() + dockGeometry.height();
            break;
        case Position::Right:
            varList[0] = 2;
            varList[1] = dockGeometry.width();
            varList[2] = dockGeometry.y();
            varList[3] = dockGeometry.y() + dockGeometry.height();
            break;
        }

        // 加上边距(时尚模式为10+10, 高效模式为0)
        varList[1] = varList[1].toDouble() + WINDOW_MARGIN * ratio * 2;
        qDebug() << "Dock strut: " << varList;

        if (parent()->windowHandle()->handle()) {
            QGuiApplication::platformNativeInterface()->setWindowProperty(parent()->windowHandle()->handle(),"_d_dwayland_dockstrut", varList);
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

        XcbMisc::instance()->set_strut_partial(static_cast<xcb_window_t>(parent()->winId()), orientation,
                                               static_cast<uint>(strut + WINDOW_MARGIN * ratio), // 设置窗口与屏幕边缘距离，需要乘缩放
                                               static_cast<uint>(strutStart),                   // 设置任务栏起点坐标（上下为x，左右为y）
                                               static_cast<uint>(strutEnd));                    // 设置任务栏终点坐标（上下为x，左右为y）
    }
}

void MultiScreenWorker::onRequestUpdatePosition(const Position &fromPos, const Position &toPos)
{
    qInfo() << "request change pos from: " << fromPos << " to: " << toPos;
    // 更新要切换到的屏幕
    if (!DIS_INS->canDock(DIS_INS->screen(m_ds.current()), m_position))
        m_ds.updateDockedScreen(getValidScreen(m_position));

    qInfo() << "update allow screen: " << m_ds.current();

    // 无论什么模式,都先显示
    changeDockPosition(m_ds.last(), m_ds.current(), fromPos, toPos);
}

void MultiScreenWorker::onRequestUpdateMonitorInfo()
{
    // resetDockScreen 调用太频繁，未在合适的时机调用时出现 wayland set_position 数值异常
    // 推测是 qt 未能正确传递窗口位置到 waylandserver， 此处修改经过测试可以较好的规避此问题。
    // for test ： DOCK_RESET_NOW=true dde-dock
    if (qEnvironmentVariableIsSet("DOCK_RESET_NOW"))
        resetDockScreen(); // m_monitorUpdateTimer timeout will call resetDockScreen

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
        parent()->setFixedSize(dockRect(m_ds.current()).size());
        parent()->setGeometry(dockRect(m_ds.current()));
        return;
    }

    m_ds.updateDockedScreen(m_delayScreen);

    // 检查当前屏幕的当前位置是否允许显示,不允许需要更新显示信息(这里应该在函数外部就处理好,不应该走到这里)
    // 检查边缘是否允许停靠
    QScreen *curScreen = DIS_INS->screen(m_delayScreen);
    if (curScreen && DIS_INS->canDock(curScreen, m_position)) {
        if (m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide) {
            displayAnimation(m_ds.current(), AniAction::Show);
        } else if (m_hideMode == HideMode::KeepShowing) {
            changeDockPosition(m_ds.last(), m_ds.current(), m_position, m_position);
        }
    }
}

MainWindow *MultiScreenWorker::parent()
{
    return static_cast<MainWindow *>(m_parent);
}

void MultiScreenWorker::initMembers()
{
    if (Utils::IS_WAYLAND_DISPLAY) {
        m_parent->createWinId();
    }
    m_monitorUpdateTimer->setInterval(100);
    m_monitorUpdateTimer->setSingleShot(true);

    m_delayWakeTimer->setSingleShot(true);

    m_delayDisplay->setSingleShot(true);
    // 延时显示配置
    DConfig *config = DConfig::create("org.deepin.dde.dock", "org.deepin.dde.dock");
    if (config->isValid() && config->keyList().contains("delayDisplay"))
        m_delayDisplay->setInterval(config->value("delayDisplay").toInt());
    delete config;

    setStates(LauncherDisplay, m_launcherInter->isValid() ? m_launcherInter->visible() : false);

    // init check
    checkDaemonDockService();
    checkXEventMonitorService();
}

void MultiScreenWorker::initConnection()
{
    connect(DIS_INS, &DisplayManager::primaryScreenChanged, this, &MultiScreenWorker::primaryScreenChanged);
    connect(DIS_INS, &DisplayManager::screenInfoChanged, this, &MultiScreenWorker::requestUpdateMonitorInfo);

    connect(m_launcherInter, static_cast<void (DBusLuncher::*)(bool) const>(&DBusLuncher::VisibleChanged), this, [ = ](bool value) { setStates(LauncherDisplay, value); });

    /** FIXME
     * 这里关联的信号有时候收不到是因为 qt-dbus-factory 中的 changed 的信号有时候会发不出来，
     * qt-dbus-factory 中的 DBusExtendedAbstractInterface::internalPropGet 在同步调用情况下，会将缓存中的数据写入属性中，
     * 导致后面 onPropertyChanged 中的判断认为属性值没变，就没有发出 changed 信号。
     * 建议：前端仅在初始化时主动获取一次 dbus 中的值存储在成员变量中，并建立 changed 信号连接，后面所有用到那个值的地方，均获取成员变量;
     * 或去修改 qt-dbus-factory，取消 DBusExtendedAbstractInterface::internalPropGet 中将数据写入属性值，
     * 但是 qt-dbus-factory 修改涉及面较广，需要大量测试确认没有问题，再合入。
     */
#if 0
    //    connect(m_dockInter, &DBusDock::PositionChanged, this, &MultiScreenWorker::onPositionChanged);
    //    connect(m_dockInter, &DBusDock::DisplayModeChanged, this, &MultiScreenWorker::onDisplayModeChanged);
    //    connect(m_dockInter, &DBusDock::HideModeChanged, this, &MultiScreenWorker::hideModeChanged);
    //    connect(m_dockInter, &DBusDock::HideStateChanged, this, &MultiScreenWorker::hideStateChanged);
#else
    QDBusConnection::sessionBus().connect("com.deepin.dde.daemon.Dock",
                                          "/com/deepin/dde/daemon/Dock",
                                          "org.freedesktop.DBus.Properties",
                                          "PropertiesChanged",
                                          "sa{sv}as",
                                          this, SLOT(handleDBusSignal(QDBusMessage)));
#endif

    connect(this, &MultiScreenWorker::requestUpdateFrontendGeometry, this, &MultiScreenWorker::onRequestUpdateFrontendGeometry);
    connect(this, &MultiScreenWorker::requestUpdatePosition, this, &MultiScreenWorker::onRequestUpdatePosition);
    connect(this, &MultiScreenWorker::requestNotifyWindowManager, this, &MultiScreenWorker::onRequestNotifyWindowManager);
    connect(this, &MultiScreenWorker::requestUpdateMonitorInfo, this, &MultiScreenWorker::onRequestUpdateMonitorInfo);

    connect(m_delayWakeTimer, &QTimer::timeout, this, &MultiScreenWorker::onRequestDelayShowDock);

    //　更新任务栏内容展示方式
    connect(this, &MultiScreenWorker::requestUpdateLayout, this, &MultiScreenWorker::onRequestUpdateLayout);

    // 刷新所有显示的内容，布局，方向，大小，位置等
    connect(m_monitorUpdateTimer, &QTimer::timeout, this, &MultiScreenWorker::updateDisplay);

    // 判断用户是否希望隐藏的任务栏显示
    connect(m_delayDisplay, &QTimer::timeout, this, [ = ] {
        setStates(CheckDockShouldDisplay, false);
        tryToShowDock(m_delayDisplayPos.x(), m_delayDisplayPos.y());
    });
}

void MultiScreenWorker::initUI()
{
    // 设置界面大小
    parent()->setFixedSize(dockRect(m_ds.current()).size());
    parent()->move(dockRect(m_ds.current()).topLeft());
    parent()->panel()->setFixedSize(dockRect(m_ds.current(), m_position, HideMode::KeepShowing, m_displayMode).size());
    parent()->panel()->move(0, 0);

    onPositionChanged(static_cast<Position>(dockInter()->position()));
    onDisplayModeChanged(static_cast<DisplayMode>(dockInter()->displayMode()));
    onHideModeChanged(static_cast<HideMode>(dockInter()->hideMode()));
    onHideStateChanged(static_cast<HideState>(dockInter()->hideState()));
    onOpacityChanged(m_dockInter->opacity());

    // 初始化透明度
    QTimer::singleShot(0, this, [ = ] {onOpacityChanged(m_dockInter->opacity());});
}

void MultiScreenWorker::initDBus()
{
    if (m_dockInter->isValid()) {
        m_position = static_cast<Dock::Position >(m_dockInter->position());
        m_hideMode = static_cast<Dock::HideMode >(m_dockInter->hideMode());
        m_hideState = static_cast<Dock::HideState >(m_dockInter->hideState());
        m_displayMode = static_cast<Dock::DisplayMode >(m_dockInter->displayMode());
        m_opacity = m_dockInter->opacity();

        DockItem::setDockPosition(m_position);
        qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));
        DockItem::setDockDisplayMode(m_displayMode);
        qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(m_displayMode));
    }
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
    initDBus();
    initDisplayData();
}

/**
 * @brief MultiScreenWorker::displayAnimation
 * 任务栏显示或隐藏过程的动画。
 * @param screen 任务栏要显示的屏幕
 * @param pos 任务栏显示的位置（上：0，右：1，下：2，左：3）
 * @param act 显示（隐藏）任务栏
 * @return void
 */
void MultiScreenWorker::displayAnimation(const QString &screen, const Position &pos, AniAction act)
{
    if (!testState(AutoHide) || qApp->property("DRAG_STATE").toBool()
            || testState(ChangePositionAnimationStart)
            || testState(HideAnimationStart)
            || testState(ShowAnimationStart))
        return;

    m_currentHideState = act ? HideState::Hide : HideState::Show;
    QRect mainwindowRect = parent()->geometry();
    QRect dockShowRect = getDockShowGeometry(screen, pos, m_displayMode);
    QRect dockHideRect = getDockHideGeometry(screen, pos, m_displayMode);

    /** FIXME
     * 在高分屏2.75倍缩放的情况下，parent()->geometry()返回的任务栏高度有问题（实际是40,返回是39）
     * 在这里增加判断，当返回值在范围（38，42）开区间内，均认为任务栏显示位置正确，直接返回，不执行动画
     * 也就是在实际值基础上下浮动1像素的误差范围
     * 正常屏幕情况下是没有这个问题的
     */
    switch (act) {
    case AniAction::Show:
        if (pos == Position::Top || pos == Position::Bottom) {
            if (qAbs(dockShowRect.height() - mainwindowRect.height()) <= 1
                    && mainwindowRect.contains(dockShowRect.center())) {
                emit requestNotifyWindowManager();
                return;
            }
        } else if (pos == Position::Left || pos == Position::Right) {
            if (qAbs(dockShowRect.width() - mainwindowRect.width()) <= 1
                    && mainwindowRect.contains(dockShowRect.center())) {
                emit requestNotifyWindowManager();
                return;
            }
        }
        break;
    case AniAction::Hide:
        if (dockHideRect == mainwindowRect) {
            emit requestNotifyWindowManager();
            return;
        }
        break;
    }

    QVariantAnimation *ani = new QVariantAnimation(this);
    ani->setEasingCurve(QEasingCurve::InOutCubic);

#ifndef DISABLE_SHOW_ANIMATION
    const bool composite = DWindowManagerHelper::instance()->hasComposite(); // 判断是否开启特效模式
    const int duration = composite ? ANIMATIONTIME : 0;
#else
    const int duration = 0;
#endif
    ani->setDuration(duration);

    ani->setStartValue(dockHideRect);
    ani->setEndValue(dockShowRect);

    switch (act) {
    case AniAction::Show:
        ani->setDirection(QAbstractAnimation::Forward);
        connect(ani, &QVariantAnimation::finished, this, &MultiScreenWorker::showAniFinished);
        connect(this, &MultiScreenWorker::requestStopShowAni, ani, &QVariantAnimation::stop);
        break;

    case AniAction::Hide:
        ani->setDirection(QAbstractAnimation::Backward); // 隐藏时动画反向走
        connect(ani, &QVariantAnimation::finished, this, &MultiScreenWorker::hideAniFinished);
        connect(this, &MultiScreenWorker::requestStopHideAni, ani, &QVariantAnimation::stop);
        break;
    }

    connect(ani, &QVariantAnimation::valueChanged, this, static_cast<void (MultiScreenWorker::*)(const QVariant &value)>(&MultiScreenWorker::updateParentGeometry));

    connect(ani, &QVariantAnimation::stateChanged, this, [ = ](QAbstractAnimation::State newState, QAbstractAnimation::State oldState) {
        // 更新动画是否正在进行的信号值
        switch (act) {
        case AniAction::Show:
            if (newState == QVariantAnimation::Running && oldState == QVariantAnimation::Stopped) {
                if (m_hideMode == HideMode::KeepShowing || duration)
                    setStates(ShowAnimationStart);
                else
                    setStates(DockIsShowing);
            }
            if (newState == QVariantAnimation::Stopped && oldState == QVariantAnimation::Running) {
                if (m_hideMode == HideMode::KeepShowing || duration)
                    setStates(ShowAnimationStart, false);
                else // 如果不是一直显示的状态，则让其延时修改状态，防止在resetDock的时候重复改变其高度引起任务栏闪烁导致无法唤醒
                    QTimer::singleShot(ANIMATIONTIME, [ = ] { setStates(DockIsShowing, false); });
            }
            break;
        case AniAction::Hide:
            if (newState == QVariantAnimation::Running && oldState == QVariantAnimation::Stopped) {
                setStates(HideAnimationStart);
            }
            if (newState == QVariantAnimation::Stopped && oldState == QVariantAnimation::Running) {
                setStates(HideAnimationStart, false);
            }
            break;
        }
    });

    parent()->panel()->setFixedSize(dockRect(m_ds.current(), m_position, HideMode::KeepShowing, m_displayMode).size());
    parent()->panel()->move(0, 0);

    emit requestStopShowAni();
    emit requestStopHideAni();
    emit requestUpdateLayout();

    ani->start(QVariantAnimation::DeleteWhenStopped);
}

/**
 * @brief MultiScreenWorker::displayAnimation
 * 任务栏显示或隐藏过程的动画。
 * @param screen 任务栏要显示的屏幕
 * @param act 显示（隐藏）任务栏
 * @return void
 */
void MultiScreenWorker::displayAnimation(const QString &screen, AniAction act)
{
    return displayAnimation(screen, m_position, act);
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

    // 更新屏幕信息
    m_ds.updateDockedScreen(toScreen);

    // TODO: 考虑切换过快的情况,这里需要停止上一次的动画,可增加信号控制,暂时无需要
    qInfo() << "from: " << fromScreen << "  to: " << toScreen;

    QSequentialAnimationGroup *group = new QSequentialAnimationGroup(this);

    QVariantAnimation *ani1 = new QVariantAnimation(group);
    QVariantAnimation *ani2 = new QVariantAnimation(group);

    //　初始化动画信息
    ani1->setEasingCurve(QEasingCurve::InOutCubic);
    ani2->setEasingCurve(QEasingCurve::InOutCubic);

    const bool composite = DWindowManagerHelper::instance()->hasComposite();
#ifndef DISABLE_SHOW_ANIMATION
    const int duration = composite ? ANIMATIONTIME : 0;
#else
    const int duration = 0;
#endif

    ani1->setDuration(duration);
    ani2->setDuration(duration);

    //　隐藏
    ani1->setStartValue(getDockShowGeometry(fromScreen, fromPos, m_displayMode));
    ani1->setEndValue(getDockHideGeometry(fromScreen, fromPos, m_displayMode));
    qDebug() << fromScreen << "hide from :" << getDockShowGeometry(fromScreen, fromPos, m_displayMode);
    qDebug() << fromScreen << "hide to   :" << getDockHideGeometry(fromScreen, fromPos, m_displayMode);

    //　显示
    ani2->setStartValue(getDockHideGeometry(toScreen, toPos, m_displayMode));
    ani2->setEndValue(getDockShowGeometry(toScreen, toPos, m_displayMode));
    qDebug() << toScreen << "show from :" << getDockHideGeometry(toScreen, toPos, m_displayMode);
    qDebug() << toScreen << "show to   :" << getDockShowGeometry(toScreen, toPos, m_displayMode);

    group->addAnimation(ani1);
    group->addAnimation(ani2);

    // 隐藏时固定一下内容大小
    connect(ani1, &QVariantAnimation::stateChanged, this, [ = ](QAbstractAnimation::State newState, QAbstractAnimation::State oldState) {
        if (newState == QVariantAnimation::Running && oldState == QVariantAnimation::Stopped) {
            parent()->panel()->setFixedSize(dockRect(fromScreen, fromPos, HideMode::KeepShowing, m_displayMode).size());
            parent()->panel()->move(0, 0);
        }
    });

    connect(ani1, &QVariantAnimation::valueChanged, this, [ = ](const QVariant & value) {
        updateParentGeometry(value, fromPos);
    });

    // 显示时固定一下内容大小
    connect(ani2, &QVariantAnimation::stateChanged, this, [ = ](QAbstractAnimation::State newState, QAbstractAnimation::State oldState) {
        // 位置发生变化时需要更新位置属性,且要在隐藏动画之后,显示动画之前
        DockItem::setDockPosition(m_position);
        qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));

        if (newState == QVariantAnimation::Running && oldState == QVariantAnimation::Stopped) {
            parent()->panel()->setFixedSize(dockRect(toScreen, toPos, HideMode::KeepShowing, m_displayMode).size());
            parent()->panel()->move(0, 0);
        }
    });

    connect(ani2, &QVariantAnimation::valueChanged, this, [ = ](const QVariant & value) {
        updateParentGeometry(value, toPos);
    });

    // 如果更改了显示位置，在显示之前应该更新一下界面布局方向
    if (fromPos != toPos)
        connect(ani1, &QVariantAnimation::finished, this, [ = ] {
            const auto display = QX11Info::display();
            if (!display) {
                qWarning() << "QX11Info::display() is " << display;
            } else {
                // 先清除原先的窗管任务栏区域
                //XcbMisc::instance()->clear_strut_partial(xcb_window_t(parent()->winId()));
            }

            // 隐藏后需要通知界面更新布局方向
            emit requestUpdateLayout();
        });

    connect(group, &QVariantAnimation::finished, this, [ = ] {
        setStates(ChangePositionAnimationStart, false);

        // 结束之后需要根据确定需要再隐藏
        emit showAniFinished();
        emit requestUpdateFrontendGeometry();
        emit requestNotifyWindowManager();
    });

    setStates(ChangePositionAnimationStart);

    group->start(QVariantAnimation::DeleteWhenStopped);
}

/**
 * @brief getValidScreen        获取一个当前任务栏可以停靠的屏幕，优先使用主屏
 * @return
 */
QString MultiScreenWorker::getValidScreen(const Position &pos)
{
    //TODO 考虑在主屏幕名变化时自动更新，是不是就不需要手动处理了
    m_ds.updatePrimary(DIS_INS->primary());

    if (m_screenMonitor->needUsedLastScreen())
        return m_screenMonitor->lastScreen();

    if (DIS_INS->canDock(DIS_INS->screen(m_ds.current()), pos))
        return m_ds.current();

    if (DIS_INS->canDock(qApp->primaryScreen(), pos))
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
            || testState(DockIsShowing))
        return;

    m_ds.updateDockedScreen(getValidScreen(position()));
    // 更新任务栏自身信息
    /**
      *注意这里要先对parent()进行setFixedSize，在分辨率切换过程中，setGeometry可能会导致其大小未改变
      */
    if (m_ds.current().isEmpty()) // wayland下当没有屏幕时连接时，qt会虚拟一个空的屏幕，空屏幕这里不处理，不然会导致崩溃
        return;
    parent()->setFixedSize(dockRect(m_ds.current()).size());
    parent()->setGeometry(dockRect(m_ds.current()));
    qDebug() << "update dock geometry: " << dockRect(m_ds.current());
    parent()->panel()->setFixedSize(dockRect(m_ds.current()).size());
    parent()->panel()->move(0, 0);
}

/**
 * @brief checkDaemonDockService
 * 避免com.deepin.dde.daemon.Dock服务比dock晚启动，导致dock启动后的状态错误
 */
void MultiScreenWorker::checkDaemonDockService()
{
    auto connectionInit = [ = ](DBusDock * dockInter) {
        connect(dockInter, &DBusDock::ServiceRestarted, this, [ = ] {
            resetDockScreen();

            emit requestUpdateFrontendGeometry();
        });
        connect(dockInter, &DBusDock::OpacityChanged, this, &MultiScreenWorker::onOpacityChanged);
        connect(dockInter, &DBusDock::WindowSizeEfficientChanged, this, &MultiScreenWorker::onWindowSizeChanged);
        connect(dockInter, &DBusDock::WindowSizeFashionChanged, this, &MultiScreenWorker::onWindowSizeChanged);
    };

    const QString serverName = "com.deepin.dde.daemon.Dock";
    QDBusConnectionInterface *ifc = QDBusConnection::sessionBus().interface();

    if (!ifc->isServiceRegistered(serverName)) {
        connect(ifc, &QDBusConnectionInterface::serviceOwnerChanged, this, [ = ](const QString & name, const QString & oldOwner, const QString & newOwner) {
            Q_UNUSED(oldOwner)
            if (name == serverName && !newOwner.isEmpty()) {
                FREE_POINT(m_dockInter);

                m_dockInter = new DBusDock(serverName, "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
                Q_EMIT notifyDaemonInterfaceUpdate();
                // connect
                connectionInit(m_dockInter);

                // reinit data
                reInitDisplayData();

                // operation
                onPositionChanged(static_cast<Position>(dockInter()->position()));
                onDisplayModeChanged(static_cast<DisplayMode>(dockInter()->displayMode()));
                onHideModeChanged(static_cast<HideMode>(dockInter()->hideMode()));
                onHideStateChanged(static_cast<HideState>(dockInter()->hideState()));
                onOpacityChanged(m_dockInter->opacity());

                disconnect(ifc);
            }
        });
    } else {
        connectionInit(m_dockInter);
    }
}

bool MultiScreenWorker::isCursorOut(int x, int y)
{
    const int realDockSize = int((m_displayMode == DisplayMode::Fashion ? m_dockInter->windowSizeFashion() + 2 * 10 /*上下的边距各10像素*/ : m_dockInter->windowSizeEfficient()) * qApp->devicePixelRatio());
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
 * 避免com.deepin.api.XEventMonitor服务比dock晚启动，导致dock启动后的状态错误
 */
void MultiScreenWorker::checkXEventMonitorService()
{
    auto connectionInit = [ = ](XEventMonitor * eventInter, XEventMonitor * extralEventInter, XEventMonitor * touchEventInter) {
        connect(eventInter, &XEventMonitor::CursorMove, this, &MultiScreenWorker::onRegionMonitorChanged);
        connect(eventInter, &XEventMonitor::ButtonPress, this, [ = ](int i, int x, int y, const QString & key) {
            Q_UNUSED(i);
            setStates(MousePress, true);

            if(m_registerKey == key)
                CheckShouldDisplay(x, y, key);
        });
        connect(eventInter, &XEventMonitor::ButtonRelease, this, [ = ](int i, int x, int y, const QString & key) {
            Q_UNUSED(i);
            setStates(MousePress, false);

            if(m_registerKey == key)
                CheckShouldDisplay(x, y, key);

            // 在鼠标拖动任务栏在最大尺寸的时候，如果当前任务栏为隐藏模式（一直隐藏或智能隐藏），在松开鼠标的那一刻，判断当前
            // 鼠标位置是否在任务栏外面(isCursorOut(x,y)==true)，此时需要触发onExtralRegionMonitorChanged函数，
            // 目的是为了让其触发隐藏
            QTimer::singleShot(50, this, [ = ] {
                if (isCursorOut(x, y)
                        && (m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide)) {
                    qApp->setProperty("DRAG_STATE", false);
                    onExtralRegionMonitorChanged(0, 0, m_extralRegisterKey);
                }
            });
        });

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

            if(m_delayDisplay->isActive()) {
                qWarning()<<"CursorOut  area....";
                m_delayDisplay->stop();
                setStates(CheckDockShouldDisplay, false);
            }
        });

        // 触屏时，后端只发送press、release消息，有move消息则为鼠标，press置false
        connect(touchEventInter, &XEventMonitor::CursorMove, this, [ = ] {
            setStates(Utils::IS_WAYLAND_DISPLAY? MousePress : TouchPress, false);
        });
        connect(touchEventInter, &XEventMonitor::ButtonPress, this, &MultiScreenWorker::onTouchPress);
        connect(touchEventInter, &XEventMonitor::ButtonRelease, this, &MultiScreenWorker::onTouchRelease);
    };

    const QString serverName = "com.deepin.api.XEventMonitor";
    QDBusConnectionInterface *ifc = QDBusConnection::sessionBus().interface();

    if (!ifc->isServiceRegistered(serverName)) {
        connect(ifc, &QDBusConnectionInterface::serviceOwnerChanged, this, [ = ](const QString & name, const QString & oldOwner, const QString & newOwner) {
            Q_UNUSED(oldOwner)
            if (name == serverName && !newOwner.isEmpty()) {
                FREE_POINT(m_eventInter);
                FREE_POINT(m_extralEventInter);
                FREE_POINT(m_touchEventInter);

                m_eventInter = new XEventMonitor(serverName, "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus());
                m_extralEventInter = new XEventMonitor(serverName, "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus());
                m_touchEventInter = new XEventMonitor(serverName, "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus());
                // connect
                connectionInit(m_eventInter, m_extralEventInter, m_touchEventInter);

                disconnect(ifc);
            }
        });
    } else {
        connectionInit(m_eventInter, m_extralEventInter, m_touchEventInter);
    }
}

/**
 * @brief MultiScreenWorker::getDockShowMinGeometry
 * @param screenName
 * @return 如果dock在当前screenName对应的屏幕上显示，返回其在最小 高/宽 度(也就是40,最大值为100)的rect区域
 */
QRect MultiScreenWorker::getDockShowMinGeometry(const QString &screenName)
{
    QRect rect;
    const double ratio = qApp->devicePixelRatio();
    const int margin = static_cast<int>((m_displayMode == DisplayMode::Fashion ? 10 : 0) * qApp->devicePixelRatio());
    const int dockSize = 40;

    for (auto s : DIS_INS->screens()) {
        if (s->name() == screenName) {
            const QRect screenRect = getScreenRect(s);
            switch (m_position) {
            case Position::Top:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(dockSize);
                break;
            case Position::Bottom:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + screenRect.height() / ratio - margin - dockSize));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(dockSize);
                break;
            case Position::Left:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(dockSize);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            case Position::Right:
                rect.setX(static_cast<int>(screenRect.x() + screenRect.width() / ratio - margin - dockSize));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(dockSize);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            }
        }
    }
    return rect;
}

bool MultiScreenWorker::launcherVisible()
{
    return m_launcherInter->isValid() ? m_launcherInter->visible() : false;
}

void MultiScreenWorker::setLauncherVisble(bool isVisible)
{
    if (m_launcherInter->isValid()) {
        isVisible ? m_launcherInter->Show() : m_launcherInter->Hide();
    }
}

/**
 * @brief 获取任务栏显示时的参数。目前多屏情况下缩放保持一致，如果后续缩放规则修改，这里需要重新调整
 *
 * @param screenName    当前屏幕名字
 * @param pos           任务栏位置
 * @param displaymode   任务栏显示模式
 * @param withoutScale  是否考虑缩放（true:获取的是真实值; false:获取的是前端认为的值(默认)）
 * @return QRect        任务栏参数
 */
QRect MultiScreenWorker::getDockShowGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale)
{
    QRect rect;
    const double ratio = withoutScale ? 1 : qApp->devicePixelRatio();
    const int margin = static_cast<int>((displaymode == DisplayMode::Fashion ? 10 : 0) * (withoutScale ? qApp->devicePixelRatio() : 1));
    const int dockSize = static_cast<int>((displaymode == DisplayMode::Fashion ? m_dockInter->windowSizeFashion() : m_dockInter->windowSizeEfficient()) * (withoutScale ? qApp->devicePixelRatio() : 1));

    for (auto s : DIS_INS->screens()) {
        if (s->name() == screenName) {
            const QRect screenRect = getScreenRect(s);
            switch (pos) {
            case Position::Top:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(dockSize);
                break;
            case Position::Bottom:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + screenRect.height() / ratio - margin - dockSize));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(dockSize);
                break;
            case Position::Left:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(dockSize);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            case Position::Right:
                rect.setX(static_cast<int>(screenRect.x() + screenRect.width() / ratio - margin - dockSize));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(dockSize);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            }
        }
    }
    return rect;
}

/**
 * @brief 获取任务栏隐藏时的参数。目前多屏情况下缩放保持一致，如果后续缩放规则修改，这里需要重新调整
 *
 * @param screenName    当前屏幕名字
 * @param pos           任务栏位置
 * @param displaymode   任务栏显示模式
 * @param withoutScale  是否考虑缩放（true:获取的是真实值; false:获取的是前端认为的值(默认)）
 * @return QRect        任务栏参数
 */
QRect MultiScreenWorker::getDockHideGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale)
{
    QRect rect;
    const double ratio = withoutScale ? 1 : qApp->devicePixelRatio();
    const int margin = static_cast<int>((displaymode == DisplayMode::Fashion ? 10 : 0) * (withoutScale ? qApp->devicePixelRatio() : 1));

    for (auto s : DIS_INS->screens()) {
        if (s->name() == screenName) {
            const QRect screenRect = getScreenRect(s);
            switch (pos) {
            case Position::Top:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(0);
                break;
            case Position::Bottom:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + screenRect.height() / ratio - margin));
                rect.setWidth(static_cast<int>(screenRect.width() / ratio - 2 * margin));
                rect.setHeight(0);
                break;
            case Position::Left:
                rect.setX(static_cast<int>(screenRect.x() + margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(0);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            case Position::Right:
                rect.setX(static_cast<int>(screenRect.x() + screenRect.width() / ratio - margin));
                rect.setY(static_cast<int>(screenRect.y() + margin));
                rect.setWidth(0);
                rect.setHeight(static_cast<int>(screenRect.height() / ratio - 2 * margin));
                break;
            }
        }
    }
    return rect;
}

QScreen *MultiScreenWorker::screenByName(const QString &screenName)
{
    foreach (QScreen *screen, qApp->screens()) {
        if (screen->name() == screenName)
            return screen;
    }
    return nullptr;
}

bool MultiScreenWorker::onScreenEdge(const QString &screenName, const QPoint &point)
{
    QScreen *screen = screenByName(screenName);
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
    case HideMode::KeepHidden:
    {
        // 这时候鼠标如果在任务栏上,就不能隐藏
        if (!parent()->geometry().contains(QCursor::pos()))
            displayAnimation(m_ds.current(), AniAction::Hide);
    }
        break;
    case HideMode::SmartHide: {
        if (m_hideState == HideState::Show)
        {
            displayAnimation(m_ds.current(), AniAction::Show);
        } else if (m_hideState == HideState::Hide) {
            displayAnimation(m_ds.current(), AniAction::Hide);
        }
    }
        break;
    case HideMode::KeepShowing:
        displayAnimation(m_ds.current(), AniAction::Show);
        break;
    }

}

/**
 * @brief tryToShowDock 根据xEvent监控区域信号的x，y坐标处理任务栏唤醒显示
 * @param eventX        监控信号x坐标
 * @param eventY        监控信号y坐标
 */
void MultiScreenWorker::tryToShowDock(int eventX, int eventY)
{
    if (qApp->property("DRAG_STATE").toBool() || testState(ChangePositionAnimationStart)) {
        qWarning() << "dock is draging or animation is running";
        return;
    }

    if (hideMode() != HideMode::KeepShowing && testState(CheckDockShouldDisplay)) {
//        qDebug() << "dock is hidden, check whether the user wants it to display";
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
    if (onScreenEdge(m_ds.current(), QPoint(eventX, eventY))) {
        toScreen = m_ds.current();
    }

    // 过滤重复坐标
    static QPoint lastPos(0, 0);
    if (lastPos == QPoint(eventX, eventY)) {
        return;
    }
    lastPos = QPoint(eventX, eventY);

    // 任务栏显示状态，但需要切换屏幕
    if (toScreen != m_ds.current()) {
        if (!m_delayWakeTimer->isActive()) {
            m_delayScreen = toScreen;
            m_delayWakeTimer->start(Utils::SettingValue("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", MonitorsSwitchTime, 2000).toInt());
        }
    } else {
        // 任务栏隐藏状态，但需要显示
        if (hideMode() == HideMode::KeepShowing) {
            parent()->setFixedSize(dockRect(m_ds.current()).size());
            parent()->setGeometry(dockRect(m_ds.current()));
            return;
        }

        if (testState(ShowAnimationStart)) {
            qDebug() << "animation is running";
            return;
        }

        const QRect boundRect = parent()->visibleRegion().boundingRect();
        qDebug() << "boundRect:" << boundRect;
        if ((m_hideMode == HideMode::KeepHidden || m_hideMode == HideMode::SmartHide)
                && (boundRect.size().isEmpty())) {
            displayAnimation(m_ds.current(), AniAction::Show);
        }
    }
}

/**
 * @brief 屏幕监视
 * @param ds
 */
#define TIMESPAN 1500
ScreenChangeMonitor::ScreenChangeMonitor(DockScreen *ds, QObject *parent)
    : QObject (parent)
{
    #define DATETOSTRING(date) date.toString("hh:mm:ss:zzz")

    connect(DIS_INS, &DisplayManager::primaryScreenChanged, this, [ this, ds ](QScreen *primaryScreen) {
        Q_ASSERT(primaryScreen);

        // 在screenAdded之后，又会发送一次主屏幕的变更的信息
        qInfo() << "primary screen changed, primary Screen" << primaryScreen->name();
        if (changedInSeconds())
            return;

        m_lastScreenName = ds->current();
        m_changePrimaryName = primaryScreen->name();
        m_changeTime = QDateTime::currentDateTime();
        qInfo() << "primary changed from: " << m_lastScreenName << "to: " << DATETOSTRING(m_changeTime);
    });
    connect(DIS_INS, &DisplayManager::screenAdded, this, [ this ](QScreen *newScreen) {
        if (newScreen) {
            m_newScreenName = newScreen->name();
            m_newTime = QDateTime::currentDateTime();
            qInfo() <<"screen added:" << m_newScreenName << ", time:" << DATETOSTRING(m_newTime);
        }
    });
    connect(DIS_INS, &DisplayManager::screenRemoved, this, [ this ](QScreen *rmScreen) {
        if (rmScreen) {
            m_removeScreenName = rmScreen->name();
            m_removeTime = QDateTime::currentDateTime();
            qInfo() <<"screen removed:" << m_removeScreenName << ", time:" << DATETOSTRING(m_removeTime);
        }
    });
}

ScreenChangeMonitor::~ScreenChangeMonitor()
{
}

bool ScreenChangeMonitor::changedInSeconds()
{
    // 如果上一次的主屏幕的名称为空，则认为未发送屏幕变化的信号
    if (m_changePrimaryName.isEmpty() || !m_changeTime.isValid())
        return false;

    qint64 nowMsec = QDateTime::currentDateTime().toMSecsSinceEpoch();
    qint64 preChangeMsec = m_changeTime.toMSecsSinceEpoch();
    return (nowMsec - preChangeMsec <= TIMESPAN);
}

bool ScreenChangeMonitor::needUsedLastScreen() const
{
    if (m_lastScreenName.isEmpty())
        return false;

    // 判断是否在关闭显示器的时候由系统发出来的先禁用显示器再开启显示器
    if (!m_changeTime.isValid() || !m_removeTime.isValid() || !m_newTime.isValid())
        return false;

    // 如果禁用的掉的显示器和开启的显示器的名称不一样，则认为他们是正常的插拔操作
    if (m_removeScreenName != m_newScreenName)
        return false;

    // 如果先插入显示器再拔掉显示器，则认为他们也是正常的插拔操作
    if (m_removeTime > m_newTime)
        return false;

    // 只有在remove显示器在add显示器之前且时间小于1.5秒，才认为他们是由关闭显示器引起的
    qint64 changeMsec = m_changeTime.toMSecsSinceEpoch();

    // 获取现在的时间，如果现在的时间大于这三个时间一定的值，则认为是手动操作
    qint64 nowMsec = QDateTime::currentDateTime().toMSecsSinceEpoch();
    if (nowMsec - changeMsec > TIMESPAN)
        return false;

    qint64 addedMsec = m_newTime.toMSecsSinceEpoch();
    qint64 removeMsec = m_removeTime.toMSecsSinceEpoch();
    return ((addedMsec - removeMsec <= TIMESPAN) && qAbs(removeMsec - changeMsec) <= TIMESPAN);
}

const QString ScreenChangeMonitor::lastScreen()
{
    return m_lastScreenName;
}
