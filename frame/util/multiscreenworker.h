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
#include "item/dockitem.h"

#include "xcb/xcb_misc.h"

#include <com_deepin_dde_daemon_dock.h>
#include <com_deepin_daemon_display.h>
#include <com_deepin_daemon_display_monitor.h>
#include <com_deepin_api_xeventmonitor.h>
#include <com_deepin_dde_launcher.h>

#include <DWindowManagerHelper>

#include <QObject>

#define WINDOWMARGIN ((m_displayMode == Dock::Efficient) ? 0 : 10)
#define ANIMATIONTIME 300

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

class DockScreen : QObject
{
    Q_OBJECT
public:
    explicit DockScreen(const QString &current, const QString &last, const QString &primary)
        : m_currentScreen(current)
        , m_lastScreen(last)
        , m_primary(primary)
    {}
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

    MultiScreenWorker(QWidget *parent, DWindowManagerHelper *helper);
    ~MultiScreenWorker();

    void initShow();

    DBusDock *dockInter() {return m_dockInter;}

    /**
     * @brief lastScreen
     * @return                      任务栏上次所在的屏幕
     */
    inline const QString &lastScreen() {return m_ds.last();/*return m_lastScreen;*/}
    /**
     * @brief deskScreen
     * @return                      任务栏目标屏幕.可以理解为任务栏当前所在屏幕
     */
    inline const QString &deskScreen() {return m_ds.current();/*return m_currentScreen;*/}
    /**
     * @brief position
     * @return                      任务栏所在方向(上下左右)
     */
    inline const Position &position() {return m_position;}
    /**
     * @brief displayMode
     * @return                      任务栏显示模式(时尚模式,高效模式)
     */
    inline const DisplayMode &displayMode() {return m_displayMode;}
    /**
     * @brief hideMode
     * @return                      任务栏状态(一直显示,一直隐藏,智能隐藏)
     */
    inline const HideMode &hideMode() {return m_hideMode;}
    /**
     * @brief hideState
     * @return                      任务栏的智能隐藏时的一个状态值(1显示,2隐藏,其他不处理)
     */
    inline const HideState &hideState() {return m_hideState;}
    /**
     * @brief opacity
     * @return                      任务栏透明度
     */
    inline quint8 opacity() {return m_opacity * 255;}
    /**
     * @brief dockRect
     * @param screenName            屏幕名
     * @param pos                   任务栏位置
     * @param hideMode              模式
     * @param displayMode           状态
     * @return                      按照给定的数据计算出任务栏所在位置
     */
    QRect dockRect(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode);
    /**
     * @brief dockRect
     * @param screenName        屏幕名
     * @return                  按照当前屏幕的当前属性给出任务栏所在区域
     */
    QRect dockRect(const QString &screenName);
    /**
     * @brief realDockRect      给出不计算缩放情况的区域信息(和后端接口保持一致)
     * @param screenName        屏幕名
     * @param pos               任务栏位置
     * @param hideMode          模式
     * @param displayMode       状态
     * @return
     */
    QRect dockRectWithoutScale(const QString &screenName, const Position &pos, const HideMode &hideMode, const DisplayMode &displayMode);

signals:
    void opacityChanged(const quint8 value) const;
    void displayModeChanegd();

    // 更新监视区域
    void requestUpdateRegionMonitor();
    void requestUpdateFrontendGeometry();                       //!!! 给后端的区域不能为是或宽度为0的区域,否则会带来HideState死循环切换的bug
    void requestNotifyWindowManager();
    void requestUpdatePosition(const Position &fromPos, const Position &toPos);
    void requestUpdateLayout();                                 //　界面需要根据任务栏更新布局的方向
    void requestUpdateDragArea();                               //　更新拖拽区域
    void requestUpdateMonitorInfo();                            //　屏幕信息发生变化，需要更新任务栏大小，拖拽区域，所在屏幕，监控区域，通知窗管，通知后端，
    void requestDelayShowDock(const QString &screenName);       //　延时唤醒任务栏

public slots:
    void onAutoHideChanged(bool autoHide);
    /**
     * @brief updateDaemonDockSize
     * @param dockSize              这里的高度是通过qt获取的，不能使用后端的接口数据
     */
    void updateDaemonDockSize(int dockSize);
    void onDragStateChanged(bool draging);

    void handleDbusSignal(QDBusMessage);

    void updateTouchRegisterRegion(const QRect &rect);

private slots:
    // Region Monitor
    void onRegionMonitorChanged(int x, int y, const QString &key);
    void onExtralRegionMonitorChanged(int x, int y, const QString &key);

    // Display Monitor
    void onMonitorListChanged(const QList<QDBusObjectPath> &mons);
    void monitorAdded(const QString &path);
    void monitorRemoved(const QString &path);

    // Animation
    void showAniFinished();
    void hideAniFinished();

    void onWindowSizeChanged(uint value);
    void primaryScreenChanged();
    void updateParentGeometry(const QVariant &value, const Position &pos);
    void updateParentGeometry(const QVariant &value);
    void delayShowDock();

    // 任务栏属性变化
    void onPositionChanged();
    void onDisplayModeChanged();
    void onHideModeChanged();
    void onHideStateChanged();
    void onOpacityChanged(const double value);

    /**
     * @brief onRequestUpdateRegionMonitor  更新监听区域信息
     * 触发时机:屏幕大小,屏幕坐标,屏幕数量,发生变化
     *          任务栏位置发生变化
     *          任务栏'模式'发生变化
     */
    void onRequestUpdateRegionMonitor();

    // 通知后端任务栏所在位置
    void onRequestUpdateFrontendGeometry();

