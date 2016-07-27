#include "constants.h"
#include "wireditem.h"
#include "util/imageutil.h"

#include <QPainter>

WiredItem::WiredItem(const QUuid &deviceUuid)
    : DeviceItem(NetworkDevice::Wired, deviceUuid),

      m_connected(false),
      m_itemTips(new QLabel(this))
{

    m_itemTips->setVisible(false);
    m_itemTips->setStyleSheet("color:white;"
                              "padding:5px 10px;");

    connect(m_networkManager, &NetworkManager::networkStateChanged, this, &WiredItem::reloadIcon);
    connect(m_networkManager, &NetworkManager::activeConnectionChanged, this, &WiredItem::activeConnectionChanged);
}

QWidget *WiredItem::itemApplet()
{
    m_itemTips->setText(tr("Unknow"));

    do {
        if (!m_connected)
        {
            m_itemTips->setText(tr("Disconnect"));
            break;
        }

        const QJsonObject info = m_networkManager->deviceConnInfo(m_deviceUuid);
        if (!info.contains("Ip4"))
            break;
        const QJsonObject ipv4 = info.value("Ip4").toObject();
        if (!ipv4.contains("Address"))
            break;
        m_itemTips->setText(ipv4.value("Address").toString());
    } while (false);

    return m_itemTips;
}

void WiredItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void WiredItem::resizeEvent(QResizeEvent *e)
{
    DeviceItem::resizeEvent(e);

    reloadIcon();
}

void WiredItem::reloadIcon()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    if (displayMode == Dock::Fashion)
    {
        const int size = std::min(width(), height()) * 0.8;

        if (m_connected)
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-connected.svg", size);
        else
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-disconnected.svg", size);
    } else {
        if (m_connected)
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-connected-small.svg", 16);
        else
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-disconnected-small.svg", 16);
    }
}

void WiredItem::activeConnectionChanged(const QUuid &uuid)
{
    if (uuid != m_deviceUuid)
        return;

    m_connected = m_networkManager->activeConnSet().contains(m_deviceUuid);
    update();
}
