/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKPLUGINSMANAGER_H
#define DOCKPLUGINSMANAGER_H

#include <QMap>
#include <QObject>
#include <QStringList>

#include "widgets/dockitem.h"
#include "interfaces/dockconstants.h"
#include "controller/dockmodedata.h"
#include "../../widgets/plugin/dockpluginssettingwindow.h"

class QFileSystemWatcher;
class DockPluginProxy;
class DockPluginsManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginsManager(QObject *parent = 0);

signals:
    void itemInsert(DockItem *baseItem, DockItem *targetItem);
    void itemAppend(DockItem * item);
    void itemRemoved(DockItem * item);

public slots:
    void initAll();
    void onPluginsSetting(int y);
    void onDockModeChanged(Dock::DockMode newMode,
                           Dock::DockMode oldMode);

private slots:
    void watchedFileChanged(const QString & file);
    void watchedDirectoryChanged(const QString & directory);
    void notePluginModeChanged();

private:
    DockItem * sysPluginItem(QString id);
    DockPluginProxy * loadPlugin(const QString & path);
    void handleSysPluginAdd(DockItem *item, QString uuid);
    void handleNormalPluginAdd(DockItem *item, QString uuid);
    void unloadPlugin(const QString & path);
    void initSettingWindow();
    void onPluginItemAdded(DockItem *item, QString uuid);
    void onPluginItemRemoved(DockItem *item, QString);

private:
    DockPluginsSettingWindow *m_settingWindow = NULL;
    QMap<DockItem *, QString> m_sysPlugins;
    QMap<DockItem *, QString> m_normalPlugins;
    QMap<QString, DockPluginProxy*> m_proxies;
    QFileSystemWatcher * m_watcher = NULL;
    QStringList m_searchPaths;
    DockModeData *m_dockModeData = DockModeData::instance();
    Dock::DockMode m_newMode;
    Dock::DockMode m_oldMode;
};

#endif // DOCKPLUGINSMANAGER_H
