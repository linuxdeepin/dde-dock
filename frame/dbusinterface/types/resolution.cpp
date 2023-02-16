// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "resolution.h"

#include <QDebug>

void registerResolutionMetaType()
{
    qRegisterMetaType<Resolution>("Resolution");
    qDBusRegisterMetaType<Resolution>();
}

Resolution::Resolution()
{

}

bool Resolution::operator!=(const Resolution &other) const
{
    return m_width != other.m_width || m_height != other.m_height || m_rate != other.m_rate;
}

bool Resolution::operator==(const Resolution &other) const
{
    return !(other != *this);
}

QDBusArgument &operator<<(QDBusArgument &arg, const Resolution &value)
{
    arg.beginStructure();
    arg << quint32(value.id()) << quint16(value.width()) << quint16(value.height()) << value.rate();
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, Resolution &value)
{
    quint32 id;
    quint16 w, h;
    double rate;

    arg.beginStructure();
    arg >> id >> w >> h >> rate;
    arg.endStructure();

    value.setId(id);
    value.setWidth(w);
    value.setHeight(h);
    value.setRate(rate);

    return arg;
}
