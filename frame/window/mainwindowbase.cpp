// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainwindowbase.h"
#include "dragwidget.h"
#include "multiscreenworker.h"
#include "dockscreen.h"
#include "touchsignalmanager.h"
#include "displaymanager.h"
#include "menuworker.h"
#include "docksettings.h"

#include <DStyle>
#include <DWindowManagerHelper>
#include <DSysInfo>
#include <DPlatformTheme>

#include <QScreen>
#include <QX11Info>
#include <qpa/qplatformscreen.h>
#include <qpa/qplatformnativeinterface.h>

#define DRAG_AREA_SIZE (5)

// 任务栏圆角最小的时候，任务栏的高度值
#define MIN_RADIUS_WINDOWSIZE 46
// 任务栏圆角最小值和最大值的差值
#define MAX_MIN_RADIUS_DIFFVALUE 6
// 最小圆角值
#define MIN_RADIUS 12

#define DOCK_SCREEN DockScreen::instance()
#define DIS_INS DisplayManager::instance()

DGUI_USE_NAMESPACE

MainWindowBase::MainWindowBase(MultiScreenWorker *multiScreenWorker, QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_platformWindowHandle(this)
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_position(Dock::Position::Bottom)
    , m_dragWidget(new DragWidget(this))
    , m_multiScreenWorker(multiScreenWorker)
    , m_updateDragAreaTimer(new QTimer(this))
    , m_shadowMaskOptimizeTimer(new QTimer(this))
    , m_isShow(false)
    , m_order(0)
{
    initUi();
    initAttribute();
    initConnection();
    initMember();
}

MainWindowBase::~MainWindowBase()
{
}

void MainWindowBase::setOrder(int order)
{
    m_order = order;
}

int MainWindowBase::order() const
{
    return m_order;
}

void MainWindowBase::initAttribute()
{
    setAttribute(Qt::WA_TranslucentBackground);
    setAttribute(Qt::WA_X11DoNotAcceptFocus);

    Qt::WindowFlags flags = Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint | Qt::Window;
    //1 确保这两行代码的先后顺序，否则会导致任务栏界面不再置顶
    setWindowFlags(windowFlags() | flags | Qt::WindowDoesNotAcceptFocus);

    if (Utils::IS_WAYLAND_DISPLAY) {
        setWindowFlag(Qt::FramelessWindowHint, false); // 会导致设置圆角为0时无效
        setAttribute(Qt::WA_NativeWindow);
        windowHandle()->setProperty("_d_dwayland_window-type", "dock");
    }

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

    m_dragWidget->setMouseTracking(true);
    m_dragWidget->setFocusPolicy(Qt::NoFocus);

    if ((Dock::Top == m_position) || (Dock::Bottom == m_position))
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    else
        m_dragWidget->setCursor(Qt::SizeHorCursor);
}

void MainWindowBase::initConnection()
{
    connect(DWindowManagerHelper::instance(), &DWindowManagerHelper::hasCompositeChanged, m_shadowMaskOptimizeTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_shadowMaskOptimizeTimer, &QTimer::timeout, this, &MainWindowBase::adjustShadowMask, Qt::QueuedConnection);

    connect(m_dragWidget, &DragWidget::dragFinished, this, [ = ] {
        Utils::setIsDraging(false);
    });

    // -拖拽任务栏改变高度或宽度-------------------------------------------------------------------------------
    connect(m_updateDragAreaTimer, &QTimer::timeout, this, &MainWindowBase::resetDragWindow);
    //TODO 后端考虑删除这块，目前还不能删除，调整任务栏高度的时候，任务栏外部区域有变化
    connect(m_updateDragAreaTimer, &QTimer::timeout, m_multiScreenWorker, &MultiScreenWorker::onRequestUpdateRegionMonitor);

    connect(m_dragWidget, &DragWidget::dragPointOffset, this, &MainWindowBase::onMainWindowSizeChanged);
    connect(m_dragWidget, &DragWidget::dragFinished, this, &MainWindowBase::resetDragWindow);   //　更新拖拽区域
    connect(TouchSignalManager::instance(), &TouchSignalManager::touchMove, m_dragWidget, &DragWidget::onTouchMove);
    connect(TouchSignalManager::instance(), &TouchSignalManager::middleTouchPress, this, &MainWindowBase::touchRequestResizeDock);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &MainWindowBase::onThemeTypeChanged);
    connect(m_multiScreenWorker, &MultiScreenWorker::opacityChanged, this, &MainWindowBase::setMaskAlpha, Qt::QueuedConnection);

    onThemeTypeChanged(DGuiApplicationHelper::instance()->themeType());
    QMetaObject::invokeMethod(this, &MainWindowBase::onCompositeChanged);
}

