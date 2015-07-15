#include "dockpluginproxy.h"
#include "pluginitemwrapper.h"
#include "Controller/dockmodedata.h"

DockPluginProxy::DockPluginProxy(DockPluginInterface * plugin, QObject * parent) :
    QObject(parent),
    m_plugin(plugin)
{
}

DockPluginInterface * DockPluginProxy::plugin()
{
    return m_plugin;
}

Dock::DockMode DockPluginProxy::dockMode()
{
    return DockModeData::instance()->getDockMode();
}

void DockPluginProxy::itemAddedEvent(QString uuid)
{
    qDebug() << "Item added on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = new PluginItemWrapper(m_plugin, uuid);
    m_items << item;

    emit itemAdded(item);
}

void DockPluginProxy::itemRemovedEvent(QString uuid)
{
    qDebug() << "Item removed on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = getItem(uuid);
    if (item) {
        m_items.takeAt(m_items.indexOf(item));

        emit itemRemoved(item);
    }
}

void DockPluginProxy::itemSizeChangedEvent(QString uuid)
{
    qDebug() << "Item size changed on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = getItem(uuid);
    item->adjustSize();
    if (item) {
        emit item->widthChanged();
    }
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
