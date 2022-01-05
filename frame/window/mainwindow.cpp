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
#include "mainpanelcontrol.h"
#include "dockitemmanager.h"
#include "menuworker.h"

#include <DStyle>
#include <DPlatformWindowHandle>
#include <DSysInfo>
#include <DPlatformTheme>
#include <DDBusSender>

#include <QDebug>
#include <QEvent>
#include <QResizeEvent>
#include <QScreen>
#include <QGuiApplication>
#include <QX11Info>
#include <QtConcurrent>
#include <qpa/qplatformwindow.h>

#include <X11/X.h>
#include <X11/Xutil.h>

#include <com_deepin_dde_daemon_dock.h>

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"

#define MAINWINDOW_MAX_SIZE       DOCK_MAX_SIZE
#define MAINWINDOW_MIN_SIZE       (40)
#define DRAG_AREA_SIZE (5)

#define DRAG_STATE_PROP "DRAG_STATE"

using org::kde::StatusNotifierWatcher;

// let startdde know that we've already started.
void MainWindow::RegisterDdeSession()
{
    QString envName("DDE_SESSION_PROCESS_COOKIE_ID");

    QByteArray cookie = qgetenv(envName.toUtf8().data());
    qunsetenv(envName.toUtf8().data());

    if (!cookie.isEmpty()) {
        QDBusPendingReply<bool> r = DDBusSender()
                .interface("com.deepin.SessionManager")
                .path("/com/deepin/SessionManager")
                .service("com.deepin.SessionManager")
                .method("Register")
                .arg(QString(cookie))
                .call();

        qDebug() << Q_FUNC_INFO << r.value();
    }
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
    , m_updateDragAreaTimer(new QTimer(this))
{
    setAttribute(Qt::WA_TranslucentBackground);
    setAttribute(Qt::WA_X11DoNotAcceptFocus);

    Qt::WindowFlags flags = Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint | Qt::Window;
    //1 确保这两行代码的先后顺序，否则会导致任务栏界面不再置顶
    setWindowFlags(windowFlags() | flags | Qt::WindowDoesNotAcceptFocus);

    if (DGuiApplicationHelper::isXWindowPlatform()) {
        const auto display = QX11Info::display();
        if (!display) {
            qWarning() << "QX11Info::display() is " << display;
        } else {
            //2 确保这两行代码的先后顺序，否则会导致任务栏界面不再置顶
            XcbMisc::instance()->set_window_type(xcb_window_t(this->winId()), XcbMisc::Dock);
        }
    }

    setMouseTracking(true);
    setAcceptDrops(true);

    DPlatformWindowHandle::enableDXcbForWindow(this, true);
    m_platformWindowHandle.setEnableBlurWindow(true);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 5));
    m_platformWindowHandle.setShadowColor(QColor(0, 0, 0, 0.3 * 255));

    m_mainPanel->setDisplayMode(m_multiScreenWorker->displayMode());

    initMember();
    initSNIHost();
    initComponents();
    initConnections();

    resetDragWindow();

    for (auto item : DockItemManager::instance()->itemList())
        m_mainPanel->insertItem(-1, item);

    m_dragWidget->setMouseTracking(true);
    m_dragWidget->setFocusPolicy(Qt::NoFocus);

    if (!Utils::IS_WAYLAND_DISPLAY) {
        if ((Top == m_multiScreenWorker->position()) || (Bottom == m_multiScreenWorker->position())) {
            m_dragWidget->setCursor(Qt::SizeVerCursor);
        } else {
            m_dragWidget->setCursor(Qt::SizeHorCursor);
        }
    }
}

/**
 * @brief MainWindow::launch
 * 任务栏初次启动时调用此方法，里面是做了一些初始化操作
 */
void MainWindow::launch()
{
    if (!qApp->property("CANSHOW").toBool())
        return;

    m_launched = true;
    m_multiScreenWorker->initShow();
    m_shadowMaskOptimizeTimer->start();
    QTimer::singleShot(0, this, [ this ] { this->setVisible(true); });
}

