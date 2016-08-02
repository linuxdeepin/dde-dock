#include "stretchitem.h"

#include <QPaintEvent>

StretchItem::StretchItem(QWidget *parent)
    : DockItem(Stretch, parent)
{
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
}

void StretchItem::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);
}
