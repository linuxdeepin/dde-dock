/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     quezhiyong <quezhiyong@uniontech.com>
 *
 * Maintainer: quezhiyong <quezhiyong@uniontech.com>
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

#ifndef TOUCHSCREENINFOLISTV2_H
#define TOUCHSCREENINFOLISTV2_H

#include <QString>
#include <QList>
#include <QDBusMetaType>

struct TouchscreenInfo_V2 {
    qint32 id;
    QString name;
    QString deviceNode;
    QString serialNumber;
    QString UUID;

    bool operator ==(const TouchscreenInfo_V2& info);
};

typedef QList<TouchscreenInfo_V2> TouchscreenInfoList_V2;

Q_DECLARE_METATYPE(TouchscreenInfo_V2)
Q_DECLARE_METATYPE(TouchscreenInfoList_V2)

QDBusArgument &operator<<(QDBusArgument &arg, const TouchscreenInfo_V2 &info);
const QDBusArgument &operator>>(const QDBusArgument &arg, TouchscreenInfo_V2 &info);

void registerTouchscreenInfoList_V2MetaType();

#endif // !TOUCHSCREENINFOLISTV2_H
