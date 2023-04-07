// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "../appitem.h"
#include "appdragwidget.h"
#include "utils.h"
#include "displaymanager.h"

#include "org_deepin_dde_xeventmonitor1.h"

#define SPLIT_NONE 0
#define SPLIT_LEFT 1
#define SPLIT_RIGHT 2

using XEventMonitor = ::org::deepin::dde::XEventMonitor1;

AppDragWidget::AppDragWidget(QWidget *parent)
    : QGraphicsView(parent)
    , m_object(new AppGraphicsObject)
    , m_scene(new QGraphicsScene(this))
    , m_followMouseTimer(new QTimer(this))
    , m_animScale(new QPropertyAnimation(m_object.get(), "scale", this))
    , m_animRotation(new QPropertyAnimation(m_object.get(), "rotation", this))
    , m_animOpacity(new QPropertyAnimation(m_object.get(), "opacity", this))
    , m_animGroup(new QParallelAnimationGroup(this))
    , m_goBackAnim(new QPropertyAnimation(this, "pos", this))
    , m_dockPosition(Dock::Position::Bottom)
    , m_popupWindow(new DockPopupWindow(nullptr))
    , m_distanceMultiple(Utils::SettingValue("com.deepin.dde.dock.distancemultiple", "/com/deepin/dde/dock/distancemultiple/", "distance-multiple", 1.5).toDouble())
    , m_item(nullptr)
    , m_dockScreen(nullptr)
{
    m_popupWindow->setRadius(18);

    m_scene->addItem(m_object.get());
    setScene(m_scene);

    setAttribute(Qt::WA_TranslucentBackground);
    if (Utils::IS_WAYLAND_DISPLAY) {
        setWindowFlags(Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint | Qt::Window | Qt::FramelessWindowHint);
        setAttribute(Qt::WA_NativeWindow);
        initWaylandEnv();
    } else {
        setWindowFlags(Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint);
    }
    viewport()->setAutoFillBackground(false);
    setFrameShape(QFrame::NoFrame);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setMouseTracking(true);

    setAcceptDrops(true);

    initAnimations();

    if (!Utils::IS_WAYLAND_DISPLAY) {
        m_followMouseTimer->setInterval(16);
        connect(m_followMouseTimer, &QTimer::timeout, this, &AppDragWidget::onFollowMouse);
        m_followMouseTimer->start();
        QTimer::singleShot(0, this, &AppDragWidget::onFollowMouse);
    }
}

void AppDragWidget::execFinished()
{
    if (!m_bDragDrop)
        return;

    dropHandler(QCursor::pos());
}

void AppDragWidget::mouseMoveEvent(QMouseEvent *event)
{
    QGraphicsView::mouseMoveEvent(event);
    // hide widget when receiving mouseMoveEvent because this means drag-and-drop has been finished
    if (m_goBackAnim->state() != QPropertyAnimation::State::Running
            && m_animGroup->state() != QParallelAnimationGroup::Running) {
        hide();
    }
}

void AppDragWidget::dragEnterEvent(QDragEnterEvent *event)
{
    event->accept();
    m_bDragDrop = true;
}

void AppDragWidget::dragMoveEvent(QDragMoveEvent *event)
{
    if (Utils::IS_WAYLAND_DISPLAY)
        return QGraphicsView::dragMoveEvent(event);

    if (m_bDragDrop)
        moveHandler(QCursor::pos());
}

/**获取应用的左上角坐标
 * @brief AppDragWidget::topleftPoint
 * @return 返回应用左上角坐标
 */
const QPoint AppDragWidget::topleftPoint() const
{
    QPoint p;
    const QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    return p;
}

/**拖动从任务栏移除应用时浮窗坐标
 * @brief AppDragWidget::popupMarkPoint
 * @param pos　任务栏所在位置
 * @return　拖动从任务栏移除应用时浮窗坐标
 */
const QPoint AppDragWidget::popupMarkPoint(Dock::Position pos)
{
    QPoint p(topleftPoint());
    QRect r = rect();
    //关闭特效,原本的图标设置小,然后隐藏,需要手动设置大小保证tips位置正确
    if (!DWindowManagerHelper::instance()->hasComposite()) {
        r.setWidth(m_iconSize.width() + 3);
        r.setHeight(m_iconSize.height() + 3);
    }

    switch (pos) {
    case Top:
        p += QPoint(r.width() / 2, r.height());
        break;
    case Bottom:
        if (!DWindowManagerHelper::instance()->hasComposite()) {
            p += QPoint(0 , -r.height() / 2);
        } else {
            p += QPoint(r.width() / 2, 0);
        }
        break;
    case Left:
        p += QPoint(r.width(), r.height() / 2);
        break;
    case Right:
        p += QPoint(0, r.height() / 2);
        break;
    }
    return p;
}