void MainWindowBase::initMember()
{
    //INFO 这里要大于动画的300ms，否则可能动画过程中这个定时器就被触发了
    m_updateDragAreaTimer->setInterval(500);
    m_updateDragAreaTimer->setSingleShot(true);
    m_shadowMaskOptimizeTimer->setSingleShot(true);
    m_shadowMaskOptimizeTimer->setInterval(100);
}

int MainWindowBase::getBorderRadius() const
{
    if (!DWindowManagerHelper::instance()->hasComposite() || m_multiScreenWorker->displayMode() != DisplayMode::Fashion)
        return 0;

    int size = ((m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) ? height() : width());
    return qMin(MAX_MIN_RADIUS_DIFFVALUE, qMax(size - MIN_RADIUS_WINDOWSIZE, 0)) + MIN_RADIUS;
}

QRect MainWindowBase::getAnimationRect(const QRect &sourceRect, const Dock::Position &pos) const
{
    if (!Utils::IS_WAYLAND_DISPLAY
            || m_multiScreenWorker->hideMode() == HideMode::KeepShowing
            || m_multiScreenWorker->displayMode() == Dock::DisplayMode::Fashion)
        return sourceRect;

    // 在wayland环境下，如果任务栏状态为智能隐藏或者一直隐藏，那么在高效模式下，任务栏距离边缘的距离如果为0
    // 会导致在WindowManager类中设置窗管的_d_dwayland_dockstrut属性失效，因此，此处将动画位置距离边缘设置为1
    // 此时就不会出现_d_dwayland_dockstrut属性失效的情况（1个像素并不影响动画效果）
    // 在时尚模式下无需做这个设置，因为时尚模式下距离边缘的距离为10
    QRect animationRect = sourceRect;
    switch (pos) {
    case Dock::Position::Bottom: {
        animationRect.setTop(animationRect.top() - 1);
        animationRect.setHeight(sourceRect.height());
        break;
    }
    case Dock::Position::Left: {
        animationRect.setLeft(1);
        animationRect.setWidth(sourceRect.width());
        break;
    }
    case Dock::Position::Top: {
        animationRect.setTop(1);
        animationRect.setHeight(sourceRect.height());
        break;
    }
    case Dock::Position::Right: {
        animationRect.setLeft(animationRect.left() - 1);
        animationRect.setWidth(sourceRect.width());
        break;
    }
    }
    return animationRect;
}

/**
 * @brief MainWindow::onMainWindowSizeChanged 任务栏拖拽过程中会不停调用此方法更新自身大小
 * @param offset 拖拽时的坐标偏移量
 */
