#include <QPluginLoader>

#include "dockpluginproxy.h"
#include "pluginitemwrapper.h"
#include "Controller/dockmodedata.h"

DockPluginProxy::DockPluginProxy(QPluginLoader * loader, DockPluginInterface * plugin) :
    QObject(),
    m_loader(loader),
    m_plugin(plugin)
{
}

DockPluginProxy::~DockPluginProxy()
{
    foreach (AbstractDockItem * item, m_items.values()) {
        emit itemRemoved(item);
    }
    m_items.clear();

    m_loader->unload();
    m_loader->deleteLater();

    qDebug() << "Plugin unloaded: " << m_loader->fileName();
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

    if (!m_items.contains(uuid)) {
        if (m_plugin->getItem(uuid)) {
            AbstractDockItem * item = new PluginItemWrapper(m_plugin, uuid);
            m_items[uuid] = item;

            emit itemAdded(item);
        }
    }
}

void DockPluginProxy::itemRemovedEvent(QString uuid)
{
    qDebug() << "Item removed on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = m_items.value(uuid);
    if (item) {
        m_items.take(uuid);

        emit itemRemoved(item);
    }
}

void DockPluginProxy::itemSizeChangedEvent(QString uuid)
{
    qDebug() << "Item size changed on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = m_items.value(uuid);
    if (item) {
        item->adjustSize();

        emit item->widthChanged();
    }
}

void DockPluginProxy::appletSizeChangedEvent(QString uuid)
{
    qWarning() << "Applet size changed on plugin " << m_plugin->name() << uuid;

    AbstractDockItem * item = m_items.value(uuid);
    if (item) {
        item->resizePreview();
    }
}
