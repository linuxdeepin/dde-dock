#include <QPropertyAnimation>
#include <QHBoxLayout>
#include <QMimeData>
#include <QImage>
#include <QDebug>
#include <QDrag>
#include "movablelayout.h"

class MovableSpacingItem : public QWidget
{
    Q_OBJECT
    Q_PROPERTY(QSize size READ size WRITE setFixedSize)
public:
    explicit MovableSpacingItem(int duration, QEasingCurve::Type easingType, const QSize &targetSize, QWidget *parent = 0);

signals:
    void growFinish();
    void declineFinish();

public slots:
    void StartGrow(bool immediately = false);
    void StartDecline();

private:
    QSize m_targetSize;
    QPropertyAnimation *m_growAnimation;
    QPropertyAnimation *m_declineAnimation;
};

MovableSpacingItem::MovableSpacingItem(int duration, QEasingCurve::Type easingType, const QSize &targetSize, QWidget *parent)
    : QWidget(parent), m_targetSize(targetSize)
{
    setFixedSize(0, 0);
    m_growAnimation = new QPropertyAnimation(this, "size");
    m_growAnimation->setDuration(duration);
    m_growAnimation->setEasingCurve(easingType);
    m_declineAnimation = new QPropertyAnimation(this, "size");
    m_declineAnimation->setDuration(duration);
    m_declineAnimation->setEasingCurve(easingType);
    connect(m_declineAnimation, &QPropertyAnimation::finished, this, &MovableSpacingItem::declineFinish);
}

void MovableSpacingItem::StartGrow(bool immediately)
{
    if (immediately) {
        setFixedSize(m_targetSize);
        return;
    }

    m_declineAnimation->stop();

    m_growAnimation->setStartValue(this->size());
    m_growAnimation->setEndValue(m_targetSize);

    m_growAnimation->start();
}

void MovableSpacingItem::StartDecline()
{
    m_growAnimation->stop();

    m_declineAnimation->setStartValue(this->size());
    m_declineAnimation->setEndValue(QSize(0, 0));

    m_declineAnimation->start();
}

#include "movablelayout.moc"
///////////////////////////////////////////////////////////////////////////

const int INVALID_MOVE_RADIUS = 5;
const int MAX_SPACINGITEM_COUNT = 2;
const int ANIMATION_DURATION = 300;
const QSize DEFAULT_SPACING_ITEM_SIZE = QSize(48, 48);
const QEasingCurve::Type ANIMATION_CURVE = QEasingCurve::OutCubic;

MovableLayout::MovableLayout(QWidget *parent)
    : QFrame(parent),
      m_lastHoverIndex(-1),
      m_draginItem(nullptr),
      m_defaultSpacingItemSize(DEFAULT_SPACING_ITEM_SIZE),
      m_animationDuration(ANIMATION_DURATION),
      m_animationCurve(ANIMATION_CURVE)
{
    setAttribute(Qt::WA_TranslucentBackground);

    m_layout = new QHBoxLayout(this);
    m_layout->setSpacing(0);
    m_layout->setContentsMargins(0, 0, 0, 0);
    setAcceptDrops(true);
}

MovableLayout::MovableLayout(QBoxLayout::Direction direction, QWidget *parent)
    : QFrame(parent),
      m_lastHoverIndex(-1),
      m_draginItem(nullptr),
      m_defaultSpacingItemSize(DEFAULT_SPACING_ITEM_SIZE),
      m_animationDuration(ANIMATION_DURATION),
      m_animationCurve(ANIMATION_CURVE)
{
    setAttribute(Qt::WA_TranslucentBackground);

    m_layout = new QHBoxLayout(this);
    m_layout->setDirection(direction);
    m_layout->setSpacing(0);
    m_layout->setContentsMargins(0, 0, 0, 0);
    setAcceptDrops(true);
}

int MovableLayout::indexOf(QWidget * const widget, int from) const
{
    return m_widgetList.indexOf(widget, from);
}

QWidget *MovableLayout::widget(int index) const
{
    return m_widgetList.at(index);
}

QList<QWidget *> MovableLayout::widgets() const
{
    return m_widgetList;
}

void MovableLayout::addWidget(QWidget *widget)
{
    m_widgetList.append(widget);

    switch (m_layout->direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::TopToBottom:
        m_layout->takeAt(m_layout->count() - 1);//remove strect
        m_layout->addWidget(widget);
        m_layout->addStretch(1);
        break;
    case QBoxLayout::RightToLeft:
    case QBoxLayout::BottomToTop:
        m_layout->takeAt(0);//remove strect
        m_layout->addWidget(widget);
        m_layout->insertStretch(0, 1);
        break;
    default:
        break;
    }
}

void MovableLayout::insertWidget(int index, QWidget *widget)
{
    m_widgetList.insert(index, widget);
    m_layout->insertWidget(index, widget);
}

void MovableLayout::removeWidget(int index)
{
    m_layout->removeWidget(m_widgetList.takeAt(index));
}

