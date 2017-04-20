#include "dockpluginloader.h"
#include "dockpluginscontroller.h"

#include <QDebug>

DockPluginLoader::DockPluginLoader(QObject *parent)
    : QThread(parent)
{

}

void DockPluginLoader::run()
{
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

        emit pluginFounded(pluginsDir.absoluteFilePath(file));

        msleep(500);
    }

    emit finished();
}
