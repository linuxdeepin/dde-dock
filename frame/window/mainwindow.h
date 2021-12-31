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
#include "statusnotifierwatcher_interface.h"
#include "mainpanelcontrol.h"
#include "multiscreenworker.h"
#include "touchsignalmanager.h"
#include "imageutil.h"
#include "utils.h"

#include <DPlatformWindowHandle>
#include <DWindowManagerHelper>
#include <DBlurEffectWidget>
#include <DGuiApplicationHelper>

#include <QWidget>

DWIDGET_USE_NAMESPACE

using XEventMonitor = ::com::deepin::api::XEventMonitor;

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
    explicit DragWidget(QWidget *parent = nullptr)
        : QWidget(parent)
    {
        setObjectName("DragWidget");
        m_dragStatus = false;
    }

public slots:
    void onTouchMove(double scaleX, double scaleY)
    {
        Q_UNUSED(scaleX);
        Q_UNUSED(scaleY);

        static QPoint lastPos;
        QPoint curPos = QCursor::pos();
        if (lastPos == curPos) {
            return;
        }
        lastPos = curPos;
        qApp->postEvent(this, new QMouseEvent(QEvent::MouseMove, mapFromGlobal(curPos)
                                                      , QPoint(), curPos, Qt::LeftButton, Qt::LeftButton
                                                      , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
    }

signals:
    void dragPointOffset(QPoint);
    void dragFinished();

private:
    void mousePressEvent(QMouseEvent *event) override
    {
        // qt转发的触屏按下信号不进行响应
        if (event->source() == Qt::MouseEventSynthesizedByQt) {
            return;
        }
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
        if (Utils::IS_WAYLAND_DISPLAY)
            updateCursor();

        QApplication::setOverrideCursor(cursor());
    }

    void leaveEvent(QEvent *) override
    {
        QApplication::setOverrideCursor(Qt::ArrowCursor);
    }

    void updateCursor()
    {
        QString theme = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-name", "bloom").toString();
        int cursorSize = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-size", 24).toInt();
        Position position = static_cast<Dock::Position>(qApp->property("position").toInt());

        static QString lastTheme;
        static int lastPosition = -1;
        static int lastCursorSize = -1;
        if (theme != lastTheme || position != lastPosition || cursorSize != lastCursorSize) {
            lastTheme = theme;
            lastPosition = position;
            lastCursorSize = cursorSize;
            const char* cursorName = (position == Bottom || position == Top) ? "v_double_arrow" : "h_double_arrow";
            QCursor *newCursor = ImageUtil::loadQCursorFromX11Cursor(theme.toStdString().c_str(), cursorName, cursorSize);
            if (!newCursor)
                return;

            setCursor(*newCursor);
            static QCursor *lastCursor = nullptr;
            if (lastCursor)
                delete lastCursor;

            lastCursor = newCursor;
        }
    }
};

class MainWindow : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = nullptr);

    void setEffectEnabled(const bool enabled);
    void setComposite(const bool hasComposite);
    void setGeometry(const QRect &rect);
    void sendNotifications();

    friend class MainPanelControl;

    MainPanelControl *panel() {return m_mainPanel;}

public slots:
    void launch();
    void callShow();
    void relaodPlugins();

private:
    using QWidget::show;
    void mousePressEvent(QMouseEvent *e) override;
    void keyPressEvent(QKeyEvent *e) override;
    void enterEvent(QEvent *e) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void moveEvent(QMoveEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

    void initMember();
    void initSNIHost();
    void initComponents();
    void initConnections();

    void resizeDockIcon();

signals:
    void panelGeometryChanged();

public slots:
    void RegisterDdeSession();
    void resizeDock(int offset, bool dragging);
    void resetDragWindow(); // 任务栏调整高度或宽度后需调用此函数

private slots:
    void compositeChanged();
    void adjustShadowMask();

    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);
    void onMainWindowSizeChanged(QPoint offset);
    void themeTypeChanged(DGuiApplicationHelper::ColorType themeType);
    void touchRequestResizeDock();

private:
    MainPanelControl *m_mainPanel;                      // 任务栏
    DPlatformWindowHandle m_platformWindowHandle;
    DWindowManagerHelper *m_wmHelper;
    MultiScreenWorker *m_multiScreenWorker;             // 多屏幕管理
    MenuWorker *m_menuWorker;
    QTimer *m_shadowMaskOptimizeTimer;
    QDBusConnectionInterface *m_dbusDaemonInterface;
    org::kde::StatusNotifierWatcher *m_sniWatcher;      // DBUS状态通知
    DragWidget *m_dragWidget;

    QString m_sniHostService;

    bool m_launched;
    QString m_registerKey;
    QStringList m_registerKeys;

    QTimer *m_updateDragAreaTimer;
};

#endif // MAINWINDOW_H
