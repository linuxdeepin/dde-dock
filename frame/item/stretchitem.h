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

#ifndef STRETCHITEM_H
#define STRETCHITEM_H

#include "dockitem.h"

class StretchItem : public DockItem
{
    Q_OBJECT

public:
    explicit StretchItem(QWidget *parent = 0);

    inline ItemType itemType() const {return Stretch;}

private:
    void mousePressEvent(QMouseEvent *e);
};

#endif // STRETCHITEM_H
