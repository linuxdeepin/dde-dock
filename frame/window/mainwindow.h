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

#include "xcb/xcb_misc.h"
#include "dbus/sni/statusnotifierwatcher_interface.h"
#include "panel/mainpanelcontrol.h"
#include "util/multiscreenworker.h"

#include <com_deepin_api_xeventmonitor.h>

#include <DPlatformWindowHandle>
#include <DWindowManagerHelper>
#include <DBlurEffectWidget>
#include <DGuiApplicationHelper>

#include <QWidget>

DWIDGET_USE_NAMESPACE

using XEventMonitor = ::com::deepin::api::XEventMonitor;

class MainPanel;
class MainPanelControl;
class QTimer;
class MenuWorker;
class DragWidget : public QWidget
{
    Q_OBJECT

private:
    bool m_dragStatus;
    QPoint m_resizePoint;

public:
    DragWidget(QWidget *parent = nullptr) : QWidget(parent)
    {
        setObjectName("DragWidget");
        m_dragStatus = false;
    }

signals:
    void dragPointOffset(QPoint);
    void dragFinished();

private:
    void mousePressEvent(QMouseEvent *event) override
    {
        if (event->button() == Qt::LeftButton) {
            m_resizePoint = event->globalPos();
            m_dragStatus = true;
            this->grabMouse();
        }
    }

    void mouseMoveEvent(QMouseEvent *) override
    {
        if (m_dragStatus) {
            QPoint offset = QPoint(QCursor::pos() - m_resizePoint);
            emit dragPointOffset(offset);
        }
    }

    void mouseReleaseEvent(QMouseEvent *) override
    {
        if (!m_dragStatus)
            return;

        m_dragStatus =  false;
        releaseMouse();
        emit dragFinished();
    }

    void enterEvent(QEvent *) override
    {
        QApplication::setOverrideCursor(cursor());
    }

    void leaveEvent(QEvent *) override
    {
        QApplication::setOverrideCursor(Qt::ArrowCursor);
    }
};

class MainWindow : public DBlurEffectWidget, public MainPanelDelegate
{
    Q_OBJECT

    enum Flag{
        Motion = 1 << 0,
        Button = 1 << 1,
        Key    = 1 << 2
    };

public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow() override;
    void setEffectEnabled(const bool enabled);
    void setComposite(const bool hasComposite);

    friend class MainPanel;
    friend class MainPanelControl;

    MainPanelControl *panel() {return m_mainPanel;}

public slots:
    void launch();

private:
    using QWidget::show;
    void showEvent(QShowEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    void keyPressEvent(QKeyEvent *e) override;
    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void moveEvent(QMoveEvent *event) override;

    void initSNIHost();
    void initComponents();
    void initConnections();

    bool appIsOnDock(const QString &appDesktop) override;
    void getTrayVisableItemCount();

signals:
    void panelGeometryChanged();

public slots:
    void resetDragWindow();

private slots:
    void compositeChanged();
    void adjustShadowMask();

    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);
    void onMainWindowSizeChanged(QPoint offset);
    void onDragFinished();
    void themeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    MainPanelControl *m_mainPanel;

    DPlatformWindowHandle m_platformWindowHandle;
    DWindowManagerHelper *m_wmHelper;
    MultiScreenWorker *m_multiScreenWorker;
    MenuWorker *m_menuWorker;
    XEventMonitor *m_eventInter;

    QTimer *m_shadowMaskOptimizeTimer;

    QDBusConnectionInterface *m_dbusDaemonInterface;
    org::kde::StatusNotifierWatcher *m_sniWatcher;
    QString m_sniHostService;
    DragWidget *m_dragWidget;

    bool m_launched;
    int m_dockSize;
    QString m_registerKey;
    QStringList m_registerKeys;
};

#endif // MAINWINDOW_H
