// Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MULTISCREENWORKER_H
#define MULTISCREENWORKER_H

#include "constants.h"
#include "utils.h"
#include "dockitem.h"
#include "xcb_misc.h"
#include "dbusutil.h"

#include "org_deepin_dde_xeventmonitor1.h"
#include "org_deepin_dde_launcher1.h"
#include "org_deepin_dde_appearance1.h"

#include <DWindowManagerHelper>

#include <QObject>
#include <QFlag>

#define WINDOWMARGIN ((m_displayMode == Dock::Efficient) ? 0 : 10)
#define ANIMATIONTIME 300
#define FREE_POINT(p) if (p) {\
        delete p;\
        p = nullptr;\
    }\

DGUI_USE_NAMESPACE
/**
 * 多屏功能这部分看着很复杂，其实只需要把握住一个核心：及时更新数据！
 * 之前测试出的诸多问题都是在切换任务栏位置，切换屏幕，主屏更改，分辨率更改等情况发生后
 * 任务栏的鼠标唤醒区域或任务栏的大小没更新或者更新时的大小还是按照原来的屏幕信息计算而来的，
 */

#define DRAG_AREA_SIZE (5)
#define DOCKSPACE (WINDOWMARGIN * 2)

using XEventMonitor = ::org::deepin::dde::XEventMonitor1;
using DBusLuncher = ::org::deepin::dde::Launcher1;
using Appearance = org::deepin::dde::Appearance1;

using namespace Dock;
class QVariantAnimation;
class QWidget;
class QTimer;
class MainWindow;
class QGSettings;
class TrayMainWindow;
class MenuWorker;

class MultiScreenWorker : public QObject
{
    Q_OBJECT

public:
    enum Flag {
        Motion = 1 << 0,
        Button = 1 << 1,
        Key    = 1 << 2
    };

    enum RunState {
        ShowAnimationStart = 0x1,           // 单次显示动画正在运行状态
        HideAnimationStart = 0x2,           // 单次隐藏动画正在运行状态
        ChangePositionAnimationStart = 0x4, // 任务栏切换位置动画运行状态
        AutoHide = 0x8,                     // 和MenuWorker保持一致,未设置此state时表示菜单已经打开
        MousePress = 0x10,                  // 当前鼠标是否被按下
        TouchPress = 0x20,                  // 当前触摸屏下是否按下
        LauncherDisplay = 0x40,             // 启动器是否显示

        // 如果要添加新的状态，可以在上面添加
        RunState_Mask = 0xffffffff,
    };

    typedef QFlags<RunState> RunStates;

    explicit MultiScreenWorker(QObject *parent = Q_NULLPTR);
    ~MultiScreenWorker() override;

    void updateDaemonDockSize(const int &dockSize);

    inline bool testState(RunState state) { return (m_state & state); }
    void setStates(RunStates state, bool on = true);

    inline const Position &position() { return m_position; }
    inline const DisplayMode &displayMode() { return m_displayMode; }
    inline const HideMode &hideMode() { return m_hideMode; }
    inline const HideState &hideState() { return m_hideState; }
    inline quint8 opacity() { return m_opacity * 255; }

signals:
    void opacityChanged(const quint8 value) const;
    void displayModeChanged(const Dock::DisplayMode &);

    // 更新监视区域
    void requestUpdateRegionMonitor();                          // 更新监听区域
    void requestUpdateFrontendGeometry();                       //!!! 给后端的区域不能为是或宽度为0的区域,否则会带来HideState死循环切换的bug
    void requestNotifyWindowManager();
    void requestUpdatePosition(const Position &fromPos, const Position &toPos);
    void requestUpdateMonitorInfo();                            //　屏幕信息发生变化，需要更新任务栏大小，拖拽区域，所在屏幕，监控区域，通知窗管，通知后端，

    // 用来通知WindowManager的信号
    void requestUpdateDockGeometry(const Dock::HideMode &hideMode);
    void positionChanged(const Dock::Position &position);

    void requestPlayAnimation(const QString &screenName, const Position &position, const Dock::AniAction &animation, bool containMouse = false, bool updatePos = false);
    void requestChangeDockPosition(const QString &fromScreen, const QString &toScreen, const Position &fromPos, const Position &toPos);

    // 服务重新启动的信号
    void serviceRestart();

public slots:
    void onAutoHideChanged(const bool autoHide);
    void onRequestUpdateRegionMonitor();

private slots:
    // Region Monitor
    void onRegionMonitorChanged(int x, int y, const QString &key);
    void onExtralRegionMonitorChanged(int x, int y, const QString &key);

    void updateDisplay();

    void onWindowSizeChanged(uint value);
    void onPrimaryScreenChanged();

    // 任务栏属性变化
    void onPositionChanged(int position);
    void onDisplayModeChanged(int displayMode);
    void onHideModeChanged(int hideMode);
    void onHideStateChanged(int state);
    void onOpacityChanged(const double value);

    void onRequestUpdatePosition(const Position &fromPos, const Position &toPos);
    void onRequestUpdateMonitorInfo();
    void onRequestDelayShowDock();

    // 触摸手势操作
    void onTouchPress(int type, int x, int y, const QString &key);
    void onTouchRelease(int type, int x, int y, const QString &key);

    void onDelayAutoHideChanged();

private:
    // 初始化数据信息
    void initMembers();
    void initDockMode();
    void initConnection();

    void initDisplayData();
    void reInitDisplayData();

    void tryToShowDock(int eventX, int eventY);
    void changeDockPosition(QString fromScreen, QString toScreen, const Position &fromPos, const Position &toPos);

    void resetDockScreen();

    void checkXEventMonitorService();

    QString getValidScreen(const Position &pos);

    bool isCursorOut(int x, int y);

    bool onScreenEdge(const QString &screenName, const QPoint &point);
    const QPoint rawXPosition(const QPoint &scaledPos);
    static bool isCopyMode();

private:
    // monitor screen
    XEventMonitor *m_eventInter;
    XEventMonitor *m_extralEventInter;
    XEventMonitor *m_touchEventInter;

    // DBus interface
    DBusLuncher *m_launcherInter;
    Appearance *m_appearanceInter;

    // update monitor info
    QTimer *m_monitorUpdateTimer;
    QTimer *m_delayWakeTimer;                   // sp3需求，切换屏幕显示延时，默认2秒唤起任务栏

    // 任务栏属性
    double m_opacity;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;
    uint m_windowFashionSize;
    uint m_windowEfficientSize;

    /***************不和其他流程产生交互,尽量不要动这里的变量***************/
    QString m_registerKey;
    QString m_extralRegisterKey;
    QString m_touchRegisterKey;                 // 触控屏唤起任务栏监控区域key
    QPoint m_touchPos;                          // 触屏按下坐标
    QList<MonitRect> m_monitorRectList;         // 监听唤起任务栏区域
    QList<MonitRect> m_extralRectList;          // 任务栏外部区域,随m_monitorRectList一起更新
    QList<MonitRect> m_touchRectList;           // 监听触屏唤起任务栏区域
    QString m_delayScreen;                      // 任务栏将要切换到的屏幕名
    RunStates m_state;
    /*****************************************************************/
};

#endif // MULTISCREENWORKER_H
