// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "../appitem.h"
#include "appdragwidget.h"
#include "utils.h"

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
    , m_removeTips(new TipsWidget(this))
    , m_popupWindow(new DockPopupWindow(nullptr))
    , m_distanceMultiple(Utils::SettingValue("com.deepin.dde.dock.distancemultiple", "/com/deepin/dde/dock/distancemultiple/", "distance-multiple", 1.5).toDouble())
    , m_item(nullptr)
    , m_cursorPosition(-1, -1)
{
    m_removeTips->setText(tr("Remove"));
    m_removeTips->setObjectName("AppRemoveTips");
    m_removeTips->setVisible(false);
    m_removeTips->installEventFilter(this);

    m_popupWindow->setShadowBlurRadius(20);
    m_popupWindow->setRadius(18);
    m_popupWindow->setShadowYOffset(2);
    m_popupWindow->setShadowXOffset(0);
    m_popupWindow->setArrowWidth(18);
    m_popupWindow->setArrowHeight(10);

    m_scene->addItem(m_object.get());
    setScene(m_scene);

    setWindowFlags(Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint);
    setAttribute(Qt::WA_TranslucentBackground);
    if (Utils::IS_WAYLAND_DISPLAY) {
        setWindowFlags(Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint | Qt::Window);
        setAttribute(Qt::WA_NativeWindow);
    } else {
        setAcceptDrops(true);
    }
    viewport()->setAutoFillBackground(false);
    setFrameShape(QFrame::NoFrame);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setMouseTracking(true);

    initAnimations();

    m_followMouseTimer->setInterval(16);
    connect(m_followMouseTimer, &QTimer::timeout, this, &AppDragWidget::onFollowMouse);
    m_followMouseTimer->start();
    QTimer::singleShot(0, this, &AppDragWidget::onFollowMouse);
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
    if (Utils::IS_WAYLAND_DISPLAY) {
        QGraphicsView::dragEnterEvent(event);
    } else {
        event->accept();
        m_bDragDrop = true;
    }
}

void AppDragWidget::dragMoveEvent(QDragMoveEvent *event)
{
    if (Utils::IS_WAYLAND_DISPLAY) {
        QGraphicsView::dragMoveEvent(event);
    } else {
        showRemoveTips();
        if (isRemoveItem() && m_bDragDrop) {
            emit requestRemoveItem();
        }
    }
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
            p += QPoint(size().width() / 2 , -r.height() / 2);
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
    if (Utils::IS_WAYLAND_DISPLAY) {
        QGraphicsView::dropEvent(event);
    } else {
        m_followMouseTimer->stop();
        m_bDragDrop = false;

        if (isRemoveAble(QCursor::pos())) {
            if (DWindowManagerHelper::instance()->hasComposite()) {
                showRemoveAnimation();
            } else {
                hide();
            }
            AppItem *appItem = static_cast<AppItem *>((Utils::IS_WAYLAND_DISPLAY && m_item) ? m_item : event->source());
            appItem->undock();
            m_popupWindow->setVisible(false);
        } else {
            if (DWindowManagerHelper::instance()->hasComposite()) {
                showGoBackAnimation();
            } else {
                hide();
            }
        }
    }
}

void AppDragWidget::hideEvent(QHideEvent *event)
{
    deleteLater();
    if (Utils::IS_WAYLAND_DISPLAY)
        QGraphicsView::hideEvent(event);
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
    if (isRemoveAble(QCursor::pos())) {
        m_object->setOpacity(opacity);
        m_animOpacity->setStartValue(opacity);
    } else {
        m_object->setOpacity(1.0);
        m_animOpacity->setStartValue(1.0);
    }
}

