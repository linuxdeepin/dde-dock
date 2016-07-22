#include "networkdevice.h"

#include <QDebug>
#include <QJsonObject>

NetworkDevice::NetworkDevice(const NetworkType type, const QJsonObject &info)
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

NetworkDevice::NetworkType NetworkDevice::type() const
{
    return m_type;
}

const QString NetworkDevice::path() const
{
    return m_objectPath;
}

NetworkDevice::NetworkType NetworkDevice::deviceType(const QString &type)
{
    if (type == "bt")
        return NetworkDevice::Bluetooth;
    if (type == "generic")
        return NetworkDevice::Generic;
    if (type == "wired")
        return NetworkDevice::Wired;
    if (type == "wireless")
        return NetworkDevice::Wireless;

    Q_ASSERT(false);

    return NetworkDevice::None;
}
