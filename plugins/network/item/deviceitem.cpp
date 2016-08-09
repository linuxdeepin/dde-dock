#include "deviceitem.h"

DeviceItem::DeviceItem(const QUuid &deviceUuid)
    : QWidget(nullptr),
      m_deviceUuid(deviceUuid),

      m_networkManager(NetworkManager::instance(this))
{

}

QSize DeviceItem::sizeHint() const
{
    return QSize(26, 26);
}

const QUuid DeviceItem::uuid() const
{
    return m_deviceUuid;
}

const QString DeviceItem::itemCommand() const
{
    return QString();
}

QWidget *DeviceItem::itemPopup()
{
    return nullptr;
}

QWidget *DeviceItem::itemApplet()
{
    return nullptr;
}
