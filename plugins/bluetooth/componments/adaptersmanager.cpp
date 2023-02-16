// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "adaptersmanager.h"
#include "adapter.h"
#include "device.h"

#include <QDBusInterface>
#include <QDBusReply>
#include <QJsonDocument>
#include <QJsonArray>

AdaptersManager::AdaptersManager(QObject *parent)
    : QObject(parent)
    , m_bluetoothInter(new DBusBluetooth("org.deepin.dde.Bluetooth1",
                                         "/org/deepin/dde/Bluetooth1",
                                         QDBusConnection::sessionBus(),
                                         this))
{
    connect(m_bluetoothInter, &DBusBluetooth::AdapterAdded, this, &AdaptersManager::onAddAdapter);
    connect(m_bluetoothInter, &DBusBluetooth::AdapterRemoved, this, &AdaptersManager::onRemoveAdapter);
    connect(m_bluetoothInter, &DBusBluetooth::AdapterPropertiesChanged, this, &AdaptersManager::onAdapterPropertiesChanged);
    connect(m_bluetoothInter, &DBusBluetooth::DeviceAdded, this, &AdaptersManager::onAddDevice);
    connect(m_bluetoothInter, &DBusBluetooth::DeviceRemoved, this, &AdaptersManager::onRemoveDevice);
    connect(m_bluetoothInter, &DBusBluetooth::DevicePropertiesChanged, this, &AdaptersManager::onDevicePropertiesChanged);

#ifdef QT_DEBUG
    connect(m_bluetoothInter, &DBusBluetooth::RequestAuthorization, this, [](const QDBusObjectPath & in0) {
        qDebug() << "request authorization: " << in0.path();
    });

    connect(m_bluetoothInter, &DBusBluetooth::RequestPasskey, this, [](const QDBusObjectPath & in0) {
        qDebug() << "request passkey: " << in0.path();
    });

    connect(m_bluetoothInter, &DBusBluetooth::RequestPinCode, this, [](const QDBusObjectPath & in0) {
        qDebug() << "request pincode: " << in0.path();
    });

    connect(m_bluetoothInter, &DBusBluetooth::DisplayPasskey, this, [](const QDBusObjectPath & in0, uint in1, uint in2) {
        qDebug() << "request display passkey: " << in0.path() << in1 << in2;
    });

    connect(m_bluetoothInter, &DBusBluetooth::DisplayPinCode, this, [](const QDBusObjectPath & in0, const QString & in1) {
        qDebug() << "request display pincode: " << in0.path() << in1;
    });
#endif

    QDBusInterface inter("org.deepin.dde.Bluetooth1",
                         "/org/deepin/dde/Bluetooth1",
                         "org.deepin.dde.Bluetooth1",
                         QDBusConnection::sessionBus());
    QDBusReply<QString> reply = inter.call(QDBus::Block, "GetAdapters");
    const QString replyStr = reply.value();
    QJsonDocument doc = QJsonDocument::fromJson(replyStr.toUtf8());
    QJsonArray arr = doc.array();

    for (int index = 0; index < arr.size(); index++) {
        auto *adapter = new Adapter(this);
        adapterAdd(adapter, arr[index].toObject());
    }
}

void AdaptersManager::setAdapterPowered(const Adapter *adapter, const bool &powered)
{
    if (!adapter)
        return;

    QDBusObjectPath path(adapter->id());
    QDBusPendingCall call = m_bluetoothInter->SetAdapterPowered(path, powered);

    if (!powered) {
        QDBusPendingCall clearUnpairedCall = m_bluetoothInter->ClearUnpairedDevice();
        QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(clearUnpairedCall, this);
        connect(watcher, &QDBusPendingCallWatcher::finished, [ = ] {
            if (clearUnpairedCall.isError())
                qWarning() << clearUnpairedCall.error().message();
        });
    }
}

void AdaptersManager::connectDevice(const Device *device, Adapter *adapter)
{
    if (device) {
        QDBusObjectPath path(device->id());
        switch (device->state()) {
        case Device::StateUnavailable: {
            m_bluetoothInter->ConnectDevice(path, QDBusObjectPath(adapter->id()));
            qDebug() << "connect to device: " << device->alias();
        }
            break;
        case Device::StateAvailable:
            break;
        case Device::StateConnected: {
            m_bluetoothInter->DisconnectDevice(path);
            qDebug() << "disconnect device: " << device->alias();
        }
            break;
        }
    }
}

int AdaptersManager::adaptersCount()
{
    return m_adapters.size();
}

