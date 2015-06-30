#include <QDir>
#include <QPluginLoader>
#include <QDebug>

#include "systraymanager.h"

SystrayManager::SystrayManager(QObject *parent)
    : QObject(parent),
      m_plugin(0)
{
    this->loadPlugin();
}

QList<AbstractDockItem*> SystrayManager::trayIcons()
{
    return m_plugin->items();
}

void SystrayManager::loadPlugin()
{
    QPluginLoader loader("/home/hualet/project/linuxdeepin/dde-workspace-2015/dde-dock-systray-plugin/build/libdock-systray-plugin.so");
    QObject *plugin = loader.instance();
    if (plugin) {
        m_plugin = qobject_cast<DockPluginInterface*>(plugin);
    } else {
        qWarning() << "Failed to load systray plugin.";
        qWarning() << loader.errorString();
    }
}
