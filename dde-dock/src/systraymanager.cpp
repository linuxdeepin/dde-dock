#include <QDir>
#include <QPluginLoader>
#include <QDebug>

#include "systraymanager.h"
#include "pluginitemwrapper.h"

static QString SystrayPluginPath = "/usr/share/dde-dock/plugins/libdock-systray-plugin.so";

SystrayManager::SystrayManager(QObject *parent)
    : QObject(parent),
      m_plugin(0)
{
    this->loadPlugin();
}

QList<AbstractDockItem*> SystrayManager::trayIcons()
{
    QList<AbstractDockItem*> result;

    if (m_plugin) {
        QStringList uuids = m_plugin->uuids();

        foreach (QString uuid, uuids) {
            result << new PluginItemWrapper(m_plugin, uuid);
        }
    }

    return result;
}

void SystrayManager::loadPlugin()
{
    if (QFile::exists(SystrayPluginPath)) {
        QPluginLoader loader(SystrayPluginPath);
        QObject *plugin = loader.instance();
        if (plugin) {
            m_plugin = qobject_cast<DockPluginInterface*>(plugin);
            m_plugin->init();
        } else {
            qWarning() << "Failed to load systray plugin.";
            qWarning() << loader.errorString();
        }
    } else {
        qWarning() << "libdock-systray-plugin.so file not found!";
    }
}
