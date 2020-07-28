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

#include "mainwindow.h"
#include "panel/mainpanelcontrol.h"
#include "controller/dockitemmanager.h"
#include "util/utils.h"
#include "util/docksettings.h"

#include <DStyle>
#include <DPlatformWindowHandle>

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <qpa/qplatformwindow.h>

#include <X11/X.h>
#include <X11/Xutil.h>

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"

#define MAINWINDOW_MAX_SIZE       DOCK_MAX_SIZE
#define MAINWINDOW_MIN_SIZE       (40)
#define DRAG_AREA_SIZE (5)

using org::kde::StatusNotifierWatcher;

class DragWidget : public QWidget
{
    Q_OBJECT

private:
    bool m_dragStatus;
    QPoint m_resizePoint;

public:
    explicit DragWidget(QWidget *parent) : QWidget(parent)
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

    void enterEvent(QEvent *) override
    {
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
    : DBlurEffectWidget(parent)
    , m_launched(false)
    , m_mainPanel(new MainPanelControl(this))
    , m_platformWindowHandle(this)
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_eventInter(new XEventMonitor("com.deepin.api.XEventMonitor", "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus()))
    , m_launcherInter(new DBusLuncher("com.deepin.dde.Launcher","/com/deepin/dde/Launcher",QDBusConnection::sessionBus()))
    , m_positionUpdateTimer(new QTimer(this))
    , m_expandDelayTimer(new QTimer(this))
    , m_leaveDelayTimer(new QTimer(this))
    , m_shadowMaskOptimizeTimer(new QTimer(this))
    , m_panelShowAni(new QVariantAnimation(this))
    , m_panelHideAni(new QVariantAnimation(this))
    , m_xcbMisc(XcbMisc::instance())
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    , m_sniWatcher(new StatusNotifierWatcher(SNI_WATCHER_SERVICE, SNI_WATCHER_PATH, QDBusConnection::sessionBus(), this))
    , m_dragWidget(new DragWidget(this))
    , m_primaryScreenChanged(false)
{
    setAccessibleName("mainwindow");
    m_mainPanel->setAccessibleName("mainpanel");
    setAttribute(Qt::WA_TranslucentBackground);
    setMouseTracking(true);
    setAcceptDrops(true);

    DPlatformWindowHandle::enableDXcbForWindow(this, true);
    m_platformWindowHandle.setEnableBlurWindow(true);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setWindowRadius(0);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 5));
    m_platformWindowHandle.setShadowColor(QColor(0, 0, 0, 0.3 * 255));

    m_settings = &DockSettings::Instance();
    m_xcbMisc->set_window_type(winId(), XcbMisc::Dock);
    m_size = m_settings->m_mainWindowSize;
    m_mainPanel->setDisplayMode(m_settings->displayMode());
    initSNIHost();
    initComponents();
    initConnections();

    resizeMainPanelWindow();

    m_mainPanel->setDelegate(this);
    for (auto item : DockItemManager::instance()->itemList())
        m_mainPanel->insertItem(-1, item);

    m_dragWidget->setMouseTracking(true);
    m_dragWidget->setFocusPolicy(Qt::NoFocus);

    m_dockPosition = m_settings->position();

    if ((Top == m_dockPosition) || (Bottom == m_dockPosition)) {
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    } else {
        m_dragWidget->setCursor(Qt::SizeHorCursor);
    }

    connect(m_panelShowAni, &QVariantAnimation::valueChanged, [ this ](const QVariant & value) {
        if (m_panelShowAni->state() != QPropertyAnimation::Running)
            return;
        // dock的宽度或高度值
        int val = value.toInt();
        // 当前dock尺寸
        const QRect windowRect = m_settings->windowRect(m_dockPosition, false);

        switch (m_dockPosition) {
        case Dock::Top:
            m_mainPanel->move(0, val - windowRect.height());
            QWidget::move(windowRect.topLeft());
            break;
        case Dock::Bottom:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.left(), windowRect.bottom() - val);
            break;
        case Dock::Left:
            m_mainPanel->move(val - windowRect.width(), 0);
            QWidget::move(windowRect.topLeft());
            break;
        case Dock::Right:
            m_mainPanel->move(0, 0);
            QWidget::move(windowRect.right() - val, windowRect.top());
            break;
        }

        if (m_dockPosition == Dock::Top || m_dockPosition == Dock::Bottom) {
            QWidget::setFixedHeight(val);
        } else {
            QWidget::setFixedWidth(val);
        }
    });

    connect(m_panelHideAni, &QVariantAnimation::valueChanged, [ this ](const QVariant & value) {
        if (m_panelHideAni->state() != QPropertyAnimation::Running)
            return;

        // dock的宽度或高度
        int val = value.toInt();
        // dock隐藏后的rect
        const QRect windowRect = m_settings->windowRect(m_dockPosition, false, true);
        const int margin = m_settings->dockMargin();
        switch (m_dockPosition) {
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
        }

        if (m_dockPosition == Dock::Top || m_dockPosition == Dock::Bottom) {
            QWidget::setFixedHeight(val);
        } else {
            QWidget::setFixedWidth(val);
        }
    });

    connect(m_panelShowAni, &QVariantAnimation::finished, [ this ]() {
        const QRect windowRect = m_settings->windowRect(m_dockPosition);
        QWidget::setFixedSize(windowRect.size());
        QWidget::move(windowRect.topLeft());
        m_mainPanel->move(QPoint(0, 0));
        qDebug() << "Show animation finished:" << frameGeometry();
        qDebug() << "Show animation finished not frame:" << geometry();
        QWidget::update();
    });

    connect(m_panelHideAni, &QVariantAnimation::finished, [ this ]() {
        // 动画完成更新dock位置
        m_dockPosition = m_settings->position();
        // 动画完成更新dock设置
        m_settings->posChangedUpdateSettings();

        const QRect windowRect = m_settings->windowRect(m_dockPosition, true);
        QWidget::setFixedSize(windowRect.size());
        QWidget::move(windowRect.topLeft());
        m_mainPanel->move(QPoint(0, 0));

        qDebug() << "Hide animation finished" << frameGeometry();
        qDebug() << "Hide animation finished not frame:" << geometry();
        QWidget::update();
    });

    updateRegionMonitorWatch();
}

MainWindow::~MainWindow()
{
    delete m_xcbMisc;
}

void MainWindow::launch()
{
    m_launched = true;
    qApp->processEvents();
    QWidget::move(m_settings->windowRect(m_dockPosition).topLeft());
    setVisible(true);
    updatePanelVisible();
    resetPanelEnvironment();
    // 用于更新mainwindow圆角
    m_shadowMaskOptimizeTimer->start();
}

bool MainWindow::event(QEvent *e)
{
//    switch (e->type()) {
//    case QEvent::Move:
//        if (!e->spontaneous())
//            QTimer::singleShot(100, this, &MainWindow::positionCheck);
//        break;
//    default:;
//    }

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
        m_settings->showDockSettingsMenu();
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
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::primaryScreenChanged, [&]() {
        m_primaryScreenChanged = true;
        updatePosition();
        m_primaryScreenChanged = false;
    });
    connect(m_settings, &DockSettings::positionChanged, this, &MainWindow::positionChanged);
    connect(m_settings, &DockSettings::autoHideChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry, Qt::DirectConnection);
    connect(m_settings, &DockSettings::trayCountChanged, this, &MainWindow::getTrayVisableItemCount, Qt::DirectConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, [this] { resetPanelEnvironment(); });
    connect(m_settings, &DockSettings::windowHideModeChanged, m_leaveDelayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowVisibleChanged, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::displayModeChanegd, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(&DockSettings::Instance(), &DockSettings::opacityChanged, this, &MainWindow::setMaskAlpha);
    connect(m_settings, &DockSettings::displayModeChanegd, this, &MainWindow::updateDisplayMode, Qt::QueuedConnection);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition, Qt::QueuedConnection);
    connect(m_expandDelayTimer, &QTimer::timeout, this, &MainWindow::expand, Qt::QueuedConnection);
    connect(m_leaveDelayTimer, &QTimer::timeout, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_shadowMaskOptimizeTimer, &QTimer::timeout, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    connect(m_panelHideAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_panelShowAni, &QPropertyAnimation::finished, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::panelGeometryChanged);
    connect(m_panelShowAni, &QPropertyAnimation::finished, this, &MainWindow::panelGeometryChanged);

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &MainWindow::compositeChanged, Qt::QueuedConnection);
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    if (m_dbusDaemonInterface && m_dbusDaemonInterface->isValid())
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
    connect(m_eventInter, &XEventMonitor::CursorMove, this, &MainWindow::onRegionMonitorChanged);
    connect(m_settings, &DockSettings::requestUpdateRegionWatch, this, &MainWindow::updateRegionMonitorWatch);
}

