/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

#ifndef ADAPTERSMANAGER_H
#define ADAPTERSMANAGER_H

#include "adapter.h"

#include <com_deepin_daemon_bluetooth.h>
using  DBusBluetooth = com::deepin::daemon::Bluetooth;

class AdaptersManager : public QObject
{
    Q_OBJECT
public:
    explicit AdaptersManager(QObject *parent = nullptr);

    QMap<QString, const Adapter *> adapters() const;
    const Adapter *adapterById(const QString &id);

    void setAdapterPowered(const Adapter *adapter, const bool &powered);
    void connectAllPairedDevice(const Adapter *adapter);
    void connectDevice(Device *deviceId);
    bool defaultAdapterInitPowerState();
    int adaptersCount();

signals:
    void adapterIncreased(Adapter *adapter);
    void adapterDecreased(Adapter *adapter);

private slots:
    void onAdapterPropertiesChanged(const QString &json);
    void onDevicePropertiesChanged(const QString &json);

    void addAdapter(const QString &json);
    void removeAdapter(const QString &json);

    void addDevice(const QString &json);
    void removeDevice(const QString &json);

private:
    void adapterAdd(Adapter *adapter, const QJsonObject &adpterObj);
    void inflateAdapter(Adapter *adapter, const QJsonObject &adapterObj);

private:
    DBusBluetooth *m_bluetoothInter;
    QMap<QString, const Adapter *> m_adapters;

    bool m_defaultAdapterState;
};

#endif // ADAPTERSMANAGER_H
