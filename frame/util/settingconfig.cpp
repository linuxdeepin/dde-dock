/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
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