//const QPoint MainWindow::x11GetWindowPos()
//{
//    const auto disp = QX11Info::display();

//    unsigned int unused;
//    int x;
//    int y;
//    Window unused_window;

//    XGetGeometry(disp, winId(), &unused_window, &x, &y, &unused, &unused, &unused, &unused);
//    XFlush(disp);

//    return QPoint(x, y);
//}

//void MainWindow::x11MoveWindow(const int x, const int y)
//{
//    const auto disp = QX11Info::display();

//    XMoveWindow(disp, winId(), x, y);
//    XFlush(disp);
//}

//void MainWindow::x11MoveResizeWindow(const int x, const int y, const int w, const int h)
//{
//    const auto disp = QX11Info::display();

//    XMoveResizeWindow(disp, winId(), x, y, w, h);
//    XFlush(disp);
//}

void MainWindow::positionChanged()
{
    // paly hide animation and disable other animation
    qDebug() << "start positionChange:" << frameGeometry();
    clearStrutPartial();

    //　需要在narrow之前执行，保证动画结束后能后更新界面布局的方向
    connect(m_panelHideAni, &QVariantAnimation::finished, this, &MainWindow::newPositionExpand);

    narrow();
}

void MainWindow::updatePosition()
{
    // all update operation need pass by timer
    //    Q_ASSERT(sender() == m_positionUpdateTimer);

    //clearStrutPartial();
    updateGeometry();
}

