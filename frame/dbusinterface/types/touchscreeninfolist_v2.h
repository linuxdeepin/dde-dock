// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
