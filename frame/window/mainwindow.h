﻿/*
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

#include <QWidget>
#include <QTimer>
#include <QRect>

#include <DPlatformWindowHandle>
#include <DWindowManagerHelper>

class MainPanel;
class DBusDockAdaptors;
class MainWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();

    friend class MainPanel;

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

    void setFixedSize(const QSize &size);
    void internalAnimationMove(int x, int y);
    void initSNIHost();
    void initComponents();
    void initConnections();

    const QPoint x11GetWindowPos();
    void x11MoveWindow(const int x, const int y);
    void x11MoveResizeWindow(const int x, const int y, const int w, const int h);

signals:
    void panelGeometryChanged();

private slots:
    void positionChanged(const Position prevPos);
    void updatePosition();
    void updateGeometry();
    void clearStrutPartial();
    void setStrutPartial();
    void compositeChanged();
    void internalMove() { internalMove(m_posChangeAni->currentValue().toPoint()); }
    void internalMove(const QPoint &p);

    void expand();
    void narrow(const Position prevPos);
    void resetPanelEnvironment(const bool visible, const bool resetPosition = true);
    void updatePanelVisible();

    void adjustShadowMask();
    void positionCheck();

    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);

private:
    bool m_launched;
    bool m_updatePanelVisible;
    MainPanel *m_mainPanel;

    DPlatformWindowHandle m_platformWindowHandle;
    DWindowManagerHelper *m_wmHelper;

    QTimer *m_positionUpdateTimer;
    QTimer *m_expandDelayTimer;
    QTimer *m_leaveDelayTimer;
    QTimer *m_shadowMaskOptimizeTimer;
    QVariantAnimation *m_sizeChangeAni;
    QVariantAnimation *m_posChangeAni;
    QPropertyAnimation *m_panelShowAni;
    QPropertyAnimation *m_panelHideAni;

    XcbMisc *m_xcbMisc;
    DockSettings *m_settings;

    QDBusConnectionInterface *m_dbusDaemonInterface;
    org::kde::StatusNotifierWatcher *m_sniWatcher;
    QString m_sniHostService;
};

#endif // MAINWINDOW_H
