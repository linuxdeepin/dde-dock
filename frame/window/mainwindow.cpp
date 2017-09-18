/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "mainwindow.h"
#include "panel/mainpanel.h"

#include <QDebug>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <DPlatformWindowHandle>

#include <X11/X.h>
#include <X11/Xutil.h>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),

      m_updatePanelVisible(true),

      m_mainPanel(new MainPanel(this)),

      m_platformWindowHandle(this),
      m_wmHelper(DWindowManagerHelper::instance()),

      m_positionUpdateTimer(new QTimer(this)),
      m_expandDelayTimer(new QTimer(this)),
      m_sizeChangeAni(new QVariantAnimation(this)),
      m_posChangeAni(new QPropertyAnimation(this, "pos")),
      m_panelShowAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_panelHideAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_xcbMisc(XcbMisc::instance())

{
    setAccessibleName("dock-mainwindow");
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowDoesNotAcceptFocus);
    setAttribute(Qt::WA_TranslucentBackground);
    setAcceptDrops(true);

    m_platformWindowHandle.setEnableBlurWindow(false);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setWindowRadius(0);
    m_platformWindowHandle.setBorderWidth(0);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 0));

    m_settings = new DockSettings(this);
    m_xcbMisc->set_window_type(winId(), XcbMisc::Dock);

    initComponents();
    initConnections();

    m_mainPanel->setFixedSize(m_settings->windowSize());

    updatePanelVisible();
    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &MainWindow::compositeChanged, Qt::QueuedConnection);
    connect(m_mainPanel, &MainPanel::geometryChanged, this, &MainWindow::panelGeometryChanged);
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, this, &MainWindow::adjustShadowMask);
    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::adjustShadowMask);
    connect(m_panelShowAni, &QPropertyAnimation::finished, this, &MainWindow::adjustShadowMask);
//    connect(m_posChangeAni, &QPropertyAnimation::valueChanged,
//            this, [=](const QVariant &v) { const QPoint p = v.toPoint(); x11MoveWindow(p.x(), p.y()); });
}

MainWindow::~MainWindow()
{
    delete m_xcbMisc;
}

QRect MainWindow::panelGeometry()
{
    QRect rect = m_mainPanel->geometry();
    rect.moveTopLeft(m_mainPanel->mapToGlobal(QPoint(0,0)));
    return rect;
}

void MainWindow::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    adjustShadowMask();
}

void MainWindow::mousePressEvent(QMouseEvent *e)
{
    e->ignore();

    if (e->button() == Qt::RightButton)
        m_settings->showDockSettingsMenu();
}

void MainWindow::keyPressEvent(QKeyEvent *e)
{
    switch (e->key())
    {
#ifdef QT_DEBUG
    case Qt::Key_Escape:        qApp->quit();       break;
#endif
    default:;
    }
}

void MainWindow::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    if (m_settings->hideState() != Show)
        m_expandDelayTimer->start();
}

void MainWindow::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_expandDelayTimer->stop();
    updatePanelVisible();
}

void MainWindow::dragEnterEvent(QDragEnterEvent *e)
{
    QWidget::dragEnterEvent(e);

    if (m_settings->hideState() != Show)
        m_expandDelayTimer->start();
}

void MainWindow::setFixedSize(const QSize &size)
{
    const QPropertyAnimation::State state = m_sizeChangeAni->state();

//    qDebug() << state;
    if (state == QPropertyAnimation::Stopped && this->size() == size)
        return;

    if (state == QPropertyAnimation::Running)
        return m_sizeChangeAni->setEndValue(size);


//    qDebug() << Q_FUNC_INFO << size;

    m_sizeChangeAni->setStartValue(this->size());
    m_sizeChangeAni->setEndValue(size);
    m_sizeChangeAni->start();
}

void MainWindow::internalMove(int x, int y)
{
    const QPropertyAnimation::State state = m_posChangeAni->state();
    const QPoint p = m_posChangeAni->endValue().toPoint();
    const QPoint tp = QPoint(x, y);

    if (state == QPropertyAnimation::Stopped && p == tp)
        return;

//    if (state == QPropertyAnimation::Running && m_posChangeAni->endValue() != tp)
//        return m_posChangeAni->setEndValue(QPoint(x, y));

    m_posChangeAni->setStartValue(pos());
    m_posChangeAni->setEndValue(tp);
    m_posChangeAni->start();
}

