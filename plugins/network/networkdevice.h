#ifndef NETWORKDEVICE_H
#define NETWORKDEVICE_H

#include <QUuid>
#include <QDBusObjectPath>

class NetworkDevice
{
public:
    enum DeviceType {
        Invaild,
        Generic,
        Bluetooth,
        Wired,
        Wireless,
    };

public:
    explicit NetworkDevice(const DeviceType type, const QJsonObject &info);
    bool operator==(const QUuid &uuid) const;

    static DeviceType deviceType(const QString &type);

private:
    DeviceType m_type;

    QUuid m_uuid;
    QString m_objectPath;

};

#endif // NETWORKDEVICE_H
