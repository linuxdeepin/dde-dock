// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef RESOLUTION_H
#define RESOLUTION_H

#include <QDBusMetaType>

class Resolution
{
public:
    friend QDBusArgument &operator<<(QDBusArgument &arg, const Resolution &value);
    friend const QDBusArgument &operator>>(const QDBusArgument &arg, Resolution &value);

    explicit Resolution();

    bool operator!=(const Resolution &other) const;
    bool operator==(const Resolution &other) const;

    int id() const { return m_id; }
    int width() const { return m_width; }
    int height() const { return m_height; }
    double rate() const { return m_rate; }

private:
    void setId(const int id) { m_id = id; }
    void setWidth(const int w) { m_width = w; }
    void setHeight(const int h) { m_height = h; }
    void setRate(const double rate) { m_rate = rate; }

private:
    int m_id;
    int m_width;
    int m_height;
    double m_rate;
};


Q_DECLARE_METATYPE(Resolution)

void registerResolutionMetaType();

#endif // RESOLUTION_H
