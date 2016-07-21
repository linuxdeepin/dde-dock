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

const NetworkManager::NetworkStates NetworkManager::states() const
{
    return m_states;
}

NetworkManager::NetworkManager(QObject *parent)
    : QObject(parent),

      m_states(Offline),

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

    m_deviceList.clear();
    for (auto infoList(obj.constBegin()); infoList != obj.constEnd(); ++infoList)
    {
        Q_ASSERT(infoList.value().isArray());
        const NetworkDevice::DeviceType deviceType = NetworkDevice::deviceType(infoList.key());

        for (auto device : infoList.value().toArray())
            m_deviceList.append(NetworkDevice(deviceType, device.toObject()));
    }
}

void NetworkManager::reloadActiveConnections()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->activeConnections().toUtf8());
    Q_ASSERT(doc.isObject());
    const QJsonObject obj = doc.object();

    NetworkStates states = Offline;
    m_activeConnList.clear();
    for (auto info(obj.constBegin()); info != obj.constEnd(); ++info)
    {
        Q_ASSERT(info.value().isObject());
        const QJsonObject infoObj = info.value().toObject();

        const QUuid uuid = infoObj.value("Uuid").toString();
        // if uuid not in device list, its a wireless connection
        const bool isWireless = std::find(m_deviceList.cbegin(), m_deviceList.cend(), uuid) == m_deviceList.cend();

        if (isWireless)
            states |= wirelessConnection;
        else
            states |= WiredConnection;

        m_activeConnList.append(uuid);
    }

    if (m_states == states)
        return;
    m_states = states;
    emit networkStateChanged(m_states);
}
