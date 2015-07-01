#include <QDir>
#include <QPluginLoader>
#include <QDebug>

#include "systraymanager.h"

static QString SystrayPluginPath = "/usr/share/dde-dock/plugins/libdock-systray-plugin.so";

SystrayManager::SystrayManager(QObject *parent)
    : QObject(parent),
      m_plugin(0)
{
    this->loadPlugin();
}

QList<AbstractDockItem*> SystrayManager::trayIcons()
{
    if (m_plugin) {
        return m_plugin->items();
    } else {
        return QList<AbstractDockItem*>();
    }
}

void SystrayManager::loadPlugin()
{
    if (QFile::exists(SystrayPluginPath)) {
        QPluginLoader loader(SystrayPluginPath);
        QObject *plugin = loader.instance();
        if (plugin) {
            m_plugin = qobject_cast<DockPluginInterface*>(plugin);
        } else {
            qWarning() << "Failed to load systray plugin.";
            qWarning() << loader.errorString();
        }
    } else {
        qWarning() << "libdock-systray-plugin.so file not found!";
    }
}
