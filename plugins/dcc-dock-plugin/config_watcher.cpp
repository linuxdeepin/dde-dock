/*
 * Copyright (C) 2011 ~ 2021 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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
#include "config_watcher.h"
#include "utils.h"

#include <QGSettings>
#include <QListView>
#include <QStandardItem>
#include <QStandardItemModel>
#include <QVariant>
#include <QWidget>

#include <DConfig>

using namespace dcc_dock_plugin;

DCORE_USE_NAMESPACE
/**
 * @brief GSettingWatcher::GSettingWatcher 用于监听处于 \a baseSchemasId + "." + \a module 配置下的配置项内容变化，并将变化应用到绑定的控件上
 */
ConfigWatcher::ConfigWatcher(const QString &fileName, QObject *parent)
    : QObject(parent)
    , m_config(new DConfig(fileName, QString(), this))
{
    if (m_config->isValid()) {
        connect(m_config, &DConfig::valueChanged, this, &ConfigWatcher::onStatusModeChanged);
    } else {
        qWarning() << "config parse failed:" << fileName;
    }
}

ConfigWatcher::~ConfigWatcher()
{
    m_map.clear();
}

void ConfigWatcher::bind(const QString &key, QWidget *binder)
{
    m_map.insert(key, binder);

    setStatus(key, binder);
    // 自动解绑
    connect(binder, &QObject::destroyed, this, [=] {
        m_map.remove(m_map.key(binder), binder);
    });
}

void ConfigWatcher::setStatus(const QString &key, QWidget *binder)
{
    if (!binder || !m_config->isValid() || !m_config->keyList().contains(key))
        return;

    const QString setting = m_config->value(key).toString();

    if ("Enabled" == setting) {
        binder->setEnabled(true);
    } else if ("Disabled" == setting) {
        binder->setEnabled(false);
    }

    binder->setVisible("Hidden" != setting);
}

void ConfigWatcher::onStatusModeChanged(const QString &key)
{
    if (!m_map.isEmpty() && m_map.contains(key)) {
        for (auto it = m_map.begin(); it != m_map.end(); ++it) {
            if (key == it.key()) {
                setStatus(key, it.value());
            }
        }
    }
}
