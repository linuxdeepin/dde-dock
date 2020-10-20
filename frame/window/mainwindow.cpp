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
#include "util/menuworker.h"

#include <DStyle>
#include <DPlatformWindowHandle>
#include <DSysInfo>
#include <DPlatformTheme>

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <qpa/qplatformwindow.h>
#include <QGSettings>

#include <X11/X.h>
#include <X11/Xutil.h>

#include <com_deepin_dde_daemon_dock.h>

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"

#define MAINWINDOW_MAX_SIZE       DOCK_MAX_SIZE
#define MAINWINDOW_MIN_SIZE       (40)
#define DRAG_AREA_SIZE (5)

using org::kde::StatusNotifierWatcher;
using DBusDock = com::deepin::dde::daemon::Dock;

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
    , m_mainPanel(new MainPanelControl(this))
    , m_platformWindowHandle(this)
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_multiScreenWorker(new MultiScreenWorker(this, m_wmHelper))
    , m_menuWorker(new MenuWorker(m_multiScreenWorker->dockInter(), this))
    , m_shadowMaskOptimizeTimer(new QTimer(this))
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    , m_sniWatcher(new StatusNotifierWatcher(SNI_WATCHER_SERVICE, SNI_WATCHER_PATH, QDBusConnection::sessionBus(), this))
    , m_dragWidget(new DragWidget(this))
    , m_launched(false)
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

    m_mainPanel->setDisplayMode(m_multiScreenWorker->displayMode());

    initSNIHost();
    initComponents();
    initConnections();

    resetDragWindow();

    m_mainPanel->setDelegate(this);
    for (auto item : DockItemManager::instance()->itemList())
        m_mainPanel->insertItem(-1, item);

    m_dragWidget->setMouseTracking(true);
    m_dragWidget->setFocusPolicy(Qt::NoFocus);

    if ((Top == m_multiScreenWorker->position()) || (Bottom == m_multiScreenWorker->position())) {
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    } else {
        m_dragWidget->setCursor(Qt::SizeHorCursor);
    }
}

MainWindow::~MainWindow()
{

}

void MainWindow::launch()
{
    if (!qApp->property("CANSHOW").toBool())
        return;

    m_launched = true;
    qApp->processEvents();
    setVisible(true);
    m_multiScreenWorker->initShow();
    m_shadowMaskOptimizeTimer->start();
}

void MainWindow::callShow()
{
    static bool flag = false;
    if (flag) {
        return;
    }
    flag = true;

    qApp->setProperty("CANSHOW", true);

    launch();

    m_callShowEventLoop.exec();
}