void AppDragWidget::dropEvent(QDropEvent *event)
{
    if (Utils::IS_WAYLAND_DISPLAY)
        return dropEvent(event);

    m_followMouseTimer->stop();
    dropHandler(QCursor::pos());
}

void AppDragWidget::hideEvent(QHideEvent *event)
{
    deleteLater();
    if (Utils::IS_WAYLAND_DISPLAY)
        QGraphicsView::hideEvent(event);
}

void AppDragWidget::dropHandler(const QPoint &pos)
{
    m_bDragDrop = false;

    if (canSplitWindow(pos)) {
        if (DWindowManagerHelper::instance()->hasComposite()) {
            showRemoveAnimation();
        } else {
            hide();
        }
        m_popupWindow->setVisible(false);
        Q_EMIT requestSplitWindow(splitPosition());
    } else {
        if (DWindowManagerHelper::instance()->hasComposite()) {
            showGoBackAnimation();
        } else {
            hide();
        }
    }
}

void AppDragWidget::moveHandler(const QPoint &pos)
{
    if (canSplitWindow(pos)) {
        QRect screenGeometry = splitGeometry(pos);
        if (screenGeometry.isValid() && screenGeometry != m_lastMouseGeometry) {
            qDebug() << "change area:" << screenGeometry;
            Q_EMIT requestChangedArea(screenGeometry);
            m_lastMouseGeometry = screenGeometry;
        }
    }
}

void AppDragWidget::moveCurrent(const QPoint &destPos)
{
    if (DWindowManagerHelper::instance()->hasComposite()) {
        move(destPos.x() - width() / 2, destPos.y() - height() / 2);
    } else {
         // 窗口特效未开启时会隐藏m_object绘制的图标，移动的图标为QDrag绘制的图标
        move(destPos.x(), destPos.y());
    }
}

void AppDragWidget::setAppPixmap(const QPixmap &pix)
{
    if (DWindowManagerHelper::instance()->hasComposite()) {
        // QSize(3, 3) to fix pixmap be cliped
        setFixedSize(pix.size() + QSize(3, 3));
    } else {
        setFixedSize(QSize(10, 10));
    }

    m_iconSize = pix.size();
    m_object->setAppPixmap(pix);
    m_object->setTransformOriginPoint(pix.rect().center());
}

void AppDragWidget::setDockInfo(Dock::Position dockPosition, const QRect &dockGeometry)
{
    m_dockPosition = dockPosition;
    m_dockGeometry = dockGeometry;
}

void AppDragWidget::setOriginPos(const QPoint position)
{
    m_originPoint = position;
}

void AppDragWidget::setPixmapOpacity(qreal opacity)
{
    if (canSplitWindow(QCursor::pos())) {
        m_object->setOpacity(opacity);
        m_animOpacity->setStartValue(opacity);
    } else {
        m_object->setOpacity(1.0);
        m_animOpacity->setStartValue(1.0);
    }
}

void AppDragWidget::initAnimations()
{
    m_animScale->setDuration(300);
    m_animScale->setStartValue(1.0);
    m_animScale->setEndValue(0.0);

    m_animRotation->setDuration(300);
    m_animRotation->setStartValue(0);
    m_animRotation->setEndValue(90);

    m_animOpacity->setDuration(300);
    m_animOpacity->setStartValue(1.0);
    m_animOpacity->setEndValue(0.0);

    m_animGroup->addAnimation(m_animScale);
    m_animGroup->addAnimation(m_animRotation);
    m_animGroup->addAnimation(m_animOpacity);

    connect(m_animGroup, &QParallelAnimationGroup::stateChanged,
            this, &AppDragWidget::onRemoveAnimationStateChanged);
    connect(m_goBackAnim, &QPropertyAnimation::finished, this, &AppDragWidget::hide);
}

/**显示移除动画
 * @brief AppDragWidget::showRemoveAnimation
 */
void AppDragWidget::showRemoveAnimation()
{
    if (m_animGroup->state() == QParallelAnimationGroup::Running) {
        m_animGroup->stop();
    }
    m_object->resetProperty();
    m_animGroup->start();
}

