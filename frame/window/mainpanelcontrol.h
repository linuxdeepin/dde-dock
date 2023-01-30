// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MAINPANELCONTROL_H
#define MAINPANELCONTROL_H

#include "constants.h"

#include <DWidget>

#include <com_deepin_daemon_gesture.h>
using namespace Dock;

DWIDGET_BEGIN_NAMESPACE
class DIconButton;
DWIDGET_END_NAMESPACE

DWIDGET_USE_NAMESPACE

class QScrollArea;
class QBoxLayout;
class QLabel;
class TrayPluginItem;
class PluginsItem;
class DockItem;
class PlaceholderItem;
class AppDragWidget;
class DesktopWidget;
class OverflowItem;

class MainPanelControl : public QWidget
{
    Q_OBJECT
public:
    explicit MainPanelControl(QWidget *parent = nullptr);

    void setPositonValue(Position position);
    void setDisplayMode(DisplayMode dislayMode);
    void resizeDockIcon();
    void updatePluginsLayout();
    void setToggleDesktopInterval(int ms);

public slots:
    void insertItem(const int index, DockItem *item);
    void removeItem(DockItem *item);
    void itemUpdated(DockItem *item);
    void setKwinAppItemMinimizedGeometry(DockItem *item, const QRect);

signals:
    void itemMoved(DockItem *sourceItem, DockItem *targetItem);
    void itemAdded(const QString &appDesktop, int idx);
    void updateLayout();

private:
    void initUI();
    void updateAppAreaSonWidgetSize();
    void updateMainPanelLayout();
    void updateDisplayMode();
    void moveAppSonWidget();

    void addFixedAreaItem(int index, QWidget *wdg);
    void removeFixedAreaItem(QWidget *wdg);
    void addAppAreaItem(int index, QWidget *wdg);
    void removeAppAreaItem(QWidget *wdg);
    void addTrayAreaItem(int index, QWidget *wdg);
    void removeTrayAreaItem(QWidget *wdg);
    void addPluginAreaItem(int index, QWidget *wdg);
    void removePluginAreaItem(QWidget *wdg);

    // 拖拽相关
    void startDrag(DockItem *);
    DockItem *dropTargetItem(DockItem *sourceItem, QPoint point);
    void moveItem(DockItem *sourceItem, DockItem *targetItem);
    void handleDragMove(QDragMoveEvent *e, bool isFilter);
    void calcuDockIconSize(int appItemSize, int maxcount ,int showtype, int traySize);
    void resizeDesktopWidget();
    bool checkNeedShowDesktop();
    bool appIsOnDock(const QString &appDesktop);

    int getItemIndex(DockItem *targetItem) const;

    // about resize
    void resizeLayout();

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
    // state
    int m_maxcount = -1;
    int m_showtype = -1;
    int m_appItemSize = 0;
    QSize m_trayareaSize = QSize(0, 0);
    QSize m_pluginareaSize = QSize(0, 0);
 
    // widgets
    QBoxLayout *m_mainPanelLayout;

    QWidget *m_fixedAreaWidget;      // 固定区域
    QBoxLayout *m_fixedAreaLayout;   //
    QLabel *m_fixedSpliter;          // 固定区域与应用区域间的分割线
    QWidget *m_appAreaWidget;        // 应用区域
                                     //
    QWidget *m_appAreaSonWidget;     // 子应用区域，所在位置根据显示模式手动指定
    QBoxLayout *m_appAreaSonLayout;  //
    QLabel *m_appSpliter;            // 应用区域与托盘区域间的分割线

    QWidget *m_appOverflowWidget;    // 溢出区域
    QBoxLayout *m_appOverflowLayout;
    OverflowItem *m_overflowBtn;     // 溢出按钮
    QBoxLayout *m_overflowAreaLayout;// 
    QLabel *m_overflowSpliter;       // 溢出区域和应用区域分割线

    QWidget *m_trayAreaWidget;       // 托盘区域
    QBoxLayout *m_trayAreaLayout;    //
    QLabel *m_traySpliter;           // 托盘区域与插件区域间的分割线
    QWidget *m_pluginAreaWidget;     // 插件区域
    QBoxLayout *m_pluginLayout;      //
    DesktopWidget *m_desktopWidget;  // 桌面预览区域

    Position m_position;
    QPointer<PlaceholderItem> m_placeholderItem;
    QString m_draggingMimeKey;
    AppDragWidget *m_appDragWidget;
    DisplayMode m_dislayMode;
    QPoint m_mousePressPos;
    TrayPluginItem *m_tray;
    int m_dragIndex = -1;  // 记录应用区域被拖拽图标的位置

    PluginsItem *m_trashItem;  // 垃圾箱插件（需要特殊处理一下）
};

#endif  // MAINPANELCONTROL_H
