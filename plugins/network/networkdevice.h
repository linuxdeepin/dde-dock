#ifndef NETWORKDEVICE_H
#define NETWORKDEVICE_H

#include "networkdevice.h"

#include <QUuid>
#include <QDBusObjectPath>
#include <QJsonObject>

class NetworkDevice
{
public:
    enum NetworkState {
        Unknow          = 0,
        Unmanaged       = 10,
        Unavailable     = 20,
        Disconnected    = 30,
        Prepare         = 40,
        Config          = 50,
        NeedAuth        = 60,
        IpConfig        = 70,
        IpCheck         = 80,
        Secondaries     = 90,
        Activated       = 100,
        Deactivation    = 110,
        Failed          = 120,
    };

    enum NetworkType {
        None        = 0,
        Generic     = 1 << 0,
        Wired       = 1 << 1,
        Wireless    = 1 << 2,
        Bluetooth   = 1 << 3,
    };
    Q_DECLARE_FLAGS(NetworkTypes, NetworkType)

public:
    static NetworkType deviceType(const QString &type);

    explicit NetworkDevice(const NetworkType type, const QJsonObject &info);
    bool operator==(const QUuid &uuid) const;
    bool operator==(const NetworkDevice &device) const;

    NetworkState state() const;
    NetworkType type() const;
    const QUuid uuid() const;
    const QString path() const;
    const QString hwAddress() const;

private:
    NetworkType m_type;

    QUuid m_uuid;
    QString m_objectPath;
    QJsonObject m_infoObj;
};

inline uint qHash(const NetworkDevice &device)
{
    return qHash(device.uuid());
}

Q_DECLARE_OPERATORS_FOR_FLAGS(NetworkDevice::NetworkTypes)

#endif // NETWORKDEVICE_H