/**显示放弃移除后的动画
 * @brief AppDragWidget::showGoBackAnimation
 */
void AppDragWidget::showGoBackAnimation()
{
    m_goBackAnim->setDuration(300);
    m_goBackAnim->setStartValue(pos());
    m_goBackAnim->setEndValue(m_originPoint);
    m_goBackAnim->start();
}

void AppDragWidget::onRemoveAnimationStateChanged(QAbstractAnimation::State newState,
                                                  QAbstractAnimation::State oldState)
{
    Q_UNUSED(oldState);
    if (newState == QAbstractAnimation::Stopped) {
        hide();
    }
}

/** 判断应用区域图标是否被拖出任务栏
 * @brief AppDragWidget::canSplitWindow
 * @return 返回true应用移出任务栏，false应用在任务栏内
 */
bool AppDragWidget::canSplitWindow(const QPoint &pos) const
{
    switch (m_dockPosition) {
    case Dock::Position::Left:
        if ((pos.x() > m_dockGeometry.topRight().x())) {
            return true;
        }
        break;
    case Dock::Position::Top:
        if ((pos.y() > m_dockGeometry.bottomLeft().y())) {
            return true;
        }
        break;
    case Dock::Position::Right:
        if ((m_dockGeometry.topLeft().x() > pos.x())) {
            return true;
        }
        break;
    case Dock::Position::Bottom:
        if ((m_dockGeometry.topLeft().y() > pos.y())) {
            return true;
        }
        break;
    }

    return false;
}

/**
 * @brief AppDragWidget::splitPosition
 * @return 1 左分屏；2 右分屏；5 左上；6 右上；9 左下；10 右下；15全屏。这些值是窗管给的
 */
ScreenSpliter::SplitDirection AppDragWidget::splitPosition() const
{
    QPoint pos = QCursor::pos();
    QScreen *currentScreen = DisplayManager::instance()->screenAt(pos);

    if (!currentScreen)
        return ScreenSpliter::None;

    int xCenter = currentScreen->geometry().x() + currentScreen->size().width() / 2;
    // 1表示左分屏
    if (pos.x() < xCenter)
        return ScreenSpliter::Left;

    // 2表示右分屏
    if (pos.x() > xCenter)
        return ScreenSpliter::Right;

    return ScreenSpliter::None;
}

void AppDragWidget::adjustDesktopGeometry(QRect &rect) const
{
    QRect rectGeometry = m_dockGeometry;
    rectGeometry.setWidth(rectGeometry.width() * qApp->devicePixelRatio());
    rectGeometry.setHeight(rectGeometry.height() * qApp->devicePixelRatio());
    switch (m_dockPosition) {
    case Dock::Position::Left: {
        int leftX = (rectGeometry.x() + rectGeometry.width()) * qApp->devicePixelRatio();
        if (rect.x() < leftX) {
            rect.setX(leftX);
            rect.setWidth(rect.width() - (leftX - rect.x()));
        }
        break;
    }
    case Dock::Position::Top: {
        int topY = (rectGeometry.y() + rectGeometry.height()) * qApp->devicePixelRatio();
        if (rect.y() < topY) {
            rect.setY(topY);
            rect.setHeight(rect.height() - (topY - rect.y()));
        }
        break;
    }
    case Dock::Position::Right: {
        int rightX = rectGeometry.x() * qApp->devicePixelRatio();
        if (rightX < rect.x() + rect.width() * qApp->devicePixelRatio())
            rect.setWidth(rect.width() - (rect.x() + rect.width() - rightX));
        break;
    }
    case Dock::Position::Bottom: {
        int bottomY =  rectGeometry.y() * qApp->devicePixelRatio();
        if (bottomY < rect.y() + rect.height() * qApp->devicePixelRatio())
            rect.setHeight(rect.height() - (rect.y() + rect.height() - bottomY));
        break;
    }
    }
}

