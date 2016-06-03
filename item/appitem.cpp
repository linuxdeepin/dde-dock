#include "appitem.h"

#include <QPainter>

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(parent),
      m_itemEntry(new DBusDockEntry(entry.path(), this))
{

}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), Qt::cyan);
    painter.drawText(rect(), m_itemEntry->id());
}
