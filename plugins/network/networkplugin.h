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
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked);
    void refershIcon(const QString &itemKey);
    const QString itemCommand(const QString &itemKey);
    const QString itemContextMenu(const QString &itemKey);
    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);
    QWidget *itemPopupApplet(const QString &itemKey);

private slots:
    void deviceAdded(const NetworkDevice &device);
    void deviceRemoved(const NetworkDevice &device);
    void networkStateChanged(const NetworkDevice::NetworkTypes &states);
    void deviceTypesChanged(const NetworkDevice::NetworkTypes &types);
    void refershDeviceItemVisible();
    void contextMenuRequested();

private:
    NetworkManager *m_networkManager;
    QTimer *m_refershTimer;

    QList<DeviceItem *> m_deviceItemList;
};

#endif // NETWORKPLUGIN_H
