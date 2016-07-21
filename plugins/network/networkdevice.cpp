#include "networkdevice.h"

#include <QDebug>
#include <QJsonObject>

NetworkDevice::NetworkDevice(const DeviceType type, const QJsonObject &info)
    : m_type(type)
{
    m_uuid = info.value("UniqueUuid").toString();
    m_objectPath = info.value("Path").toString();

    //    qDebug() << m_uuid << m_objectPath;
}

bool NetworkDevice::operator==(const QUuid &uuid) const
{
    return m_uuid == uuid;
}

NetworkDevice::DeviceType NetworkDevice::deviceType(const QString &type)
{
    if (type == "bt")
        return Bluetooth;
    if (type == "generic")
        return Generic;
    if (type == "wired")
        return Wired;
    if (type == "wireless")
        return Wireless;

    Q_ASSERT(false);

    return Invaild;
}
