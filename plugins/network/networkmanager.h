#ifndef NETWORKMANAGER_H
#define NETWORKMANAGER_H

#include "dbus/dbusnetwork.h"
#include "networkdevice.h"

#include <QJsonObject>
#include <QJsonDocument>
#include <QJsonArray>

class NetworkManager : public QObject
{
    Q_OBJECT

public:
    static NetworkManager *instance(QObject *parent = nullptr);

    void init();

    const NetworkDevice::NetworkTypes states() const;
    const NetworkDevice::NetworkTypes types() const;
    const QSet<NetworkDevice> deviceList() const;
    const QSet<QUuid> activeConnSet() const;

    NetworkDevice::NetworkState deviceState(const QUuid &uuid) const;
    const QString deviceHwAddr(const QUuid &uuid) const;
    const QString devicePath(const QUuid &uuid) const;
    const QJsonObject deviceConnInfo(const QUuid &uuid) const;
    const QSet<NetworkDevice>::const_iterator device(const QUuid &uuid) const;

signals:
    void deviceAdded(const NetworkDevice &device) const;
    void deviceChanged(const NetworkDevice &device) const;
    void deviceRemoved(const NetworkDevice &device) const;
    void activeConnectionChanged(const QUuid &uuid) const;
    void networkStateChanged(const NetworkDevice::NetworkTypes &states) const;
    void deviceTypesChanged(const NetworkDevice::NetworkTypes &types) const;

private:
    explicit NetworkManager(QObject *parent = 0);


private slots:
    void reloadDevices();
    void reloadActiveConnections();

private:
    NetworkDevice::NetworkTypes m_states;
    NetworkDevice::NetworkTypes m_types;
    DBusNetwork *m_networkInter;

    QSet<NetworkDevice> m_deviceSet;
    QSet<QUuid> m_activeConnSet;

    static NetworkManager *INSTANCE;
};
#endif // NETWORKMANAGER_H
