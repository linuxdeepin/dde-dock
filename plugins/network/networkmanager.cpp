#include "networkmanager.h"
#include "networkdevice.h"

NetworkManager *NetworkManager::INSTANCE = nullptr;

NetworkManager *NetworkManager::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new NetworkManager(parent);

    return INSTANCE;
}

void NetworkManager::init()
{
    reloadDevices();
    reloadActiveConnections();
}

const NetworkDevice::NetworkTypes NetworkManager::states() const
{
    return m_states;
}

const NetworkDevice::NetworkTypes NetworkManager::types() const
{
    return m_types;
}

const QSet<NetworkDevice> NetworkManager::deviceList() const
{
    return m_deviceList;
}

NetworkManager::NetworkManager(QObject *parent)
    : QObject(parent),

      m_states(NetworkDevice::None),
      m_types(NetworkDevice::None),

      m_networkInter(new DBusNetwork(this))
{

    connect(m_networkInter, &DBusNetwork::DevicesChanged, this, &NetworkManager::reloadDevices);
    connect(m_networkInter, &DBusNetwork::ActiveConnectionsChanged, this, &NetworkManager::reloadActiveConnections);
}

void NetworkManager::reloadDevices()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->devices().toUtf8());
    Q_ASSERT(doc.isObject());
    const QJsonObject obj = doc.object();

    NetworkDevice::NetworkTypes types = NetworkDevice::None;
    QSet<NetworkDevice> deviceSet;
    for (auto infoList(obj.constBegin()); infoList != obj.constEnd(); ++infoList)
    {
        Q_ASSERT(infoList.value().isArray());
        const NetworkDevice::NetworkType deviceType = NetworkDevice::deviceType(infoList.key());

        types |= deviceType;

        for (auto device : infoList.value().toArray())
        {
            deviceSet.insert(NetworkDevice(deviceType, device.toObject()));
        }
    }

    const QSet<NetworkDevice> removedDeviceList = m_deviceList - deviceSet;
    for (auto dev : removedDeviceList)
        emit deviceRemoved(dev);
    const QSet<NetworkDevice> addedDeviceList = deviceSet - m_deviceList;
    for (auto dev : addedDeviceList)
        emit deviceAdded(dev);

    m_deviceList = std::move(deviceSet);
    if (m_types == types)
        return;

    m_types = types;
    qDebug() << "device type: " << m_types;
}

void NetworkManager::reloadActiveConnections()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->activeConnections().toUtf8());
    Q_ASSERT(doc.isObject());
    const QJsonObject obj = doc.object();

    NetworkDevice::NetworkTypes states = NetworkDevice::None;
    m_activeConnList.clear();
    for (auto info(obj.constBegin()); info != obj.constEnd(); ++info)
    {
        Q_ASSERT(info.value().isObject());
        const QJsonObject infoObj = info.value().toObject();

        const QUuid uuid = infoObj.value("Uuid").toString();
        // if uuid not in device list, its a wireless connection
        const bool isWireless = std::find(m_deviceList.cbegin(), m_deviceList.cend(), uuid) == m_deviceList.cend();

        if (isWireless)
            states |= NetworkDevice::Wireless;
        else
            states |= NetworkDevice::Wired;

        m_activeConnList.insert(uuid);
    }

    if (m_states == states)
        return;

    m_states = states;
    emit networkStateChanged(m_states);

    qDebug() << "network states: " << m_states;
}
