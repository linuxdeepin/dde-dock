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

PluginsItemInterface::ItemType ShutdownPlugin::pluginType(const QString &itemKey)
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

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}