void MainWindow::initComponents()
{
    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(20);
    m_positionUpdateTimer->start();

    m_expandDelayTimer->setSingleShot(true);
    m_expandDelayTimer->setInterval(m_settings->expandTimeout());

    m_sizeChangeAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_posChangeAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_panelShowAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_panelHideAni->setEasingCurve(QEasingCurve::InOutCubic);

    QTimer::singleShot(1, this, &MainWindow::compositeChanged);
}

void MainWindow::compositeChanged()
{
//    qDebug() << Q_FUNC_INFO;

    const int duration = m_wmHelper->hasComposite() ? 200 : 0;

    m_sizeChangeAni->setDuration(duration);
    m_posChangeAni->setDuration(duration);
    m_panelShowAni->setDuration(duration);
    m_panelHideAni->setDuration(duration);

    m_positionUpdateTimer->start();
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::dataChanged, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::positionChanged, this, &MainWindow::positionChanged);
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry, Qt::DirectConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, [this] { resetPanelEnvironment(true); });
    connect(m_settings, &DockSettings::windowVisibleChanged, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::autoHideChanged, this, &MainWindow::updatePanelVisible);

    connect(m_mainPanel, &MainPanel::requestRefershWindowVisible, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_mainPanel, &MainPanel::requestWindowAutoHide, m_settings, &DockSettings::setAutoHide);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition, Qt::QueuedConnection);
    connect(m_expandDelayTimer, &QTimer::timeout, this, &MainWindow::expand, Qt::QueuedConnection);

    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::updateGeometry, Qt::QueuedConnection);

//    // to fix qt animation bug, sometimes window size not change
    connect(m_sizeChangeAni, &QPropertyAnimation::valueChanged, [this] {
        const QSize size = m_sizeChangeAni->currentValue().toSize();

        QWidget::setFixedSize(size);
        m_mainPanel->setFixedSize(size);
    });

//    connect(m_posChangeAni, &QPropertyAnimation::finished, [this] {
//        QWidget::move(m_posChangeAni->endValue().toPoint());
//    });

}

const QPoint MainWindow::x11GetWindowPos()
{
    const auto disp = QX11Info::display();

    unsigned int unused;
    int x;
    int y;
    Window unused_window;

    XGetGeometry(disp, winId(), &unused_window, &x, &y, &unused, &unused, &unused, &unused);
    XFlush(disp);

    return QPoint(x, y);
}

void MainWindow::x11MoveWindow(const int x, const int y)
{
    const auto disp = QX11Info::display();

    XMoveWindow(disp, winId(), x, y);
    XFlush(disp);
}

void MainWindow::x11MoveResizeWindow(const int x, const int y, const int w, const int h)
{
    const auto disp = QX11Info::display();

    XMoveResizeWindow(disp, winId(), x, y, w, h);
    XFlush(disp);
}

void MainWindow::positionChanged(const Position prevPos)
{
    // paly hide animation and disable other animation
    m_updatePanelVisible = false;
    clearStrutPartial();
    narrow(prevPos);

    // reset position & layout and slide out
    QTimer::singleShot(200, this, [&] {
        resetPanelEnvironment(false);
        updateGeometry();
        expand();
    });

    // set strut
    QTimer::singleShot(400, this, [&] {
        setStrutPartial();
    });

    // reset to right environment when animation finished
    QTimer::singleShot(600, this, [&] {
        m_updatePanelVisible = true;
        updatePanelVisible();
    });
}

void MainWindow::updatePosition()
{
//    qDebug() << Q_FUNC_INFO;
    // all update operation need pass by timer
    Q_ASSERT(sender() == m_positionUpdateTimer);

    clearStrutPartial();
    updateGeometry();

    // make sure strut partial is set after the size/position animation;
    const int inter = qMax(m_sizeChangeAni->duration(), m_posChangeAni->duration());
    QTimer::singleShot(inter + 100, this, &MainWindow::setStrutPartial);
}

