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

// private methods
DockPluginProxy* DockPluginManager::loadPlugin(QString &path)
{
    QPluginLoader pluginLoader(path);
    QObject *plugin = pluginLoader.instance();

    if (plugin) {
        DockPluginInterface * interface = qobject_cast<DockPluginInterface*>(plugin);
        if (interface) {
            DockPluginProxy *proxy = new DockPluginProxy(interface);
            interface->init(proxy);
            qDebug() << "Plugin loaded: " << path;
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
