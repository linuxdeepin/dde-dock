/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "xcb/xcb_misc.h"
#include "dbus/dbusdisplay.h"
#include "dbus/dbusdockadaptors.h"
#include "dbus/sni/statusnotifierwatcher_interface.h"
#include "util/docksettings.h"
#include "panel/mainpanelcontrol.h"

#include <QWidget>
#include <QTimer>
#include <QRect>

#include <DPlatformWindowHandle>
#include <DWindowManagerHelper>
#include <DBlurEffectWidget>
#include <DGuiApplicationHelper>
#include <DRegionMonitor>

DWIDGET_USE_NAMESPACE


class DragWidget;
class MainPanel;
class MainPanelControl;
class DBusDockAdaptors;
class MainWindow : public DBlurEffectWidget, public MainPanelDelegate
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
    void setEffectEnabled(const bool enabled);
    void setComposite(const bool hasComposite);

    friend class MainPanel;
    friend class MainPanelControl;

public slots:
    void launch();

private:
    using QWidget::show;
    bool event(QEvent *e);
    void showEvent(QShowEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void keyPressEvent(QKeyEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void mouseMoveEvent(QMouseEvent *e);

    void initSNIHost();
    void initComponents();
    void initConnections();
    void resizeMainWindow();
    void resizeMainPanelWindow();

    const QPoint x11GetWindowPos();
    void x11MoveWindow(const int x, const int y);
    void x11MoveResizeWindow(const int x, const int y, const int w, const int h);
    bool appIsOnDock(const QString &appDesktop);
    void onRegionMonitorChanged();
    void updateRegionMonitorWatch();
    void getTrayVisableItemCount();

signals:
    void panelGeometryChanged();
    void loaderPlugins();

private slots:
    void positionChanged(const Position prevPos, const Position nextPos);
    void updatePosition();
    void updateGeometry();
    void clearStrutPartial();
    void setStrutPartial();
    void compositeChanged();
    void internalMove(const QPoint &p);
    void updateDisplayMode();

    void expand();
    void narrow(const Position prevPos);
    void resetPanelEnvironment(const bool visible, const bool resetPosition = true);
    void updatePanelVisible();

    void adjustShadowMask();
    void positionCheck();

    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);
    void onMainWindowSizeChanged(QPoint offset);
    void onDragFinished();
    void themeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    bool m_launched;
    MainPanelControl *m_mainPanel;

    DPlatformWindowHandle m_platformWindowHandle;
    DWindowManagerHelper *m_wmHelper;
    DRegionMonitor *m_regionMonitor;

    QTimer *m_positionUpdateTimer;
    QTimer *m_expandDelayTimer;
    QTimer *m_leaveDelayTimer;
    QTimer *m_shadowMaskOptimizeTimer;
    QVariantAnimation *m_panelShowAni;
    QVariantAnimation *m_panelHideAni;

    XcbMisc *m_xcbMisc;
    DockSettings *m_settings;

    QDBusConnectionInterface *m_dbusDaemonInterface;
    org::kde::StatusNotifierWatcher *m_sniWatcher;
    QString m_sniHostService;
    QSize m_size;
    DragWidget *m_dragWidget;
    Position m_curDockPos;
    Position m_newDockPos;
};

#endif // MAINWINDOW_H
