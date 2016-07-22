#include "networkplugin.h"

#define WIRED_ITEM      "wired"
#define WIRELESS_ITEM   "wireless"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent),

      m_networkManager(NetworkManager::instance(this)),

      m_wiredItem(new WiredItem)
{
    m_wiredItem->setVisible(false);
}

const QString NetworkPlugin::pluginName() const
{
    return "network";
}

void NetworkPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    connect(m_networkManager, &NetworkManager::networkStateChanged, this, &NetworkPlugin::networkStateChanged);

    m_networkManager->init();
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == WIRED_ITEM)
        return m_wiredItem;

    return nullptr;
}

void NetworkPlugin::networkStateChanged(const NetworkDevice::NetworkTypes &states)
{
    // has wired connection
    if (states.testFlag(NetworkDevice::Wired))
        m_proxyInter->itemAdded(this, WIRED_ITEM);
    else
        m_proxyInter->itemRemoved(this, WIRED_ITEM);
}