/**
 * @brief MainWindow::callShow
 * 此方法是被外部进程通过DBus调用的。
 * @note 当任务栏以-r参数启动时，其不会显示界面，需要在外部通过DBus调用此接口之后才会显示界面，
 * 这里是以前为了优化任务栏的启动速度做的处理，当任务栏启动时，此时窗管进程可能还未启动完全，
 * 部分设置未初始化完等，导致任务栏显示的界面异常，所以留下此接口，被startdde延后调用
 */
void MainWindow::callShow()
{
    static bool flag = false;
    if (flag) {
        return;
    }
    flag = true;

    qApp->setProperty("CANSHOW", true);

    launch();

    // 预留200ms提供给窗口初始化再通知startdde，不影响启动速度
    QTimer::singleShot(200, this, &MainWindow::RegisterDdeSession);
}

/**
 * @brief MainWindow::relaodPlugins
 * 需要重新加载插件时，此接口会被调用，目前是用于任务栏的安全模式退出时调用
 */
void MainWindow::relaodPlugins()
{
    if (qApp->property("PLUGINSLOADED").toBool()) {
        return;
    }

    DockItemManager::instance()->startLoadPlugins();
    qApp->setProperty("PLUGINSLOADED", true);
}

/**
 * @brief MainWindow::mousePressEvent
 * @param e
 * @note 右键显示任务栏的菜单
 */
void MainWindow::mousePressEvent(QMouseEvent *e)
{
    e->ignore();
    if (e->button() == Qt::RightButton) {
        QTimer::singleShot(10, this, [this]{
            m_menuWorker->showDockSettingsMenu();
        });
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

    if (!qApp->property(DRAG_STATE_PROP).toBool())
        m_updateDragAreaTimer->start();
}

void MainWindow::resizeEvent(QResizeEvent *event)
{
    if (!qApp->property(DRAG_STATE_PROP).toBool())
        m_updateDragAreaTimer->start();

    // 任务栏大小、位置、模式改变都会触发resize，发射大小改变信号，供依赖项目更新位置
    Q_EMIT panelGeometryChanged();

    m_mainPanel->updatePluginsLayout();
    m_shadowMaskOptimizeTimer->start();

    return DBlurEffectWidget::resizeEvent(event);
}

void MainWindow::dragEnterEvent(QDragEnterEvent *e)
{
    QWidget::dragEnterEvent(e);
}

void MainWindow::initMember()
{
    //INFO 这里要大于动画的300ms，否则可能动画过程中这个定时器就被触发了
    m_updateDragAreaTimer->setInterval(500);
    m_updateDragAreaTimer->setSingleShot(true);
}

/**
 * @brief MainWindow::initSNIHost
 * @note 将Dock注册到StatusNotifierWatcher服务上
 */
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

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::frameMarginsChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(&m_platformWindowHandle, &DPlatformWindowHandle::windowRadiusChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &MainWindow::onDbusNameOwnerChanged);

    connect(DockItemManager::instance(), &DockItemManager::itemInserted, m_mainPanel, &MainPanelControl::insertItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemRemoved, m_mainPanel, &MainPanelControl::removeItem, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::itemUpdated, m_mainPanel, &MainPanelControl::itemUpdated, Qt::DirectConnection);
    connect(DockItemManager::instance(), &DockItemManager::trayVisableCountChanged, this, &MainWindow::resizeDockIcon, Qt::QueuedConnection);
    connect(DockItemManager::instance(), &DockItemManager::requestWindowAutoHide, m_menuWorker, &MenuWorker::setAutoHide);
    connect(m_mainPanel, &MainPanelControl::itemMoved, DockItemManager::instance(), &DockItemManager::itemMoved, Qt::DirectConnection);
    connect(m_mainPanel, &MainPanelControl::itemAdded, DockItemManager::instance(), &DockItemManager::itemAdded, Qt::DirectConnection);

    // -拖拽任务栏改变高度或宽度-------------------------------------------------------------------------------
    connect(m_updateDragAreaTimer, &QTimer::timeout, this, &MainWindow::resetDragWindow);
    //TODO 后端考虑删除这块，目前还不能删除，调整任务栏高度的时候，任务栏外部区域有变化
    connect(m_updateDragAreaTimer, &QTimer::timeout, m_multiScreenWorker, &MultiScreenWorker::onRequestUpdateRegionMonitor);

    connect(m_dragWidget, &DragWidget::dragPointOffset, this, [ = ] { qApp->setProperty(DRAG_STATE_PROP, true); });
    connect(m_dragWidget, &DragWidget::dragFinished, this, [ = ] { qApp->setProperty(DRAG_STATE_PROP, false); });

    connect(m_dragWidget, &DragWidget::dragPointOffset, this, &MainWindow::onMainWindowSizeChanged);
    connect(m_dragWidget, &DragWidget::dragFinished, this, &MainWindow::resetDragWindow);   //　更新拖拽区域
    // ----------------------------------------------------------------------------------------------------

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &MainWindow::themeTypeChanged);

    connect(m_menuWorker, &MenuWorker::autoHideChanged, m_multiScreenWorker, &MultiScreenWorker::onAutoHideChanged);

    connect(m_multiScreenWorker, &MultiScreenWorker::opacityChanged, this, &MainWindow::setMaskAlpha, Qt::QueuedConnection);
    connect(m_multiScreenWorker, &MultiScreenWorker::displayModeChanegd, this, &MainWindow::adjustShadowMask, Qt::QueuedConnection);

    connect(m_multiScreenWorker, &MultiScreenWorker::requestUpdateDockEntry, DockItemManager::instance(), &DockItemManager::requestUpdateDockItem);

    // 响应后端触控屏拖拽任务栏高度长按信号
    connect(TouchSignalManager::instance(), &TouchSignalManager::middleTouchPress, this, &MainWindow::touchRequestResizeDock);
    connect(TouchSignalManager::instance(), &TouchSignalManager::touchMove, m_dragWidget, &DragWidget::onTouchMove);
}

