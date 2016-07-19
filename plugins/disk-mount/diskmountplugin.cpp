#include "diskmountplugin.h"

DiskMountPlugin::DiskMountPlugin(QObject *parent)
    : QObject(parent)
{

}

const QString DiskMountPlugin::pluginName() const
{
    return "disk-mount";
}

void DiskMountPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
}

QWidget *DiskMountPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}
