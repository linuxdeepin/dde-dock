#include "networkplugin.h"
#include "item/wireditem.h"
#include "item/wirelessitem.h"

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
        {
            return deviceItem;
        }

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
    DeviceItem *item = nullptr;
    switch (device.type())
    {
    case NetworkDevice::Wired:      item = new WiredItem(device.uuid());        break;
    case NetworkDevice::Wireless:   item = new WirelessItem(device.uuid());     break;
    default:;
    }

    if (!item)
        return;

    m_deviceItemList.append(item);
    refershDeviceItemVisible();
}

void NetworkPlugin::deviceRemoved(const NetworkDevice &device)
{
    const auto item = std::find_if(m_deviceItemList.begin(), m_deviceItemList.end(),
                                   [&] (DeviceItem *dev) {return device == dev->uuid();});

    if (item == m_deviceItemList.cend())
        return;

    m_proxyInter->itemRemoved(this, (*item)->uuid().toString());
    (*item)->deleteLater();
    m_deviceItemList.erase(item);
}

void NetworkPlugin::networkStateChanged(const NetworkDevice::NetworkTypes &states)
{
    Q_UNUSED(states)

    refershDeviceItemVisible();
}

void NetworkPlugin::deviceTypesChanged(const NetworkDevice::NetworkTypes &types)
{
    Q_UNUSED(types)

    refershDeviceItemVisible();
}

void NetworkPlugin::refershDeviceItemVisible()
{
    const NetworkDevice::NetworkTypes types = m_networkManager->types();
    const bool hasWirelessDevice = types.testFlag(NetworkDevice::Wireless);

    for (auto item : m_deviceItemList)
    {
        switch (item->type())
        {
        case NetworkDevice::Wireless:
            m_proxyInter->itemAdded(this, item->uuid().toString());
            break;

        case NetworkDevice::Wired:
            if (item->state() == NetworkDevice::Activated || !hasWirelessDevice)
                m_proxyInter->itemAdded(this, item->uuid().toString());
            else
                m_proxyInter->itemRemoved(this, item->uuid().toString());
            break;

        default:;
        }
    }
}
