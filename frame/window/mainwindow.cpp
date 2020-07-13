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
#include "panel/mainpanelcontrol.h"
#include "controller/dockitemmanager.h"
#include "util/utils.h"
#include "util/imageutil.h"

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <QCursor>
#include <qpa/qplatformwindow.h>
#include <DStyle>
#include <DPlatformWindowHandle>
#include <QGSettings>

#include <X11/X.h>
#include <X11/Xutil.h>
#include <X11/Xcursor/Xcursor.h>

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"
#define SessionManagerService "com.deepin.SessionManager"
#define SessionManagerPath "/com/deepin/SessionManager"

#define WINDOWS_DOCK 0x00000081
#define MAINWINDOW_MAX_SIZE       DOCK_MAX_SIZE
#define MAINWINDOW_MIN_SIZE       (40)
#define DRAG_AREA_SIZE (5)

using org::kde::StatusNotifierWatcher;
using namespace com::deepin;

class DragWidget : public QWidget
{
    Q_OBJECT

private:
    bool m_dragStatus;
    QPoint m_resizePoint;

public:
    DragWidget(QWidget *parent) : QWidget(parent)
    {
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

    void mouseMoveEvent(QMouseEvent *event) override
    {
        if (m_dragStatus) {
            QPoint offset = QPoint(QCursor::pos() - m_resizePoint);
            emit dragPointOffset(offset);
        }
    }

    void mouseReleaseEvent(QMouseEvent *event) override
    {
        if (!m_dragStatus)
            return;

        m_dragStatus =  false;
        releaseMouse();
        emit dragFinished();
    }

    void updateCursor()
    {
        QGSettings gsetting("com.deepin.xsettings", "/com/deepin/xsettings/");
        QString theme = gsetting.get("gtk-cursor-theme-name").toString();
        Position position = DockSettings::Instance().position();
        int cursorSize = gsetting.get("gtk-cursor-theme-size").toInt();

        static QCursor *lastCursor = nullptr;
        static QString lastTheme;
        static int lastPosition = -1;
        static int lastCursorSize = -1;
        if (theme != lastTheme || position != lastPosition || cursorSize != lastCursorSize) {
            lastTheme = theme;
            lastPosition = position;
            lastCursorSize = cursorSize;
            const char* cursorName = (position == Bottom || position == Top) ? "v_double_arrow" : "h_double_arrow";
            QCursor *newCursor = ImageUtil::loadQCursorFromX11Cursor(theme.toStdString().c_str(), cursorName, cursorSize);
            setCursor(*newCursor);
            if (lastCursor != nullptr)
                delete lastCursor;

            lastCursor = newCursor;
        }
    }

    void enterEvent(QEvent *) override
    {
        updateCursor();
        if (QApplication::overrideCursor() && QApplication::overrideCursor()->shape() != cursor()) {
            QApplication::setOverrideCursor(cursor());
        }
    }

    void leaveEvent(QEvent *) override
    {
        QApplication::restoreOverrideCursor();
    }
};

const QPoint rawXPosition(const QPoint &scaledPos)
{
    QScreen const *screen = Utils::screenAtByScaled(scaledPos);

    return screen ? screen->geometry().topLeft() +
           (scaledPos - screen->geometry().topLeft()) *
           screen->devicePixelRatio()
           : scaledPos;
}

const QPoint scaledPos(const QPoint &rawXPos)
{
    QScreen const *screen = Utils::screenAt(rawXPos);

    return screen
           ? screen->geometry().topLeft() +
           (rawXPos - screen->geometry().topLeft()) / screen->devicePixelRatio()
           : rawXPos;
}

MainWindow::MainWindow(QWidget *parent)
    : DBlurEffectWidget(parent),

      m_launched(false),

      m_mainPanel(new MainPanelControl(this)),

      m_platformWindowHandle(this),
      m_wmHelper(DWindowManagerHelper::instance()),
      m_regionMonitor(new DRegionMonitor(this)),

      m_positionUpdateTimer(new QTimer(this)),
      m_expandDelayTimer(new QTimer(this)),
      m_leaveDelayTimer(new QTimer(this)),
      m_shadowMaskOptimizeTimer(new QTimer(this)),

