// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MULTISCREENWORKER_H
#define MULTISCREENWORKER_H

#include "constants.h"
#include "utils.h"
#include "dockitem.h"
#include "xcb_misc.h"

#include <com_deepin_dde_daemon_dock.h>
#include <com_deepin_api_xeventmonitor.h>
#include <com_deepin_dde_launcher.h>

#include <QObject>
#include <QFlag>
#include <QDebug>

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
using XEventMonitor = ::com::deepin::api::XEventMonitor;
using DBusLuncher = ::com::deepin::dde::Launcher;

using namespace Dock;
class QVariantAnimation;
class QWidget;
class QTimer;
class MainWindow;
class QGSettings;
class ScreenChangeMonitor;

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
    inline const QString &current() const {return m_currentScreen;}
    inline const QString &last() const {return m_lastScreen;}
    inline const QString &primary() const {return m_primary;}

    void updateDockedScreen(const QString &screenName)
    {
        qInfo() << "update docked screen" << screenName;
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
        DockIsShowing = 0x80,               // 任务栏正在显示
        CheckDockShouldDisplay = 0x100,     // 任务栏已经隐藏时，正在判断是否需要显示


        // 如果要添加新的状态，可以在上面添加
        RunState_Mask = 0xffffffff,
    };

    typedef QFlags<RunState> RunStates;

    MultiScreenWorker(QWidget *parent);

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
    QRect getDockShowMinGeometry(const QString &screenName);

    bool launcherVisible();
    void setLauncherVisble(bool isVisible);

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
    void notifyDaemonInterfaceUpdate();

public slots:
    void onAutoHideChanged(bool autoHide);
    void updateDaemonDockSize(int dockSize);
    void onRequestUpdateRegionMonitor();

private slots:
    void handleDBusSignal(QDBusMessage);

    // Region Monitor
    void onRegionMonitorChanged(int x, int y, const QString &key);
    void onExtralRegionMonitorChanged(int x, int y, const QString &key);
    void CheckShouldDisplay(int x, int y, const QString &key);

    // Animation
    void showAniFinished();
    void hideAniFinished();

    void updateDisplay();

    void onWindowSizeChanged(uint value);
    void primaryScreenChanged(QScreen *screen);
    void updateParentGeometry(const QVariant &value, const Position &pos);
    void updateParentGeometry(const QVariant &value);

    // 任务栏属性变化
    void onPositionChanged(const Position &position);
    void onDisplayModeChanged(const DisplayMode &displayMode);
    void onHideModeChanged(const HideMode &hideMode);
    void onHideStateChanged(const Dock::HideState &state);
    void onOpacityChanged(const double value);

    // 通知后端任务栏所在位置
    void onRequestUpdateFrontendGeometry();

    void onRequestUpdateLayout();
    void onRequestNotifyWindowManager();
    void onRequestUpdatePosition(const Position &fromPos, const Position &toPos);
    void onRequestUpdateMonitorInfo();
    void onRequestDelayShowDock();

    // 触摸手势操作
    void onTouchPress(int type, int x, int y, const QString &key);
    void onTouchRelease(int type, int x, int y, const QString &key);

    void onDelayAutoHideChanged();

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
    void changeDockPosition(QString fromScreen, QString toScreen, const Position &fromPos, const Position &toPos);

    QString getValidScreen(const Position &pos);
    void resetDockScreen();

    void checkDaemonDockService();
    void checkXEventMonitorService();

    QRect dockRectWithoutScale(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode);

    QRect getDockShowGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);
    QRect getDockHideGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);
    bool isCursorOut(int x, int y);

    QScreen *screenByName(const QString &screenName);
    bool onScreenEdge(const QString &screenName, const QPoint &point);
    const QPoint rawXPosition(const QPoint &scaledPos);
    static bool isCopyMode();
    QRect getScreenRect(QScreen *s);

private:
    QWidget *m_parent;

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
    QTimer *m_delayDisplay;                     // 任务栏隐藏时，鼠标移动到屏幕边缘，延时唤起任务栏
    QPoint m_delayDisplayPos;

    DockScreen m_ds;                            // 屏幕名称信息
    ScreenChangeMonitor *m_screenMonitor;       // 用于监视屏幕是否为系统先拔再插

    // 任务栏属性
    double m_opacity;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    HideState m_currentHideState;
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

/**
 * 在控制中心设置了自动关闭显示器，等一段时间后，显示器会自动关闭。在多显示器的情况下，此时系统会自动禁用主屏
 * 马上又开启主屏（底层系统是这么设计的），唤醒后，在禁用主屏后，后端会改变主屏为另外一个显示器，再开启主屏，
 * 后端又改变主屏为原来的主屏，这样就会导致的问题是：本来任务栏在主屏，由于瞬间触发了删除主屏幕的操作，引起了
 * 任务栏显示到副屏的问题，开启主屏后，由于任务栏已经在副屏了，就再也回不到主屏，因此，增加了这个类用于专门来
 * 处理这种异常情况
 */
class ScreenChangeMonitor : public QObject
{
    Q_OBJECT

public:
    ScreenChangeMonitor(DockScreen *ds, QObject *parent);
    ~ScreenChangeMonitor();

    bool needUsedLastScreen() const;
    const QString lastScreen();

private:
    bool changedInSeconds();

private:
    QString m_lastScreenName;

    QString m_changePrimaryName;
    QDateTime m_changeTime;

    QString m_newScreenName;
    QDateTime m_newTime;

    QString m_removeScreenName;
    QDateTime m_removeTime;
};

#endif // MULTISCREENWORKER_H
