#include "wirelessapplet.h"
#include "accesspointwidget.h"

#include <QJsonDocument>

#define WIDTH           300
#define MAX_HEIGHT      200
#define ITEM_HEIGHT     30

WirelessApplet::WirelessApplet(const QSet<NetworkDevice>::const_iterator &deviceIter, QWidget *parent)
    : QScrollArea(parent),

      m_device(*deviceIter),
      m_activeAP(),

      m_updateAPTimer(new QTimer(this)),

      m_centeralLayout(new QVBoxLayout),
      m_centeralWidget(new QWidget),
      m_controlPanel(new DeviceControlWidget(this)),
      m_networkInter(new DBusNetwork(this))
{
    setFixedHeight(WIDTH);

    m_updateAPTimer->setSingleShot(true);
    m_updateAPTimer->setInterval(100);

    m_centeralWidget->setFixedWidth(WIDTH);
    m_centeralWidget->setLayout(m_centeralLayout);

    m_centeralLayout->addWidget(m_controlPanel);
    m_centeralLayout->setSpacing(0);
    m_centeralLayout->setMargin(0);

    setWidget(m_centeralWidget);
    setFrameStyle(QFrame::NoFrame);
    setFixedWidth(300);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);

    connect(m_networkInter, &DBusNetwork::AccessPointAdded, this, &WirelessApplet::APAdded);
    connect(m_networkInter, &DBusNetwork::AccessPointRemoved, this, &WirelessApplet::APRemoved);
    connect(m_networkInter, &DBusNetwork::AccessPointPropertiesChanged, this, &WirelessApplet::APPropertiesChanged);
    connect(m_networkInter, &DBusNetwork::DevicesChanged, this, &WirelessApplet::deviceStateChanegd);

    connect(m_controlPanel, &DeviceControlWidget::deviceEnableChanged, this, &WirelessApplet::deviceEnableChanged);

    connect(m_updateAPTimer, &QTimer::timeout, this, &WirelessApplet::updateAPList);
}

NetworkDevice::NetworkState WirelessApplet::wirelessState() const
{
    return m_device.state();
}

int WirelessApplet::activeAPStrgength() const
{
    return m_activeAP.strength();
}

void WirelessApplet::init()
{
    setDeviceInfo();
    loadAPList();
    onActiveAPChanged();
}

void WirelessApplet::APAdded(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    AccessPoint ap(info);
    if (m_apList.contains(ap))
        return;

    m_apList.append(ap);
    m_updateAPTimer->start();
}

void WirelessApplet::APRemoved(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    m_apList.removeOne(AccessPoint(info));
    m_updateAPTimer->start();
}

void WirelessApplet::setDeviceInfo()
{
    // set device enable state
    m_controlPanel->setDeviceEnabled(m_networkInter->IsDeviceEnabled(m_device.dbusPath()));
    m_controlPanel->setDeviceName(m_device.vendor());
}

void WirelessApplet::loadAPList()
{
    const QString data = m_networkInter->GetAccessPoints(m_device.dbusPath());
    const QJsonDocument doc = QJsonDocument::fromJson(data.toUtf8());
    Q_ASSERT(doc.isArray());

    for (auto item : doc.array())
    {
        Q_ASSERT(item.isObject());

        AccessPoint ap(item.toObject());
        if (!m_apList.contains(ap))
            m_apList.append(ap);
    }

    m_updateAPTimer->start();
}

void WirelessApplet::APPropertiesChanged(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    QJsonDocument doc = QJsonDocument::fromJson(info.toUtf8());
    Q_ASSERT(doc.isObject());
    const AccessPoint ap(doc.object());

    auto it = std::find_if(m_apList.begin(), m_apList.end(),
                           [&] (const AccessPoint &a) {return a == ap;});

    if (it == m_apList.end())
        return;

    if (*it > ap)
    {
        *it = ap;
        m_activeAP = ap;
        m_updateAPTimer->start();

        emit activeAPChanged();
    }
}

void WirelessApplet::updateAPList()
{
    Q_ASSERT(sender() == m_updateAPTimer);

    // remove old items
    while (QLayoutItem *item = m_centeralLayout->takeAt(1))
    {
        delete item->widget();
        delete item;
    }

    int avaliableAPCount = 0;

    if (m_networkInter->IsDeviceEnabled(m_device.dbusPath()))
    {
        // sort ap list by strength
        std::sort(m_apList.begin(), m_apList.end(), std::greater<AccessPoint>());

        for (auto ap : m_apList)
        {
            AccessPointWidget *apw = new AccessPointWidget(ap);
            m_centeralLayout->addWidget(apw);

            ++avaliableAPCount;
        }
    }

    const int contentHeight = avaliableAPCount * ITEM_HEIGHT + m_controlPanel->height();
    m_centeralWidget->setFixedHeight(contentHeight);
    setFixedHeight(std::min(contentHeight, MAX_HEIGHT));
}

void WirelessApplet::deviceEnableChanged(const bool enable)
{
    m_networkInter->EnableDevice(m_device.dbusPath(), enable);
    m_updateAPTimer->start();
}

void WirelessApplet::deviceStateChanegd()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->devices().toUtf8());
    Q_ASSERT(doc.isObject());
    const QJsonObject obj = doc.object();

    for (auto infoList(obj.constBegin()); infoList != obj.constEnd(); ++infoList)
    {
        Q_ASSERT(infoList.value().isArray());

        if (infoList.key() != "wireless")
            continue;

        for (auto wireless : infoList.value().toArray())
        {
            const QJsonObject info = wireless.toObject();
            if (info.value("Path") == m_device.path())
            {
                const NetworkDevice prevInfo = m_device;
                m_device = NetworkDevice(NetworkDevice::Wireless, info);

                if (prevInfo.state() != m_device.state())
                    emit wirelessStateChanged(m_device.state());
                if (prevInfo.activeAp() != m_device.activeAp())
                    onActiveAPChanged();

                break;
            }
        }
    }

}

void WirelessApplet::onActiveAPChanged()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->GetAccessPoints(m_device.dbusPath()).value().toUtf8());
    Q_ASSERT(doc.isArray());

    for (auto dev : doc.array())
    {
        Q_ASSERT(dev.isObject());
        const QJsonObject obj = dev.toObject();

        if (obj.value("Path").toString() != m_device.activeAp())
            continue;

        m_activeAP = AccessPoint(obj);
        break;
    }

    emit activeAPChanged();
}
