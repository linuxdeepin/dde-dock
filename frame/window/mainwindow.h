// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "xcb_misc.h"
#include "multiscreenworker.h"
#include "touchsignalmanager.h"
#include "imageutil.h"
#include "utils.h"
#include "mainwindowbase.h"

#include <DPlatformWindowHandle>

#include <QWidget>

DWIDGET_USE_NAMESPACE

class MainPanelControl;
class QTimer;
class MenuWorker;
class QScreen;

class MainWindow : public MainWindowBase
{
    Q_OBJECT

public:
    explicit MainWindow(MultiScreenWorker *multiScreenWorker, QWidget *parent = nullptr);
    void setGeometry(const QRect &rect);

    friend class MainPanelControl;

    // 以下接口是实现基类的接口
    // 用来更新子区域的位置，一般用于在执行动画的过程中，根据当前的位置来更新里面panel的大小
    DockWindowType windowType() const override;
    void setPosition(const Dock::Position &position) override;
    void setDisplayMode(const Dock::DisplayMode &displayMode) override;
    void updateParentGeometry(const Dock::Position &pos, const QRect &rect) override;
    QSize suitableSize(const Dock::Position &pos, const int &screenSize, const double &deviceRatio) const override;
    void resetPanelGeometry() override;
    void serviceRestart() override;
    void animationFinished(bool showOrHide) override;

private:
    using QWidget::show;
    void initConnections();
    void resizeDockIcon();

private:
    MainPanelControl *m_mainPanel;                      // 任务栏
    MultiScreenWorker *m_multiScreenWorker;             // 多屏幕管理

    QString m_sniHostService;

    QString m_registerKey;
    QStringList m_registerKeys;
    bool m_needUpdateUi;
};

#endif // MAINWINDOW_H
