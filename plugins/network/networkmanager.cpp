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
    return m_deviceSet;
}

const QSet<QUuid> NetworkManager::activeConnSet() const
{
    return m_activeConnSet;
}

const QString NetworkManager::deviceHwAddr(const QUuid &uuid) const
{
    const auto item = std::find_if(m_deviceSet.cbegin(), m_deviceSet.cend(),
                                   [&] (const NetworkDevice &dev) {return dev == uuid;});

    if (item == m_deviceSet.cend())
        return QString();

    return item->hwAddress();
}

const QJsonObject NetworkManager::deviceInfo(const QUuid &uuid) const
{
    const QString addr = deviceHwAddr(uuid);
    if (addr.isEmpty())
        return QJsonObject();

    const QJsonDocument infos = QJsonDocument::fromJson(m_networkInter->GetActiveConnectionInfo().value().toUtf8());
    Q_ASSERT(infos.isArray());

    for (auto info : infos.array())
    {
        Q_ASSERT(info.isObject());
        const QJsonObject obj = info.toObject();
        if (obj.contains("HwAddress") && obj.value("HwAddress").toString() == addr)
            return obj;
    }

    return QJsonObject();
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
            deviceSet.insert(NetworkDevice(deviceType, device.toObject()));
    }

    const QSet<NetworkDevice> removedDeviceList = m_deviceSet - deviceSet;
    for (auto dev : removedDeviceList)
        emit deviceRemoved(dev);
    for (auto dev : deviceSet)
    {
        if (m_deviceSet.contains(dev))
            emit deviceChanged(dev);
        else
            emit deviceAdded(dev);
    }

    m_deviceSet = std::move(deviceSet);
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
    QSet<QUuid> activeConnList;
    for (auto info(obj.constBegin()); info != obj.constEnd(); ++info)
    {
        Q_ASSERT(info.value().isObject());
        const QJsonObject infoObj = info.value().toObject();

        const QUuid uuid = infoObj.value("Uuid").toString();
        // if uuid not in device list, its a wireless connection
        const bool isWireless = std::find(m_deviceSet.cbegin(), m_deviceSet.cend(), uuid) == m_deviceSet.cend();

        if (isWireless)
            states |= NetworkDevice::Wireless;
        else
            states |= NetworkDevice::Wired;

        activeConnList.insert(uuid);
    }

    const QSet<QUuid> removedConnList = m_activeConnSet - activeConnList;
    m_activeConnSet = std::move(activeConnList);

    for (auto uuid : removedConnList)
        emit activeConnectionChanged(uuid);

    for (auto uuid : m_activeConnSet)
        emit activeConnectionChanged(uuid);

    if (m_states == states)
        return;

    m_states = states;
    emit networkStateChanged(m_states);

    qDebug() << "network states: " << m_states;
}