void AdaptersManager::onAdapterPropertiesChanged(const QString &json)
{
    const QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    const QJsonObject obj = doc.object();
    const QString id = obj["Path"].toString();
    QDBusObjectPath dPath(id);

    if (!m_adapters.contains(id)) {
        return;
    }

    Adapter *adapter = const_cast<Adapter *>(m_adapters[id]);
    if (adapter) {
        inflateAdapter(adapter, obj);
    }
}

void AdaptersManager::onDevicePropertiesChanged(const QString &json)
{
    const QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    const QJsonObject obj = doc.object();
    for (const Adapter *constAdapter : m_adapters) {
        auto adapter = const_cast<Adapter *>(constAdapter);
        if (adapter)
            adapter->updateDevice(obj);
    }
}

void AdaptersManager::onAddAdapter(const QString &json)
{
    const QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    auto adapter = new Adapter(this);
    adapterAdd(adapter, doc.object());
}

void AdaptersManager::onRemoveAdapter(const QString &json)
{
    QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    QJsonObject obj = doc.object();
    const QString id = obj["Path"].toString();

    if (!m_adapters.contains(id)) {
        return;
    }

    const Adapter *result = m_adapters[id];
    Adapter *adapter = const_cast<Adapter *>(result);
    if (adapter) {
        m_adapters.remove(id);
        m_adapterIds.removeOne(id);
        emit adapterDecreased(adapter);
        adapter->deleteLater();
    }
}

void AdaptersManager::onAddDevice(const QString &json)
{
    const QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    const QJsonObject obj = doc.object();
    const QString adapterId = obj["AdapterPath"].toString();
    const QString deviceId = obj["Path"].toString();

    if (!m_adapters.contains(adapterId)) {
        return;
    }

    const Adapter *result = m_adapters[adapterId];
    Adapter *adapter = const_cast<Adapter *>(result);
    if (adapter) {
        const Device *result1 = adapter->deviceById(deviceId);
        Device *device = const_cast<Device *>(result1);
        if (!device) {
            adapter->addDevice(obj);
        }
    }
}

void AdaptersManager::onRemoveDevice(const QString &json)
{
    QJsonDocument doc = QJsonDocument::fromJson(json.toUtf8());
    QJsonObject obj = doc.object();
    const QString adapterId = obj["AdapterPath"].toString();
    const QString deviceId = obj["Path"].toString();

    if (!m_adapters.contains(adapterId)) {
        return;
    }

    const Adapter *result = m_adapters[adapterId];
    Adapter *adapter = const_cast<Adapter *>(result);
    if (adapter) {
        adapter->removeDevice(deviceId);
    }
}

void AdaptersManager::adapterAdd(Adapter *adapter, const QJsonObject &adpterObj)
{
    if (!adapter)
        return;

    inflateAdapter(adapter, adpterObj);
    QDBusObjectPath dPath(adpterObj["Path"].toString());
    QDBusPendingCall call = m_bluetoothInter->GetDevices(dPath);
    QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(call, this);
    connect(watcher, &QDBusPendingCallWatcher::finished, watcher, &QDBusPendingCallWatcher::deleteLater);
    connect(watcher, &QDBusPendingCallWatcher::finished, [ this, adapter, call ] {
        if (!call.isError()) {
            QDBusReply<QString> reply = call.reply();
            const QString replyStr = reply.value();
            QJsonDocument doc = QJsonDocument::fromJson(replyStr.toUtf8());
            adapter->initDevicesList(doc);
            emit this->adapterIncreased(adapter);
        } else {
            qWarning() << call.error().message();
        }
    });

    QString id = adapter->id();
    if (!id.isEmpty() && (!m_adapters.contains(id) || !m_adapters[id])) {
        m_adapters[id] = adapter;
        m_adapterIds << id;
    }
}

void AdaptersManager::inflateAdapter(Adapter *adapter, const QJsonObject &adapterObj)
{
    if (!adapter)
        return;

    const QString path = adapterObj["Path"].toString();
    const QString alias = adapterObj["Alias"].toString();
    const bool powered = adapterObj["Powered"].toBool();
    const bool discovering = adapterObj["Discovering"].toBool();

    adapter->setId(path);
    adapter->setName(alias);
    adapter->setPowered(powered);
    adapter->setDiscover(discovering);
}

void AdaptersManager::adapterRefresh(const Adapter *adapter)
{
    QDBusObjectPath dPath(adapter->id());
    m_bluetoothInter->RequestDiscovery(dPath);
}

QList<const Adapter *> AdaptersManager::adapters()
{
    QList<const Adapter *> allAdapter = m_adapters.values();
    std::sort(allAdapter.begin(), allAdapter.end(), [ & ](const Adapter *adapter1, const Adapter *adapter2) {
        return m_adapterIds.indexOf(adapter1->id()) < m_adapterIds.indexOf(adapter2->id());
    });
    return allAdapter;
}
