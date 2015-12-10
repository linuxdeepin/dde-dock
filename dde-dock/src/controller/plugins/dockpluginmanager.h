#ifndef DOCKPLUGINMANAGER_H
#define DOCKPLUGINMANAGER_H

#include <QMap>
#include <QObject>
#include <QStringList>

#include "interfaces/dockconstants.h"
#include "widgets/abstractdockitem.h"
#include "controller/dockmodedata.h"
#include "pluginssettingframe.h"

class QFileSystemWatcher;
class DockPluginProxy;
class DockPluginManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginManager(QObject *parent = 0);

signals:
    void itemInsert(AbstractDockItem *baseItem, AbstractDockItem *targetItem);
    void itemAppend(AbstractDockItem * item);
    void itemRemoved(AbstractDockItem * item);

public slots:
    void initAll();
    void onPluginsSetting(int y);
    void onDockModeChanged(Dock::DockMode newMode,
                           Dock::DockMode oldMode);

private slots:
    void watchedFileChanged(const QString & file);
    void watchedDirectoryChanged(const QString & directory);

private:
    AbstractDockItem * sysPluginItem(QString id);
    DockPluginProxy * loadPlugin(const QString & path);
    void handleSysPluginAdd(AbstractDockItem *item, QString uuid);
    void handleNormalPluginAdd(AbstractDockItem *item, QString uuid);
    void unloadPlugin(const QString & path);
    void initSettingWindow();
    void onPluginItemAdded(AbstractDockItem *item, QString uuid);
    void onPluginItemRemoved(AbstractDockItem *item, QString);

private:
    PluginsSettingFrame *m_settingFrame = NULL;
    QMap<AbstractDockItem *, QString> m_sysPlugins;
    QMap<AbstractDockItem *, QString> m_normalPlugins;
    QMap<QString, DockPluginProxy*> m_proxies;
    QFileSystemWatcher * m_watcher = NULL;
    QStringList m_searchPaths;
    DockModeData *m_dockModeData = DockModeData::instance();

    const QString SYSTRAY_PLUGIN_ID = "composite_item_key";
    const QString DATETIME_PLUGIN_ID = "id_datetime";
    const QString SHUTDOWN_PLUGIN_ID = "shutdown";
};

#endif // DOCKPLUGINMANAGER_H
