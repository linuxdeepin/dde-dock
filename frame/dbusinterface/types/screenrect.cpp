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

#include "screenrect.h"

ScreenRect::ScreenRect()
    : x(0),
      y(0),
      w(0),
      h(0)
{

}

QDebug operator<<(QDebug debug, const ScreenRect &rect)
{
    debug << QString("ScreenRect(%1, %2, %3, %4)").arg(rect.x)
                                                    .arg(rect.y)
                                                    .arg(rect.w)
                                                    .arg(rect.h);

    return debug;
}

ScreenRect::operator QRect() const
{
    return QRect(x, y, w, h);
}

QDBusArgument &operator<<(QDBusArgument &arg, const ScreenRect &rect)
{
    arg.beginStructure();
    arg << rect.x << rect.y << rect.w << rect.h;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, ScreenRect &rect)
{
    arg.beginStructure();
    arg >> rect.x >> rect.y >> rect.w >> rect.h;
    arg.endStructure();

    return arg;
}

void registerScreenRectMetaType()
{
    qRegisterMetaType<ScreenRect>("ScreenRect");
    qDBusRegisterMetaType<ScreenRect>();
}
