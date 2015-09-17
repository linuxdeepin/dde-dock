#include <QPluginLoader>

#include "dockpluginproxy.h"
#include "pluginitemwrapper.h"
#include "controller/dockmodedata.h"

DockPluginProxy::DockPluginProxy(QPluginLoader * loader, DockPluginInterface * plugin) :
    QObject(),
    m_loader(loader),
    m_plugin(plugin)
{
}

DockPluginProxy::~DockPluginProxy()
{
    foreach (QString id, m_items.keys()) {
        emit itemRemoved(m_items.take(id), id);
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

bool DockPluginProxy::isSystemPlugin()
{
    return m_loader->metaData()["MetaData"].toObject()["sys_plugin"].toBool();
}

void DockPluginProxy::itemAddedEvent(QString id)
{
    if (m_plugin->getItem(id)) {
        qDebug() << "Item added on plugin " << m_plugin->getPluginName() << id;

        AbstractDockItem * item = new PluginItemWrapper(m_plugin, id);
        m_items[id] = item;

        emit itemAdded(item, id);
    }
}

void DockPluginProxy::itemRemovedEvent(QString id)
{
    qDebug() << "Item removed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item) {
        m_items.take(id);

        emit itemRemoved(item, id);
    }
}

void DockPluginProxy::itemSizeChangedEvent(QString id)
{
    qDebug() << "Item size changed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item) {
        item->adjustSize();

        emit item->widthChanged();
    }
}

void DockPluginProxy::appletSizeChangedEvent(QString id)
{
    qWarning() << "Applet size changed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item)
        item->needPreviewUpdate();
}
