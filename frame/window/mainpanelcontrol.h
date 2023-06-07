// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MAINPANELCONTROL_H
#define MAINPANELCONTROL_H

#include "constants.h"
#include "dbusutil.h"

#include <QWidget>
#include <sys/types.h>

using namespace Dock;

class QBoxLayout;
class QLabel;
class DockTrayWindow;
class PluginsItem;
class DockItem;
class PlaceholderItem;
class AppDragWidget;
class DesktopWidget;
class RecentAppHelper;
class ToolAppHelper;
class MultiWindowHelper;

class MainPanelControl : public QWidget
{
    Q_OBJECT

public:
    explicit MainPanelControl(QWidget *parent = nullptr);

    void setPositonValue(Position position);
    void setDisplayMode(DisplayMode dislayMode);
    void resizeDockIcon();

    QSize suitableSize(const Position &position, int screenSize, double deviceRatio) const;

public slots:
    void insertItem(const int index, DockItem *item);
    void removeItem(DockItem *item);
    void itemUpdated(DockItem *item);

signals:
    void itemMoved(DockItem *sourceItem, DockItem *targetItem);
    void itemAdded(const QString &appDesktop, int idx);
    void requestUpdate();

private:
    void initUI();
    void initConnection();
    void updateAppAreaSonWidgetSize();
    void updateMainPanelLayout();
    void updateDisplayMode();
    void moveAppSonWidget();
    void updateModeChange();

    void addFixedAreaItem(int index, QWidget *wdg);
    void removeFixedAreaItem(QWidget *wdg);
    void removeAppAreaItem(QWidget *wdg);
    int getScreenSize() const;
    int trayAreaSize(qreal ratio) const;

    // 拖拽相关
    void startDrag(DockItem *);
    DockItem *dropTargetItem(DockItem *sourceItem, QPoint point);
    void moveItem(DockItem *sourceItem, DockItem *targetItem);
    void handleDragMove(QDragMoveEvent *e, bool isFilter);
    void calcuDockIconSize(int w, int h);
    bool checkNeedShowDesktop();
    bool appIsOnDock(const QString &appDesktop);
    void dockRecentApp(DockItem *dockItem);
    PluginsItem *trash() const;

private Q_SLOTS:
    void onRecentVisibleChanged(bool visible);
    void onDockAppVisibleChanged(bool visible);
    void onToolVisibleChanged(bool visible);
    void onTrayRequestUpdate();

protected:
    void dragMoveEvent(QDragMoveEvent *e) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *e) override;
    void dropEvent(QDropEvent *) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void enterEvent(QEvent *event) override;
    void mousePressEvent(QMouseEvent *e) override;
    void resizeEvent(QResizeEvent *event) override;
    void paintEvent(QPaintEvent *event) override;

private:
    QBoxLayout *m_mainPanelLayout;

    QWidget *m_fixedAreaWidget;     // 固定区域
    QBoxLayout *m_fixedAreaLayout;  // 固定区域布局
    QLabel *m_fixedSpliter;         // 固定区域与应用区域间的分割线
    QWidget *m_appAreaWidget;       // 应用区域
    QWidget *m_appAreaSonWidget;    // 子应用区域，所在位置根据显示模式手动指定
    QBoxLayout *m_appAreaSonLayout; // 子应用区域布局
    QLabel *m_appSpliter;           // 应用区域与托盘区域间的分割线
    QWidget *m_recentAreaWidget;    // 最近打开应用
    QBoxLayout *m_recentLayout;
    QLabel *m_recentSpliter;        // 最近打开应用区域分割线
    QWidget *m_toolAreaWidget;      // 工具区域，用来存放多开窗口和回收站等
    QBoxLayout *m_toolAreaLayout;   // 工具区域的布局
    QWidget *m_multiWindowWidget;   // 多开窗口区域，用来存放多开窗口
    QBoxLayout *m_multiWindowLayout;// 用来存放多开窗口的布局
    QWidget *m_toolSonAreaWidget;   // 工具区域，用来存放回收站等工具
    QBoxLayout *m_toolSonLayout;    // 工具区域布局

    Position m_position;
    QPointer<PlaceholderItem> m_placeholderItem;
    QString m_draggingMimeKey;
    AppDragWidget *m_appDragWidget;
    DisplayMode m_displayMode;
    QPoint m_mousePressPos;
    DockTrayWindow *m_tray;
    int m_dragIndex = -1;           // 记录应用区域被拖拽图标的位置

    RecentAppHelper *m_recentHelper;
    ToolAppHelper *m_toolHelper;
    MultiWindowHelper *m_multiHelper;
    bool m_showRecent;
};

#endif // MAINPANELCONTROL_H
