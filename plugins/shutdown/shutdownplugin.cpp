#include "shutdownplugin.h"

#include <QIcon>

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent)
{
}

const QString ShutdownPlugin::pluginName() const
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

    displayModeChanged(qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>());
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    if (displayMode == Dock::Fashion)
        m_icon.addFile(":/icons/resources/icons/fashion.svg");
    else
        m_icon.addFile(":/icons/resources/icons/normal.svg");
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
