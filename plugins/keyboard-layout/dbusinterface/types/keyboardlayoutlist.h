// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef KEYBOARDLAYOUTLIST_H
#define KEYBOARDLAYOUTLIST_H

#include <QMap>
#include <QString>
#include <QObject>
#include <QDBusMetaType>

typedef QMap<QString, QString> KeyboardLayoutList;

void registerKeyboardLayoutListMetaType();

#endif // KEYBOARDLAYOUTLIST_H
