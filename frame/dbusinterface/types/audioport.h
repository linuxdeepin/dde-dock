/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

#ifndef AUDIOPORT_H
#define AUDIOPORT_H

#include <QDBusMetaType>
#include <QString>
#include <QDBusArgument>
#include <QDebug>

class AudioPort
{
public:
    QString name;
    QString description;
    uchar availability; // 0 for Unknown, 1 for Not Available, 2 for Available.

    friend QDebug operator<<(QDebug argument, const AudioPort &port) {
        argument << port.description;

        return argument;
    }

    friend QDBusArgument &operator<<(QDBusArgument &argument, const AudioPort &port) {
        argument.beginStructure();
        argument << port.name << port.description << port.availability;
        argument.endStructure();

        return argument;
    }

    friend const QDBusArgument &operator>>(const QDBusArgument &argument, AudioPort &port) {
        argument.beginStructure();
        argument >> port.name >> port.description >> port.availability;
        argument.endStructure();

        return argument;
    }

    bool operator==(const AudioPort what) const {
        return what.name == name && what.description == description && what.availability == availability;
    }

    bool operator!=(const AudioPort what) const {
        return what.name != name || what.description != description || what.availability != availability;
    }
};

Q_DECLARE_METATYPE(AudioPort)

void registerAudioPortMetaType();

#endif // AUDIOPORT_H
