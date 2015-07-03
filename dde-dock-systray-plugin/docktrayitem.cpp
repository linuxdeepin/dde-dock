#include <QHBoxLayout>

#include "docktrayitem.h"

DockTrayItem * DockTrayItem::fromWinId(WId winId)
{
    DockTrayItem *item = new DockTrayItem;

    QWindow *win = QWindow::fromWinId(winId);
    QWidget *child = QWidget::createWindowContainer(win, item);
    child->setFixedSize(item->size());
    child->setParent(item);

    return item;
}

DockTrayItem::DockTrayItem(QWidget *parent)
    : QWidget(parent)
{
    setFixedSize(16, 16);
}

DockTrayItem::~DockTrayItem()
{

}
