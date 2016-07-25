#include "wirelessapplet.h"

#include <QJsonDocument>

#define WIDTH   300

WirelessApplet::WirelessApplet(const QString &devicePath, QWidget *parent)
    : QScrollArea(parent),
      m_devicePath(devicePath),

      m_centeralLayout(new QVBoxLayout),
      m_centeralWidget(new QWidget),
      m_controlPanel(new DeviceControlWidget(this)),
      m_networkInter(new DBusNetwork(this))
{
    setFixedHeight(WIDTH);

    m_centeralWidget->setFixedWidth(WIDTH);
    m_centeralWidget->setLayout(m_centeralLayout);

    m_centeralLayout->addWidget(m_controlPanel);

    setWidget(m_centeralWidget);
    setFrameStyle(QFrame::NoFrame);
    setFixedWidth(300);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    setDeviceInfo();

    connect(m_networkInter, &DBusNetwork::AccessPointPropertiesChanged, this, &WirelessApplet::APChanged);
}

void WirelessApplet::setDeviceInfo()
{
    // set device enable state
    m_controlPanel->setDeviceEnabled(m_networkInter->IsDeviceEnabled(QDBusObjectPath(m_devicePath)));

    // set device name
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
            if (info.value("Path") == m_devicePath)
            {
                m_controlPanel->setDeviceName(info.value("Vendor").toString());
                break;
            }
        }
    }
}

void WirelessApplet::APChanged(const QString &devPath, const QString &info)
{
    if (devPath != m_devicePath)
        return;

    qDebug() << info;
}
