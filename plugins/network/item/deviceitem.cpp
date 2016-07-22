#include "deviceitem.h"

DeviceItem::DeviceItem(const NetworkDevice::NetworkType type, const QUuid &deviceUuid)
    : QWidget(nullptr),
      m_type(type),
      m_deviceUuid(deviceUuid),

      m_networkManager(NetworkManager::instance(this))
{

}

QSize DeviceItem::sizeHint() const
{
    return QSize(24, 24);
}

const QUuid DeviceItem::uuid() const
{
    return m_deviceUuid;
}

NetworkDevice::NetworkType DeviceItem::type() const
{
    return m_type;
}