void MainWindow::updateGeometry()
{
    // DockDisplayMode and DockPosition MUST be set before invoke setFixedSize method of MainPanel

    //为了防止当后端发送错误值，然后发送正确值时，任务栏没有移动在相应的位置
    //当ｑｔ没有获取到屏幕资源时候，move函数会失效。可以直接return
    if (m_settings->primaryRect().width() == 0 || m_settings->primaryRect().height() == 0) {
        return;
    }

    setStrutPartial();

    m_mainPanel->setDisplayMode(m_settings->displayMode());
    m_mainPanel->setPositonValue(m_dockPosition);

    bool isHide = m_settings->hideState() == Hide && !testAttribute(Qt::WA_UnderMouse);

    const QRect windowRect = m_settings->windowRect(m_dockPosition, isHide);

    internalMove(windowRect.topLeft());

    if (!m_primaryScreenChanged || m_settings->hideState() != Hide) {
        QWidget::move(windowRect.topLeft());
        QWidget::setFixedSize(m_settings->m_mainWindowSize);
    }

    resizeMainPanelWindow();

    m_mainPanel->update();
}

void MainWindow::getTrayVisableItemCount()
{
    m_mainPanel->getTrayVisableItemCount();
}

void MainWindow::newPositionExpand()
{
    // set strut
    setStrutPartial();

    // reset to right environment when animation finished
    m_mainPanel->setPositonValue(m_dockPosition);
    resetPanelEnvironment();

    if ((Top == m_dockPosition) || (Bottom == m_dockPosition)) {
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    } else {
        m_dragWidget->setCursor(Qt::SizeHorCursor);
    }

    disconnect(m_panelHideAni, &QVariantAnimation::finished, this, &MainWindow::newPositionExpand);

    updatePanelVisible();
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
    //resetPanelEnvironment(true);

    if (m_settings->hideMode() != Dock::KeepShowing)
        return;

    const auto ratio = devicePixelRatioF();
    const int maxScreenHeight = m_settings->screenRawHeight();
    const int maxScreenWidth = m_settings->screenRawWidth();
    const QPoint &p = rawXPosition(m_settings->windowRect(m_dockPosition).topLeft());
    const QSize &s = m_settings->windowSize();
    const QRect &primaryRawRect = m_settings->currentRawRect();

    XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
    uint strut = 0;
    uint strutStart = 0;
    uint strutEnd = 0;

    QRect strutArea(0, 0, maxScreenWidth, maxScreenHeight);
    switch (m_dockPosition) {
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
    //    const QRect pr = m_settings->currentRawRect();
    //    for (auto *screen : qApp->screens()) {
    //        const QRect sr = screen->geometry();
    //        if (sr == pr)
    //            continue;

    //        if (sr.intersects(strutArea))
    //            ++count;
    //    }
    //    if (count > 0) {
    //        qWarning() << "strutArea is intersects with another screen.";
    //        qWarning() << maxScreenHeight << maxScreenWidth << m_dockPosition << p << s;
    //        return;
    //    }

    m_xcbMisc->set_strut_partial(winId(), orientation, strut + m_settings->dockMargin() * ratio, strutStart, strutEnd);
}

void MainWindow::expand()
{
    qDebug() << "expand started";
    if (m_panelHideAni->state() == QPropertyAnimation::Running) {
        m_panelHideAni->stop();
        emit m_panelHideAni->finished();
    }

    const auto showAniState = m_panelShowAni->state();

    int startValue = 0;
    int endValue = 0;

    resetPanelEnvironment();
    if (showAniState != QPropertyAnimation::Running && pos() != m_panelShowAni->currentValue()) {
        const QRect windowRect = m_settings->windowRect(m_dockPosition);

        startValue = (m_dockPosition == Top || m_dockPosition == Bottom) ? height() : width();
        endValue = (m_dockPosition == Top || m_dockPosition == Bottom) ? windowRect.height() : windowRect.width();

        qDebug() << "expand     " << "start value:" << startValue
                 << "end value:" << endValue;

        if (startValue > DOCK_MAX_SIZE || endValue > DOCK_MAX_SIZE) {
            return;
        }

        if (startValue > endValue)
            return;

        m_panelShowAni->setStartValue(startValue);
        m_panelShowAni->setEndValue(endValue);
        m_panelShowAni->start();
        qDebug() << "show ani start";
        m_shadowMaskOptimizeTimer->start();
        m_settings->posChangedUpdateSettings();
    }
}

void MainWindow::narrow()
{
    qDebug() << "narrow started";
    int startValue = (m_dockPosition == Top || m_dockPosition == Bottom) ? height() : width();

    qDebug() << "narrow     " << "start value:" << startValue;
    m_panelShowAni->stop();
    m_panelHideAni->setStartValue(startValue);
    m_panelHideAni->setEndValue(0);
    m_panelHideAni->start();
}

void MainWindow::resetPanelEnvironment()
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
        if (!m_registerKey.isEmpty()) {
            m_eventInter->UnregisterArea(m_registerKey);
            qDebug() << "register area clear";
            //清空registerKey
            m_registerKey.clear();
        }
        return expand();
    }

    if (m_registerKey.isEmpty()) {
        updateRegionMonitorWatch();
    }

    const Dock::HideState state = m_settings->hideState();

    if (state == Hide && m_settings->autoHide()) {
        QRectF r(pos(), size());
        const int margin = m_settings->dockMargin();
        switch (m_dockPosition) {
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
        if (!r.contains(QCursor::pos())) {
            qDebug() << "hide narrow";
            return narrow();
        }
    }

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
    // internalMove();
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
    QSize size = m_settings->windowSize();
    const QRect windowRect = m_settings->windowRect(m_dockPosition, false);
    internalMove(windowRect.topLeft());
    resizeMainPanelWindow();
    QWidget::setFixedSize(size);
}

