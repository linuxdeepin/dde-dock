#include "docklayout.h"

DockLayout::DockLayout(QWidget *parent) :
    QWidget(parent)
{
    this->setAcceptDrops(true);
}

void DockLayout::setParent(QWidget *parent)
{
    this->setParent(parent);
}

void DockLayout::addItem(AppItem *item)
{
    insertItem(item,appList.count());
}

void DockLayout::insertItem(AppItem *item, int index)
{
    item->setParent(this);
    int appCount = appList.count();
    index = index > appCount ? appCount : (index < 0 ? 0 : index);

    appList.insert(index,item);
    connect(item,SIGNAL(mouseRelease(int,int,AppItem*)),this,SLOT(slotItemRelease(int,int,AppItem*)));
    connect(item, SIGNAL(dragStart(AppItem*)),this,SLOT(slotItemDrag(AppItem*)));
    connect(item,SIGNAL(dragEntered(QDragEnterEvent*,AppItem*)),this,SLOT(slotItemEntered(QDragEnterEvent*,AppItem*)));
    connect(item,SIGNAL(dragExited(QDragLeaveEvent*,AppItem*)),this,SLOT(slotItemExited(QDragLeaveEvent*,AppItem*)));

    relayout();
}

void DockLayout::removeItem(int index)
{
    delete appList.takeAt(index);
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

void DockLayout::setMargin(qreal margin)
{
    this->leftMargin = margin;
    this->rightMargin = margin;
    this->topMargin = margin;
    this->bottomMargin = margin;
}

void DockLayout::setMargin(DockLayout::MarginEdge edge, qreal margin)
{
    switch(edge)
    {
    case DockLayout::LeftMargin:
        this->leftMargin = margin;
        break;
    case DockLayout::RightMargin:
        this->rightMargin = margin;
        break;
    case DockLayout::TopMargin:
        this->topMargin = margin;
        break;
    case DockLayout::BottomMargin:
        this->bottomMargin = margin;
        break;
    default:
        break;
    }
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

    appList.at(0)->move(0,0);

    for (int i = 1; i < appList.count(); i ++)
    {
        AppItem * frontItem = appList.at(i - 1);
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
        AppItem *fromItem = appList.at(i - 1);
        AppItem *toItem = appList.at(i);
        toItem->move(fromItem->x() - itemSpacing - toItem->width(),0);
    }
}

void DockLayout::sortTopToBottom()
{

}

void DockLayout::sortBottomToTop()
{

}

int DockLayout::indexOf(AppItem *item)
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
    case TopToBottom:
        sortTopToBottom();
        break;
    case BottomToTop:
        sortBottomToTop();
        break;
    default:
        break;
    }
}

void DockLayout::dragoutFromLayout(int index)
{
    AppItem * tmpItem = appList.takeAt(index);
    tmpItem->setVisible(false);
    tmpAppMap.insert(tmpItem,index);

    if (index == appList.count())//note,target hast been remove before
    {
        return;//at the end of list
    }

    //move follow item,note,target hast been remove before
    AppItem * followItem = appList.at(index);
    followItem->setNextPos(followItem->x() - tmpItem->width() - itemSpacing,0);
    //move last item
    for (int i = index + 1; i < appList.count(); i ++)
    {
        AppItem * frontItem = appList.at(i - 1);
        AppItem * targetItem = appList.at(i);
        targetItem->setNextPos(frontItem->getNextPos().x() + frontItem->width() + itemSpacing,0);
    }

    for (int i = index; i < appList.count(); i ++)
    {
        AppItem *button= appList.at(i);
        QPropertyAnimation *animation = new QPropertyAnimation(button, "pos");
        animation->setStartValue(button->pos());
        animation->setEndValue(button->getNextPos());
        animation->setDuration(500 + i * 100);
        animation->setEasingCurve(QEasingCurve::InOutBack);

        animation->start();
    }

}

void DockLayout::dragEnterEvent(QDragEnterEvent *event)
{
    event->setDropAction(Qt::MoveAction);
    event->accept();
}

void DockLayout::dropEvent(QDropEvent *event)
{
    AppItem * tmpItem = tmpAppMap.firstKey();
    tmpAppMap.remove(tmpItem);
    tmpItem->setVisible(true);
    if (indexOf(tmpItem) == -1)
    {
        insertItem(tmpItem,lastHoverIndex);
    }

    emit itemDropped();
}

void DockLayout::slotItemDrag(AppItem *item)
{
//    qWarning() << "Item draging..."<<x<<y<<item;
    int tmpIndex = indexOf(item);
    if (tmpIndex != -1)
    {
        dragoutFromLayout(tmpIndex);

        emit dragStarted();
    }
}

void DockLayout::slotItemRelease(int x, int y, AppItem *item)
{
    //outside frame,destroy it
    //inside frame,insert it
    item->setVisible(true);
    if (indexOf(item) == -1)
    {
        insertItem(item,lastHoverIndex);
    }
}

void DockLayout::slotItemEntered(QDragEnterEvent * event,AppItem *item)
{
    this->lastHoverIndex = indexOf(item);
}

void DockLayout::slotItemExited(QDragLeaveEvent *event,AppItem *item)
{

}