void MainWindowBase::onMainWindowSizeChanged(QPoint offset)
{
    QScreen *screen = DIS_INS->screen(DOCK_SCREEN->current());
    if (!screen)
        return;

    const QRect rect = getDockGeometry(screen, position(), displayMode(), Dock::HideState::Show);
    QRect newRect;
    switch (m_multiScreenWorker->position()) {
    case Top: {
        newRect.setX(rect.x());
        newRect.setY(rect.y());
        newRect.setWidth(rect.width());
        newRect.setHeight(qBound(DOCK_MIN_SIZE, rect.height() + offset.y(), DOCK_MAX_SIZE));
    }
        break;
    case Bottom: {
        newRect.setX(rect.x());
        newRect.setY(rect.y() + rect.height() - qBound(DOCK_MIN_SIZE, rect.height() - offset.y(), DOCK_MAX_SIZE));
        newRect.setWidth(rect.width());
        newRect.setHeight(qBound(DOCK_MIN_SIZE, rect.height() - offset.y(), DOCK_MAX_SIZE));
    }
        break;
    case Left: {
        newRect.setX(rect.x());
        newRect.setY(rect.y());
        newRect.setWidth(qBound(DOCK_MIN_SIZE, rect.width() + offset.x(), DOCK_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
        break;
    case Right: {
        newRect.setX(rect.x() + rect.width() - qBound(DOCK_MIN_SIZE, rect.width() - offset.x(), DOCK_MAX_SIZE));
        newRect.setY(rect.y());
        newRect.setWidth(qBound(DOCK_MIN_SIZE, rect.width() - offset.x(), DOCK_MAX_SIZE));
        newRect.setHeight(rect.height());
    }
        break;
    }

    Utils::setIsDraging(true);

    setFixedSize(newRect.size());
    move(newRect.topLeft());
    resetPanelGeometry();

    Q_EMIT requestUpdate();
}

void MainWindowBase::updateDragGeometry()
{
    switch (position()) {
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

    m_dragWidget->raise();
    if ((Top == position()) || (Bottom == position())) {
        m_dragWidget->setCursor(Qt::SizeVerCursor);
    } else {
        m_dragWidget->setCursor(Qt::SizeHorCursor);
    }
}

void MainWindowBase::resetDragWindow()
{
    updateDragGeometry();
    QScreen *screen = DIS_INS->screen(DOCK_SCREEN->current());
    if (!screen)
        return;

    QRect currentRect = getDockGeometry(screen, position(), displayMode(), Dock::HideState::Show);

    // 这个时候屏幕有可能是隐藏的，不能直接使用this->width()这种去设置任务栏的高度，而应该保证原值
    int dockSize = 0;
    if (m_multiScreenWorker->position() == Position::Left
            || m_multiScreenWorker->position() == Position::Right) {
        dockSize = this->width() == 0 ? currentRect.width() : this->width();
    } else {
        dockSize = this->height() == 0 ? currentRect.height() : this->height();
    }

    /** FIX ME
     * 作用：限制dockSize的值在40～100之间。
     * 问题1：如果dockSize为39，会导致dock的mainwindow高度变成99，显示的内容高度却是39。
     * 问题2：dockSize的值在这里不应该为39，但在高分屏上开启缩放后，拉高任务栏操作会概率出现。
     * 暂时未分析出原因，后面再修改。
     */
    dockSize = qBound(DOCK_MIN_SIZE, dockSize, DOCK_MAX_SIZE);

    // 通知窗管和后端更新数据
    m_multiScreenWorker->updateDaemonDockSize(dockSize);                                // 1.先更新任务栏高度
    m_multiScreenWorker->requestUpdateFrontendGeometry();                               // 2.再更新任务栏位置,保证先1再2
    m_multiScreenWorker->requestNotifyWindowManager();
    m_multiScreenWorker->requestUpdateRegionMonitor();                                  // 界面发生变化，应更新监控区域
}

void MainWindowBase::touchRequestResizeDock()
{
    const QPoint touchPos(QCursor::pos());
    QRect dockRect = geometry();
    // 隐藏状态返回
    if (width() == 0 || height() == 0)
        return;

    int resizeHeight = Utils::SettingValue("com.deepin.dde.dock.touch", QByteArray(), "resizeHeight", 7).toInt();

    QRect touchRect;
    // 任务栏屏幕 内侧边线 内外resizeHeight距离矩形区域内长按可拖动任务栏高度
    switch (position()) {
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

    if (!touchRect.contains(touchPos))
        return;

    qApp->postEvent(m_dragWidget, new QMouseEvent(QEvent::MouseButtonPress, m_dragWidget->mapFromGlobal(touchPos)
                                                  , QPoint(), touchPos, Qt::LeftButton, Qt::NoButton
                                                  , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
}

void MainWindowBase::adjustShadowMask()
{
    if (!m_isShow || m_shadowMaskOptimizeTimer->isActive())
        return;

    m_platformWindowHandle.setWindowRadius(getBorderRadius());
}

void MainWindowBase::onCompositeChanged()
{
    setMaskColor(AutoColor);

    setMaskAlpha(m_multiScreenWorker->opacity());
    m_platformWindowHandle.setBorderWidth(0);

    m_shadowMaskOptimizeTimer->start();
}

void MainWindowBase::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    if (DWindowManagerHelper::instance()->hasComposite()) {
        if (themeType == DGuiApplicationHelper::DarkType) {
            QColor color = Qt::black;
            color.setAlpha(255 * 0.3);
            m_platformWindowHandle.setBorderColor(color);
        } else {
            m_platformWindowHandle.setBorderColor(QColor(QColor::Invalid));
        }
    }
}

void MainWindowBase::setDisplayMode(const Dock::DisplayMode &displayMode)
{
    m_displayMode = displayMode;
    adjustShadowMask();
    m_platformWindowHandle.setShadowOffset(QPoint(0, (displayMode == Dock::DisplayMode::Fashion ? 5 : 0)));
}

void MainWindowBase::setPosition(const Dock::Position &position)
{
    m_position = position;
}

QRect MainWindowBase::getDockGeometry(QScreen *screen, const Dock::Position &pos, const Dock::DisplayMode &displaymode, const Dock::HideState &hideState, bool withoutScale) const
{
    QList<MainWindowBase const *> topMainWindows;           // 所有的顶层主窗口列表
    QList<MainWindowBase const *> lessOrderMainWindows;     // 所有在当前窗口之前的主窗口
    QWidgetList topWidgets = qApp->topLevelWidgets();
    for (QWidget *widget : topWidgets) {
        MainWindowBase *currentWindow = qobject_cast<MainWindowBase *>(widget);
        if (!currentWindow || !currentWindow->isVisible())
            continue;

        topMainWindows << currentWindow;
        if (currentWindow->order() < order())
            lessOrderMainWindows << currentWindow;
    }

    if (!topMainWindows.contains(this))
        return QRect();

    // 对当前窗口前面的所有窗口按照order进行排序
    std::sort(lessOrderMainWindows.begin(), lessOrderMainWindows.end(), [](MainWindowBase const *window1, MainWindowBase const *window2) {
        return window1->order() < window2->order();
    });
    QRect rect;
    const double ratio = withoutScale ? 1 : qApp->devicePixelRatio();
    const int margin = static_cast<int>((displaymode == DisplayMode::Fashion ? 10 : 0) * (withoutScale ? qApp->devicePixelRatio() : 1));
    int dockSize = 0;
    if (hideState == Dock::HideState::Show)
        dockSize = windowSize() * (withoutScale ? qApp->devicePixelRatio() : 1);

    // 拿到当前显示器缩放之前的分辨率
    QRect screenRect = screen->handle()->geometry();
    // 计算所有的窗口的总尺寸
    int totalSize = 0;
    switch (pos) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        // 计算任务栏的总的尺寸
        int width = 0;
        for (MainWindowBase const *mainWindow : topMainWindows) {
            QSize windowSize = mainWindow->suitableSize(pos, screenRect.width(), ratio);
            totalSize += windowSize.width() + mainWindow->dockSpace();
            if (mainWindow == this)
                width = windowSize.width();
        }

        // 计算第一个窗口的X坐标
        int x = screenRect.x() + (static_cast<int>((screenRect.width() / ratio) - totalSize) / 2);
        // 计算当前的X坐标
        for (MainWindowBase const *mainWindow : lessOrderMainWindows) {
            x += mainWindow->suitableSize(pos, screenRect.width(), ratio).width() + mainWindow->dockSpace();
        }
        int y = 0;
        if (pos == Dock::Position::Top)
            y = (screenRect.y() + static_cast<int>(margin / ratio));
        else
            y = (screenRect.y() + static_cast<int>(screenRect.height() / ratio - margin / ratio)) - dockSize;
        rect.setX(x);
        rect.setY(y);
        rect.setWidth(width);
        rect.setHeight(dockSize);
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        int height = 0;
        for (MainWindowBase const *mainWindow : topMainWindows) {
            QSize windowSize = mainWindow->suitableSize(pos, screenRect.height(), ratio);
            totalSize += windowSize.height() + mainWindow->dockSpace();
            if (mainWindow == this)
                height = windowSize.height();
        }
        int x = 0;
        if (pos == Dock::Position::Left)
            x = screenRect.x() + static_cast<int>(margin / ratio);
        else
            x = screenRect.x() + static_cast<int>(screenRect.width() /ratio - margin / ratio) - dockSize;

        int y = screenRect.y() + static_cast<int>(((screenRect.height() / ratio) - totalSize) / 2);
        // 计算y坐标
        for (MainWindowBase const *mainWindow : lessOrderMainWindows)
            y += mainWindow->suitableSize(pos, screenRect.height(), ratio).height() + mainWindow->dockSpace();

        rect.setX(x);
        rect.setY(y);
        rect.setWidth(dockSize);
        rect.setHeight(height);
        break;
    }
    }
    return rect;
}

QVariantAnimation *MainWindowBase::createAnimation(QScreen *screen, const Dock::Position &pos, const Dock::AniAction &act)
{
    /** FIXME
     * 在高分屏2.75倍缩放的情况下，mainWindowGeometry返回的任务栏高度有问题（实际是40,返回是39）
     * 在这里增加判断，当返回值在范围（38，42）开区间内，均认为任务栏显示位置正确，直接返回，不执行动画
     * 也就是在实际值基础上下浮动1像素的误差范围
     * 正常屏幕情况下是没有这个问题的
     */
    QRect mainwindowRect = geometry();
    const QRect dockShowRect = getAnimationRect(getDockGeometry(screen, pos, m_multiScreenWorker->displayMode(), Dock::HideState::Show), pos);
    const QRect &dockHideRect = getAnimationRect(getDockGeometry(screen, pos, m_multiScreenWorker->displayMode(), Dock::HideState::Hide), pos);
    if (act == Dock::AniAction::Show) {
        if (pos == Position::Top || pos == Position::Bottom) {
            if (qAbs(dockShowRect.height() - mainwindowRect.height()) <= 1
                    && mainwindowRect.contains(dockShowRect.center()))
                return nullptr;
        } else if (pos == Position::Left || pos == Position::Right) {
            if (qAbs(dockShowRect.width() - mainwindowRect.width()) <= 1
                    && mainwindowRect.contains(dockShowRect.center()))
                return nullptr;
        }
    }
    if (act == Dock::AniAction::Hide && dockHideRect.size() == mainwindowRect.size())
        return nullptr;

    // 开始播放动画
    QVariantAnimation *ani = new QVariantAnimation(nullptr);
    ani->setEasingCurve(QEasingCurve::InOutCubic);
#ifndef DISABLE_SHOW_ANIMATION
    const bool composite = DWindowManagerHelper::instance()->hasComposite(); // 判断是否开启特效模式
    const int duration = composite ? ANIMATIONTIME : 0;
#else
    const int duration = 0;
#endif
    ani->setDuration(duration);

    connect(ani, &QVariantAnimation::valueChanged, this, [ = ](const QVariant &value) {
        if ((!m_multiScreenWorker->testState(MultiScreenWorker::ShowAnimationStart)
                && !m_multiScreenWorker->testState(MultiScreenWorker::HideAnimationStart)
                && !m_multiScreenWorker->testState(MultiScreenWorker::ChangePositionAnimationStart))
                || ani->state() != QVariantAnimation::State::Running)
            return;

        updateParentGeometry(pos, value.value<QRect>());
    });

    switch (act) {
    case Dock::AniAction::Show: {
        ani->setStartValue(dockHideRect);
        ani->setEndValue(dockShowRect);
        connect(ani, &QVariantAnimation::finished, this, [ = ]{
            updateParentGeometry(pos, dockShowRect);
        });
        break;
    }
    case Dock::AniAction::Hide: {
        ani->setStartValue(dockShowRect);
        ani->setEndValue(dockHideRect);
        connect(ani, &QVariantAnimation::finished, this, [ = ]{
            updateParentGeometry(pos, dockHideRect);
        });
        break;
    }
    }

    return ani;
}

Dock::DisplayMode MainWindowBase::displayMode() const
{
    return m_displayMode;
}

Dock::Position MainWindowBase::position() const
{
    return m_position;
}

int MainWindowBase::windowSize() const
{
    if (m_displayMode == Dock::DisplayMode::Efficient)
        return DockSettings::instance()->getWindowSizeEfficient();

    return DockSettings::instance()->getWindowSizeFashion();
}

bool MainWindowBase::isDraging() const
{
    return m_dragWidget->isDraging();
}

int MainWindowBase::dockSpace() const
{
    return DOCKSPACE;
}

void MainWindowBase::initUi()
{
    DPlatformWindowHandle::enableDXcbForWindow(this, true);
    m_platformWindowHandle.setEnableBlurWindow(true);
    m_platformWindowHandle.setTranslucentBackground(true);
    m_platformWindowHandle.setShadowOffset(QPoint(0, 0));
    QColor shadorColor = Qt::black;
    shadorColor.setAlpha(static_cast<int>(0.3 * 255));
    m_platformWindowHandle.setShadowColor(shadorColor);
}

void MainWindowBase::resizeEvent(QResizeEvent *event)
{
    updateDragGeometry();

    m_shadowMaskOptimizeTimer->start();

    if (!isDraging())
        m_updateDragAreaTimer->start();
}

void MainWindowBase::moveEvent(QMoveEvent *)
{
    updateDragGeometry();

    if (!isDraging())
        m_updateDragAreaTimer->start();
}

void MainWindowBase::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    if (QApplication::overrideCursor() && QApplication::overrideCursor()->shape() != Qt::ArrowCursor)
        QApplication::restoreOverrideCursor();
}

void MainWindowBase::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::RightButton && geometry().contains(QCursor::pos())) {
        m_multiScreenWorker->onAutoHideChanged(false);
        MenuWorker menuWorker;
        menuWorker.exec();
        m_multiScreenWorker->onAutoHideChanged(true);
    }

    DBlurEffectWidget::mousePressEvent(event);
}

void MainWindowBase::showEvent(QShowEvent *event)
{
    if (!m_isShow) {
        m_isShow = true;
        m_shadowMaskOptimizeTimer->start();
    }

    DBlurEffectWidget::showEvent(event);
}