void MainWindow::showEvent(QShowEvent *e)
{
    QWidget::showEvent(e);

    if (m_callShowEventLoop.isRunning()) {
        m_callShowEventLoop.quit();
    }

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
    if (e->button() == Qt::RightButton && m_menuWorker->menuEnable()) {
        m_menuWorker->showDockSettingsMenu();
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
    Q_UNUSED(e);
    //重写mouseMoveEvent 解决bug12866  leaveEvent事件失效
}

void MainWindow::moveEvent(QMoveEvent *event)
{
    Q_UNUSED(event);
}

void MainWindow::resizeEvent(QResizeEvent *event)
{
    // 任务栏大小、位置、模式改变都会触发resize，发射大小改变信号，供依赖项目更新位置
    Q_EMIT panelGeometryChanged();

    adjustShadowMask();

    return DBlurEffectWidget::resizeEvent(event);
}

void MainWindow::dragEnterEvent(QDragEnterEvent *e)
{
    QWidget::dragEnterEvent(e);
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
    m_shadowMaskOptimizeTimer->setSingleShot(true);
    m_shadowMaskOptimizeTimer->setInterval(100);

    QTimer::singleShot(1, this, &MainWindow::compositeChanged);

    themeTypeChanged(DGuiApplicationHelper::instance()->themeType());
}

void MainWindow::compositeChanged()
{
    const bool composite = m_wmHelper->hasComposite();
    setComposite(composite);

    m_shadowMaskOptimizeTimer->start();
}

void MainWindow::initConnections()
{
    connect(m_shadowMaskOptimizeTimer, &QTimer::timeout, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &MainWindow::compositeChanged, Qt::QueuedConnection);
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &MainWindow::onDbusNameOwnerChanged);

    connect(DockItemManager::instance(), &DockItemManager::itemInserted, m_mainPanel, &MainPanelControl::insertItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemRemoved, m_mainPanel, &MainPanelControl::removeItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemUpdated, m_mainPanel, &MainPanelControl::itemUpdated, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::requestWindowAutoHide, m_menuWorker, &MenuWorker::setAutoHide);
    connect(m_mainPanel, &MainPanelControl::itemMoved, DockItemManager::instance(), &DockItemManager::itemMoved, Qt::DirectConnection);
    connect(m_mainPanel, &MainPanelControl::itemAdded, DockItemManager::instance(), &DockItemManager::itemAdded, Qt::DirectConnection);

    connect(m_dragWidget, &DragWidget::dragPointOffset, m_multiScreenWorker, [ = ] {m_multiScreenWorker->onDragStateChanged(true);});
    connect(m_dragWidget, &DragWidget::dragFinished, m_multiScreenWorker, [ = ] {m_multiScreenWorker->onDragStateChanged(false);});

    connect(m_dragWidget, &DragWidget::dragPointOffset, this, &MainWindow::onMainWindowSizeChanged);
    connect(m_dragWidget, &DragWidget::dragFinished, this, &MainWindow::onDragFinished);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &MainWindow::themeTypeChanged);

    connect(m_menuWorker, &MenuWorker::trayCountChanged, this, &MainWindow::getTrayVisableItemCount, Qt::DirectConnection);
    connect(m_menuWorker, &MenuWorker::autoHideChanged, m_multiScreenWorker, &MultiScreenWorker::onAutoHideChanged);

    connect(m_multiScreenWorker, &MultiScreenWorker::opacityChanged, this, &MainWindow::setMaskAlpha, Qt::QueuedConnection);
    connect(m_multiScreenWorker, &MultiScreenWorker::displayModeChanegd, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    //　更新拖拽区域
    connect(m_multiScreenWorker, &MultiScreenWorker::requestUpdateDragArea, this, &MainWindow::resetDragWindow);

    // 响应后端触控屏拖拽任务栏高度长按信号
    connect(TouchSignalManager::instance(), &TouchSignalManager::middleTouchPress, this, &MainWindow::touchRequestResizeDock);
    connect(TouchSignalManager::instance(), &TouchSignalManager::touchMove, m_dragWidget, [ this ]() {
        static QPoint lastPos;
        QPoint curPos = QCursor::pos();
        if (lastPos == curPos) {
            return;
        }
        lastPos = curPos;
        qApp->postEvent(m_dragWidget, new QMouseEvent(QEvent::MouseMove, m_dragWidget->mapFromGlobal(curPos)
                                                      , QPoint(), curPos, Qt::LeftButton, Qt::LeftButton
                                                      , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
    });
}

void MainWindow::getTrayVisableItemCount()
{
    m_mainPanel->getTrayVisableItemCount();
}

void MainWindow::adjustShadowMask()
{
    if (!m_launched)
        return;

    if (m_shadowMaskOptimizeTimer->isActive())
        return;

    const bool composite = m_wmHelper->hasComposite();
    const bool isFasion = m_multiScreenWorker->displayMode() == Fashion;

    DStyleHelper dstyle(style());
    int radius = dstyle.pixelMetric(DStyle::PM_TopLevelWindowRadius);

    if (Dtk::Core::DSysInfo::isCommunityEdition()) {
        auto theme = DGuiApplicationHelper::instance()->systemTheme();
        radius = theme->windowRadius();
    }

    int newRadius = composite && isFasion ? radius : 0;
    m_platformWindowHandle.setWindowRadius(newRadius);

    QPainterPath clipPath;
    clipPath.addRect(QRect(QPoint(0, 0), this->geometry().size()));

    m_platformWindowHandle.setClipPath(newRadius != 0 ? QPainterPath() : clipPath);

    // 临时解决方案: 在控制中心修改圆角之后, 会过几秒才会刷新显示 ???
    if (Dtk::Core::DSysInfo::isCommunityEdition())
        update();
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

    setMaskAlpha(m_multiScreenWorker->opacity());

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

void MainWindow::resetDragWindow()
{
    switch (m_multiScreenWorker->position()) {
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
    }

    QRect rect = m_multiScreenWorker->dockRect(m_multiScreenWorker->deskScreen()
                                               , m_multiScreenWorker->position()
                                               , HideMode::KeepShowing
                                               , m_multiScreenWorker->displayMode());

    // 这个时候屏幕有可能是隐藏的，不能直接使用this->width()这种去设置任务栏的高度，而应该保证原值
    int dockSize = 0;
    if (m_multiScreenWorker->position() == Position::Left
            || m_multiScreenWorker->position() == Position::Right) {
        dockSize = this->width() == 0 ? rect.width() : this->width();
    } else {
        dockSize = this->height() == 0 ? rect.height() : this->height();
    }

    /** FIX ME
     * 作用：限制dockSize的值在40～100之间。
     * 问题1：如果dockSize为39，会导致dock的mainwindow高度变成99，显示的内容高度却是39。
     * 问题2：dockSize的值在这里不应该为39，但在高分屏上开启缩放后，拉高任务栏操作会概率出现。
     * 暂时未分析出原因，后面再修改。
     */
    dockSize = qBound(MAINWINDOW_MIN_SIZE, dockSize, MAINWINDOW_MAX_SIZE);

    // 通知窗管和后端更新数据
    m_multiScreenWorker->updateDaemonDockSize(dockSize);                                // 1.先更新任务栏高度
    m_multiScreenWorker->requestUpdateFrontendGeometry();                               // 2.再更新任务栏位置,保证先1再2
    m_multiScreenWorker->requestNotifyWindowManager();
    m_multiScreenWorker->requestUpdateRegionMonitor();                                  // 界面发生变化，应更新监控区域

    if ((Top == m_multiScreenWorker->position()) || (Bottom == m_multiScreenWorker->position())) {
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    } else {
        m_dragWidget->setCursor(Qt::SizeHorCursor);
    }
}

void MainWindow::onMainWindowSizeChanged(QPoint offset)
{
    const QRect &rect = m_multiScreenWorker->dockRect(m_multiScreenWorker->deskScreen()
                                                      , m_multiScreenWorker->position()
                                                      , HideMode::KeepShowing,
                                                      m_multiScreenWorker->displayMode());
    QRect newRect;
    switch (m_multiScreenWorker->position()) {
    case Top: {
        newRect.setX(rect.x());
        newRect.setY(rect.y());
        newRect.setWidth(rect.width());
        newRect.setHeight(qBound(MAINWINDOW_MIN_SIZE, rect.height() + offset.y(), MAINWINDOW_MAX_SIZE));
    }
    break;
    case Bottom: {
        newRect.setX(rect.x());
        newRect.setY(rect.y() + rect.height() - qBound(MAINWINDOW_MIN_SIZE, rect.height() - offset.y(), MAINWINDOW_MAX_SIZE));
        newRect.setWidth(rect.width());
        newRect.setHeight(qBound(MAINWINDOW_MIN_SIZE, rect.height() - offset.y(), MAINWINDOW_MAX_SIZE));
    }
    break;
    case Left: {
        newRect.setX(rect.x());
        newRect.setY(rect.y());
        newRect.setWidth(qBound(MAINWINDOW_MIN_SIZE, rect.width() + offset.x(), MAINWINDOW_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
    break;
    case Right: {
        newRect.setX(rect.x() + rect.width() - qBound(MAINWINDOW_MIN_SIZE, rect.width() - offset.x(), MAINWINDOW_MAX_SIZE));
        newRect.setY(rect.y());
        newRect.setWidth(qBound(MAINWINDOW_MIN_SIZE, rect.width() - offset.x(), MAINWINDOW_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
    break;
    }

    // 更新界面大小
    m_mainPanel->setFixedSize(newRect.size());
    setFixedSize(newRect.size());
    move(newRect.topLeft());
}

void MainWindow::onDragFinished()
{
    qDebug() << "drag finished";
    resetDragWindow();
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

void MainWindow::touchRequestResizeDock()
{
    const QPoint touchPos(QCursor::pos());
    QRect dockRect = m_multiScreenWorker->dockRect(m_multiScreenWorker->deskScreen()
                                                   , m_multiScreenWorker->position()
                                                   , HideMode::KeepShowing
                                                   , m_multiScreenWorker->displayMode());

    // 隐藏状态返回
    if (width() == 0 || height() == 0) {
        return;
    }

    QGSettings settings("com.deepin.dde.dock.touch", QByteArray(), this);
    int resizeHeight = settings.get("resizeHeight").toInt();

    QRect touchRect;
    // 任务栏屏幕 内侧边线 内外resizeHeight距离矩形区域内长按可拖动任务栏高度
    switch (m_multiScreenWorker->position()) {
    case Position::Top:
        touchRect = QRect(dockRect.x(), dockRect.y() + dockRect.height() - resizeHeight, dockRect.width(), resizeHeight * 2);
        break;
    case Position::Bottom:
        touchRect = QRect(dockRect.x(), dockRect.y() - resizeHeight, dockRect.width(), resizeHeight * 2);
        break;
    case Position::Left:
        touchRect = QRect(dockRect.x() + dockRect.width() - resizeHeight, dockRect.y(), resizeHeight * 2, dockRect.height());
        break;
    case Position::Right:
        touchRect = QRect(dockRect.x() - resizeHeight, dockRect.y(), resizeHeight * 2, dockRect.height());
        break;
    }

    if (!touchRect.contains(touchPos)) {
        return;
    }
    qApp->postEvent(m_dragWidget, new QMouseEvent(QEvent::MouseButtonPress, m_dragWidget->mapFromGlobal(touchPos)
                                                  , QPoint(), touchPos, Qt::LeftButton, Qt::NoButton
                                                  , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
}

void MainWindow::setGeometry(const QRect &rect)
{
    static QRect lastRect;
    if (lastRect == rect) {
        return;
    }
    lastRect = rect;
    DBlurEffectWidget::setGeometry(rect);
    emit panelGeometryChanged();
}

#include "mainwindow.moc"
