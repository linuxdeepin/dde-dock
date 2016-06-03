#include "dockitem.h"

DockItem::DockItem(QWidget *parent)
    : QWidget(parent)
{

}

void DockItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
}
