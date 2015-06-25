#include "docklayout.h"

DockLayout::DockLayout(QWidget *parent) :
    QWidget(parent)
{
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
    connect(item, SIGNAL(mouseMove(int,int,AppItem*)),this,SLOT(slotItemDrag(int,int,AppItem*)));

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

}

void DockLayout::sortTopToBottom()
{

}

void DockLayout::sortBottomToTop()
{

}

int DockLayout::indexOf(AppItem *item)
{
    return    appList.indexOf(item);
}

int DockLayout::indexOf(int x, int y)
{
    //TODO
    return 0;
}

void DockLayout::slotItemDrag(int x, int y, AppItem *item)
{
    qWarning() << "Item draging..."<<x<<y<<item;
}
