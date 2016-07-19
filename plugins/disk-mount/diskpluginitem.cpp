#include "diskpluginitem.h"

#include <QPainter>
#include <QDebug>

DiskPluginItem::DiskPluginItem(QWidget *parent)
    : QWidget(parent),
      m_displayMode(Dock::Efficient)
{
}

void DiskPluginItem::setDockDisplayMode(const Dock::DisplayMode mode)
{
    m_displayMode = mode;

    updateIcon();
}

void DiskPluginItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void DiskPluginItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    updateIcon();
}

QSize DiskPluginItem::sizeHint() const
{
    return QSize(24, 24);
}

void DiskPluginItem::updateIcon()
{
    if (m_displayMode == Dock::Efficient)
        m_icon = QPixmap(":/icons/resources/icon_small.png");
    else
        m_icon = QPixmap(":/icons/resources/icon.png");

    update();
}
