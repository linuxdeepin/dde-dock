/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "containeritem.h"

#include <QPainter>

ContainerItem::ContainerItem(QWidget *parent)
    : DockItem(parent),
      m_dropping(false),
      m_popupTips(new QLabel),
      m_containerWidget(new ContainerWidget(this))
{
    m_containerWidget->setVisible(false);
    m_popupTips->setText(tr("Click to display hidden icon"));
    m_popupTips->setVisible(false);

    setAcceptDrops(true);
}

void ContainerItem::setDropping(const bool dropping)
{
    if (dropping)
        showPopupApplet(m_containerWidget);
    else
        hidePopup();

    m_dropping = dropping;
    update();
}

void ContainerItem::addItem(DockItem * const item)
{
    m_containerWidget->addWidget(item);
    item->setVisible(true);
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

    if (!m_containerWidget->itemCount() && !m_dropping)
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

QWidget *ContainerItem::popupTips()
{
    return m_popupTips;
}
