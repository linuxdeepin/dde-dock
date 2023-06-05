// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef SETTINGS_H
#define SETTINGS_H

#include <DConfig>

#include <QObject>
#include <QString>

DCORE_USE_NAMESPACE

// Dconfig 配置类
class Settings
{
public:
    Settings();
    ~Settings();

    static DConfig *ConfigPtr(const QString &name, const QString &subpath = QString(), QObject *parent = nullptr);
    static const QVariant ConfigValue(const QString &name, const QString &subPath = QString(), const QString &key = QString(), const QVariant &fallback = QVariant());
    static bool ConfigSaveValue(const QString &name, const QString &subPath, const QString &key, const QVariant &value);
};

#endif // SETTINGS_H
