#include "placeholderitem.h"

PlaceholderItem::PlaceholderItem(QWidget *parent)
    : DockItem(Placeholder, parent)
{
    setBaseSize(0, 0);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
}

void PlaceholderItem::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);
}

void PlaceholderItem::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);
}
