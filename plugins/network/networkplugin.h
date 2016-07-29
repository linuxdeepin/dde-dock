#ifndef NETWORKPLUGIN_H
#define NETWORKPLUGIN_H

#include "pluginsiteminterface.h"
#include "networkmanager.h"
#include "item/deviceitem.h"

class NetworkPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "network.json")

public:
    explicit NetworkPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);
    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemPopupApplet(const QString &itemKey);

private slots:
    void deviceAdded(const NetworkDevice &device);
    void deviceRemoved(const NetworkDevice &device);
    void networkStateChanged(const NetworkDevice::NetworkTypes &states);
    void deviceTypesChanged(const NetworkDevice::NetworkTypes &types);
    void refershDeviceItemVisible();

private:
    NetworkManager *m_networkManager;

    QList<DeviceItem *> m_deviceItemList;
};

#endif // NETWORKPLUGIN_H
