#include "containeritem.h"

#include <QPainter>

ContainerItem::ContainerItem(QWidget *parent)
    : DockItem(parent),
      m_icon(":/indicator/resources/arrow_up_normal.png"),
      m_containerWidget(new ContainerWidget(this))
{
    m_containerWidget->setVisible(false);

    setAcceptDrops(true);
}

void ContainerItem::addItem(DockItem * const item)
{
    m_containerWidget->addWidget(item);
}

void ContainerItem::removeItem(DockItem * const item)
{
    m_containerWidget->removeWidget(item);
}

bool ContainerItem::contains(DockItem * const item)
{
    if (m_containerWidget->itemList().contains(item))
    {
        item->setParent(m_containerWidget);
        return true;
    }

    return false;
}

void ContainerItem::dragEnterEvent(QDragEnterEvent *e)
{
    if (m_containerWidget->allowDragEnter(e))
        return e->accept();
}

void ContainerItem::dragMoveEvent(QDragMoveEvent *e)
{
    Q_UNUSED(e);

    return;
}

void ContainerItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (DockDisplayMode == Dock::Fashion)
        return;

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void ContainerItem::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton && m_containerWidget->itemCount())
        return showPopupApplet(m_containerWidget);

    return DockItem::mouseReleaseEvent(e);
}

QSize ContainerItem::sizeHint() const
{
    return QSize(24, 24);
}
