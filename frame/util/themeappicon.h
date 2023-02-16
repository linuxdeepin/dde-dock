// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef THEMEAPPICON_H
#define THEMEAPPICON_H

#include <QObject>
#include <QIcon>

class ThemeAppIcon
{
public:
    explicit ThemeAppIcon();
    ~ThemeAppIcon();

    static QIcon getIcon(const QString &name);
    static bool getIcon(QPixmap &pix, const QString iconName, const int size, bool reObtain = false);

private:
    static bool createCalendarIcon(const QDate &date, const QString &fileName);
};

#endif // THEMEAPPICON_H
