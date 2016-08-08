#include "constants.h"
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

void ContainerItem::dragEnterEvent(QDragEnterEvent *e)
{
    if (!e->mimeData()->hasFormat(DOCK_PLUGIN_MIME))
        return;

    e->accept();
}

void ContainerItem::dragMoveEvent(QDragMoveEvent *e)
{
    Q_UNUSED(e);

    return;
}

void ContainerItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

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
