/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *             zhaolong <zhaolong@uniontech.com>
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

#ifndef AREALIST_H
#define AREALIST_H

#include <QDBusMetaType>
#include <QRect>
#include <QList>

struct MonitRect {
    int x1;
    int y1;
    int x2;
    int y2;

    bool operator ==(const MonitRect& rect);
};

typedef QList<MonitRect> AreaList;

Q_DECLARE_METATYPE(MonitRect)
Q_DECLARE_METATYPE(AreaList)

QDBusArgument &operator<<(QDBusArgument &arg, const MonitRect &rect);
const QDBusArgument &operator>>(const QDBusArgument &arg, MonitRect &rect);

void registerAreaListMetaType();

#endif // AREALIST_H
