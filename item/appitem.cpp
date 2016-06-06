#include "appitem.h"

#include <QPainter>

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(App, parent),
      m_itemEntry(new DBusDockEntry(entry.path(), this))
{
    qDebug() << m_itemEntry->data();
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    const QRect itemRect = rect();
    const int iconSize = std::min(itemRect.width(), itemRect.height());

    QRect iconRect;
    iconRect.setWidth(iconSize);
    iconRect.setHeight(iconSize);
    iconRect.moveTopLeft(itemRect.center() - iconRect.center());

    QPainter painter(this);
    painter.fillRect(rect(), Qt::cyan);
    painter.fillRect(iconRect, Qt::yellow);
    painter.drawText(rect(), m_itemEntry->id());
}
