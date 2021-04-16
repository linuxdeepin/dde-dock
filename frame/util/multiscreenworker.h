/*
 * Copyright (C) 2018 ~ 2028 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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

#ifndef MULTISCREENWORKER_H
#define MULTISCREENWORKER_H

#include "constants.h"
#include "monitor.h"
#include "utils.h"
#include "dockitem.h"
#include "xcb_misc.h"

#include <com_deepin_dde_daemon_dock.h>
#include <com_deepin_daemon_display.h>
#include <com_deepin_daemon_display_monitor.h>
#include <com_deepin_api_xeventmonitor.h>
#include <com_deepin_dde_launcher.h>

#include <DWindowManagerHelper>

#include <QObject>

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
using DBusDock = com::deepin::dde::daemon::Dock;
using DisplayInter = com::deepin::daemon::Display;
using MonitorInter = com::deepin::daemon::display::Monitor;
using XEventMonitor = ::com::deepin::api::XEventMonitor;
using DBusLuncher = ::com::deepin::dde::Launcher;

using namespace Dock;
class QVariantAnimation;
class QWidget;
class QTimer;
class MainWindow;
class QGSettings;

/**
 * @brief The DockScreen class
 * 保存任务栏的屏幕信息
 */
class DockScreen
{
public:
    explicit DockScreen(const QString &primary)
        : m_currentScreen(primary)
        , m_lastScreen(primary)
        , m_primary(primary)
    {}
    inline const QString &current() {return m_currentScreen;}
    inline const QString &last() {return m_lastScreen;}
    inline const QString &primary() {return m_primary;}

    void updateDockedScreen(const QString &screenName)
    {
        m_lastScreen = m_currentScreen;
        m_currentScreen = screenName;
    }

    void updatePrimary(const QString &primary)
    {
        m_primary = primary;
        if (m_currentScreen.isEmpty()) {
            updateDockedScreen(primary);
        }
    }

private:
    QString m_currentScreen;
    QString m_lastScreen;
    QString m_primary;
};

class MultiScreenWorker : public QObject
{
    Q_OBJECT
public:
    enum Flag {
        Motion = 1 << 0,
        Button = 1 << 1,
        Key    = 1 << 2
    };

    enum AniAction {
        Show = 0,
        Hide
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

     Q_DECLARE_FLAGS(RunStates, RunState)

    MultiScreenWorker(QWidget *parent, DWindowManagerHelper *helper);
    ~MultiScreenWorker();

    void initShow();

    DBusDock *dockInter() { return m_dockInter; }

    inline bool testState(RunState state) { return (m_state & state); }
    void setStates(RunStates state, bool on = true);

    inline const QString &lastScreen() { return m_ds.last(); }
    inline const QString &deskScreen() { return m_ds.current(); }
    inline const Position &position() { return m_position; }
    inline const DisplayMode &displayMode() { return m_displayMode; }
    inline const HideMode &hideMode() { return m_hideMode; }
    inline const HideState &hideState() { return m_hideState; }
    inline quint8 opacity() { return m_opacity * 255; }

    QRect dockRect(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode);
    QRect dockRect(const QString &screenName);

signals:
    void opacityChanged(const quint8 value) const;
    void displayModeChanegd();

    // 更新监视区域
    void requestUpdateRegionMonitor();                          // 更新监听区域
    void requestUpdateFrontendGeometry();                       //!!! 给后端的区域不能为是或宽度为0的区域,否则会带来HideState死循环切换的bug
    void requestNotifyWindowManager();
    void requestUpdatePosition(const Position &fromPos, const Position &toPos);
    void requestUpdateLayout();                                 //　界面需要根据任务栏更新布局的方向
    void requestUpdateDragArea();                               //　更新拖拽区域
    void requestUpdateMonitorInfo();                            //　屏幕信息发生变化，需要更新任务栏大小，拖拽区域，所在屏幕，监控区域，通知窗管，通知后端，

    void requestStopShowAni();
    void requestStopHideAni();

    void requestUpdateDockEntry();

public slots:
    void onAutoHideChanged(bool autoHide);
    void updateDaemonDockSize(int dockSize);
    void onRequestUpdateRegionMonitor();
    void handleDbusSignal(QDBusMessage);

private slots:
    // Region Monitor
    void onRegionMonitorChanged(int x, int y, const QString &key);
    void onExtralRegionMonitorChanged(int x, int y, const QString &key);

    // Animation
    void showAniFinished();
    void hideAniFinished();

    void onWindowSizeChanged(uint value);
    void primaryScreenChanged();
    void updateParentGeometry(const QVariant &value, const Position &pos);
    void updateParentGeometry(const QVariant &value);

    // 任务栏属性变化
    void onPositionChanged();
    void onDisplayModeChanged();
    void onHideModeChanged();
    void onHideStateChanged();
    void onOpacityChanged(const double value);

    // 通知后端任务栏所在位置
    void onRequestUpdateFrontendGeometry();

    void onRequestNotifyWindowManager();
    void onRequestUpdatePosition(const Position &fromPos, const Position &toPos);
    void onRequestUpdateMonitorInfo();
    void onRequestDelayShowDock();

    void onTouchPress(int type, int x, int y, const QString &key);
    void onTouchRelease(int type, int x, int y, const QString &key);

private:
    MainWindow *parent();
    // 初始化数据信息
    void initMembers();
    void initDBus();
    void initConnection();
    void initUI();
    void initDisplayData();
    void reInitDisplayData();

    void displayAnimation(const QString &screen, const Position &pos, AniAction act);
    void displayAnimation(const QString &screen, AniAction act);

    void tryToShowDock(int eventX, int eventY);
    void changeDockPosition(QString lastScreen, QString deskScreen, const Position &fromPos, const Position &toPos);

    void updateDockScreenName(const QString &screenName);
    QString getValidScreen(const Position &pos);
    void resetDockScreen();

    void checkDaemonDockService();
    void checkXEventMonitorService();

    QRect dockRectWithoutScale(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode);

    QRect getDockShowGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);
    QRect getDockHideGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);

    QScreen *screenByName(const QString &screenName);
    bool onScreenEdge(const QString &screenName, const QPoint &point);
    const QPoint rawXPosition(const QPoint &scaledPos);

private:
    QWidget *m_parent;
    DWindowManagerHelper *m_wmHelper;

    // monitor screen
    XEventMonitor *m_eventInter;
    XEventMonitor *m_extralEventInter;
    XEventMonitor *m_touchEventInter;

    // DBus interface
    DBusDock *m_dockInter;
    DBusLuncher *m_launcherInter;

    // update monitor info
    QTimer *m_monitorUpdateTimer;
    QTimer *m_delayWakeTimer;                   // sp3需求，切换屏幕显示延时，默认2秒唤起任务栏

    DockScreen m_ds;                            // 屏幕名称信息

    // 任务栏属性
    double m_opacity;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;

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
