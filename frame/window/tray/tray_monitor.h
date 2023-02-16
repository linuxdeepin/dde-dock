// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRAYMONITOR_H
#define TRAYMONITOR_H

#include <QObject>

#include "dbustraymanager.h"
#include "statusnotifierwatcher_interface.h"

class PluginsItemInterface;

using namespace org::kde;

class TrayMonitor : public QObject
{
    Q_OBJECT

public:
    explicit TrayMonitor(QObject *parent = nullptr);

    QList<quint32> trayWinIds() const;
    QStringList sniServices() const;
    QStringList indicatorNames() const;
    QList<PluginsItemInterface *> systemTrays() const;

Q_SIGNALS:
    void requestUpdateIcon(quint32);
    void xEmbedTrayAdded(quint32);
    void xEmbedTrayRemoved(quint32);

    void sniTrayAdded(const QString &);
    void sniTrayRemoved(const QString &);

    void systemTrayAdded(PluginsItemInterface *);
    void systemTrayRemoved(PluginsItemInterface *);

    void indicatorFounded(const QString &);

public Q_SLOTS:
    void onTrayIconsChanged();
    void onSniItemsChanged();

    void startLoadIndicators();

private:
    DBusTrayManager *m_trayInter;
    StatusNotifierWatcher *m_sniWatcher;

    QList<quint32> m_trayWids;
    QStringList m_sniServices;
    QStringList m_indicatorNames;
    QList<PluginsItemInterface *> m_systemTrays;
};

#endif // TRAYMONITOR_H
