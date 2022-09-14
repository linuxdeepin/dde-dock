// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
ConfigWatcher::ConfigWatcher(const QString &appId, const QString &fileName, QObject *parent)
    : QObject(parent)
    , m_config(DConfig::create(appId, fileName, QString(), this))
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
