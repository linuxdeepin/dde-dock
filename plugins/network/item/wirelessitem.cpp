#include "wirelessitem.h"

#include <QPainter>

WirelessItem::WirelessItem(const QUuid &uuid)
    : DeviceItem(NetworkDevice::Wireless, uuid)
{

}

QWidget *WirelessItem::itemApplet()
{
    return nullptr;
}

void WirelessItem::paintEvent(QPaintEvent *e)
{
    DeviceItem::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), Qt::red);
}

const QPixmap WirelessItem::icon(const QString &key)
{
    if (!m_icons.contains(key))
    {

    }

    return m_icons.value(key);
}