void MovableLayout::removeWidget(QWidget *widget)
{
    m_layout->removeWidget(widget);
    m_widgetList.removeAll(widget);
}

void MovableLayout::addSpacingItem(QWidget *souce, MovableLayout::MoveDirection md, const QSize &size)
{
    int layoutCount = m_layout->count();
    for (int i = 0; i < layoutCount; i ++) {
        if (m_layout->itemAt(i)->widget() == souce) {
            MovableSpacingItem * item;
            int tmpIndex = i;
            switch (md) {
            case MoveLeftToRight:
            case MoveTopToBottom:
                tmpIndex = i == layoutCount - 1 ? layoutCount - 1 : i + 1;
                item = qobject_cast<MovableSpacingItem *>(m_layout->itemAt(tmpIndex)->widget());
                if (item) {
                    //note to other
                    emit spacingItemAdded();
                    item->StartGrow();
                }
                else {
                    insertSpacingItemToLayout(tmpIndex, size);
                }
                break;
            case MoveRightToLeft:
            case MoveBottomToTop:
                tmpIndex = i ;
                item = qobject_cast<MovableSpacingItem *>(m_layout->itemAt(tmpIndex)->widget());
                if (item) {
                    //note to other
                    emit spacingItemAdded();
                    item->StartGrow();
                }
                else {
                    insertSpacingItemToLayout(tmpIndex, size);
                }
                break;
            default:
                return;
            }

            break;
        }
        else {
            continue;
        }
    }
}

int MovableLayout::count() const
{
    return m_widgetList.count();
}

QBoxLayout::Direction MovableLayout::direction() const
{
    return m_layout->direction();
}

void MovableLayout::setDirection(QBoxLayout::Direction direction)
{
    m_layout->setDirection(direction);
}

int MovableLayout::getLayoutSpacing() const
{
    return m_layout->spacing();
}

void MovableLayout::setLayoutSpacing(int spacing)
{
    m_layout->setSpacing(spacing);
}