void MainWindow::updateGeometry()
{
    const Position position = m_settings->position();
    QSize size = m_settings->windowSize();

//    qDebug() << Q_FUNC_INFO << position << size;

    m_mainPanel->setFixedSize(size);
    m_mainPanel->updateDockPosition(position);
    m_mainPanel->updateDockDisplayMode(m_settings->displayMode());

    if (m_settings->hideState() == Hide)
    {
        m_sizeChangeAni->stop();
        m_posChangeAni->stop();
        switch (position)
        {
        case Top:
        case Bottom:    size.setHeight(1);      break;
        case Left:
        case Right:     size.setWidth(1);       break;
        }
        QWidget::setFixedSize(size);
    }
    else
    {
        setFixedSize(size);
    }

    //    const QRect primaryRect = m_settings->primaryRect();
    //    const int offsetX = (primaryRect.width() - size.width()) / 2;
    //    const int offsetY = (primaryRect.height() - size.height()) / 2;

    //    switch (position)
    //    {
    //    case Top:
    //        move(primaryRect.topLeft().x() + offsetX, primaryRect.y());                   break;
    //    case Left:
    //        move(primaryRect.topLeft().x(), primaryRect.y() + offsetY);                   break;
    //    case Right:
    //        move(primaryRect.right() - size.width() + 1, primaryRect.y() + offsetY);      break;
    //    case Bottom:
    //        move(primaryRect.x() + offsetX, primaryRect.bottom() - size.height() + 1);    break;
    //    default:
    //        Q_ASSERT(false);
    //    }
    const QRect windowRect = m_settings->windowRect(position, m_settings->hideState() == Hide);
    internalMove(windowRect.x(), windowRect.y());
    m_mainPanel->update();
}

void MainWindow::clearStrutPartial()
{
    m_xcbMisc->clear_strut_partial(winId());
}

void MainWindow::setStrutPartial()
{
//    qDebug() << Q_FUNC_INFO;
    // first, clear old strut partial
    clearStrutPartial();

    // reset env
    resetPanelEnvironment(true);

    if (m_settings->hideMode() != Dock::KeepShowing)
        return;

    const Position side = m_settings->position();
    const int maxScreenHeight = m_settings->screenHeight();
    const int maxScreenWidth = m_settings->screenWidth();

    XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
    uint strut = 0;
    uint strutStart = 0;
    uint strutEnd = 0;

    const QPoint p = m_posChangeAni->endValue().toPoint();
    const QRect r = QRect(p, m_settings->windowSize());
    QRect strutArea(0, 0, maxScreenWidth, maxScreenHeight);

    qDebug() << "screen info: " << r << p << strutArea;

    switch (side)
    {
    case Position::Top:
        orientation = XcbMisc::OrientationTop;
        strut = r.bottom() + 1;
        strutStart = r.left();
        strutEnd = r.right();
        strutArea.setLeft(strutStart);
        strutArea.setRight(strutEnd);
        strutArea.setBottom(r.bottom());
        break;
    case Position::Bottom:
        orientation = XcbMisc::OrientationBottom;
        strut = maxScreenHeight - p.y();
        strutStart = r.left();
        strutEnd = r.right();
        strutArea.setLeft(strutStart);
        strutArea.setRight(strutEnd);
        strutArea.setTop(p.y());
        break;
    case Position::Left:
        orientation = XcbMisc::OrientationLeft;
        strut = r.right() + 1;
        strutStart = r.top();
        strutEnd = r.bottom();
        strutArea.setTop(r.top());
        strutArea.setBottom(r.bottom());
        strutArea.setRight(r.right());
        break;
    case Position::Right:
        orientation = XcbMisc::OrientationRight;
        strut = maxScreenWidth - r.left();
        strutStart = r.top();
        strutEnd = r.bottom();
        strutArea.setTop(r.top());
        strutArea.setBottom(r.bottom());
        strutArea.setLeft(r.left());
        break;
    default:
        Q_ASSERT(false);
    }

    // pass if strut area is intersect with other screen
    int count = 0;
    const QRect pr = m_settings->primaryRect();
    for (auto *screen : qApp->screens())
    {
        const QRect sr = screen->geometry();
        if (sr == pr)
            continue;

        if (sr.intersects(strutArea))
            ++count;
    }
    if (count > 0)
        return;

    m_xcbMisc->set_strut_partial(winId(), orientation, strut * qApp->devicePixelRatio(), strutStart, strutEnd);
}

void MainWindow::expand()
{
//    qDebug() << "expand";
    const QPoint finishPos(0, 0);

    if (m_mainPanel->pos() == finishPos && m_mainPanel->size() == this->size() && m_panelHideAni->state() == QPropertyAnimation::Stopped)
        return;
    m_panelHideAni->stop();

    resetPanelEnvironment(true);

    if (m_panelShowAni->state() == QPropertyAnimation::Running)
        return m_panelShowAni->setEndValue(finishPos);

    const QSize size = m_settings->windowSize();

    QPoint startPos(0, 0);
    switch (m_settings->position())
    {
    case Top:       startPos.setY(-size.height());     break;
    case Bottom:    startPos.setY(size.height());      break;
    case Left:      startPos.setX(-size.width());      break;
    case Right:     startPos.setX(size.width());       break;
    }

    m_panelShowAni->setStartValue(startPos);
    m_panelShowAni->setEndValue(finishPos);
    m_panelShowAni->start();
}

