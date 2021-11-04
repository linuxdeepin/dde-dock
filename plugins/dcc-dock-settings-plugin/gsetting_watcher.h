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
#ifndef GSETTINGWATCHER_H
#define GSETTINGWATCHER_H

#include <QObject>
#include <QHash>
#include <QMap>

class QGSettings;
class QListView;
class QStandardItem;

namespace dcc_dock_plugin {
class GSettingWatcher : public QObject
{
    Q_OBJECT

public:
    GSettingWatcher(const QString &baseSchemasId, const QString &module, QObject *parent = nullptr);
    ~GSettingWatcher();

    void bind(const QString &key, QWidget *binder);

private:
    void setStatus(const QString &key, QWidget *binder);
    void onStatusModeChanged(const QString &key);

private:
    QMultiHash<QString, QWidget *> m_map;
    QGSettings *m_gsettings;
};
}

#endif // GSETTINGWATCHER_H
