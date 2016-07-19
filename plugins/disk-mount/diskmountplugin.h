#ifndef DISKMOUNTPLUGIN_H
#define DISKMOUNTPLUGIN_H

#include "pluginsiteminterface.h"
#include "dbus/dbusdiskmount.h"

class DiskMountPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "disk-mount.json")

public:
    explicit DiskMountPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);

    QWidget *itemWidget(const QString &itemKey);

private:
    DBusDiskMount *m_diskInter;
};

#endif // DISKMOUNTPLUGIN_H
