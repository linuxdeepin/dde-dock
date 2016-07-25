#ifndef NETWORKDEVICE_H
#define NETWORKDEVICE_H

#include "networkdevice.h"

#include <QUuid>
#include <QDBusObjectPath>
#include <QJsonObject>

class NetworkDevice
{
public:
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
