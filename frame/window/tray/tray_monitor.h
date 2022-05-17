/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef TRAYMONITOR_H
#define TRAYMONITOR_H

#include <QObject>

#include "dbustraymanager.h"
#include "statusnotifierwatcher_interface.h"

using namespace org::kde;
class TrayMonitor : public QObject
{
    Q_OBJECT

public:
    explicit TrayMonitor(QObject *parent = nullptr);

Q_SIGNALS:
    void requestUpdateIcon(quint32);
    void xEmbedTrayAdded(quint32);
    void xEmbedTrayRemoved(quint32);

    void sniTrayAdded(const QString &);
    void sniTrayRemoved(const QString &);

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
};

#endif // TRAYMONITOR_H
