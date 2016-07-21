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
    enum NetworkState {
        Offline             = 0,
        WiredConnection     = 1 << 0,
        wirelessConnection  = 1 << 1,
    };
    Q_DECLARE_FLAGS(NetworkStates, NetworkState)

    static NetworkManager *instance(QObject *parent = nullptr);

    void init();

    const NetworkStates states() const;

signals:
    void networkStateChanged(const NetworkStates &states) const;

private:
    explicit NetworkManager(QObject *parent = 0);

private slots:
    void reloadDevices();
    void reloadActiveConnections();

private:
    NetworkStates m_states;
    DBusNetwork *m_networkInter;

    QList<NetworkDevice> m_deviceList;
    QList<QUuid> m_activeConnList;

    static NetworkManager *INSTANCE;
};

Q_DECLARE_OPERATORS_FOR_FLAGS(NetworkManager::NetworkStates)

#endif // NETWORKMANAGER_H
