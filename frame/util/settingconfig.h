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
