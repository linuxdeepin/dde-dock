#include <QDir>
#include <QPluginLoader>

#include "dockpluginmanager.h"
#include "dockpluginproxy.h"
#include "dockplugininterface.h"

DockPluginManager::DockPluginManager(QObject *parent) :
    QObject(parent)
{
    m_searchPaths << "/usr/share/dde-dock/plugins/";

    foreach (QString path, m_searchPaths) {
        QDir pluginsDir(path);

        foreach (QString fileName, pluginsDir.entryList(QDir::Files)) {
            QString pluginPath = pluginsDir.absoluteFilePath(fileName);
            DockPluginProxy * proxy = loadPlugin(pluginPath);
            m_proxies[pluginPath] = proxy;
        }
    }
}

QList<DockPluginProxy*> DockPluginManager::getAll()
{
    return m_proxies.values();
}

// public slots
void DockPluginManager::onDockModeChanged(Dock::DockMode newMode,
                                          Dock::DockMode oldMode)
{
    qDebug() << "DockPluginManager::onDockModeChanged " << newMode << oldMode;

    foreach (DockPluginProxy * proxy, m_proxies) {
        DockPluginInterface * plugin = proxy->plugin();
        plugin->changeMode(newMode, oldMode);
    }
}

// private methods
DockPluginProxy* DockPluginManager::loadPlugin(QString &path)
{
    QPluginLoader pluginLoader(path);
    QObject *plugin = pluginLoader.instance();

    if (plugin) {
        DockPluginInterface * interface = qobject_cast<DockPluginInterface*>(plugin);
        if (interface) {
            qDebug() << "Plugin loaded: " << path;

            DockPluginProxy *proxy = new DockPluginProxy(interface);
            interface->init(proxy);

            return proxy;
        } else {
            qWarning() << "Load plugin failed(failed to convert) " << path;
            return NULL;
        }
    } else {
        qWarning() << "Load plugin failed" << pluginLoader.errorString();
        return NULL;
    }
}
