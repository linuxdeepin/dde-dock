#include <QDir>
#include <QLibrary>
#include <QPluginLoader>
#include <QFileSystemWatcher>

#include "dockpluginmanager.h"
#include "dockpluginproxy.h"
#include "dockplugininterface.h"

DockPluginManager::DockPluginManager(QObject *parent) :
    QObject(parent)
{
    m_searchPaths << "/usr/share/dde-dock/plugins/";

    m_watcher = new QFileSystemWatcher(this);
    m_watcher->addPaths(m_searchPaths);

    foreach (QString path, m_searchPaths) {
        QDir pluginsDir(path);

        foreach (QString fileName, pluginsDir.entryList(QDir::Files)) {
            QString pluginPath = pluginsDir.absoluteFilePath(fileName);

            this->loadPlugin(pluginPath);
        }
    }

    connect(m_watcher, &QFileSystemWatcher::fileChanged, this, &DockPluginManager::watchedFileChanged);
    connect(m_watcher, &QFileSystemWatcher::directoryChanged, this, &DockPluginManager::watchedDirectoryChanged);
}

void DockPluginManager::initAll()
{
    foreach (DockPluginProxy * proxy, m_proxies.values()) {
        proxy->plugin()->init(proxy);
    }
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
DockPluginProxy * DockPluginManager::loadPlugin(const QString &path)
{
    // check the file type
    if (!QLibrary::isLibrary(path)) return NULL;

    QPluginLoader * pluginLoader = new QPluginLoader(path);

    // check the apiVersion the plugin uses
    double apiVersion = pluginLoader->metaData()["MetaData"].toObject()["api_version"].toDouble();
    if (apiVersion != PLUGIN_API_VERSION) return NULL;


    QObject *plugin = pluginLoader->instance();

    if (plugin) {
        DockPluginInterface * interface = qobject_cast<DockPluginInterface*>(plugin);

        if (interface) {
            qDebug() << "Plugin loaded: " << path;

            DockPluginProxy * proxy = new DockPluginProxy(pluginLoader, interface);
            if (proxy) {
                m_proxies[path] = proxy;
                m_watcher->addPath(path);

                connect(proxy, &DockPluginProxy::itemAdded, this, &DockPluginManager::itemAdded);
                connect(proxy, &DockPluginProxy::itemRemoved, this, &DockPluginManager::itemRemoved);

                return proxy;
            }
        } else {
            qWarning() << "Load plugin failed(failed to convert) " << path;

            return NULL;
        }
    } else {
        qWarning() << "Load plugin failed" << pluginLoader->errorString();

        return NULL;
    }
}

void DockPluginManager::unloadPlugin(const QString &path)
{
    if (m_proxies.contains(path)) {
        DockPluginProxy * proxy = m_proxies.take(path);
        delete proxy;
    }
}

// private slots
void DockPluginManager::watchedFileChanged(const QString & file)
{
    qDebug() << "DockPluginManager::watchedFileChanged" << file;
    this->unloadPlugin(file);

    if (QFile::exists(file)) {
        DockPluginProxy * proxy = loadPlugin(file);

        if (proxy) proxy->plugin()->init(proxy);
    }
}

void DockPluginManager::watchedDirectoryChanged(const QString & directory)
{
    qDebug() << "DockPluginManager::watchedDirectoryChanged" << directory;
    // we just need to take care of the situation that new files pop up in
    // our watched directory.
    QDir targetDir(directory);
    foreach (QString fileName, targetDir.entryList(QDir::Files)) {
        QString absPath = targetDir.absoluteFilePath(fileName);
        if (!m_proxies.contains(absPath)) {
            DockPluginProxy * proxy = loadPlugin(absPath);

            if (proxy) proxy->plugin()->init(proxy);
        }
    }
}
