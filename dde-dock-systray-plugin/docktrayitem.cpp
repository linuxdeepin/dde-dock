#include <QHBoxLayout>

#include "docktrayitem.h"

DockTrayItem * DockTrayItem::fromWinId(WId winId, QWidget *parent)
{
    DockTrayItem *item = new DockTrayItem(parent);

    QWindow *win = QWindow::fromWinId(winId);
    QWidget *child = QWidget::createWindowContainer(win, item);

    QHBoxLayout *layout = new QHBoxLayout(item);
    layout->addWidget(child);
    item->setLayout(layout);

    return item;
}

DockTrayItem::DockTrayItem(QWidget *parent)
    : AbstractDockItem(parent)
{
    setFixedSize(32, 32);
}

DockTrayItem::~DockTrayItem()
{

}

void DockTrayItem::setTitle(const QString &)
{

}

void DockTrayItem::setIcon(const QString &, int)
{

}

void DockTrayItem::setMoveable(bool)
{

}

bool DockTrayItem::moveable()
{
    return false;
}

void DockTrayItem::setActived(bool)
{

}

bool DockTrayItem::actived()
{
    return false;
}

void DockTrayItem::setIndex(int value)
{
    m_itemIndex = value;
}

int DockTrayItem::index()
{
    return m_itemIndex;
}

QWidget * DockTrayItem::getContents()
{
    return NULL;
}
