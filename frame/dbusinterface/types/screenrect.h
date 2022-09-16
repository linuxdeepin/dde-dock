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

#ifndef SCREENRECT_H
#define SCREENRECT_H

#include <QRect>
#include <QDBusArgument>
#include <QDebug>
#include <QDBusMetaType>

struct ScreenRect
{
public:
    ScreenRect();
    operator QRect() const;

    friend QDebug operator<<(QDebug debug, const ScreenRect &rect);
    friend const QDBusArgument &operator>>(const QDBusArgument &arg, ScreenRect &rect);
    friend QDBusArgument &operator<<(QDBusArgument &arg, const ScreenRect &rect);

private:
    qint16 x;
    qint16 y;
    quint16 w;
    quint16 h;
};

Q_DECLARE_METATYPE(ScreenRect)

void registerScreenRectMetaType();

#endif // SCREENRECT_H
