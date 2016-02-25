/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "mainwidget.h"
#include "xcb_misc.h"
#include "controller/stylemanager.h"

#include <QApplication>

const int ENTER_DELAY_INTERVAL = 600;
MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    this->setWindowFlags(Qt::FramelessWindowHint | Qt::WindowDoesNotAcceptFocus);
    this->setAttribute(Qt::WA_TranslucentBackground);
    //the attribute "Qt::WA_X11DoNotAcceptFocus" will not tack effect, not know the reason
    //this->setAttribute(Qt::WA_X11DoNotAcceptFocus);

    initHideStateManager();

#ifdef NEW_DOCK_LAYOUT
    m_mainPanel = new DockPanel(this);
    connect(m_mainPanel, &DockPanel::startShow, this, &MainWidget::showDock);
    connect(m_mainPanel, &DockPanel::panelHasHidden, this, &MainWidget::hideDock);
    connect(m_mainPanel, &DockPanel::sizeChanged, this, &MainWidget::onPanelSizeChanged);
#else
    m_mainPanel = new Panel(this);
    connect(m_mainPanel, &Panel::startShow, this, &MainWidget::showDock);
    connect(m_mainPanel, &Panel::panelHasHidden, this, &MainWidget::hideDock);
    connect(m_mainPanel, &Panel::sizeChanged, this, &MainWidget::onPanelSizeChanged);
#endif

    connect(m_dmd, &DockModeData::dockModeChanged, this, &MainWidget::onDockModeChanged);

    //For init
    m_display = new DBusDisplay(this);
    updatePosition(m_display->primaryRect());

    DockUIDbus *dockUIDbus = new DockUIDbus(this);
    Q_UNUSED(dockUIDbus)

    XcbMisc::instance()->set_window_type(winId(), XcbMisc::Dock);

    connect(m_display, &DBusDisplay::PrimaryChanged, this, &MainWidget::updateGeometry);
    connect(m_display, &DBusDisplay::ScreenHeightChanged, this, &MainWidget::updateGeometry);
    connect(m_display, &DBusDisplay::ScreenWidthChanged, this, &MainWidget::updateGeometry);
}

void MainWidget::onDockModeChanged()
{
    updatePosition(m_display->primaryRect());
//    QMetaObject::invokeMethod(this, "updatePosition", Qt::QueuedConnection, Q_ARG(QRect, m_display->primaryRect()));
}

void MainWidget::updatePosition(const QRect &rec)
{
    // sometimes rec's width or height is ZERO, we need to ignore these wrong data.
    if (!rec.width() || !rec.height()) {
        return;
    }

//    qDebug() << "move to " << rec;

    if (m_hasHidden) {
        //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
#ifdef NEW_DOCK_LAYOUT
        this->setFixedSize(m_mainPanel->sizeHint().width(), 1);
#else
        this->setFixedSize(m_mainPanel->width(), 1);
#endif
        this->move(rec.x() + (rec.width() - width()) / 2,
                   rec.y() + rec.height() - 1);//1 pixel for grab mouse enter event to show panel
    } else {
#ifdef NEW_DOCK_LAYOUT
        this->setFixedSize(m_mainPanel->sizeHint().width(), m_dmd->getDockHeight());
#else
        this->setFixedSize(m_mainPanel->width(), m_dmd->getDockHeight());
#endif

        move(rec.x() + (rec.width() - width()) / 2,
             rec.y() + rec.height() - height() /*- 10*/);
    }

//    qDebug() << "size = " << width() << ", " << height();
//    qDebug() << "move to " << this->x() << ", " << this->y() << " end";

    updateXcbStructPartial();
}

void MainWidget::updateXcbStructPartial()
{
    int tmpHeight = 0;
    DBusDockSetting dds;
    if (dds.GetHideMode() == Dock::KeepShowing) {
        // qApp's screenHeight is wrong. its a bug, use dbus data instead.
//        int maxMonitorHeight = qApp->desktop()->size().height();
//        int max = 0;
//        for (QScreen *screen : qApp->screens())
//        {
//            QRect screenRect = screen->geometry();
//            max = qMax(max, screenRect.y() + screenRect.height());
//        }

//        qDebug() << "max = " << max;

        int maxMonitorHeight = m_display->screenHeight();
        tmpHeight = maxMonitorHeight - y();
    }

    // sometimes screen height is wrong, we need to ignore wrong data.
    if (tmpHeight && tmpHeight < m_dmd->getDockHeight()) {
        return;
    }

//    qDebug() << "maxHeight dbus: " << m_display->screenHeight();
//    qDebug() << "set structPartial: " << x() << ", " << width() << ", " << tmpHeight;

    XcbMisc::instance()->set_strut_partial(winId(),
                                           XcbMisc::OrientationBottom,
                                           tmpHeight,
                                           x(),
                                           x() + width() - 1);
// The line below causes deepin-wm to regard dde-dock as a normal window
// while previewing windows. https://github.com/fasheng/arch-deepin/issues/249
//    this->setVisible(true);
}

void MainWidget::updateGeometry()
{
    QRect primaryRect = m_display->primaryRect();

    qDebug() << "change screen, primary is: " << m_display->primary();

    for (const QScreen *screen : qApp->screens()) {
        if (screen->name() == m_display->primary()) {
            primaryRect = screen->geometry();
            connect(screen, &QScreen::geometryChanged, this, &MainWidget::updatePosition);
        } else {
            disconnect(screen, &QScreen::geometryChanged, this, &MainWidget::updatePosition);
        }
    }

    updatePosition(primaryRect);
}

void MainWidget::initHideStateManager()
{
    m_dhsm = new DBusHideStateManager(this);
    m_dhsm->SetState(Dock::HideStateHiding);
}

void MainWidget::enterEvent(QEvent *)
{
    if (height() == 1) {
        QTimer *st = new QTimer(this);
        connect(st, &QTimer::timeout, this, [ = ] {
            //make sure the panel will show by mouse-enter
            if (geometry().contains(QCursor::pos()))
            {
                qDebug() << "MouseEntered, show dock...";
                showDock();
                m_mainPanel->startShow();
            }
            sender()->deleteLater();
        });
        st->start(ENTER_DELAY_INTERVAL);
    }
}

void MainWidget::leaveEvent(QEvent *)
{
    if (!this->geometry().contains(QCursor::pos())) {
        m_dhsm->UpdateState();
    }
}

void MainWidget::showDock()
{
    m_hasHidden = false;
    updatePosition(m_display->primaryRect());
}

void MainWidget::hideDock()
{
    m_hasHidden = true;
    updatePosition(m_display->primaryRect());
}

void MainWidget::onPanelSizeChanged()
{
    updatePosition(m_display->primaryRect());
}

MainWidget::~MainWidget()
{
    qDebug() << "dde-dock destroyed";
}

void MainWidget::loadResources()
{
    m_mainPanel->loadResources();
}

DockUIDbus::DockUIDbus(MainWidget *parent):
    QDBusAbstractAdaptor(parent),
    m_parent(parent)
{
    QDBusConnection::sessionBus().registerObject(DBUS_PATH, parent);
}

DockUIDbus::~DockUIDbus()
{

}

qulonglong DockUIDbus::Xid()
{
    return m_parent->winId();
}

QString DockUIDbus::currentStyleName()
{
    return StyleManager::instance()->currentStyle();
}

QStringList DockUIDbus::styleNameList()
{
    return StyleManager::instance()->styleNameList();
}

void DockUIDbus::applyStyle(const QString &styleName)
{
    StyleManager::instance()->applyStyle(styleName);
}
