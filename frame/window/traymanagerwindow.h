// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRAYMANAGERWINDOW_H
#define TRAYMANAGERWINDOW_H

#include "constants.h"
#include "dbusutil.h"

#include "org_deepin_dde_timedate1.h"

#include <QPainterPath>
#include <QWidget>

namespace Dtk { namespace Widget { class DBlurEffectWidget; } }

using namespace Dtk::Widget;

using Timedate = org::deepin::dde::Timedate1;

class QuickPluginWindow;
class QBoxLayout;
class TrayGridView;
class TrayModel;
class TrayDelegate;
class SystemPluginWindow;
class QLabel;
class QDropEvent;
class DateTimeDisplayer;
class QPainterPath;

class TrayManagerWindow : public QWidget
{
    Q_OBJECT

public:
    explicit TrayManagerWindow(QWidget *parent = nullptr);
    ~TrayManagerWindow() override;

    void updateBorderRadius(int borderRadius);
    void updateLayout();
    void setPositon(Dock::Position position);
    void setDisplayMode(Dock::DisplayMode displayMode);
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

Q_SIGNALS:
    void requestUpdate();

protected:
    void resizeEvent(QResizeEvent *event) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragMoveEvent(QDragMoveEvent *e) override;
    void dropEvent(QDropEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *event) override;
    void paintEvent(QPaintEvent *event) override;

private:
    void initUi();
    void initConnection();

    void resetChildWidgetSize();
    void resetMultiDirection();
    void resetSingleDirection();
    QColor maskColor(uint8_t alpha) const;

    int appDatetimeSize(const Dock::Position &position) const;
    QPainterPath roundedPaths();
    void updateItemLayout(int dockSize);
    int pathRadius() const;

private Q_SLOTS:
    void onTrayCountChanged();
    void updateHighlightArea(const QRect &rect);

private:
    QWidget *m_appPluginDatetimeWidget;
    SystemPluginWindow *m_systemPluginWidget;
    QWidget *m_appPluginWidget;
    QuickPluginWindow *m_quickIconWidget;
    DateTimeDisplayer *m_dateTimeWidget;
    QBoxLayout *m_appPluginLayout;
    QBoxLayout *m_mainLayout;
    TrayGridView *m_trayView;
    TrayModel *m_model;
    TrayDelegate *m_delegate;
    Dock::Position m_position;
    Dock::DisplayMode m_displayMode;
    QLabel *m_splitLine;
    bool m_singleShow;                              // 用于记录当前日期时间和插件区域是显示一行还是显示多行
    int m_borderRadius;                             // 圆角的值
    uint m_windowFashionSize;
    QPainterPath m_highlightArea;
};

#endif // PLUGINWINDOW_H
