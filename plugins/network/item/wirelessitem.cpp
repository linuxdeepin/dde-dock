
#include "wirelessitem.h"
#include "util/imageutil.h"

#include <QPainter>

WirelessItem::WirelessItem(const QUuid &uuid)
    : DeviceItem(NetworkDevice::Wireless, uuid),
      m_applet(nullptr)
{
    QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);
}

QWidget *WirelessItem::itemApplet()
{
    return m_applet;
}

void WirelessItem::paintEvent(QPaintEvent *e)
{
    DeviceItem::paintEvent(e);

    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    const int iconSize = std::min(width(), height()) * 0.8;
    const QPixmap pixmap = iconPix(displayMode, iconSize);

    QPainter painter(this);
    if (displayMode == Dock::Fashion)
    {
        const QPixmap pixmap = backgroundPix(iconSize);
        painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap);
    }
    painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap);
}

void WirelessItem::resizeEvent(QResizeEvent *e)
{
    DeviceItem::resizeEvent(e);

    m_icons.clear();
}

const QPixmap WirelessItem::iconPix(const Dock::DisplayMode displayMode, const int size)
{
    const QString key = QString("wireless-%1%2")
                                .arg(8)
                                .arg(displayMode == Dock::Fashion ? "" : "-symbolic");

    return cachedPix(key, size);
}

const QPixmap WirelessItem::backgroundPix(const int size)
{
    return cachedPix("wireless-background", size);
}

const QPixmap WirelessItem::cachedPix(const QString &key, const int size)
{
    if (!m_icons.contains(key))
        m_icons.insert(key, ImageUtil::loadSvg(":/wireless/resources/wireless/" + key + ".svg", size));

    return m_icons.value(key);
}

void WirelessItem::init()
{
    const auto devInfo = m_networkManager->device(m_deviceUuid);

    m_applet = new WirelessApplet(devInfo, this);
    m_applet->setVisible(false);
}
