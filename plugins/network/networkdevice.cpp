#include "networkdevice.h"

#include <QDebug>
#include <QJsonObject>

NetworkDevice::NetworkDevice(const NetworkType type, const QJsonObject &info)
    : m_type(type),
      m_infoObj(info)
{
    m_uuid = info.value("UniqueUuid").toString();
    m_objectPath = info.value("Path").toString();
}

bool NetworkDevice::operator==(const QUuid &uuid) const
{
    return m_uuid == uuid;
}

bool NetworkDevice::operator==(const NetworkDevice &device) const
{
    return m_uuid == device.m_uuid;
}

NetworkDevice::NetworkState NetworkDevice::state() const
{
    return NetworkState(m_infoObj.value("State").toInt());
}

NetworkDevice::NetworkType NetworkDevice::type() const
{
    return m_type;
}

const QUuid NetworkDevice::uuid() const
{
    return m_uuid;
}

const QString NetworkDevice::path() const
{
    return m_objectPath;
}

const QString NetworkDevice::hwAddress() const
{
    return std::move(m_infoObj.value("HwAddress").toString());
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

