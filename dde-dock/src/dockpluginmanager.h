#ifndef DOCKPLUGINMANAGER_H
#define DOCKPLUGINMANAGER_H

#include <QObject>
#include <QMap>
#include <QStringList>

#include "dockconstants.h"

class DockPluginProxy;
class DockPluginManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginManager(QObject *parent = 0);

    QList<DockPluginProxy*> getAll();

public slots:
    void onDockModeChanged(Dock::DockMode newMode,
                           Dock::DockMode oldMode);

private:
    QStringList m_searchPaths;
    QMap<QString, DockPluginProxy*> m_proxies;

    DockPluginProxy* loadPlugin(QString &path);
};

#endif // DOCKPLUGINMANAGER_H
