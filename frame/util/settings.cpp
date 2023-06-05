// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "settings.h"

#include <QSharedPointer>
#include <QDebug>

Settings::Settings()
{

}

Settings::~Settings()
{

}

DConfig *Settings::ConfigPtr(const QString &name, const QString &subpath, QObject *parent)
{
    DConfig *config = DConfig::create("dde-dock", name, subpath, parent);
    if (!config)
        return nullptr;

    if (config->isValid())
        return config;

    delete config;
    qDebug() << "Cannot find dconfigs, name:" << name;
    return nullptr;
}

const QVariant Settings::ConfigValue(const QString &name, const QString &subPath, const QString &key, const QVariant &fallback)
{
    QSharedPointer<DConfig> config(ConfigPtr(name, subPath));
    if (config && config->isValid() && config->keyList().contains(key)) {
        QVariant v = config->value(key);
        return v;
    }

    qDebug() << "Cannot find dconfigs, name:" << name
             << " subPath:" << subPath << " key:" << key
             << "Use fallback value:" << fallback;
    return fallback;
}

bool Settings::ConfigSaveValue(const QString &name, const QString &subPath, const QString &key, const QVariant &value)
{
    QSharedPointer<DConfig> config(ConfigPtr(name, subPath));
    if (config && config->isValid() && config->keyList().contains(key)) {
        config->setValue(key, value);
        return true;
    }

    qDebug() << "Cannot find dconfigs, name:" << name
             << " subPath:" << subPath << " key:" << key;
    return false;
}