void MainWindow::narrow(const Position prevPos)
{
//    qDebug() << "narrow" << prevPos;
    //    const QSize size = m_settings->windowSize();
    const QSize size = m_mainPanel->size();

    QPoint finishPos(0, 0);
    switch (prevPos)
    {
    case Top:       finishPos.setY(-size.height());     break;
    case Bottom:    finishPos.setY(size.height());      break;
    case Left:      finishPos.setX(-size.width());      break;
    case Right:     finishPos.setX(size.width());       break;
    }

    if (m_panelHideAni->state() == QPropertyAnimation::Running)
        return m_panelHideAni->setEndValue(finishPos);

    m_panelShowAni->stop();
    m_panelHideAni->setStartValue(m_mainPanel->pos());
    m_panelHideAni->setEndValue(finishPos);
    m_panelHideAni->start();

    // clear shadow
    adjustShadowMask();
}

void MainWindow::resetPanelEnvironment(const bool visible)
{
    // reset environment
    m_sizeChangeAni->stop();
    m_posChangeAni->stop();

    const Position position = m_settings->position();
    const QRect r(m_settings->windowRect(position));

//    qDebug() << Q_FUNC_INFO << r;
    m_sizeChangeAni->setEndValue(r.size());
    m_mainPanel->setFixedSize(r.size());
    QWidget::setFixedSize(r.size());
    m_posChangeAni->setEndValue(r.topLeft());
    QWidget::move(r.topLeft());
//    x11MoveWindow(r.topLeft().x(), r.topLeft().y());

    QPoint finishPos(0, 0);
    if (!visible)
    {
        switch (position)
        {
        case Top:       finishPos.setY(-r.height());     break;
        case Bottom:    finishPos.setY(r.height());      break;
        case Left:      finishPos.setX(-r.width());      break;
        case Right:     finishPos.setX(r.width());       break;
        }
    }

    m_mainPanel->move(finishPos);
}

void MainWindow::updatePanelVisible()
{
//    qDebug() << m_updatePanelVisible;

    if (!m_updatePanelVisible)
        return;
    if (m_settings->hideMode() == KeepShowing)
        return expand();

    const Dock::HideState state = m_settings->hideState();

//    qDebug() << state;

    do
    {
        if (state != Hide)
            break;

        if (!m_settings->autoHide())
            break;

        QRect r(pos(), size());
        if (r.contains(QCursor::pos()))
            break;

        return narrow(m_settings->position());

    } while (false);

    return expand();
}

void MainWindow::adjustShadowMask()
{
//    qDebug() << Q_FUNC_INFO << m_mainPanel->pos() << m_panelHideAni->state() << m_panelShowAni->state() << m_wmHelper->hasComposite();
    if (m_mainPanel->pos() != QPoint(0, 0) ||
        m_panelHideAni->state() == QPropertyAnimation::Running ||
        m_panelShowAni->state() == QPauseAnimation::Running ||
        !m_wmHelper->hasComposite())
    {
        m_platformWindowHandle.setShadowRadius(0);
        m_platformWindowHandle.setClipPath(QPainterPath());
        return;
    }

    const QRect r = QRect(QPoint(), rect().size());
    const int radius = 5;

    QPainterPath path;
    if (m_settings->displayMode() == DisplayMode::Fashion)
    {
        switch (m_settings->position())
        {
        case Top:
            path.addRoundedRect(0, -radius, r.width(), r.height() + radius, radius, radius);
            break;
        case Bottom:
            path.addRoundedRect(0, 0, r.width(), r.height() + radius, radius, radius);
            break;
        case Left:
            path.addRoundedRect(-radius, 0, r.width() + radius, r.height(), radius, radius);
            break;
        case Right:
            path.addRoundedRect(0, 0, r.width() + radius, r.height(), radius, radius);
        default:;
        }
    } else {
        path.addRect(r);
    }

    m_platformWindowHandle.setShadowRadius(60);
    m_platformWindowHandle.setClipPath(path);
}
