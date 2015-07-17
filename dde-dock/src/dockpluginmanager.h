#ifndef DOCKPLUGINMANAGER_H
#define DOCKPLUGINMANAGER_H

#include <QObject>
#include <QMap>
#include <QStringList>

#include "dockconstants.h"
#include "abstractdockitem.h"

class QFileSystemWatcher;
class DockPluginProxy;
class DockPluginManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginManager(QObject *parent = 0);

    void initAll();

signals:
    void itemAdded(AbstractDockItem * item);
    void itemRemoved(AbstractDockItem * item);

public slots:
    void onDockModeChanged(Dock::DockMode newMode,
                           Dock::DockMode oldMode);

private:
    QStringList m_searchPaths;
    QMap<QString, DockPluginProxy*> m_proxies;
    QFileSystemWatcher * m_watcher;

    DockPluginProxy * loadPlugin(const QString & path);
    void unloadPlugin(const QString & path);

private slots:
    void watchedFileChanged(const QString & file);
    void watchedDirectoryChanged(const QString & directory);
};

#endif // DOCKPLUGINMANAGER_H
