/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             zhaolong <zhaolong@uniontech.com>
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
    void resizeEvent(QResizeEvent *event) override;
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