    void onRequestNotifyWindowManager();
    void onRequestUpdatePosition(const Position &fromPos, const Position &toPos);
    void onRequestUpdateMonitorInfo();
    void onRequestDelayShowDock(const QString &screenName);

    void updateMonitorDockedInfo();

    void onTouchPress(int type, int x, int y, const QString &key);
    void onTouchRelease(int type, int x, int y, const QString &key);

private:
    // 初始化数据信息
    void initMembers();
    void initConnection();
    void initUI();
    /**
     * @brief showAni   任务栏显示动画
     * @param screen    显示到目标屏幕上
     */
    void showAni(const QString &screen);
    /**
     * @brief tryToShowDock 根据xEvent监控区域信号的x，y坐标处理任务栏唤醒显示
     * @param eventX        监控信号x坐标
     * @param eventY        监控信号y坐标
     */
    void tryToShowDock(int eventX, int eventY);
    /**
     * @brief hideAni   任务栏隐藏动画
     * @param screen    从目标屏幕上隐藏
     */
    void hideAni(const QString &screen);
    /**
     * @brief changeDockPosition    做一个动画操作
     * @param lastScreen            上次任务栏所在的屏幕
     * @param deskScreen            任务栏要移动到的屏幕
     * @param fromPos               任务栏上次的方向
     * @param toPos                 任务栏打算移动到的方向
     */
    void changeDockPosition(QString lastScreen, QString deskScreen, const Position &fromPos, const Position &toPos);
    /**
     * @brief updateDockScreenName  将任务栏所在屏幕信息进行更新,在任务栏切换屏幕显示后,这里应该被调用
     * @param screenName            目标屏幕
     */
    void updateDockScreenName(const QString &screenName);
    /**
     * @brief getValidScreen        获取一个当前任务栏可以停靠的屏幕，优先使用主屏
     * @return
     */
    QString getValidScreen(const Position &pos);
    /**
     * @brief autosetDockScreen     检查一下当前屏幕所在边缘是够允许任务栏停靠，不允许的情况需要更换下一块屏幕
     */
    void autosetDockScreen();

    void checkDaemonDockService();

    MainWindow *parent();
    /**
     * @brief getDockShowGeometry       获取任务栏显示时的位置
     * @param screenName                当前屏幕名
     * @param pos                       任务栏位置
     * @param displaymode               显示模式
     * @param real                      接口是否计算算法,false为计算,true为不计算(与后端接口大小的形式保持一致.比如,后端给出的屏幕大小是不计算缩放的)
     * @return
     */
    QRect getDockShowGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);
    QRect getDockHideGeometry(const QString &screenName, const Position &pos, const DisplayMode &displaymode, bool withoutScale = false);

    Monitor *monitorByName(const QList<Monitor *> &list, const QString &screenName);
    QScreen *screenByName(const QString &screenName);
    qreal scaleByName(const QString &screenName);
    bool onScreenEdge(const QString &screenName, const QPoint &point);
    bool onScreenEdge(const QPoint &point);
    bool contains(const MonitRect &rect, const QPoint &pos);
    bool contains(const QList<MonitRect> &rectList, const QPoint &pos);
    const QPoint rawXPosition(const QPoint &scaledPos);
    const QPoint scaledPos(const QPoint &rawXPos);
    QList<Monitor *> validMonitorList(const QMap<Monitor *, MonitorInter *> &map);

private:
    QWidget *m_parent;
    DWindowManagerHelper *m_wmHelper;
    XcbMisc *m_xcbMisc;

    // monitor screen
    XEventMonitor *m_eventInter;
    XEventMonitor *m_extralEventInter;
    // 触控屏唤起任务栏监控区域接口
    XEventMonitor *m_touchEventInter;

    // DBus interface
    DBusDock *m_dockInter;
    DisplayInter *m_displayInter;
    DBusLuncher* m_launcherInter;

    // update monitor info
    QTimer *m_monitorUpdateTimer;
    QTimer *m_delayTimer;               // sp3需求,切换屏幕显示延时2秒唤起任务栏

    // animation
    QVariantAnimation *m_showAni;
    QVariantAnimation *m_hideAni;

    // 屏幕名称信息
    DockScreen m_ds;

    // 任务栏属性
    double m_opacity;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;
    /**
     * @brief m_monitorInfo 屏幕信息(注意需要保证内容实时更新,且有效)
     */
    QMap<Monitor *, MonitorInter *> m_monitorInfo;
    /***************不和其他流程产生交互,尽量不要动这里的变量***************/
    int m_screenRawHeight;
    int m_screenRawWidth;
    QString m_registerKey;
    QString m_extralRegisterKey;
    QString m_touchRegisterKey;                 // 触控屏唤起任务栏监控区域key
    bool m_aniStart;                            // changeDockPosition是否正在运行中
    bool m_draging;                             // 鼠标是否正在调整任务栏的宽度或高度
    bool m_autoHide;                            // 和MenuWorker保持一致,为false时表示菜单已经打开
    bool m_btnPress;                            // 鼠标按下时移动到唤醒区域不应该响应唤醒
    QList<MonitRect> m_monitorRectList;         // 监听唤起任务栏区域
    QList<MonitRect> m_extralRectList;          // 任务栏外部区域,随m_monitorRectList一起更新
    QList<MonitRect> m_touchRectList;           // 监听触屏唤起任务栏区域
    QString m_delayScreen;                      // 任务栏将要切换到的屏幕名

    bool m_touchPress;                          // 触屏按下
    QPoint m_touchPos;                          // 触屏按下坐标
    /*****************************************************************/
};

#endif // MULTISCREENWORKER_H
