// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef GSETTINGWATCHER_H
#define GSETTINGWATCHER_H

#include <dtkcore_global.h>

#include <QObject>
#include <QHash>
#include <QMap>

class QGSettings;
class QListView;
class QStandardItem;

DCORE_BEGIN_NAMESPACE
class DConfig;
DCORE_END_NAMESPACE

namespace dcc_dock_plugin {
class ConfigWatcher : public QObject
{
    Q_OBJECT

public:
    ConfigWatcher(const QString &appId, const QString &fileName, QObject *parent = nullptr);
    ~ConfigWatcher();

    void bind(const QString &key, QWidget *binder);

private:
    void setStatus(const QString &key, QWidget *binder);
    void onStatusModeChanged(const QString &key);

private:
    QMultiHash<QString, QWidget *> m_map;
    DTK_CORE_NAMESPACE::DConfig *m_config;
};
}

#endif // GSETTINGWATCHER_H