      m_panelShowAni(new QVariantAnimation(this)),
      m_panelHideAni(new QVariantAnimation(this)),
      m_dbusDaemonInterface(QDBusConnection::sessionBus().interface()),
      m_sniWatcher(new StatusNotifierWatcher(SNI_WATCHER_SERVICE, SNI_WATCHER_PATH, QDBusConnection::sessionBus(), this)),
      m_sessionManagerInter(new SessionManager(SessionManagerService, SessionManagerPath, QDBusConnection::sessionBus(), this)),
      m_dragWidget(new DragWidget(this))
{
    Qt::WindowFlags flags = Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint | Qt::Window | (Qt::WindowFlags)WINDOWS_DOCK;
    setWindowFlags(windowFlags() | flags);
    setAttribute(Qt::WA_TranslucentBackground);
    setMouseTracking(true);
    setAcceptDrops(true);

    m_platformWindowHandle.setEnableBlurWindow(true);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setWindowRadius(0);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 5));
    m_platformWindowHandle.setShadowColor(QColor(0, 0, 0, 0.3 * 255));

    m_settings = &DockSettings::Instance();
    m_size = m_settings->m_mainWindowSize;
    m_mainPanel->setDisplayMode(m_settings->displayMode());
    initSNIHost();
    initComponents();
    initConnections();

    resizeMainPanelWindow();

    m_mainPanel->setDelegate(this);
    for (auto item : DockItemManager::instance()->itemList())
        m_mainPanel->insertItem(-1, item);

    m_dragWidget->setAccessibleName("Form_dragwidget");
    m_dragWidget->setMouseTracking(true);
    m_dragWidget->setFocusPolicy(Qt::NoFocus);

    m_curDockPos = m_settings->position();
    m_newDockPos = m_curDockPos;

    ///begin: 因为klu跟随控制中心主题的临时解决方案注释掉,后期会恢复,勿删
    // if ((Top == m_curDockPos) || (Bottom == m_curDockPos)) {
    //     m_dragWidget->setCursor(Qt::SizeVerCursor);
    // } else {
    //     m_dragWidget->setCursor(Qt::SizeHorCursor);
    // }
    ///end

    connect(m_panelShowAni, &QVariantAnimation::valueChanged, [ this ](const QVariant & value) {

        if (m_panelShowAni->state() != QPropertyAnimation::Running)
            return;

        // dock的宽度或高度值
        int val = value.toInt();
        // 当前dock尺寸
        const QRectF windowRect = m_settings->windowRect(m_curDockPos, false);

        switch (m_curDockPos) {
        case Dock::Top:
            m_mainPanel->move(0, val - windowRect.height());
            QWidget::move(windowRect.left(), windowRect.top());
            break;
        case Dock::Bottom:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.left(), windowRect.bottom() - val);
            break;
        case Dock::Left:
            m_mainPanel->move(val - windowRect.width(), 0);
            QWidget::move(windowRect.left(), windowRect.top());
            break;
        case Dock::Right:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.right() - val, windowRect.top());
            break;
        default: break;
        }

        if (m_curDockPos == Dock::Top || m_curDockPos == Dock::Bottom) {
            QWidget::setFixedHeight(val);
        } else {
            QWidget::setFixedWidth(val);
        }

        m_mainPanel->setFixedSize(windowRect.width(), windowRect.height());

    });

    connect(m_panelHideAni, &QVariantAnimation::valueChanged, [ this ](const QVariant & value) {

        if (m_panelHideAni->state() != QPropertyAnimation::Running)
            return;

        // dock的宽度或高度
        int val = value.toInt();
        // dock隐藏后的rect
        const QRectF windowRect = m_settings->windowRect(m_curDockPos, false);
        const int margin = m_settings->dockMargin();

        switch (m_curDockPos) {
        case Dock::Top:
            m_mainPanel->move(0, val - windowRect.height());
            QWidget::move(windowRect.left(), windowRect.top() - margin);
            break;
        case Dock::Bottom:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.left(), windowRect.bottom() - val + margin);
            break;
        case Dock::Left:
            m_mainPanel->move(val - windowRect.width(), 0);
            QWidget::move(windowRect.left() - margin, windowRect.top());
            break;
        case Dock::Right:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.right() - val + margin, windowRect.top());
            break;
        default: break;
        }

        if (m_curDockPos == Dock::Top || m_curDockPos == Dock::Bottom) {
            QWidget::setFixedHeight(val);
        } else {
            QWidget::setFixedWidth(val);
        }

    });

    connect(m_panelShowAni, &QVariantAnimation::finished, [ this ]() {
        const QRect windowRect = m_settings->windowRect(m_curDockPos, false);

        QWidget::move(windowRect.left(), windowRect.top());
        QWidget::setFixedSize(windowRect.size());

        m_mainPanel->move(QPoint(0, 0));

        resizeMainPanelWindow();
    });

    connect(m_panelHideAni, &QVariantAnimation::finished, [ this ]() {
        m_curDockPos = m_newDockPos;
        const QRect windowRect = m_settings->windowRect(m_curDockPos, true);

        QWidget::move(windowRect.left(), windowRect.top());
        QWidget::setFixedSize(windowRect.size());
        m_mainPanel->move(QPoint(0, 0));
        if (m_settings->hideMode() != KeepShowing)
            this->setVisible(false);
    });

    connect(m_sessionManagerInter, &SessionManager::LockedChanged, this, [ this ] (bool locked) {
        this->setVisible(!locked);
    });

    updateRegionMonitorWatch();
}

