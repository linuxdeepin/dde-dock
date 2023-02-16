// Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
