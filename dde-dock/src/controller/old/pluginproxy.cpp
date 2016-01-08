#include <QPluginLoader>

#include "pluginproxy.h"
#include "pluginitemwrapper.h"
#include "controller/dockmodedata.h"

PluginProxy::PluginProxy(QPluginLoader * loader, DockPluginInterface * plugin) :
    QObject(),
    m_loader(loader),
    m_plugin(plugin)
{
}

PluginProxy::~PluginProxy()
{
    foreach (QString id, m_items.keys()) {
        emit itemRemoved(m_items.take(id), id);
    }
    m_items.clear();

    m_loader->unload();
    m_loader->deleteLater();

    qDebug() << "Plugin unloaded: " << m_loader->fileName();
}

DockPluginInterface * PluginProxy::plugin()
{
    return m_plugin;
}

Dock::DockMode PluginProxy::dockMode()
{
    return DockModeData::instance()->getDockMode();
}

bool PluginProxy::isSystemPlugin()
{
    return m_loader->metaData()["MetaData"].toObject()["sys_plugin"].toBool();
}

void PluginProxy::itemAddedEvent(QString id)
{
    if (m_plugin->getItem(id)) {
        qDebug() << "Item added on plugin " << m_plugin->getPluginName() << id;

        AbstractDockItem * item = new PluginItemWrapper(m_plugin, id);
        m_items[id] = item;

        emit itemAdded(item, id);
    }
}

void PluginProxy::itemRemovedEvent(QString id)
{
    qDebug() << "Item removed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item) {
        item->needPreviewImmediatelyHide();
        m_items.take(id);

        emit itemRemoved(item, id);
    }
}

void PluginProxy::infoChangedEvent(DockPluginInterface::InfoType type, const QString &id)
{
    switch (type) {
    case DockPluginInterface::InfoTypeItemSize:
    case DockPluginInterface::ItemSize: //Q_DECL_DEPRECATED
        itemSizeChangedEvent(id);
        break;
    case DockPluginInterface::InfoTypeAppletSize:
    case DockPluginInterface::AppletSize:   //Q_DECL_DEPRECATED
        appletSizeChangedEvent(id);
        break;
    case DockPluginInterface::InfoTypeTitle:
        emit titleChanged(id);
        break;
    case DockPluginInterface::InfoTypeConfigurable:
        emit configurableChanged(id);
        break;
    case DockPluginInterface::InfoTypeEnable:
        emit enabledChanged(id);
        break;
    default:
        break;
    }
}

void PluginProxy::itemSizeChangedEvent(QString id)
{
//    qDebug() << "Item size changed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item) {
        item->adjustSize();

        emit item->widthChanged();
    }
}

void PluginProxy::appletSizeChangedEvent(QString id)
{
//    qWarning() << "Applet size changed on plugin " << m_plugin->getPluginName() << id;

    AbstractDockItem * item = m_items.value(id);
    if (item)
        item->needPreviewUpdate();
}