QRect AppDragWidget::splitGeometry(const QPoint &pos) const
{
    QList<QScreen *> screens = DisplayManager::instance()->screens();
    for (QScreen *screen : screens) {
        QRect screenGeometry = screen->geometry();
        screenGeometry.setWidth(screenGeometry.width() * qApp->devicePixelRatio());
        screenGeometry.setHeight(screenGeometry.height() * qApp->devicePixelRatio());
        if (!screenGeometry.contains(pos))
            continue;

        // 左右分屏即可
        int centerX = screenGeometry.x() + screenGeometry.width() / 2;
        if (pos.x() < centerX) {
            // 左分屏
            QRect rectLeft = screenGeometry;
            rectLeft.setWidth(screenGeometry.width() / 2);
            adjustDesktopGeometry(rectLeft);
            return rectLeft;
        }
        if (pos.x() > centerX) {
            // 右分屏
            QRect rectRight = screenGeometry;
            rectRight.setLeft(screenGeometry.x() + screenGeometry.width() / 2);
            rectRight.setWidth(screenGeometry.width() / 2);
            adjustDesktopGeometry(rectRight);
            return rectRight;
        }
        break;
    }
    return QRect();
}

void AppDragWidget::initWaylandEnv()
{
    if (!Utils::IS_WAYLAND_DISPLAY)
        return;

    // 由于在wayland环境下无法触发drop事件，导致鼠标无法释放，所以这里暂时用XEventMonitor的方式(具体原因待查)
    XEventMonitor *extralEventInter = new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus());
    QString key = extralEventInter->RegisterFullScreen();
    connect(this, &AppDragWidget::destroyed, this, [ key, extralEventInter ] {
        extralEventInter->UnregisterArea(key);
        delete extralEventInter;
        QDrag::cancel();
    });

    connect(extralEventInter, &XEventMonitor::ButtonRelease, this, &AppDragWidget::onButtonRelease);
    connect(extralEventInter, &XEventMonitor::CursorMove,this, &AppDragWidget::onCursorMove);
}

void AppDragWidget::onButtonRelease(int, int x, int y, const QString &)
{
    if (!m_bDragDrop)
        return;

    dropHandler(QPoint(x, y));

    QDrag::cancel();
}

void AppDragWidget::onCursorMove(int x, int y, const QString &)
{
    QPoint pos = QPoint(x, y);
    moveCurrent(pos);
    moveHandler(pos);
}

void AppDragWidget::onFollowMouse()
{
    moveCurrent(QCursor::pos());
}

void AppDragWidget::enterEvent(QEvent *event)
{
    Q_UNUSED(event);

    if (m_goBackAnim->state() != QPropertyAnimation::State::Running
            && m_animGroup->state() != QParallelAnimationGroup::Running) {
        hide();
    }
}

QuickDragWidget::QuickDragWidget(QWidget *parent)
    : AppDragWidget(parent)
{

}

QuickDragWidget::~QuickDragWidget()
{
}

void QuickDragWidget::dropEvent(QDropEvent *event)
{
    Q_UNUSED(event);

    m_followMouseTimer->stop();
    m_bDragDrop = false;

    if (isRemoveAble(QCursor::pos())) {
        if (DWindowManagerHelper::instance()->hasComposite())
            showRemoveAnimation();
        else
            hide();

        m_popupWindow->setVisible(false);
    } else {
        if (DWindowManagerHelper::instance()->hasComposite())
            showGoBackAnimation();
        else
            hide();

        Q_EMIT requestDropItem(event);
    }
}

void QuickDragWidget::dragMoveEvent(QDragMoveEvent *event)
{
    AppDragWidget::dragMoveEvent(event);
    requestDragMove(event);
}

/**判断图标拖到一定高度(默认任务栏高度的1.5倍)后是否可以移除
 * @brief AppDragWidget::isRemoveAble
 * @param curPos 当前鼠标所在位置
 * @return 返回true可移除，false不可移除
 */
bool QuickDragWidget::isRemoveAble(const QPoint &curPos)
{
    const QPoint &p = curPos;
    switch (m_dockPosition) {
    case Dock::Position::Left:
        if ((p.x() - m_dockGeometry.topRight().x()) > (m_dockGeometry.width() * m_distanceMultiple)) {
            return true;
        }
        break;
    case Dock::Position::Top:
        if ((p.y() - m_dockGeometry.bottomLeft().y()) > (m_dockGeometry.height() * m_distanceMultiple)) {
            return true;
        }
        break;
    case Dock::Position::Right:
        if ((m_dockGeometry.topLeft().x() - p.x()) > (m_dockGeometry.width() * m_distanceMultiple)) {
            return true;
        }
        break;
    case Dock::Position::Bottom:
        if ((m_dockGeometry.topLeft().y() - p.y()) > (m_dockGeometry.height() * m_distanceMultiple)) {
            return true;
        }
        break;
    }
    return false;
}
