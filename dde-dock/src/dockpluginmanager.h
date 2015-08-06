#ifndef DOCKPLUGINMANAGER_H
#define DOCKPLUGINMANAGER_H

#include <QObject>
#include <QMap>
#include <QStringList>

#include "dockconstants.h"
#include "abstractdockitem.h"
#include "Controller/dockmodedata.h"

class QFileSystemWatcher;
class DockPluginProxy;
class DockPluginManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginManager(QObject *parent = 0);

    void initAll();

signals:
    void itemMove(AbstractDockItem *baseItem, AbstractDockItem *targetItem);
    void itemInsert(AbstractDockItem *baseItem, AbstractDockItem *targetItem);
    void itemAppend(AbstractDockItem * item);
    void itemRemoved(AbstractDockItem * item);

public slots:
    void onDockModeChanged(Dock::DockMode newMode,
                           Dock::DockMode oldMode);

private slots:
    void watchedFileChanged(const QString & file);
    void watchedDirectoryChanged(const QString & directory);

private:
    AbstractDockItem * sysPluginItem(QString id);
    DockPluginProxy * loadPlugin(const QString & path);
    void handleSysPluginAdd(AbstractDockItem *item, QString uuid);
    void handleNormalPluginAdd(AbstractDockItem *item);
    void unloadPlugin(const QString & path);
    void updatePluginPos(Dock::DockMode newMode, Dock::DockMode oldMode);

private:
    QMap<AbstractDockItem *, QString> m_sysPlugins;
    QMap<QString, DockPluginProxy*> m_proxies;
    QList<AbstractDockItem *> m_normalPlugins;
    QFileSystemWatcher * m_watcher = NULL;
    QStringList m_searchPaths;
    DockModeData *m_dockModeData = DockModeData::instance();

    const QString SYSTRAY_PLUGIN_ID = "composite_item_key";
    const QString DATETIME_PLUGIN_ID = "id_datetime";
};

#endif // DOCKPLUGINMANAGER_H
