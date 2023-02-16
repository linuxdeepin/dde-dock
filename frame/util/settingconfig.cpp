// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "settingconfig.h"

#include <DConfig>

DCORE_USE_NAMESPACE

SettingConfig *SettingConfig::instance()
{
    static SettingConfig instance;
    return &instance;
}

void SettingConfig::setValue(const QString &key, const QVariant &value)
{
    if (m_config->isValid() && m_config->keyList().contains(key))
        m_config->setValue(key, value);
}

QVariant SettingConfig::value(const QString &key) const
{
    if (m_config->isValid() && m_config->keyList().contains(key))
        return m_config->value(key);

    return QVariant();
}

SettingConfig::SettingConfig(QObject *parent)
    : QObject(parent)
    , m_config(new DConfig(QString("com.deepin.dde.dock.dconfig"), QString()))
{
    connect(m_config, &DConfig::valueChanged, this, &SettingConfig::onValueChanged);
}

void SettingConfig::onValueChanged(const QString &key)
{
    Q_EMIT valueChanged(key, m_config->value(key));
}
