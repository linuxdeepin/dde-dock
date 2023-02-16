// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SETTINGCONFIG_H
#define SETTINGCONFIG_H

#include <dtkcore_global.h>

#include <QObject>
#include <QVariant>

DCORE_BEGIN_NAMESPACE
class DConfig;
DCORE_END_NAMESPACE

DCORE_USE_NAMESPACE

#define SETTINGCONFIG SettingConfig::instance()

class SettingConfig : public QObject
{
    Q_OBJECT

public:
    static SettingConfig *instance();

    void setValue(const QString &key, const QVariant &value);
    QVariant value(const QString &key) const;

Q_SIGNALS:
    void valueChanged(const QString &key, const QVariant &value);

protected:
    explicit SettingConfig(QObject *parent = nullptr);

private Q_SLOTS:
    void onValueChanged(const QString &key);

private:
    DConfig *m_config;
};

#endif // SETTINGCONFIG_H
