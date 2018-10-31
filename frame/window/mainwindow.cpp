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

#include "mainwindow.h"
#include "panel/mainpanel.h"

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <qpa/qplatformwindow.h>

#include <DPlatformWindowHandle>

#include <X11/X.h>
#include <X11/Xutil.h>

const QPoint rawXPosition(const QPoint &scaledPos)
{
    QRect g = qApp->primaryScreen()->geometry();
    for (auto *screen : qApp->screens())
    {
        const QRect &sg = screen->geometry();
        if (sg.contains(scaledPos))
        {
            g = sg;
            break;
        }
    }

    return g.topLeft() + (scaledPos - g.topLeft()) * qApp->devicePixelRatio();
}

const QPoint scaledPos(const QPoint &rawXPos)
{
    QRect g = qApp->primaryScreen()->geometry();
    for (auto *screen : qApp->screens())
    {
        const QRect &sg = screen->geometry();
        if (sg.contains(rawXPos))
        {
            g = sg;
            break;
        }
    }

    return g.topLeft() + (rawXPos - g.topLeft()) / qApp->devicePixelRatio();
}

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),

      m_launched(false),
      m_updatePanelVisible(false),

      m_mainPanel(new MainPanel(this)),

      m_platformWindowHandle(this),
      m_wmHelper(DWindowManagerHelper::instance()),

      m_positionUpdateTimer(new QTimer(this)),
      m_expandDelayTimer(new QTimer(this)),
      m_leaveDelayTimer(new QTimer(this)),
      m_shadowMaskOptimizeTimer(new QTimer(this)),

      m_sizeChangeAni(new QVariantAnimation(this)),
      m_posChangeAni(new QVariantAnimation(this)),
      m_panelShowAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_panelHideAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_xcbMisc(XcbMisc::instance())

{
    setAccessibleName("dock-mainwindow");
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowDoesNotAcceptFocus);
    setAttribute(Qt::WA_TranslucentBackground);
    setMouseTracking(true);
    setAcceptDrops(true);

    DPlatformWindowHandle::enableDXcbForWindow(this, true);
    m_platformWindowHandle.setEnableBlurWindow(false);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setWindowRadius(0);
    m_platformWindowHandle.setBorderWidth(0);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 0));
    m_platformWindowHandle.setShadowRadius(0);

    m_settings = &DockSettings::Instance();
    m_xcbMisc->set_window_type(winId(), XcbMisc::Dock);

    initComponents();
    initConnections();

    m_mainPanel->setFixedSize(m_settings->panelSize());
}

MainWindow::~MainWindow()
{
    delete m_xcbMisc;
}

void MainWindow::launch()
{
    m_updatePanelVisible = false;
    m_mainPanel->setVisible(false);
    resetPanelEnvironment(false);
    setVisible(false);

    QTimer::singleShot(400, this, [&] {
        m_launched = true;
        m_mainPanel->setVisible(true);
        resetPanelEnvironment(false);
        updateGeometry();
        expand();
    });

    // set strut
    QTimer::singleShot(600, this, [&] {
        setStrutPartial();
    });

    // reset to right environment when animation finished
    QTimer::singleShot(800, this, [&] {
        m_updatePanelVisible = true;
        updatePanelVisible();
    });

    qApp->processEvents();
    QTimer::singleShot(300, this, &MainWindow::show);
}

bool MainWindow::event(QEvent *e)
{
    switch (e->type())
    {
    case QEvent::Move:
        if (!e->spontaneous())
            QTimer::singleShot(1, this, &MainWindow::positionCheck);
        break;
    default:;
    }

    return QWidget::event(e);
}

void MainWindow::showEvent(QShowEvent *e)
{
    QWidget::showEvent(e);

    m_platformWindowHandle.setEnableBlurWindow(false);
    m_platformWindowHandle.setShadowOffset(QPoint());
    m_platformWindowHandle.setShadowRadius(0);
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

    m_leaveDelayTimer->stop();
    if (m_settings->hideState() != Show && m_panelShowAni->state() != QPropertyAnimation::Running)
        m_expandDelayTimer->start();
}

void MainWindow::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_expandDelayTimer->stop();
    m_leaveDelayTimer->start();
}

void MainWindow::dragEnterEvent(QDragEnterEvent *e)
{
    QWidget::dragEnterEvent(e);

    if (m_settings->hideState() != Show) {
        m_expandDelayTimer->start();
    }
}

void MainWindow::setFixedSize(const QSize &size)
{
    const QPropertyAnimation::State state = m_sizeChangeAni->state();

    if (state == QPropertyAnimation::Stopped && this->size() == size)
        return;

    if (state == QPropertyAnimation::Running)
        return m_sizeChangeAni->setEndValue(size);

    m_sizeChangeAni->setStartValue(this->size());
    m_sizeChangeAni->setEndValue(size);
    m_sizeChangeAni->start();
}

