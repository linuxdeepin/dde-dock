#include "networkplugin.h"
#include "item/wireditem.h"

#define WIRED_ITEM      "wired"
#define WIRELESS_ITEM   "wireless"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent),

      m_networkManager(NetworkManager::instance(this))
{
}

const QString NetworkPlugin::pluginName() const
{
    return "network";
}

void NetworkPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    connect(m_networkManager, &NetworkManager::networkStateChanged, this, &NetworkPlugin::networkStateChanged);
    connect(m_networkManager, &NetworkManager::deviceAdded, this, &NetworkPlugin::deviceAdded);
    connect(m_networkManager, &NetworkManager::deviceRemoved, this, &NetworkPlugin::deviceRemoved);

    m_networkManager->init();
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    for (auto deviceItem : m_deviceItemList)
        if (deviceItem->uuid() == itemKey)
            return deviceItem;

    return nullptr;
}

QWidget *NetworkPlugin::itemPopupApplet(const QString &itemKey)
{
    for (auto deviceItem : m_deviceItemList)
        if (deviceItem->uuid() == itemKey)
            return deviceItem->itemApplet();

    return nullptr;
}

void NetworkPlugin::deviceAdded(const NetworkDevice &device)
{
    qDebug() << "add: " << device.uuid();

    DeviceItem *item = nullptr;
    switch (device.type())
    {
    case NetworkDevice::Wired:      item = new WiredItem(device.uuid());        break;
    default:;
    }

    if (!item)
        return;

    m_deviceItemList.append(item);
    m_proxyInter->itemAdded(this, device.uuid().toString());
}

void NetworkPlugin::deviceRemoved(const NetworkDevice &device)
{
    qDebug() << "remove: " << device.uuid();

    const auto item = std::find_if(m_deviceItemList.begin(), m_deviceItemList.end(),
                                   [&] (DeviceItem *dev) {return dev->uuid() == device.uuid();});

    if (item == m_deviceItemList.cend())
        return;

    m_proxyInter->itemRemoved(this, (*item)->uuid().toString());
    (*item)->deleteLater();
    m_deviceItemList.erase(item);
}

void NetworkPlugin::networkStateChanged(const NetworkDevice::NetworkTypes &states)
{
    Q_UNUSED(states)

    qDebug() << states;

//    for (auto item : m_deviceItemList)
//        if (item->ty)
//    const QList<NetworkDevice> deviceList = m_networkManager->deviceList();
//    const auto items = states.testFlag(NetworkDevice::Wireless)
//                        ? std::find_if(m_deviceList.cbegin(), m_deviceList.cend(),
//                                    [] (const NetworkDevice &dev) {return dev.type() == NetworkDevice::Wireless;})
//                        : std::find_if(m_deviceList.cbegin(), m_deviceList.cend(),
//                                      [] (const NetworkDevice &dev) {return dev.type() == NetworkDevice::Wired;});

//    for (auto dev : deviceList)
//        qDebug() << dev.uuid() << dev.type();

//    // has wired connection
//    if (states.testFlag(NetworkDevice::Wired))
//        m_proxyInter->itemAdded(this, WIRED_ITEM);
//    else
//        m_proxyInter->itemRemoved(this, WIRED_ITEM);
}
