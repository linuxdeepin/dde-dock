/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PLUGINMANAGER_H
#define PLUGINMANAGER_H

#include <QMap>
#include <QObject>
#include <QStringList>

#include "interfaces/dockconstants.h"
#include "widgets/old/abstractdockitem.h"
#include "controller/dockmodedata.h"
#include "pluginssettingframe.h"

class QFileSystemWatcher;
class PluginProxy;
class PluginManager : public QObject
{
    Q_OBJECT
public:
    explicit PluginManager(QObject *parent = 0);

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
    void notePluginModeChanged();

private:
    AbstractDockItem * sysPluginItem(QString id);
    PluginProxy * loadPlugin(const QString & path);
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
    QMap<QString, PluginProxy*> m_proxies;
    QFileSystemWatcher * m_watcher = NULL;
    QStringList m_searchPaths;
    DockModeData *m_dockModeData = DockModeData::instance();
    Dock::DockMode m_newMode;
    Dock::DockMode m_oldMode;
};

#endif // PLUGINMANAGER_H
