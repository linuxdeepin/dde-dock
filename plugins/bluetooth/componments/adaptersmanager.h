// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ADAPTERSMANAGER_H
#define ADAPTERSMANAGER_H

#include <com_deepin_daemon_bluetooth.h>
using  DBusBluetooth = com::deepin::daemon::Bluetooth;

class Adapter;
class Device;
class AdaptersManager : public QObject
{
    Q_OBJECT
public:
    explicit AdaptersManager(QObject *parent = nullptr);

    void setAdapterPowered(const Adapter *adapter, const bool &powered);
    void connectDevice(const Device *device, Adapter *adapter);
    int adaptersCount();
    void adapterRefresh(const Adapter *adapter);

signals:
    void adapterIncreased(Adapter *adapter);
    void adapterDecreased(Adapter *adapter);

private slots:
    void onAdapterPropertiesChanged(const QString &json);
    void onDevicePropertiesChanged(const QString &json);

    void onAddAdapter(const QString &json);
    void onRemoveAdapter(const QString &json);

    void onAddDevice(const QString &json);
    void onRemoveDevice(const QString &json);

private:
    void adapterAdd(Adapter *adapter, const QJsonObject &adpterObj);
    void inflateAdapter(Adapter *adapter, const QJsonObject &adapterObj);

private:
    DBusBluetooth *m_bluetoothInter;
    QMap<QString, const Adapter *> m_adapters;
};

#endif // ADAPTERSMANAGER_H
