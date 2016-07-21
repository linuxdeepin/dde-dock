#include "networkplugin.h"

#define WIRED_ITEM      "wired"
#define WIRELESS_ITEM   "wireless"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent),

      m_networkManager(NetworkManager::instance(this)),

      m_wiredItem(nullptr)
{

}

const QString NetworkPlugin::pluginName() const
{
    return "network";
}

void NetworkPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == WIRED_ITEM)
        return m_wiredItem;

    return nullptr;
}
