// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKTRAYWINDOW_H
#define DOCKTRAYWINDOW_H

#include "constants.h"
#include "dbusutil.h"

#include <QWidget>

class QBoxLayout;
class SystemPluginWindow;
class DateTimeDisplayer;
class QuickPluginWindow;
class TrayGridView;
class TrayModel;
class TrayDelegate;
class PluginsItem;
class PluginsItemInterface;
class QLabel;

class DockTrayWindow : public QWidget
{
    Q_OBJECT

public:
    explicit DockTrayWindow(QWidget *parent = nullptr);

    void setPositon(const Dock::Position &position);
    void setDisplayMode(const Dock::DisplayMode &displayMode);

    QSize suitableSize(const Dock::Position &position, const int &, const double &) const;
    QSize suitableSize() const;

    void layoutWidget();

Q_SIGNALS:
    void requestUpdate();

protected:
    void resizeEvent(QResizeEvent *event) override;
    void paintEvent(QPaintEvent *event) override;
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    void initUi();
    void initConnection();
    void initAttribute();
    void updateLayout(const Dock::Position &position);
    void resizeTool() const;
    bool pluginExists(PluginsItemInterface *itemInter) const;
    void moveToolPlugin();
    void updateToolWidget();

private Q_SLOTS:
    void onUpdateComponentSize();
    void onItemAdded(PluginsItemInterface *itemInter);
    void onItemRemove(PluginsItemInterface *itemInter);
    void onDropIcon(QDropEvent *dropEvent);

private:
    Dock::Position m_position;
    Dock::DisplayMode m_displayMode;
    QBoxLayout *m_mainBoxLayout;
    QWidget *m_showDesktopWidget;
    QWidget *m_toolWidget;
    QBoxLayout *m_toolLayout;
    QLabel *m_toolLineLabel;
    DateTimeDisplayer *m_dateTimeWidget;            // 日期时间
    SystemPluginWindow *m_systemPuginWidget;        // 固定区域-一般是右侧的电源按钮
    QuickPluginWindow *m_quickIconWidget;           // 插件区域-包括网络、蓝牙等
    TrayGridView *m_trayView;                       // 托盘区域视图
    TrayModel *m_model;                             // 托盘区域的model
    TrayDelegate *m_delegate;                       // 托盘区域的视图代理
    QWidget *m_toolFrontSpaceWidget;                // 用于显示桌面和回收站中间的间隔
    QWidget *m_toolBackSpaceWidget;                 // 用于回收站和时间日期分割线中间的间隔
    QWidget *m_dateTimeSpaceWidget;                 // 用于时间日期分割线和时间日期中间的间隔
};

#endif // DOCKTRAYWINDOW_H