/**
 * @brief MainWindow::getTrayVisableItemCount
 * 重新获取以下当前托盘区域有多少个可见的图标，并更新图标的大小
 */
void MainWindow::resizeDockIcon()
{
    m_mainPanel->resizeDockIcon();
}

/**
 * @brief MainWindow::adjustShadowMask 更新任务栏的圆角大小（时尚模式下才有圆角效果）
 */
void MainWindow::adjustShadowMask()
{
    if (!m_launched || m_shadowMaskOptimizeTimer->isActive())
        return;

    DStyleHelper dstyle(style());
    int radius = 0;
    if (m_wmHelper->hasComposite() && m_multiScreenWorker->displayMode() == DisplayMode::Fashion) {
        if (Dtk::Core::DSysInfo::isCommunityEdition()) { // 社区版圆角与专业版不同
            DPlatformTheme *theme = DGuiApplicationHelper::instance()->systemTheme();
            radius = theme->windowRadius(radius);
        } else {
            radius = dstyle.pixelMetric(DStyle::PM_TopLevelWindowRadius);
        }
    }

    m_platformWindowHandle.setWindowRadius(radius);
    m_mainPanel->updatePluginsLayout();
}

void MainWindow::onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (name == SNI_WATCHER_SERVICE && !newOwner.isEmpty()) {
        qDebug() << SNI_WATCHER_SERVICE << "SNI watcher daemon started, register dock to watcher as SNI Host";
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    }
}

/**
 * @brief MainWindow::setEffectEnabled
 * @param enabled 根据当前系统是否enabled特效来更新任务栏的外观样式
 */
void MainWindow::setEffectEnabled(const bool enabled)
{
    setMaskColor(AutoColor);

    setMaskAlpha(m_multiScreenWorker->opacity());

    m_platformWindowHandle.setBorderWidth(enabled ? 1 : 0);
}

/**
 * @brief MainWindow::setComposite
 * @param hasComposite 系统是否支持混成（也就是特效）
 */
void MainWindow::setComposite(const bool hasComposite)
{
    setEffectEnabled(hasComposite);
}

/**
 * @brief MainWindow::resetDragWindow 更新任务栏的拖拽区域
 * @note 任务栏远离屏幕的一边是支持拖拽的，由一个不可见的widget提拽支持，当任务栏的geometry发生变化的时候，此拖拽区域也需要更新其自身的geometry
 */
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

