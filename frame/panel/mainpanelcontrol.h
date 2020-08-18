/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     xuwenw <xuwenw@xuwenw.so>
 *
 * Maintainer:  <@xuwenw.so>
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

#ifndef MAINPANELCONTROL_H
#define MAINPANELCONTROL_H

#include "constants.h"

#include <QWidget>
#include <QBoxLayout>
#include <QLabel>

#include <com_deepin_daemon_gesture.h>

using namespace Dock;
using Gesture = com::deepin::daemon::Gesture;

class TrayPluginItem;
class PluginsItem;

class MainPanelDelegate
{
public:
    virtual bool appIsOnDock(const QString &appDesktop) = 0;
};

class DesktopWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DesktopWidget(QWidget *parent) : QWidget(parent){
    }
};

class DockItem;
class PlaceholderItem;
class AppDragWidget;
class MainPanelControl : public QWidget
{
    Q_OBJECT
public:
    explicit MainPanelControl(QWidget *parent = nullptr);
    ~MainPanelControl() override;

    void addFixedAreaItem(int index, QWidget *wdg);
    void addAppAreaItem(int index, QWidget *wdg);
    void addTrayAreaItem(int index, QWidget *wdg);
    void addPluginAreaItem(int index, QWidget *wdg);
    void removeFixedAreaItem(QWidget *wdg);
    void removeAppAreaItem(QWidget *wdg);
    void removeTrayAreaItem(QWidget *wdg);
    void removePluginAreaItem(QWidget *wdg);
    void setPositonValue(Position position);
    void setDisplayMode(DisplayMode m_displayMode);
    void getTrayVisableItemCount();

    MainPanelDelegate *delegate() const;
    void setDelegate(MainPanelDelegate *delegate);

signals:
    void itemMoved(DockItem *sourceItem, DockItem *targetItem);
    void itemAdded(const QString &appDesktop, int idx);

private:
    void resizeEvent(QResizeEvent *event) override;

    void init();
    void updateAppAreaSonWidgetSize();
    void updateMainPanelLayout();
    void updateDisplayMode();
    void moveAppSonWidget();

    void dragMoveEvent(QDragMoveEvent *e) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *e) override;
    void dropEvent(QDropEvent *) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void mousePressEvent(QMouseEvent *e) override;

    void startDrag(DockItem *);
    DockItem *dropTargetItem(DockItem *sourceItem, QPoint point);
    void moveItem(DockItem *sourceItem, DockItem *targetItem);
    void handleDragMove(QDragMoveEvent *e, bool isFilter);
    void paintEvent(QPaintEvent *event) override;
    void resizeDockIcon();
    void calcuDockIconSize(int w, int h, PluginsItem *trashPlugin, PluginsItem *shutdownPlugin, PluginsItem *keyboardPlugin, PluginsItem *notificationPlugin);
    void resizeDesktopWidget();
    bool checkNeedShowDesktop();

public slots:
    void insertItem(const int index, DockItem *item);
    void removeItem(DockItem *item);
    void itemUpdated(DockItem *item);

    // void
    void onGSettingsChanged(const QString &key);
    
protected:
    void showEvent(QShowEvent *event) override;
private:
    QBoxLayout *m_mainPanelLayout;
    QWidget *m_fixedAreaWidget;
    QWidget *m_appAreaWidget;
    QWidget *m_trayAreaWidget;
    QWidget *m_pluginAreaWidget;
    DesktopWidget *m_desktopWidget;
    QBoxLayout *m_fixedAreaLayout;
    QBoxLayout *m_trayAreaLayout;
    QBoxLayout *m_pluginLayout;
    QWidget *m_appAreaSonWidget;
    QBoxLayout *m_appAreaSonLayout;
    //    QBoxLayout *m_appAreaLayout;
    Position m_position;
    QPointer<PlaceholderItem> m_placeholderItem;
    MainPanelDelegate *m_delegate;
    QString m_draggingMimeKey;
    AppDragWidget *m_appDragWidget;
    DisplayMode m_dislayMode;
    QLabel *m_fixedSpliter;
    QLabel *m_appSpliter;
    QLabel *m_traySpliter;
    QPoint m_mousePressPos;
    int m_trayIconCount;
    TrayPluginItem *m_tray = nullptr;
    bool m_isHover;//判断鼠标是否移到desktop区域
    bool m_needRecoveryWin; // 判断鼠标移出desktop区域是否恢复之前窗口
    bool m_isEnableLaunch;//判断是否使能了com.deepin.dde.dock.module.launcher
};

#endif // MAINPANELCONTROL_H
