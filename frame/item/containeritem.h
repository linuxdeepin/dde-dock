/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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

#ifndef CONTAINERITEM_H
#define CONTAINERITEM_H

#include "dockitem.h"
#include "components/containerwidget.h"
#include "../widgets/tipswidget.h"

#include <QPixmap>

class ContainerItem : public DockItem
{
    Q_OBJECT

public:
    explicit ContainerItem(QWidget *parent = 0);

    inline ItemType itemType() const {return Container;}

    void setDropping(const bool dropping);
    void addItem(DockItem * const item);
    void removeItem(DockItem * const item);
    bool contains(DockItem * const item);

public slots:
    void refershIcon();

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    QSize sizeHint() const;
    QWidget *popupTips();

private:
    bool m_dropping;
    TipsWidget *m_popupTips;
    ContainerWidget *m_containerWidget;
    QPixmap m_icon;
};

#endif // CONTAINERITEM_H