MainWindow::~MainWindow()
{
}

void MainWindow::launch()
{
    setVisible(false);
    QTimer *waitServerTimer = new QTimer(this);
    waitServerTimer->start(100);
    connect(waitServerTimer, &QTimer::timeout, this, [ = ] {
        QDBusInterface interface("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock",
                                 "com.deepin.dde.daemon.Dock",
                                 QDBusConnection::sessionBus());
        if(interface.isValid()) {
            waitServerTimer->stop();
            waitServerTimer->deleteLater();
            emit loaderPlugins();
            QTimer::singleShot(400, this, [&] {
                m_launched = true;
                qApp->processEvents();
                QWidget::move(m_settings->windowRect(m_curDockPos).topLeft());
                setVisible(true);
                updatePanelVisible();
                expand();
                resetPanelEnvironment(false);
            });
        }
    });
}

bool MainWindow::event(QEvent *e)
{
    switch (e->type()) {
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

    //    connect(qGuiApp, &QGuiApplication::primaryScreenChanged,
    //    windowHandle(), [this](QScreen * new_screen) {
    //        QScreen *old_screen = windowHandle()->screen();
    //        windowHandle()->setScreen(new_screen);
    //        // 屏幕变化后可能导致控件缩放比变化，此时应该重设控件位置大小
    //        // 比如：窗口大小为 100 x 100, 显示在缩放比为 1.0 的屏幕上，此时窗口的真实大小 = 100x100
    //        // 随后窗口被移动到了缩放比为 2.0 的屏幕上，应该将真实大小改为 200x200。另外，只能使用
    //        // QPlatformWindow直接设置大小来绕过QWidget和QWindow对新旧geometry的比较。
    //        const qreal scale = devicePixelRatioF();
    //        const QPoint screenPos = new_screen->geometry().topLeft();
    //        const QPoint posInScreen = this->pos() - old_screen->geometry().topLeft();
    //        const QPoint pos = screenPos + posInScreen * scale;
    //        const QSize size = this->size() * scale;

    //        windowHandle()->handle()->setGeometry(QRect(pos, size));
    //    }, Qt::UniqueConnection);

    //    windowHandle()->setScreen(qGuiApp->primaryScreen());
}

void MainWindow::mousePressEvent(QMouseEvent *e)
{
    e->ignore();
    if (e->button() == Qt::RightButton && m_settings->m_menuVisible) {
        m_settings->showDockSettingsMenu(mapToGlobal(e->pos()));
        return;
    } else if (e->button() == Qt::LeftButton && m_settings->m_menuVisible) {
        m_settings->hideDockSettingsMenu();
        return;
    }
}

void MainWindow::keyPressEvent(QKeyEvent *e)
{
    switch (e->key()) {
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

    if (QApplication::overrideCursor() && QApplication::overrideCursor()->shape() != Qt::ArrowCursor)
        QApplication::restoreOverrideCursor();
}

void MainWindow::mouseMoveEvent(QMouseEvent *e)
{
    //重写mouseMoveEvent 解决bug12866  leaveEvent事件失效
}

void MainWindow::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);
    if (m_panelHideAni->state() == QPropertyAnimation::Running)
        return;

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

void MainWindow::initSNIHost()
{
    // registor dock as SNI Host on dbus
    QDBusConnection dbusConn = QDBusConnection::sessionBus();
    m_sniHostService = QString("org.kde.StatusNotifierHost-") + QString::number(qApp->applicationPid());
    dbusConn.registerService(m_sniHostService);
    dbusConn.registerObject("/StatusNotifierHost", this);

    if (m_sniWatcher->isValid()) {
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    } else {
        qDebug() << SNI_WATCHER_SERVICE << "SNI watcher daemon is not exist for now!";
    }
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

    m_panelShowAni->setEasingCurve(QEasingCurve::InOutCubic);
    m_panelHideAni->setEasingCurve(QEasingCurve::InOutCubic);

    QTimer::singleShot(1, this, &MainWindow::compositeChanged);

    themeTypeChanged(DGuiApplicationHelper::instance()->themeType());
}

void MainWindow::compositeChanged()
{
    const bool composite = m_wmHelper->hasComposite();
    setComposite(composite);

    // NOTE(justforlxz): On the sw platform, there is an unstable
    // display position error, disable animation solution
#ifndef DISABLE_SHOW_ANIMATION
    const int duration = composite ? 300 : 0;
#else
    const int duration = 0;
#endif

    m_panelHideAni->setDuration(duration);
    m_panelShowAni->setDuration(duration);

    m_shadowMaskOptimizeTimer->start();
}

void MainWindow::internalMove(const QPoint &p)
{
    const bool isHide = m_settings->hideState() == HideState::Hide && !testAttribute(Qt::WA_UnderMouse);
    const bool pos_adjust = m_settings->hideMode() != HideMode::KeepShowing &&
                            isHide &&
                            m_panelShowAni->state() == QVariantAnimation::Stopped;
    if (!pos_adjust) {
        m_mainPanel->move(0, 0);
        return QWidget::move(p);
    }


    QPoint rp = rawXPosition(p);
    const auto ratio = devicePixelRatioF();

    const QRect &r = m_settings->primaryRawRect();
    switch (m_curDockPos) {
    case Left:      rp.setX(r.x());             break;
    case Top:       rp.setY(r.y());             break;
    case Right:     rp.setX(r.right() - 1);     break;
    case Bottom:    rp.setY(r.bottom() - 1);    break;
    }

    int hx = height() * ratio, wx = width() * ratio;
    if (m_settings->hideMode() != HideMode::KeepShowing &&
            isHide &&
            m_panelHideAni->state() == QVariantAnimation::Stopped &&
            m_panelShowAni->state() == QVariantAnimation::Stopped) {
        switch (m_curDockPos) {
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
    //    windowHandle()->handle()->setGeometry(QRect(rp.x(), rp.y(), wx, hx));
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::dataChanged, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::regionChanged, this, &MainWindow::updateRegion);
    connect(m_settings, &DockSettings::positionChanged, this, &MainWindow::positionChanged);
    connect(m_settings, &DockSettings::autoHideChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry, Qt::DirectConnection);
    connect(m_settings, &DockSettings::trayCountChanged, this, &MainWindow::getTrayVisableItemCount, Qt::DirectConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, [this] { resetPanelEnvironment(true); });
    connect(m_settings, &DockSettings::windowHideModeChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowVisibleChanged, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::displayModeChanegd, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(&DockSettings::Instance(), &DockSettings::opacityChanged, this, &MainWindow::setMaskAlpha);
    connect(m_settings, &DockSettings::displayModeChanegd, this, &MainWindow::updateDisplayMode, Qt::QueuedConnection);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition, Qt::QueuedConnection);  //+ 6-22-1
    connect(m_expandDelayTimer, &QTimer::timeout, this, &MainWindow::expand, Qt::QueuedConnection);
    connect(m_leaveDelayTimer, &QTimer::timeout, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_shadowMaskOptimizeTimer, &QTimer::timeout, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    connect(m_panelHideAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_panelShowAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &MainWindow::compositeChanged, Qt::QueuedConnection);
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &MainWindow::onDbusNameOwnerChanged);

    connect(DockItemManager::instance(), &DockItemManager::itemInserted, m_mainPanel, &MainPanelControl::insertItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemRemoved, m_mainPanel, &MainPanelControl::removeItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemUpdated, m_mainPanel, &MainPanelControl::itemUpdated, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::requestRefershWindowVisible, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(DockItemManager::instance(), &DockItemManager::requestWindowAutoHide, m_settings, &DockSettings::setAutoHide);
    connect(m_mainPanel, &MainPanelControl::itemMoved, DockItemManager::instance(), &DockItemManager::itemMoved, Qt::DirectConnection);
    connect(m_mainPanel, &MainPanelControl::itemAdded, DockItemManager::instance(), &DockItemManager::itemAdded, Qt::DirectConnection);
    connect(m_dragWidget, &DragWidget::dragPointOffset, this, &MainWindow::onMainWindowSizeChanged);
    connect(m_dragWidget, &DragWidget::dragFinished, this, &MainWindow::onDragFinished);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &MainWindow::themeTypeChanged);
    connect(m_regionMonitor, &DRegionMonitor::cursorMove, this, &MainWindow::onRegionMonitorChanged);
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

void MainWindow::positionChanged(const Position prevPos, const Position nextPos)
{
    m_newDockPos = nextPos;
    // paly hide animation and disable other animation
    clearStrutPartial();
    narrow(prevPos);

    // set strut
    QTimer::singleShot(400, this, [&] {
        setStrutPartial();
    });

    // reset to right environment when animation finished
    QTimer::singleShot(500, this, [&] {
        m_mainPanel->setPositonValue(m_curDockPos);
        resetPanelEnvironment(true);

        ///begin: 因为klu跟随控制中心主题的临时解决方案注释掉,后期会恢复,勿删
        // if ((Top == m_curDockPos) || (Bottom == m_curDockPos))
        // {
        //     m_dragWidget->setCursor(Qt::SizeVerCursor);
        // } else
        // {
        //     m_dragWidget->setCursor(Qt::SizeHorCursor);
        // }
        ///end

        updatePanelVisible();
    });

    updateRegionMonitorWatch();
}

void MainWindow::updatePosition()
{
    // all update operation need pass by timer
    Q_ASSERT(sender() == m_positionUpdateTimer);

    //clearStrutPartial();
    setStrutPartial();
    m_mainPanel->setDisplayMode(m_settings->displayMode());
    m_mainPanel->setPositonValue(m_curDockPos);

    if (DisplayMode::Efficient == m_settings->displayMode()) {
        QWidget::setFixedSize(m_settings->m_mainWindowSize);
    }

    resizeMainPanelWindow();
    m_mainPanel->update();
    m_leaveDelayTimer->start();
}

void MainWindow::updateGeometry()
{
    // DockDisplayMode and DockPosition MUST be set before invoke setFixedSize method of MainPanel
    setStrutPartial();

    m_mainPanel->setDisplayMode(m_settings->displayMode());
    m_mainPanel->setPositonValue(m_curDockPos);

    bool isHide = m_settings->hideState() == Hide && !testAttribute(Qt::WA_UnderMouse);

    const QRect windowRect = m_settings->windowRect(m_curDockPos, isHide);

    internalMove(windowRect.topLeft());

    QWidget::move(windowRect.topLeft());
    QWidget::setFixedSize(m_settings->m_mainWindowSize);

    resizeMainPanelWindow();

    m_mainPanel->update();
}

void MainWindow::getTrayVisableItemCount()
{
    m_mainPanel->getTrayVisableItemCount();
}

void MainWindow::clearStrutPartial()
{
    //m_xcbMisc->clear_strut_partial(winId());
}

void MainWindow::setStrutPartial()
{
    // first, clear old strut partial
    clearStrutPartial();

    // reset env
    //resetPanelEnvironment(true);

    if (m_settings->hideMode() != Dock::KeepShowing)
        return;

    const auto ratio = devicePixelRatioF();
    const int maxScreenHeight = m_settings->screenRawHeight();
    const int maxScreenWidth = m_settings->screenRawWidth();
    const Position side = m_curDockPos;
    const QPoint &p = rawXPosition(m_settings->windowRect(m_curDockPos).topLeft());
    const QSize &s = m_settings->windowSize();
    const QRect &primaryRawRect = m_settings->primaryRawRect();

    XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
    uint strut = 0;
    uint strutStart = 0;
    uint strutEnd = 0;

    QRect strutArea(0, 0, maxScreenWidth, maxScreenHeight);
    switch (side) {
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

    // pass if strut area is intersect with other screen
    //优化了文件管理的代码 会导致bug 15351 需要注释一下代码
//    int count = 0;
//    const QRect pr = m_settings->primaryRect();
//    for (auto *screen : qApp->screens()) {
//        const QRect sr = screen->geometry();
//        if (sr == pr)
//            continue;

//        if (sr.intersects(strutArea))
//            ++count;
//    }
//    if (count > 0) {
//        qWarning() << "strutArea is intersects with another screen.";
//        qWarning() << maxScreenHeight << maxScreenWidth << side << p << s;
//        return;
//    }

//    m_xcbMisc->set_strut_partial(winId(), orientation, strut + m_settings->dockMargin() * ratio, strutStart, strutEnd);
}

void MainWindow::expand()
{
    qApp->processEvents();
    setVisible(true);

    if (m_panelHideAni->state() == QPropertyAnimation::Running)
    {
        m_panelHideAni->stop();
        emit m_panelHideAni->finished();
    }

    const auto showAniState = m_panelShowAni->state();

    int startValue = 2;
    int endValue = 2;

    resetPanelEnvironment(true, false);
    if (showAniState != QPropertyAnimation::Running /*&& pos() != m_panelShowAni->currentValue()*/) {
        bool isHide = m_settings->hideState() == Hide && !testAttribute(Qt::WA_UnderMouse);
        const QRectF windowRect = m_settings->windowRect(m_curDockPos, isHide);
        switch (m_curDockPos) {
        case Top:
        case Bottom:
            startValue = height();
            endValue = windowRect.height();
            break;
        case Left:
        case Right:
            startValue = width();
            endValue = windowRect.width();
            break;
        }

        if (startValue > DOCK_MAX_SIZE || endValue > DOCK_MAX_SIZE) {
            return;
        }

        if (startValue > endValue)
            return;

        m_panelShowAni->setStartValue(startValue);
        m_panelShowAni->setEndValue(endValue);
        m_panelShowAni->start();
        m_shadowMaskOptimizeTimer->start();
    }
}

void MainWindow::narrow(const Position prevPos)
{
    int startValue = 2;
    int endValue = 2;

    switch (prevPos) {
    case Top:
    case Bottom:
        startValue = height();
        endValue = 2;
        break;
    case Left:
    case Right:
        startValue = width();
        endValue = 2;
        break;
    }

    m_panelShowAni->stop();
    m_panelHideAni->setStartValue(startValue);
    m_panelHideAni->setEndValue(endValue);
    m_panelHideAni->start();
}

void MainWindow::resetPanelEnvironment(const bool visible, const bool resetPosition)
{
    if (!m_launched)
        return;

    resizeMainPanelWindow();
    updateRegionMonitorWatch();
    if (m_size != m_settings->m_mainWindowSize) {
        m_size = m_settings->m_mainWindowSize;
        setStrutPartial();
    }
}

void MainWindow::updatePanelVisible()
{
    if (m_settings->hideMode() == KeepShowing) {
        if(m_regionMonitor->registered()){
            m_regionMonitor->unregisterRegion();
        }
        return expand();
    }

    if(!m_regionMonitor->registered()){
        m_regionMonitor->registerRegion();
        m_regionMonitor->setCoordinateType(DRegionMonitor::ScaleRatio);
    }

    const Dock::HideState state = m_settings->hideState();

    do {
        if (state != Hide)
            break;

        if (!m_settings->autoHide())
            break;

        QRectF r(pos(), size());
        const int margin = m_settings->dockMargin();
        switch (m_curDockPos) {
        case Dock::Top:
            r.setY(r.y() - margin);
            break;
        case Dock::Bottom:
            r.setHeight(r.height() + margin);
            break;
        case Dock::Left:
            r.setX(r.x() - margin);
            break;
        case Dock::Right:
            r.setWidth(r.width() + margin);
            break;
        }
        if (r.contains(QCursor::pos())) {
            break;
        }

        //        const QRect windowRect = m_settings->windowRect(m_curDockPos, true);
        //        move(windowRect.topLeft());

        return narrow(m_curDockPos);

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

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_TopLevelWindowRadius);

    m_platformWindowHandle.setWindowRadius(composite && isFasion ? radius : 0);
}

void MainWindow::positionCheck()
{
    if (m_positionUpdateTimer->isActive())
        return;

    const QPoint scaledFrontPos = scaledPos(m_settings->frontendWindowRect().topLeft());

    if (QPoint(pos() - scaledFrontPos).manhattanLength() < 2)
        return;

    // this may cause some position error and animation caton
    //internalMove();
}

void MainWindow::onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (name == SNI_WATCHER_SERVICE && !newOwner.isEmpty()) {
        qDebug() << SNI_WATCHER_SERVICE << "SNI watcher daemon started, register dock to watcher as SNI Host";
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    }
}

void MainWindow::setEffectEnabled(const bool enabled)
{
    setMaskColor(AutoColor);

    setMaskAlpha(DockSettings::Instance().Opacity());

    m_platformWindowHandle.setBorderWidth(enabled ? 1 : 0);
}

void MainWindow::setComposite(const bool hasComposite)
{
    setEffectEnabled(hasComposite);
}

bool MainWindow::appIsOnDock(const QString &appDesktop)
{
    return DockItemManager::instance()->appIsOnDock(appDesktop);
}

void MainWindow::resizeMainWindow()
{
    const Position position = m_curDockPos;
    QSize size = m_settings->windowSize();
    const QRect windowRect = m_settings->windowRect(position, false);
    internalMove(windowRect.topLeft());
    resizeMainPanelWindow();
    QWidget::setFixedSize(size);
}

void MainWindow::resizeMainPanelWindow()
{
    m_mainPanel->setFixedSize(m_settings->m_mainWindowSize);

    switch (m_curDockPos) {
    case Dock::Top:
        m_dragWidget->setGeometry(0, height() - DRAG_AREA_SIZE, width(), DRAG_AREA_SIZE);
        break;
    case Dock::Bottom:
        m_dragWidget->setGeometry(0, 0, width(), DRAG_AREA_SIZE);
        break;
    case Dock::Left:
        m_dragWidget->setGeometry(width() - DRAG_AREA_SIZE, 0, DRAG_AREA_SIZE, height());
        break;
    case Dock::Right:
        m_dragWidget->setGeometry(0, 0, DRAG_AREA_SIZE, height());
        break;
    default: break;
    }
}

void MainWindow::updateDisplayMode()
{
    m_mainPanel->setDisplayMode(m_settings->displayMode());
    setStrutPartial();
    adjustShadowMask();
    updateRegionMonitorWatch();
}

void MainWindow::onMainWindowSizeChanged(QPoint offset)
{
    if (Dock::Top == m_curDockPos) {
        m_settings->m_mainWindowSize.setHeight(qBound(MAINWINDOW_MIN_SIZE, m_size.height() + offset.y(), MAINWINDOW_MAX_SIZE));
        m_settings->m_mainWindowSize.setWidth(width());
    } else if (Dock::Bottom == m_curDockPos) {
        m_settings->m_mainWindowSize.setHeight(qBound(MAINWINDOW_MIN_SIZE, m_size.height() - offset.y(), MAINWINDOW_MAX_SIZE));
        m_settings->m_mainWindowSize.setWidth(width());
    } else if (Dock::Left == m_curDockPos) {
        m_settings->m_mainWindowSize.setHeight(height());
        m_settings->m_mainWindowSize.setWidth(qBound(MAINWINDOW_MIN_SIZE, m_size.width() + offset.x(), MAINWINDOW_MAX_SIZE));
    } else {
        m_settings->m_mainWindowSize.setHeight(height());
        m_settings->m_mainWindowSize.setWidth(qBound(MAINWINDOW_MIN_SIZE, m_size.width() - offset.x(), MAINWINDOW_MAX_SIZE));
    }

    resizeMainWindow();
    m_settings->updateFrontendGeometry();
}

void MainWindow::onDragFinished()
{
    if (m_size == m_settings->m_mainWindowSize)
        return;

    m_size = m_settings->m_mainWindowSize;

    if (m_settings->displayMode() == Fashion) {
        if (Dock::Top == m_curDockPos || Dock::Bottom == m_curDockPos) {
            m_settings->m_dockInter->setWindowSizeFashion(m_settings->m_mainWindowSize.height());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.height());
        } else {
            m_settings->m_dockInter->setWindowSizeFashion(m_settings->m_mainWindowSize.width());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.width());
        }
    } else {
        if (Dock::Top == m_curDockPos || Dock::Bottom == m_curDockPos) {
            m_settings->m_dockInter->setWindowSizeEfficient(m_settings->m_mainWindowSize.height());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.height());
        } else {
            m_settings->m_dockInter->setWindowSizeEfficient(m_settings->m_mainWindowSize.width());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.width());
        }
    }


    setStrutPartial();
}