void MainWindow::internalAnimationMove(int x, int y)
{
    const QPropertyAnimation::State state = m_posChangeAni->state();
    const QPoint p = m_posChangeAni->endValue().toPoint();
    const QPoint tp = QPoint(x, y);

    if (state == QPropertyAnimation::Stopped && p == tp)
        return;

    if (state == QPropertyAnimation::Running && p != tp)
        return m_posChangeAni->setEndValue(QPoint(x, y));

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

    m_leaveDelayTimer->setSingleShot(true);
    m_leaveDelayTimer->setInterval(m_settings->narrowTimeout());

    m_shadowMaskOptimizeTimer->setSingleShot(true);
    m_shadowMaskOptimizeTimer->setInterval(100);

    m_sizeChangeAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_posChangeAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_panelShowAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_panelHideAni->setEasingCurve(QEasingCurve::InOutCubic);

    QTimer::singleShot(1, this, &MainWindow::compositeChanged);
}

void MainWindow::compositeChanged()
{
    const bool composite = m_wmHelper->hasComposite();
    const int duration = composite ? 300 : 0;

    m_sizeChangeAni->setDuration(duration);
    m_posChangeAni->setDuration(duration);
    m_panelShowAni->setDuration(duration);
    m_panelHideAni->setDuration(duration);

    m_mainPanel->setComposite(composite);

    m_shadowMaskOptimizeTimer->start();
    m_positionUpdateTimer->start();
}

void MainWindow::internalMove(const QPoint &p)
{
    const bool isHide = m_settings->hideState() == HideState::Hide && !testAttribute(Qt::WA_UnderMouse);
    const bool pos_adjust = m_settings->hideMode() != HideMode::KeepShowing &&
                             isHide &&
                             m_posChangeAni->state() == QVariantAnimation::Stopped;
    if (!pos_adjust)
        return QWidget::move(p);

    QPoint rp = rawXPosition(p);
    const auto ratio = devicePixelRatioF();

    const QRect &r = m_settings->primaryRawRect();
    switch (m_settings->position())
    {
    case Left:      rp.setX(r.x());             break;
    case Top:       rp.setY(r.y());             break;
    case Right:     rp.setX(r.right() - 1);     break;
    case Bottom:    rp.setY(r.bottom() - 1);    break;
    }

    int hx = height() * ratio, wx = width() * ratio;
    if (m_settings->hideMode() != HideMode::KeepShowing &&
        isHide &&
        m_panelHideAni->state() == QVariantAnimation::Stopped &&
        m_panelShowAni->state() == QVariantAnimation::Stopped)
    {
        switch (m_settings->position())
        {
        case Top:
        case Bottom:
            hx = 2;
            break;
        case Left:
        case Right:
            wx = 2;
        }
    }

    // using platform window to set real window position
    windowHandle()->handle()->setGeometry(QRect(rp.x(), rp.y(), wx, hx));
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::dataChanged, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::positionChanged, this, &MainWindow::positionChanged);
    connect(m_settings, &DockSettings::autoHideChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry, Qt::DirectConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, [this] { resetPanelEnvironment(true); });
    connect(m_settings, &DockSettings::windowHideModeChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowVisibleChanged, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::displayModeChanegd, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_mainPanel, &MainPanel::requestRefershWindowVisible, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_mainPanel, &MainPanel::requestWindowAutoHide, m_settings, &DockSettings::setAutoHide);
    connect(m_mainPanel, &MainPanel::geometryChanged, this, &MainWindow::panelGeometryChanged);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition, Qt::QueuedConnection);
    connect(m_expandDelayTimer, &QTimer::timeout, this, &MainWindow::expand, Qt::QueuedConnection);
    connect(m_leaveDelayTimer, &QTimer::timeout, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_shadowMaskOptimizeTimer, &QTimer::timeout, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::updateGeometry, Qt::QueuedConnection);
    connect(m_panelHideAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_panelShowAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_posChangeAni, &QVariantAnimation::valueChanged, this, static_cast<void (MainWindow::*)()>(&MainWindow::internalMove));
    connect(m_posChangeAni, &QVariantAnimation::finished, this, static_cast<void (MainWindow::*)()>(&MainWindow::internalMove), Qt::QueuedConnection);

    // to fix qt animation bug, sometimes window size not change
    connect(m_sizeChangeAni, &QPropertyAnimation::valueChanged, [=](const QVariant &value) {
        QWidget::setFixedSize(value.toSize());
    });

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &MainWindow::compositeChanged, Qt::QueuedConnection);
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
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
        resetPanelEnvironment(false, true);
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
    // all update operation need pass by timer
    Q_ASSERT(sender() == m_positionUpdateTimer);

    clearStrutPartial();
    updateGeometry();

    // make sure strut partial is set after the size/position animation;
    const int duration = qMax(m_sizeChangeAni->duration(), m_posChangeAni->duration());

    QTimer::singleShot(duration, this, &MainWindow::setStrutPartial);
    QTimer::singleShot(duration, this, &MainWindow::updatePanelVisible);
}

