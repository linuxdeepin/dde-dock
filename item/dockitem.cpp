#include "dockitem.h"

DockItem::DockItem(const ItemType type, QWidget *parent)
    : QWidget(parent),
      m_side(DockSettings::Top),
      m_type(type)
{
}

void DockItem::setDockSide(const DockSettings::DockSide side)
{
    m_side = side;

    update();
}

DockItem::ItemType DockItem::itemType() const
{
    return m_type;
}

void DockItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
}
