#include "containeritem.h"

#include <QPainter>

ContainerItem::ContainerItem(QWidget *parent)
    : DockItem(parent),
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

void ContainerItem::refershIcon()
{
    QPixmap icon;
    switch (DockPosition)
    {
    case Top:       icon = QPixmap(":/icons/resources/arrow-down.svg");     break;
    case Left:      icon = QPixmap(":/icons/resources/arrow-right.svg");    break;
    case Bottom:    icon = QPixmap(":/icons/resources/arrow-up.svg");       break;
    case Right:     icon = QPixmap(":/icons/resources/arrow-left.svg");     break;
    default:        Q_UNREACHABLE();
    }

    m_icon = icon;
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
    if (e->button() == Qt::LeftButton)
        return showPopupApplet(m_containerWidget);

    return DockItem::mouseReleaseEvent(e);
}

QSize ContainerItem::sizeHint() const
{
    return QSize(24, 24);
}
