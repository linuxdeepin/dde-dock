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

signals:
    void deviceAdded(const NetworkDevice &device) const;
    void deviceRemoved(const NetworkDevice &device) const;
    void networkStateChanged(const NetworkDevice::NetworkTypes &states) const;

private:
    explicit NetworkManager(QObject *parent = 0);

private slots:
    void reloadDevices();
    void reloadActiveConnections();

private:
    NetworkDevice::NetworkTypes m_states;
    NetworkDevice::NetworkTypes m_types;
    DBusNetwork *m_networkInter;

    QSet<NetworkDevice> m_deviceList;
    QSet<QUuid> m_activeConnList;

    static NetworkManager *INSTANCE;
};
#endif // NETWORKMANAGER_H
