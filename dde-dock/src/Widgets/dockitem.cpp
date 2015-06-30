#include "dockitem.h"

DockItem::DockItem(QWidget *parent) :
    QFrame(parent)
{
}

QWidget * DockItem::getContents()
{
    return NULL;
}

void DockItem::setTitle(const QString &title)
{
    this->itemTitle = title;
}

void DockItem::setIcon(const QString &iconPath, int size)
{
    appIcon = new AppIcon(iconPath,this);
    appIcon->resize(size,size);
    appIcon->move(this->width() / 2, this->height() / 2);
}

void DockItem::setActived(bool value)
{
    this->itemActived = value;
}

bool DockItem::actived()
{
    return this->itemActived;
}

void DockItem::setMoveable(bool value)
{
    this->itemMoveable = value;
}

bool DockItem::moveable()
{
    return this->itemMoveable;
}

void DockItem::setIndex(int value)
{
    this->itemIndex = value;
}

int DockItem::index()
{
    return this->itemIndex;
}
