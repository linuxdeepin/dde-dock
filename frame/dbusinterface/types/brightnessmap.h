// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef BRIGHTNESSMAP_H
#define BRIGHTNESSMAP_H

#include <QMap>
#include <QDBusMetaType>

typedef QMap<QString, double> BrightnessMap;

void registerBrightnessMapMetaType();

#endif // BRIGHTNESSMAP_H