void MainWindow::updateGeometry()
{
    const Position position = m_settings->position();
    QSize size = m_settings->windowSize();

    // this->setFixedSize has been overrided for size animation
    m_mainPanel->setFixedSize(m_settings->panelSize());
    m_mainPanel->updateDockPosition(position);
    m_mainPanel->updateDockDisplayMode(m_settings->displayMode());

    bool animation = true;
    bool isHide = m_settings->hideState() == Hide && !testAttribute(Qt::WA_UnderMouse);

    if (isHide) {
        m_sizeChangeAni->stop();
        m_posChangeAni->stop();
        switch (position) {
        case Top:
        case Bottom:    size.setHeight(2);      break;
        case Left:
        case Right:     size.setWidth(2);       break;
        }
        animation = false;
        m_sizeChangeAni->setEndValue(size);
        QWidget::setFixedSize(size);
    } else {
        // this->setFixedSize has been overrided for size animation
        // 如果要增大则直接设置大小, 如果要缩小则使用重写的setFixedSize函数使用动画缩小
        // 因为缩小时如果也直接设置大小则会无法正常显示panel的缩小动画
        switch (position) {
        case Dock::Position::Top:
        case Dock::Position::Bottom: {
            if (size.width() >= this->size().width()) {
                QWidget::setFixedSize(size);
            } else {
                setFixedSize(size);
            }
            break;
        }
        case Dock::Position::Left:
        case Dock::Position::Right: {
            if (size.height() >= this->size().height()) {
                QWidget::setFixedSize(size);
            } else {
                setFixedSize(size);
            }
            break;
        }
        default:
            break;
        }
    }

    const QRect windowRect = m_settings->windowRect(position, isHide);

    if (animation)
        internalAnimationMove(windowRect.x(), windowRect.y());
    else
        internalMove(windowRect.topLeft());

    m_mainPanel->update();
    m_shadowMaskOptimizeTimer->start();
}

void MainWindow::clearStrutPartial()
{
    m_xcbMisc->clear_strut_partial(winId());
}

void MainWindow::setStrutPartial()
{
    // first, clear old strut partial
    clearStrutPartial();

    // reset env
    resetPanelEnvironment(true);

    if (m_settings->hideMode() != Dock::KeepShowing)
        return;

    const auto ratio = devicePixelRatioF();
    const int maxScreenHeight = m_settings->screenRawHeight();
    const int maxScreenWidth = m_settings->screenRawWidth();
    const Position side = m_settings->position();
    const QPoint &p = rawXPosition(m_posChangeAni->endValue().toPoint());
    const QSize &s = m_settings->windowSize();
    const QRect &primaryRawRect = m_settings->primaryRawRect();

    XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
    uint strut = 0;
    uint strutStart = 0;
    uint strutEnd = 0;

    QRect strutArea(0, 0, maxScreenWidth, maxScreenHeight);
    switch (side)
    {
    case Position::Top:
        orientation = XcbMisc::OrientationTop;
        strut = p.y() + s.height() * ratio;
        strutStart = p.x();
        strutEnd = qMin(qRound(p.x() + s.width() * ratio), primaryRawRect.right());
        strutArea.setLeft(strutStart);
        strutArea.setRight(strutEnd);
        strutArea.setBottom(strut);
        break;
    case Position::Bottom:
        orientation = XcbMisc::OrientationBottom;
        strut = maxScreenHeight - p.y();
        strutStart = p.x();
        strutEnd = qMin(qRound(p.x() + s.width() * ratio), primaryRawRect.right());
        strutArea.setLeft(strutStart);
        strutArea.setRight(strutEnd);
        strutArea.setTop(p.y());
        break;
    case Position::Left:
        orientation = XcbMisc::OrientationLeft;
        strut = p.x() + s.width() * ratio;
        strutStart = p.y();
        strutEnd = qMin(qRound(p.y() + s.height() * ratio), primaryRawRect.bottom());
        strutArea.setTop(strutStart);
        strutArea.setBottom(strutEnd);
        strutArea.setRight(strut);
        break;
    case Position::Right:
        orientation = XcbMisc::OrientationRight;
        strut = maxScreenWidth - p.x();
        strutStart = p.y();
        strutEnd = qMin(qRound(p.y() + s.height() * ratio), primaryRawRect.bottom());
        strutArea.setTop(strutStart);
        strutArea.setBottom(strutEnd);
        strutArea.setLeft(p.x());
        break;
    default:
        Q_ASSERT(false);
    }

    qDebug() << "screen info: " << p << strutArea;

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
    {
        qWarning() << "strutArea is intersects with another screen.";
        qWarning() << maxScreenHeight << maxScreenWidth << side << p << s;
        return;
    }

    m_xcbMisc->set_strut_partial(winId(), orientation, strut, strutStart, strutEnd);
}

