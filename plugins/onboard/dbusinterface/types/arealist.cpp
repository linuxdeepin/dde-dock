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

#include "arealist.h"

bool MonitRect::operator ==(const MonitRect &rect)
{
    return x1 == rect.x1 && y1 == rect.y1 && x2 == rect.x2 && y2 == rect.y2;
}

QDBusArgument &operator<<(QDBusArgument &arg, const MonitRect &rect)
{
    arg.beginStructure();
    arg << rect.x1 << rect.y1 << rect.x2 << rect.y2;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, MonitRect &rect)
{
    arg.beginStructure();
    arg >> rect.x1 >> rect.y1 >> rect.x2 >> rect.y2;
    arg.endStructure();

    return arg;
}

void registerAreaListMetaType()
{
    qRegisterMetaType<MonitRect>("MonitRect");
    qDBusRegisterMetaType<MonitRect>();

    qRegisterMetaType<AreaList>("AreaList");
    qDBusRegisterMetaType<AreaList>();
}
