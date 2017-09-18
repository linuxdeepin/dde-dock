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

#include "constants.h"
#include "containerwidget.h"
#include "item/pluginsitem.h"

#include <QDebug>
#include <QDragEnterEvent>

#define ITEM_HEIGHT         30
#define ITEM_WIDTH          30

ContainerWidget::ContainerWidget(QWidget *parent)
    : QWidget(parent),

      m_centralLayout(new QHBoxLayout)
{
    m_centralLayout->addStretch();
    m_centralLayout->setSpacing(0);
    m_centralLayout->setMargin(0);

    setLayout(m_centralLayout);
    setFixedHeight(ITEM_HEIGHT);
    setFixedWidth(ITEM_WIDTH);
    setAcceptDrops(true);
}

void ContainerWidget::addWidget(QWidget * const w)
{
    w->setParent(this);
    w->setFixedSize(ITEM_WIDTH, ITEM_HEIGHT);
    m_centralLayout->addWidget(w);
    m_itemList.append(w);

    setFixedWidth(ITEM_WIDTH * std::max(1, m_itemList.size()));
}

void ContainerWidget::removeWidget(QWidget * const w)
{
    m_centralLayout->removeWidget(w);
    m_itemList.removeOne(w);

    setFixedWidth(ITEM_WIDTH * std::max(1, m_itemList.size()));
}

int ContainerWidget::itemCount() const
{
    return m_itemList.count();
}

const QList<QWidget *> ContainerWidget::itemList() const
{
    return m_itemList;
}

bool ContainerWidget::allowDragEnter(QDragEnterEvent *e)
{
    if (!e->mimeData()->hasFormat(DOCK_PLUGIN_MIME))
        return false;

    PluginsItem *pi = static_cast<PluginsItem *>(e->source());
    if (pi && pi->allowContainer())
        return true;

    return false;
}

void ContainerWidget::dragEnterEvent(QDragEnterEvent *e)
{
    if (allowDragEnter(e))
        return e->accept();
}