void MainWindow::expand()
{
    qApp->processEvents();

    const auto showAniState = m_panelShowAni->state();
    m_panelHideAni->stop();

    QPoint finishPos(0, 0);
    switch (m_settings->position())
    {
    case Left:  finishPos.setX(WINDOW_OVERFLOW);    break;
    case Top:   finishPos.setY(WINDOW_OVERFLOW);    break;
    default:;
    }

    resetPanelEnvironment(true, false);

    if (showAniState != QPropertyAnimation::Running && m_mainPanel->pos() != m_panelShowAni->currentValue())
    {
        QPoint startPos(0, 0);
        const QSize &size = m_settings->windowSize();
        switch (m_settings->position())
        {
        case Top:       startPos.setY(-size.height() + WINDOW_OVERFLOW);     break;
        case Bottom:    startPos.setY(size.height());      break;
        case Left:      startPos.setX(-size.width() + WINDOW_OVERFLOW);      break;
        case Right:     startPos.setX(size.width());       break;
        }

        m_panelShowAni->setStartValue(startPos);
        m_panelShowAni->setEndValue(finishPos);
        m_panelShowAni->start();
        m_shadowMaskOptimizeTimer->start();
        m_platformWindowHandle.setShadowRadius(0);
    }
}

void MainWindow::narrow(const Position prevPos)
{
    const QSize size = m_settings->panelSize();

    QPoint finishPos(0, 0);
    switch (prevPos)
    {
    case Top:       finishPos.setY(-size.height() + WINDOW_OVERFLOW);     break;
    case Bottom:    finishPos.setY(size.height());      break;
    case Left:      finishPos.setX(-size.width() + WINDOW_OVERFLOW);      break;
    case Right:     finishPos.setX(size.width());       break;
    }

    m_panelShowAni->stop();
    m_panelHideAni->setStartValue(m_mainPanel->pos());
    m_panelHideAni->setEndValue(finishPos);
    m_panelHideAni->start();
    m_platformWindowHandle.setShadowRadius(0);
}

void MainWindow::resetPanelEnvironment(const bool visible, const bool resetPosition)
{
    if (!m_launched)
        return;

    // reset environment
    m_sizeChangeAni->stop();
    m_posChangeAni->stop();

    const Position position = m_settings->position();
    const QRect r(m_settings->windowRect(position));

    m_sizeChangeAni->setEndValue(r.size());
    m_mainPanel->setFixedSize(m_settings->panelSize());
    QWidget::setFixedSize(r.size());
    m_posChangeAni->setEndValue(r.topLeft());
    QWidget::move(r.topLeft());

    if (!resetPosition)
        return;
    QPoint finishPos(0, 0);
    switch (position)
    {
    case Top:       finishPos.setY((visible ? WINDOW_OVERFLOW : -r.height()));     break;
    case Bottom:    finishPos.setY(visible ? 0 : r.height());      break;
    case Left:      finishPos.setX((visible ? WINDOW_OVERFLOW : -r.width()));       break;
    case Right:     finishPos.setX(visible ? 0 : r.width());       break;
    }

    m_mainPanel->move(finishPos);
    m_mainPanel->updateDockPosition(position);
}

void MainWindow::updatePanelVisible()
{
    if (!m_updatePanelVisible)
        return;
    if (m_settings->hideMode() == KeepShowing)
        return expand();

    const Dock::HideState state = m_settings->hideState();

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
    if (!m_launched)
        return;

    if (m_shadowMaskOptimizeTimer->isActive())
        return;

    const bool composite = m_wmHelper->hasComposite();
    const bool isFasion = m_settings->displayMode() == Fashion;

    m_platformWindowHandle.setWindowRadius(composite && isFasion ? 5 : 0);
}

void MainWindow::positionCheck()
{
    if (m_posChangeAni->state() == QPropertyAnimation::Running)
        return;
    if (m_positionUpdateTimer->isActive())
        return;

    const QPoint scaledFrontPos = scaledPos(m_settings->frontendWindowRect().topLeft());

    if (QPoint(pos() - scaledFrontPos).manhattanLength() < 2)
        return;

    qWarning() << "Dock position may error!!!!!";
    qDebug() << pos() << m_settings->frontendWindowRect() << m_settings->windowRect(m_settings->position(), false);

    // this may cause some position error and animation caton
    //internalMove();
}