void MainWindow::resizeDock(int offset, bool dragging)
{
    qApp->setProperty(DRAG_STATE_PROP, dragging);

    const QRect &rect = m_multiScreenWorker->getDockShowMinGeometry(m_multiScreenWorker->deskScreen());
    QRect newRect;
    switch (m_multiScreenWorker->position()) {
    case Top: {
       newRect.setX(rect.x());
       newRect.setY(rect.y());
       newRect.setWidth(rect.width());
       newRect.setHeight(qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
   }
        break;
    case Bottom: {
        newRect.setX(rect.x());
        newRect.setY(rect.y() + rect.height() - qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
        newRect.setWidth(rect.width());
        newRect.setHeight(qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
    }
        break;
    case Left: {
        newRect.setX(rect.x());
        newRect.setY(rect.y());
        newRect.setWidth(qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
        break;
    case Right: {
        newRect.setX(rect.x() + rect.width() - qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
        newRect.setY(rect.y());
        newRect.setWidth(qBound(MAINWINDOW_MIN_SIZE, offset, MAINWINDOW_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
        break;
    }

    // 更新界面大小
    m_mainPanel->setFixedSize(newRect.size());
    setFixedSize(newRect.size());
    move(newRect.topLeft());

    if (!dragging)
        resetDragWindow();
}

/**
 * @brief MainWindow::onMainWindowSizeChanged 任务栏拖拽过程中会不听调用此方法更新自身大小
 * @param offset 拖拽时的坐标偏移量
 */
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

/**
 * @brief MainWindow::themeTypeChanged 系统主题发生变化时，此方法被调用
 * @param themeType 当前系统主题
 */
void MainWindow::themeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    if (m_wmHelper->hasComposite()) {
        if (themeType == DGuiApplicationHelper::DarkType)
            m_platformWindowHandle.setBorderColor(QColor(0, 0, 0, 255 * 0.3));
        else
            m_platformWindowHandle.setBorderColor(QColor(QColor::Invalid));
    }
}

/**
 * @brief MainWindow::touchRequestResizeDock 触屏情况用手指调整任务栏高度或宽度
 */
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

    int resizeHeight = Utils::SettingValue("com.deepin.dde.dock.touch", QByteArray(), "resizeHeight", 7).toInt();

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

/**
 * @brief MainWindow::setGeometry
 * @param rect 设置任务栏的位置和大小，重写此函数时为了及时发出panelGeometryChanged信号，最终供外部DBus调用方使用
 */
void MainWindow::setGeometry(const QRect &rect)
{
    if (rect == this->geometry()) {
        return;
    }
    DBlurEffectWidget::setGeometry(rect);
    emit panelGeometryChanged();
}

/**
 * @brief 当进入安全模式时，通过此方法发送通知告知用户
 */
void MainWindow::sendNotifications()
{
    QStringList actionButton;
    actionButton << "reload" << tr("Exit Safe Mode");
    QVariantMap hints;
    hints["x-deepin-action-reload"] = QString("dbus-send,--session,--dest=com.deepin.dde.Dock,--print-reply,/com/deepin/dde/Dock,com.deepin.dde.Dock.ReloadPlugins");
    // 在进入安全模式时，执行此DBUS耗时25S左右，导致任务栏显示阻塞，所以使用线程调用
    QtConcurrent::run(QThreadPool::globalInstance(), [=] {
        DDBusSender()
                .service("com.deepin.dde.Notification")
                .path("/com/deepin/dde/Notification")
                .interface("com.deepin.dde.Notification")
                .method(QString("Notify"))
                .arg(QString("dde-control-center"))                                            // appname
                .arg(static_cast<uint>(0))                                                     // id
                .arg(QString("preferences-system"))                                            // icon
                .arg(QString(tr("Dock - Safe Mode")))                                          // summary
                .arg(QString(tr("The Dock is in safe mode, please exit to show it properly"))) // content
                .arg(actionButton)                                                             // actions
                .arg(hints)                                                                    // hints
                .arg(15000)                                                                    // timeout
                .call();
    });
}
