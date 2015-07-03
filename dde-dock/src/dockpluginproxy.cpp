#include "dockpluginproxy.h"
#include "pluginitemwrapper.h"

DockPluginProxy::DockPluginProxy(DockPluginInterface * plugin, QObject * parent) :
    QObject(parent),
    m_plugin(plugin)
{
}

void DockPluginProxy::itemAddedEvent(QString uuid)
{
    qDebug() << "Item added on plugin " << m_plugin->name();

    AbstractDockItem * item = new PluginItemWrapper(m_plugin, uuid);
    m_items << item;

    emit itemAdded(item);
}

void DockPluginProxy::itemRemovedEvent(QString uuid)
{
    qDebug() << "Item removed on plugin " << m_plugin->name();

    emit itemRemoved(getItem(uuid));
}

AbstractDockItem * DockPluginProxy::getItem(QString uuid)
{
    foreach (AbstractDockItem * item, m_items) {
        PluginItemWrapper *wrapper = qobject_cast<PluginItemWrapper*>(item);
        if (wrapper->uuid() == uuid) {
            return item;
        }
    }
    return NULL;
}
