#include "docklayout.h"
#include "abstractdockitem.h"

DockLayout::DockLayout(QWidget *parent) :
    QWidget(parent)
{
    this->setAcceptDrops(true);
}

void DockLayout::addItem(AbstractDockItem *item)
{
    insertItem(item,appList.count());
}

void DockLayout::insertItem(AbstractDockItem *item, int index)
{
    item->setParent(this);
    int appCount = appList.count();
    index = index > appCount ? appCount : (index < 0 ? 0 : index);

    appList.insert(index,item);
    connect(item, &AbstractDockItem::mouseRelease, this, &DockLayout::slotItemRelease);
    connect(item, &AbstractDockItem::dragStart, this, &DockLayout::slotItemDrag);
    connect(item, &AbstractDockItem::dragEntered, this, &DockLayout::slotItemEntered);
    connect(item, &AbstractDockItem::dragExited, this, &DockLayout::slotItemExited);

    relayout();
}

void DockLayout::removeItem(int index)
{
    delete appList.takeAt(index);
    relayout();
}

void DockLayout::moveItem(int from, int to)
{
    appList.move(from,to);
    relayout();
}

void DockLayout::setItemMoveable(int index, bool moveable)
{
    appList.at(index)->setMoveable(moveable);
}

void DockLayout::setSpacing(qreal spacing)
{
    this->itemSpacing = spacing;
}

void DockLayout::setSortDirection(DockLayout::Direction value)
{
    this->sortDirection = value;
}

void DockLayout::sortLeftToRight()
{
    if (appList.count() <= 0)
        return;

    appList.at(0)->move(itemSpacing,0);

    for (int i = 1; i < appList.count(); i ++)
    {
        AbstractDockItem * frontItem = appList.at(i - 1);
        appList.at(i)->move(frontItem->pos().x() + frontItem->width() + itemSpacing,0);
    }
}

void DockLayout::sortRightToLeft()
{
    if (appList.count()<=0)
        return;

    appList.at(0)->move(this->width() - itemSpacing - appList.at(0)->width(),0);

    for (int i = 1; i < appList.count(); i++)
    {
        AbstractDockItem *fromItem = appList.at(i - 1);
        AbstractDockItem *toItem = appList.at(i);
        toItem->move(fromItem->x() - itemSpacing - toItem->width(),0);
    }
}

bool DockLayout::hasSpacingItemInList()
{
    if (appList.count() <= 1)
        return false;
    if (appList.at(0)->x() > itemSpacing)
        return true;

    for (int i = 1; i < appList.count(); i ++)
    {
        if (appList.at(i)->x() - itemSpacing != appList.at(i - 1)->x() + appList.at(i - 1)->width())
        {
            return true;
        }
    }
    return false;
}

int DockLayout::indexOf(AbstractDockItem *item)
{
    return appList.indexOf(item);
}

int DockLayout::indexOf(int x, int y)
{
    //TODO
    return 0;
}

void DockLayout::relayout()
{
    switch (sortDirection)
    {
    case LeftToRight:
        sortLeftToRight();
        break;
    case RightToLeft:
        sortRightToLeft();
        break;
    default:
        break;
    }
}

void DockLayout::addSpacingItem()
{
    if (tmpAppMap.isEmpty())
        return;

    AbstractDockItem *tmpItem = tmpAppMap.firstKey();
    for (int i = appList.count() -1;i >= lastHoverIndex; i-- )
    {
        AbstractDockItem *targetItem = appList.at(i);
        targetItem->setNextPos(targetItem->x() + tmpItem->width() + itemSpacing,0);

        QPropertyAnimation *animation = new QPropertyAnimation(targetItem, "pos");
        animation->setStartValue(targetItem->pos());
        animation->setEndValue(targetItem->getNextPos());
        animation->setDuration(150 + i * 10);
        animation->setEasingCurve(QEasingCurve::OutCubic);

        animation->start();
    }
}

void DockLayout::dragoutFromLayout(int index)
{
    AbstractDockItem * tmpItem = appList.takeAt(index);
    tmpItem->setVisible(false);
    tmpAppMap.insert(tmpItem,index);
}

int DockLayout::getContentsWidth()
{
    int tmpWidth = appList.count() * itemSpacing;
    for (int i = 0; i < appList.count(); i ++)
    {
        tmpWidth += appList.at(i)->width();
    }
    return tmpWidth;
}

int DockLayout::getItemCount()
{
    return appList.count();
}

void DockLayout::dragEnterEvent(QDragEnterEvent *event)
{
    event->setDropAction(Qt::MoveAction);
    event->accept();
}

void DockLayout::dropEvent(QDropEvent *event)
{
    AbstractDockItem * tmpItem = tmpAppMap.firstKey();
    tmpAppMap.remove(tmpItem);
    tmpItem->setVisible(true);
    if (indexOf(tmpItem) == -1)
    {
        if (movingForward)
            insertItem(tmpItem,lastHoverIndex);
        else
            insertItem(tmpItem,lastHoverIndex + 1);
    }

    emit itemDropped();
}

void DockLayout::slotItemDrag()
{
//    qWarning() << "Item draging..."<<x<<y<<item;
    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    int tmpIndex = indexOf(item);
    if (tmpIndex != -1)
    {
        lastHoverIndex = tmpIndex;
        m_lastPost = QCursor::pos();
        dragoutFromLayout(tmpIndex);

        emit dragStarted();
    }
}

void DockLayout::slotItemRelease(int, int)
{
    //outside frame,destroy it
    //inside frame,insert it
    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    item->setVisible(true);
    if (indexOf(item) == -1)
    {
        insertItem(item,lastHoverIndex);
    }
}

void DockLayout::slotItemEntered(QDragEnterEvent *)
{
    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    int tmpIndex = indexOf(item);
    lastHoverIndex = tmpIndex;
    if (!hasSpacingItemInList())
    {
        addSpacingItem();
        return;
    }

    QPoint tmpPos = QCursor::pos();

    if (tmpPos.x() - m_lastPost.x() == 0)
        return;

    switch (sortDirection)
    {
    case LeftToRight:
        movingForward = tmpPos.x() - m_lastPost.x() < 0;
        break;
    case RightToLeft:
        movingForward = tmpPos.x() - m_lastPost.x() > 0;
        break;
    }

    m_lastPost = tmpPos;

    if (!tmpAppMap.isEmpty())
    {
        AbstractDockItem *targetItem = appList.at(tmpIndex);
        if (movingForward)
        {
            targetItem->setNextPos(QPoint(targetItem->x() + tmpAppMap.firstKey()->width() + itemSpacing,0));
        }
        else
        {
            targetItem->setNextPos(QPoint(targetItem->x() - tmpAppMap.firstKey()->width() - itemSpacing,0));
        }

        QPropertyAnimation *animation = new QPropertyAnimation(targetItem, "pos");
        animation->setStartValue(targetItem->pos());
        animation->setEndValue(targetItem->getNextPos());
        animation->setDuration(200);
        animation->setEasingCurve(QEasingCurve::OutCubic);
        animation->start();
    }
}

void DockLayout::slotItemExited(QDragLeaveEvent *)
{

}