QPoint basePos(0, 0);
void MovableLayout::mouseMoveEvent(QMouseEvent *event)
{
    //小范围内拖动无效
    if (basePos.isNull()) {
        basePos = event->pos();
        return;
    }
    else {
        if (event->pos().x() - basePos.x() > INVALID_MOVE_RADIUS ||
                event->pos().y() - basePos.y() > INVALID_MOVE_RADIUS) {
            basePos = QPoint(0, 0);
        }
        else {
            return;
        }
    }

    int index = getHoverIndextByPos(event->pos());
    if (index == -1)
        return;

    m_draginItem = m_widgetList.at(index);
    m_lastHoverIndex = index;
    storeDragingItem();

    Qt::MouseButtons btn = event->buttons();
    if(btn == Qt::LeftButton)
    {
        //drag and mimeData object will delete automatically
        QDrag* drag = new QDrag(this);
        QMimeData* mimeData = new QMimeData();
        QImage dataImg = m_draginItem->grab().toImage();
        mimeData->setImageData(QVariant(dataImg));
        drag->setMimeData(mimeData);
        drag->setHotSpot(QPoint(15,15));

        drag->setPixmap(m_draginItem->grab());

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
}

void MovableLayout::dragEnterEvent(QDragEnterEvent *event)
{
    handleDrag(event->pos());

    event->accept();
}

void MovableLayout::dragLeaveEvent(QDragLeaveEvent *event)
{
    Q_UNUSED(event)

    emit spacingItemAdded();
}

void MovableLayout::dragMoveEvent(QDragMoveEvent *event)
{
    if (m_layout->count() > m_widgetList.count() + MAX_SPACINGITEM_COUNT)
        return;

    handleDrag(event->pos());
}

void MovableLayout::dropEvent(QDropEvent *event)
{
    if (m_draginItem && event->source() == this) {
        restoreDragingItem();
    }

    emit spacingItemAdded();
    emit drop(event);
    event->accept();
}

void MovableLayout::resizeEvent(QResizeEvent *event)
{
    emit sizeChanged(event);
}

void MovableLayout::storeDragingItem()
{
    m_draginItem->setVisible(false);
    removeWidget(m_draginItem);
    insertSpacingItemToLayout(m_lastHoverIndex, m_draginItem->size(), true);
}

void MovableLayout::restoreDragingItem()
{
    bool head = true;
    switch (direction()) {
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        if (m_vMoveDirection == MoveBottomToTop)
            head = false;
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        if (m_hMoveDirection == MoveRightToLeft)
            head = false;
    }

    m_draginItem->setVisible(true);
    insertWidget(head ? m_lastHoverIndex : m_lastHoverIndex + 1, m_draginItem);
    m_draginItem = nullptr;
}

void MovableLayout::insertSpacingItemToLayout(int index, const QSize &size, bool immediately)
{
    //note to other
    emit spacingItemAdded();

    MovableSpacingItem *nItem = new MovableSpacingItem(m_animationDuration, m_animationCurve, size);
    connect(this, &MovableLayout::spacingItemAdded, nItem, &MovableSpacingItem::StartDecline);
    connect(nItem, &MovableSpacingItem::declineFinish, [=] {
        m_layout->removeWidget(nItem);
        nItem->deleteLater();
    });
    m_layout->insertWidget(index, nItem);
    nItem->StartGrow(immediately);
}

MovableLayout::MoveDirection MovableLayout::getVMoveDirection(int index, const QPoint &pos)
{
    QWidget *widget = m_widgetList.at(index);
    if (widget) {
        QPoint pp = widget->mapToParent(QPoint(0, 0));
        QRect wr = widget->geometry();
        wr.moveTo(pp);
        if (wr.contains(pos)) {
            QRect topArea(wr.x(), wr.y(), wr.width(), wr.height() / 2);
            if (topArea.contains(pos))
                return MoveTopToBottom;
            else
                return MoveBottomToTop;
        }
        else {
            return MoveUnknow;
        }
    }
    else {
        return MoveUnknow;
    }
}

MovableLayout::MoveDirection MovableLayout::getHMoveDirection(int index, const QPoint &pos)
{
    QWidget *widget = m_widgetList.at(index);
    if (widget) {
        QPoint pp = widget->mapToParent(QPoint(0, 0));
        QRect wr = widget->geometry();
        wr.moveTo(pp);
        if (wr.contains(pos)) {
            QRect leftArea(wr.x(), wr.y(), wr.width() / 2, wr.height());
            if (leftArea.contains(pos))
                return MoveLeftToRight;
            else
                return MoveRightToLeft;
        }
        else {
            return MoveUnknow;
        }
    }
    else {
        return MoveUnknow;
    }
}

bool MovableLayout::getAutoResize() const
{
    return m_autoResize;
}

void MovableLayout::setAutoResize(bool autoResize)
{
    m_autoResize = autoResize;
}

QEasingCurve::Type MovableLayout::getAnimationCurve() const
{
    return m_animationCurve;
}

void MovableLayout::setAnimationCurve(const QEasingCurve::Type &animationCurve)
{
    m_animationCurve = animationCurve;
}

int MovableLayout::getAnimationDuration() const
{
    return m_animationDuration;
}

void MovableLayout::setAnimationDuration(int animationDuration)
{
    m_animationDuration = animationDuration;
}

QSize MovableLayout::getDefaultSpacingItemSize() const
{
    return m_defaultSpacingItemSize;
}

void MovableLayout::setDefaultSpacingItemSize(const QSize &defaultSpacingItemSize)
{
    m_defaultSpacingItemSize = defaultSpacingItemSize;
}

void MovableLayout::setDuration(int v)
{
    m_animationDuration = v;
}

void MovableLayout::setEasingCurve(QEasingCurve::Type curve)
{
    m_animationCurve = curve;
}

bool MovableLayout::event(QEvent *e)
{
    if (e->type() == QEvent::LayoutRequest && getAutoResize()) {
        setFixedSize(sizeHint());
    }

    return QFrame::event(e);
}

int MovableLayout::getHoverIndextByPos(const QPoint &pos)
{
    for (int i = 0; i < m_widgetList.count(); i ++) {
        QWidget *widget = m_widgetList.at(i);
        QPoint pp = widget->mapToParent(QPoint(0, 0));
        QRect wr = widget->geometry();
        wr.moveTo(pp);
        if (wr.contains(pos)) {
            return i;
        }
    }

    return -1;
}

void MovableLayout::handleDrag(const QPoint &pos)
{
    int index = getHoverIndextByPos(pos);
    if (index != -1) {
        QSize spacingSize = m_draginItem ? m_draginItem->size() : m_defaultSpacingItemSize;
        bool shouldAddSpacing = false;
        MoveDirection d = MoveUnknow;

        switch (m_layout->direction()) {
        case QBoxLayout::LeftToRight:
        case QBoxLayout::RightToLeft:
            //the same index but direction rever
            if (index != m_lastHoverIndex || (m_hMoveDirection != getHMoveDirection(index, pos))) {
                shouldAddSpacing = true;
                d = m_hMoveDirection;
            }
            break;
        case QBoxLayout::TopToBottom:
        case QBoxLayout::BottomToTop:
            if (index != m_lastHoverIndex || (m_vMoveDirection != getVMoveDirection(index, pos))) {
                shouldAddSpacing = true;
                d = m_vMoveDirection;
            }
            break;
        }

        if (shouldAddSpacing) {
            updateCurrentHoverInfo(index, pos);
            addSpacingItem(m_widgetList.at(index), d, spacingSize);
        }

        m_hoverToSpacing = false;
    }
    else {
        m_hoverToSpacing = true;
    }
}

void MovableLayout::updateCurrentHoverInfo(int index, const QPoint &pos)
{
    m_lastHoverIndex = index;
    m_vMoveDirection = getVMoveDirection(index, pos);
    m_hMoveDirection = getHMoveDirection(index, pos);
}