void MainWindow::themeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    if (m_wmHelper->hasComposite()) {

        if (themeType == DGuiApplicationHelper::DarkType)
            m_platformWindowHandle.setBorderColor(QColor(0, 0, 0, 255 * 0.3));
        else
            m_platformWindowHandle.setBorderColor(QColor(QColor::Invalid));
    }
}

void MainWindow::onRegionMonitorChanged(const QPoint &p)
{
    if (m_settings->hideMode() == KeepShowing)
        return;

    if (!isVisible())
        setVisible(true);
}

void MainWindow::updateRegion()
{
    if ((m_settings->hideMode() == KeepHidden) || (m_settings->hideMode() == SmartHide)) {
        updateRegionMonitorWatch();
    }
}

void MainWindow::updateRegionMonitorWatch()
{
    if (m_settings->hideMode() == KeepShowing)
        return;
    int val = 5;

    bool isHide = m_settings->hideState() == Hide && !testAttribute(Qt::WA_UnderMouse);
    const QRect windowRect = m_settings->windowRect(m_curDockPos, isHide);
    const qreal scale = devicePixelRatioF();
    const int margin = m_settings->dockMargin();
    int x, y, w, h;

    if (Dock::Top == m_curDockPos) {
            x = windowRect.topLeft().x();
            y = windowRect.topLeft().y();
            w = m_settings->primaryRect().width();
            h = val+ margin;
        } else if (Dock::Bottom == m_curDockPos) {
            x = windowRect.bottomLeft().x();
            y = windowRect.bottomLeft().y() - val;
            w = m_settings->primaryRect().width();
            h = val+ margin;
        } else if (Dock::Left == m_curDockPos) {
            x = windowRect.topLeft().x();
            y = windowRect.topLeft().y();
            h = val+ margin;
            h = m_settings->primaryRect().height();
        } else {
            x = windowRect.topRight().x() - val - margin;
            y = windowRect.topRight().y();
            w = m_settings->primaryRect().width();
            h = m_settings->primaryRect().height();
        }

    m_regionMonitor->setWatchedRegion(QRegion(x * scale, y * scale, w * scale, h * scale));
}


#include "mainwindow.moc"
