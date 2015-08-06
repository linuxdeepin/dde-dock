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

    updatePluginPos(newMode, oldMode);
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

                connect(proxy, &DockPluginProxy::itemAdded, [=](AbstractDockItem *item, QString uuid){
                    if (pluginLoader->metaData()["MetaData"].toObject()["sys_plugin"].toBool())
                        handleSysPluginAdd(item, uuid);
                    else
                        handleNormalPluginAdd(item);
                });
                connect(proxy, &DockPluginProxy::itemRemoved, [=](AbstractDockItem *item){
                    m_sysPlugins.remove(item);
                    m_normalPlugins.removeAt(m_normalPlugins.indexOf(item));
                    emit itemRemoved(item);
                });

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

void DockPluginManager::updatePluginPos(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode == Dock::FashionMode && oldMode != Dock::FashionMode){
        foreach (AbstractDockItem *item, m_normalPlugins) {
            emit itemMove(NULL, item);  //Move to the front of the list
        }
    }else if (oldMode == Dock::FashionMode){
        AbstractDockItem * systrayItem = sysPluginItem(SYSTRAY_PLUGIN_ID);
        foreach (AbstractDockItem *item, m_normalPlugins) {
            emit itemMove(systrayItem, item);   //Move to the back of systray plugin
        }
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

AbstractDockItem *DockPluginManager::sysPluginItem(QString id)
{
    int si = m_sysPlugins.values().indexOf(id);

    if (si != -1)
        return m_sysPlugins.keys().at(si);
    else
        return NULL;
}

void DockPluginManager::handleSysPluginAdd(AbstractDockItem *item, QString uuid)
{
    m_sysPlugins.insert(item, uuid);

    if (uuid == SYSTRAY_PLUGIN_ID){
        if (m_dockModeData->getDockMode() == Dock::FashionMode)
            emit itemInsert(sysPluginItem(DATETIME_PLUGIN_ID), item);
        else
            emit itemInsert(NULL, item);
    }
    else
        emit itemAppend(item);
}

void DockPluginManager::handleNormalPluginAdd(AbstractDockItem *item)
{
    m_normalPlugins.append(item);

    if (m_dockModeData->getDockMode() == Dock::FashionMode)
        emit itemInsert(NULL, item);
    else
        emit itemInsert(sysPluginItem(SYSTRAY_PLUGIN_ID), item);
}
