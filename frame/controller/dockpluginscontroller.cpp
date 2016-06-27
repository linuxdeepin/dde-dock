#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dockitemcontroller.h"

#include <QDebug>
#include <QDir>

DockPluginsController::DockPluginsController(DockItemController *itemControllerInter)
    : QObject(itemControllerInter),
      m_itemControllerInter(itemControllerInter)
{
    QMetaObject::invokeMethod(this, "loadPlugins", Qt::QueuedConnection);
}

DockPluginsController::~DockPluginsController()
{
}

void DockPluginsController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = new PluginsItem(itemInter, itemKey);

    emit pluginItemInserted(item);
}

DisplayMode DockPluginsController::displayMode() const
{
    return Fashion;
}

void DockPluginsController::loadPlugins()
{
//    Q_ASSERT(m_pluginLoaderList.isEmpty());
//    Q_ASSERT(m_pluginsInterfaceList.isEmpty());

#ifdef QT_DEBUG
    const QDir pluginsDir("plugins");
#else
    const QDir pluginsDir("../lib/dde-dock/plugins");
#endif
    const QStringList plugins = pluginsDir.entryList(QDir::Files);

    for (const QString file : plugins)
    {
        if (!QLibrary::isLibrary(file))
            continue;

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-"))
            continue;

        const QString pluginFilePath = pluginsDir.absoluteFilePath(file);
        QPluginLoader *pluginLoader = new QPluginLoader(pluginFilePath, this);
        PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
        if (!interface)
        {
            pluginLoader->unload();
            pluginLoader->deleteLater();
            continue;
        }

        interface->init(this);

//        m_pluginLoaderList.append(pluginLoader);
//        m_pluginsInterfaceList.append(interface);
    }
}
