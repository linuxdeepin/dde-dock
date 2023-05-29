// SPDX-FileCopyrightText: 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once

#include <QObject>

#undef signals
#include <gio/gio.h>
#define signals Q_SIGNALS

class TrashHelper: public QObject
{
    Q_OBJECT

public:
    explicit TrashHelper(QObject * parent);
    ~TrashHelper();

    int trashItemCount();
    bool emptyTrash();

Q_SIGNALS:
    void trashAttributeChanged();

private:
    GFile * m_trash;
    GFileMonitor * m_trashMonitor;

    void onTrashMonitorChanged(GFileMonitor *monitor, GFile *file, GFile *other_file, GFileMonitorEvent event_type);
    static void slot_onTrashMonitorChanged(GFileMonitor *monitor, GFile *file, GFile *other_file, GFileMonitorEvent event_type, gpointer user_data);
};