void MainWindow::resizeMainPanelWindow()
{
    m_mainPanel->setFixedSize(m_settings->m_mainWindowSize);

    switch (m_dockPosition) {
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
        Q_UNREACHABLE();
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
    if (Dock::Top == m_dockPosition) {
        m_settings->m_mainWindowSize.setHeight(qBound(MAINWINDOW_MIN_SIZE, m_size.height() + offset.y(), MAINWINDOW_MAX_SIZE));
        m_settings->m_mainWindowSize.setWidth(width());
    } else if (Dock::Bottom == m_dockPosition) {
        m_settings->m_mainWindowSize.setHeight(qBound(MAINWINDOW_MIN_SIZE, m_size.height() - offset.y(), MAINWINDOW_MAX_SIZE));
        m_settings->m_mainWindowSize.setWidth(width());
    } else if (Dock::Left == m_dockPosition) {
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
        if (Dock::Top == m_dockPosition || Dock::Bottom == m_dockPosition) {
            m_settings->m_dockInter->setWindowSizeFashion(m_settings->m_mainWindowSize.height());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.height());
        } else {
            m_settings->m_dockInter->setWindowSizeFashion(m_settings->m_mainWindowSize.width());
            m_settings->m_dockInter->setWindowSize(m_settings->m_mainWindowSize.width());
        }
    } else {
        if (Dock::Top == m_dockPosition || Dock::Bottom == m_dockPosition) {
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

void MainWindow::onRegionMonitorChanged(int x, int y, const QString &key)
{
    if (m_registerKey != key)
        return;

    // 同一个坐标，只响应一次
    static QPoint lastPos(0, 0);
    if (lastPos == QPoint(x, y)) {
        return;
    }
    lastPos = QPoint(x, y);

    QScreen *screen = Utils::screenAt(QPoint(x, y));
    if (!screen)
        return;

    QRect screenRect = screen->geometry();
    qDebug() << y << screenRect.y() << screenRect.y() + screenRect.height() / 2;
    switch (m_dockPosition) {
    case Top:
        if (y > screenRect.y() + screenRect.height() / 2)
            return;
        break;
    case Bottom:
        if (y < screenRect.y() + screenRect.height() / 2)
            return;
        break;
    case Left:
        if (x > screenRect.x() + screenRect.width() / 2)
            return;
        break;
    case Right:
        if (x < screenRect.x() + screenRect.width() / 2)
            return;
    }

    if (screen->name() == m_settings->currentDockScreen()) {
        if (m_settings->hideMode() == KeepShowing)
            return;

        if (m_panelShowAni->state() == QPropertyAnimation::Running)
            return;

        // 一直隐藏模式不用通过时间延迟的方式调用,影响离开动画的响应
        expand();
    } else {
        // 移动Dock至相应屏相应位置
        if (m_launcherInter->IsVisible())//启动器显示,则dock不显示
            return;

        if (m_settings->setDockScreen(screen->name())) {
            if (m_settings->hideMode() == KeepShowing || m_settings->hideMode() == SmartHide) {
                narrow();
                newPositionExpand();
            } else {
                int screenWidth = screen->size().width();
                int screenHeight = screen->size().height();
                switch (m_dockPosition) {
                case Dock::Top:
                case Dock::Bottom:
                    setFixedWidth(screenWidth);
                    break;
                case Dock::Left:
                case Dock::Right:
                    setFixedHeight(screenHeight);
                    break;
                }
                expand();
            }
        }
    }
}

void MainWindow::updateRegionMonitorWatch()
{
    if (!m_registerKey.isEmpty()) {
        bool ret = m_eventInter->UnregisterArea(m_registerKey);
        qDebug() << "register area clear:" << ret;
        m_registerKey.clear();
    }

    const int flags = Motion | Button | Key;

    QList<QRect> screensRect = m_settings->monitorsRect();
    QList<MonitRect> monitorAreas;

    int val = 3;
    int x, y, w, h;

    auto func = [&](MonitRect & monitRect) {
        monitRect.x1 = x;
        monitRect.y1 = y;
        monitRect.x2 = x + w;
        monitRect.y2 = y + h;
        monitorAreas << monitRect;
    };

    if (screensRect.size()) {
        MonitRect monitRect;
        switch (m_dockPosition) {
        case Dock::Top: {
            for (QRect rect : screensRect) {
                x = rect.x();
                y = rect.y();
                w = rect.width();
                h = val;
                func(monitRect);
            }
        }
        break;
        case Dock::Bottom: {
            for (QRect rect : screensRect) {
                x = rect.x();
                y = rect.y() + rect.height() - val;
                w = rect.width();
                h = val;
                func(monitRect);
            }
        }
        break;
        case Dock::Left: {
            for (QRect rect : screensRect) {
                x = rect.x();
                y = rect.y();
                w = val;
                h = rect.height();
                func(monitRect);
            }
        }
        break;
        case Dock::Right: {
            for (QRect rect : screensRect) {
                x = rect.x() + rect.width() - val;
                y = rect.y();
                w = val;
                h = rect.height();
                func(monitRect);
            }
        }
        break;
        }
        m_registerKey = m_eventInter->RegisterAreas(monitorAreas, flags);
        qDebug() << "register key" << m_registerKey;
    } else {
        m_registerKey = m_eventInter->RegisterFullScreen();
        qDebug() << "register full screen" << m_registerKey;
    }
}

#include "mainwindow.moc"