bool AppDragWidget::isRemoveable(const Position &dockPos, const QRect &doctRect)
{
    const QPoint &p = QCursor::pos();
    switch (dockPos) {
        case Dock::Position::Left:
            if ((p.x() - doctRect.topRight().x()) > (doctRect.width() * 3)) {
                return true;
            }
            break;
        case Dock::Position::Top:
            if ((p.y() - doctRect.bottomLeft().y()) > (doctRect.height() * 3)) {
                return true;
            }
            break;
        case Dock::Position::Right:
            if ((doctRect.topLeft().x() - p.x()) > (doctRect.width() * 3)) {
                return true;
            }
            break;
        case Dock::Position::Bottom:
            if ((doctRect.topLeft().y() - p.y()) > (doctRect.height() * 3)) {
                return true;
            }
            break;
    }
    return false;
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

/**判断图标拖到一定高度(默认任务栏高度的1.5倍)后是否可以移除
 * @brief AppDragWidget::isRemoveAble
 * @param curPos 当前鼠标所在位置
 * @return 返回true可移除，false不可移除
 */
bool AppDragWidget::isRemoveAble(const QPoint &curPos)
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

/**判断应用区域图标是否被拖出任务栏
 * @brief AppDragWidget::isRemoveItem
 * @return 返回true应用移出任务栏，false应用在任务栏内
 */
bool AppDragWidget::isRemoveItem()
{
    const QPoint &p = QCursor::pos();
    switch (m_dockPosition) {
    case Dock::Position::Left:
        if ((p.x() > m_dockGeometry.topRight().x())) {
            return true;
        }
        break;
    case Dock::Position::Top:
        if ((p.y() > m_dockGeometry.bottomLeft().y())) {
            return true;
        }
        break;
    case Dock::Position::Right:
        if ((m_dockGeometry.topLeft().x() > p.x())) {
            return true;
        }
        break;
    case Dock::Position::Bottom:
        if ((m_dockGeometry.topLeft().y() > p.y())) {
            return true;
        }
        break;
    }
    return false;
}

void AppDragWidget::onFollowMouse()
{
    QPoint destPos = QCursor::pos();
    move(destPos.x() - width() / 2, destPos.y() - height() / 2);
}

void AppDragWidget::enterEvent(QEvent *event)
{
    Q_UNUSED(event);

    if (m_goBackAnim->state() != QPropertyAnimation::State::Running
            && m_animGroup->state() != QParallelAnimationGroup::Running) {
        hide();
    }
}

bool AppDragWidget::event(QEvent *event)
{
    if (Utils::IS_WAYLAND_DISPLAY && event->type() == QEvent::DeferredDelete)
        requestRemoveSelf(isRemoveAble(m_cursorPosition));

    return QGraphicsView::event(event);
}

/**显示移除应用提示窗口
 * @brief AppDragWidget::showRemoveTips
 */
void AppDragWidget::showRemoveTips()
{
    Dock::Position pos = Dock::Position::Bottom;

    DockPopupWindow *popup = m_popupWindow.data();
    if (isRemoveAble(QCursor::pos())) {
        QWidget *lastContent = popup->getContent();
        if (lastContent)
            lastContent->setVisible(false);

        switch (pos) {
        case Top:    popup->setArrowDirection(DockPopupWindow::ArrowTop);     break;
        case Bottom: popup->setArrowDirection(DockPopupWindow::ArrowBottom);  break;
        case Left:   popup->setArrowDirection(DockPopupWindow::ArrowLeft);    break;
        case Right:  popup->setArrowDirection(DockPopupWindow::ArrowRight);   break;
        }
        popup->resize(m_removeTips->sizeHint());
        popup->setContent(m_removeTips);

        const QPoint p = popupMarkPoint(pos);
        if (!popup->isVisible())
            QMetaObject::invokeMethod(popup, "show", Qt::QueuedConnection, Q_ARG(QPoint, p), Q_ARG(bool, true));
        else
            popup->show(p, true);

        m_object->setOpacity(0.5);
        m_animOpacity->setStartValue(0.5);
    } else {
        m_object->setOpacity(1.0);
        m_animOpacity->setStartValue(1.0);
        if (popup->isVisible())
            popup->setVisible(false);
    }
}

void AppDragWidget::moveEvent(QMoveEvent *event)
{
    Q_UNUSED(event);
    showRemoveTips();
    if (Utils::IS_WAYLAND_DISPLAY) {
        m_cursorPosition = QCursor::pos();
        if (isRemoveItem()) {
            emit requestRemoveItem();
        }
    }
}
