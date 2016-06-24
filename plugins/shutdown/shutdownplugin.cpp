#include "shutdownplugin.h"

#include <QIcon>

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent)
{
    m_icon.addFile(":/icons/resources/icons/fashion.svg");
}

const QString ShutdownPlugin::pluginName()
{
    return "shutdown";
}

PluginsItemInterface::PluginType ShutdownPlugin::pluginType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Simple;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, QString());
}

const QIcon ShutdownPlugin::itemIcon(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_icon;
}

int ShutdownPlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return 0;
}
