/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QDir>
#include <QLibrary>
#include <QPluginLoader>

#include "dockpluginproxy.h"
#include "dockpluginsmanager.h"
#include "interfaces/dockplugininterface.h"

const QString SYSTRAY_PLUGIN_ID = "composite_item_key";
const QString DATETIME_PLUGIN_ID = "id_datetime";
const QString SHUTDOWN_PLUGIN_ID = "shutdown";
const int DELAY_NOTE_MODE_CHANGED_INTERVAL = 500;

DockPluginsManager::DockPluginsManager(QObject *parent) :
    QObject(parent)
{
    m_settingWindow = new DockPluginsSettingWindow;

    m_searchPaths << "/usr/lib/dde-dock/plugins/";

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockPluginsManager::onDockModeChanged);
}

void DockPluginsManager::initAll()
{
    for (QString path : m_searchPaths) {
        QDir pluginsDir(path);

        for (QString fileName : pluginsDir.entryList(QDir::Files)) {
            QString pluginPath = pluginsDir.absoluteFilePath(fileName);

            this->loadPlugin(pluginPath);
        }
    }

    for (DockPluginProxy * proxy : m_proxies.values()) {
        connect(proxy, &DockPluginProxy::configurableChanged, [=](const QString &id) {
            if (proxy->plugin()->configurable(id)) {
                m_settingWindow->onPluginAdd(proxy->plugin()->enabled(id),
                                            id,
                                            proxy->plugin()->getName(id),
                                            proxy->plugin()->getIcon(id));
            }
            else {
                m_settingWindow->onPluginRemove(id);
            }
        });
        connect(proxy, &DockPluginProxy::enabledChanged, [=](const QString &id) {
            m_settingWindow->onPluginEnabledChanged(id, proxy->plugin()->enabled(id));
        });
        connect(proxy, &DockPluginProxy::titleChanged, [=](const QString &id) {
            m_settingWindow->onPluginTitleChanged(id, proxy->plugin()->getName(id));
        });

        proxy->plugin()->init(proxy);
    }

    initSettingWindow();
}

void DockPluginsManager::onPluginsSetting(int y)
{
    m_settingWindow->move(QCursor::pos().x(), y - m_settingWindow->height());
    m_settingWindow->show();
}

// public slots
void DockPluginsManager::onDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode == oldMode)
        return;

    m_newMode = newMode;
    m_oldMode = oldMode;

    //hide plugin immediately
    for (DockItem *item : m_sysPlugins.keys()) {
        item->setVisible(false);
    }
    for (DockItem *item : m_normalPlugins.keys()) {
        item->setVisible(false);
    }

    //Many plugin when changing the mode has done a lot of work
    //and cause UI block
    //so,delay to note plugin dock-mode changed after the dock updated it style
    QTimer::singleShot(DELAY_NOTE_MODE_CHANGED_INTERVAL, this, SLOT(notePluginModeChanged()));
}

DockPluginProxy * DockPluginsManager::loadPlugin(const QString &path)
{
    // check the file type
    if (!QLibrary::isLibrary(path))
        return NULL;

    QPluginLoader * pluginLoader = new QPluginLoader(path);

    // check the apiVersion the plugin uses
    double apiVersion = pluginLoader->metaData()["MetaData"].toObject()["api_version"].toDouble();
    if (apiVersion != PLUGIN_API_VERSION)
        return NULL;


    QObject *plugin = pluginLoader->instance();

    if (plugin) {
        DockPluginInterface * interface = qobject_cast<DockPluginInterface*>(plugin);

        if (interface) {
            qDebug() << "Plugin loaded: " << path;

            DockPluginProxy * proxy = new DockPluginProxy(pluginLoader, interface);
            if (proxy) {
                m_proxies[path] = proxy;
                connect(proxy, &DockPluginProxy::itemAdded, this, &DockPluginsManager::onPluginItemAdded);
                connect(proxy, &DockPluginProxy::itemRemoved, this, &DockPluginsManager::onPluginItemRemoved);
                connect(m_settingWindow, &DockPluginsSettingWindow::checkedChanged, [=](QString uuid, bool checked){
                    //NOTE:one sender, multi receiver
                    if (interface->ids().indexOf(uuid) != -1) {
                        interface->setEnabled(uuid, checked);
                    }
                });

                return proxy;
            }
        }
        else {
            qWarning() << "Load plugin failed(failed to convert) " << path;
        }
    }
    else {
        qWarning() << "Load plugin failed" << pluginLoader->errorString();
    }

    return NULL;
}

void DockPluginsManager::unloadPlugin(const QString &path)
{
    if (m_proxies.contains(path)) {
        DockPluginProxy * proxy = m_proxies.take(path);
        delete proxy;
    }
}

void DockPluginsManager::initSettingWindow()
{
    foreach (DockPluginProxy *proxy, m_proxies.values()) {
        QStringList ids = proxy->plugin()->ids();
        foreach (QString uuid, ids) {
            if (proxy->plugin()->configurable(uuid)){
                m_settingWindow->onPluginAdd(proxy->plugin()->enabled(uuid),
                                            uuid,
                                            proxy->plugin()->getName(uuid),
                                            proxy->plugin()->getIcon(uuid));
            }
        }
    }
}

void DockPluginsManager::onPluginItemAdded(DockItem *item, QString uuid)
{
    DockPluginProxy *proxy = qobject_cast<DockPluginProxy *>(sender());
    if (!proxy)
        return;

    if (proxy->isSystemPlugin())
        handleSysPluginAdd(item, uuid);
    else
        handleNormalPluginAdd(item, uuid);
}

void DockPluginsManager::onPluginItemRemoved(DockItem *item, QString)
{
    m_sysPlugins.remove(item);
    m_normalPlugins.remove(item);

    emit itemRemoved(item);
    item->setVisible(false);
    item->deleteLater();
}

// private slots
void DockPluginsManager::watchedFileChanged(const QString & file)
{
    qDebug() << "DockPluginsManager::watchedFileChanged" << file;
    this->unloadPlugin(file);

    if (QFile::exists(file)) {
        DockPluginProxy * proxy = loadPlugin(file);

        if (proxy) proxy->plugin()->init(proxy);
    }
}

void DockPluginsManager::watchedDirectoryChanged(const QString & directory)
{
    qDebug() << "DockPluginsManager::watchedDirectoryChanged" << directory;
    // we just need to take care of the situation that new files pop up in
    // our watched directory.
    QDir targetDir(directory);
    foreach (QString fileName, targetDir.entryList(QDir::Files)) {
        QString absPath = targetDir.absoluteFilePath(fileName);
        if (!m_proxies.contains(absPath)) {
            DockPluginProxy * proxy = loadPlugin(absPath);

            if (proxy) proxy->plugin()->init(proxy);
        }
    }
}

void DockPluginsManager::notePluginModeChanged()
{
    for (DockPluginProxy * proxy : m_proxies) {
        DockPluginInterface * plugin = proxy->plugin();
        plugin->changeMode(m_newMode, m_oldMode);
    }

    //make sure all plugin will show
    for (DockItem *item : m_sysPlugins.keys()) {
        item->setVisible(true);
    }
    for (DockItem *item : m_normalPlugins.keys()) {
        item->setVisible(true);
    }

    //reanchor systray-plugin
    DockItem *sysItem = m_sysPlugins.key(SYSTRAY_PLUGIN_ID);
    m_sysPlugins.remove(sysItem);
    emit itemRemoved(sysItem);
    handleSysPluginAdd(sysItem, SYSTRAY_PLUGIN_ID);
}

DockItem *DockPluginsManager::sysPluginItem(QString id)
{
    int si = m_sysPlugins.values().indexOf(id);

    if (si != -1)
        return m_sysPlugins.keys().at(si);
    else
        return NULL;
}

void DockPluginsManager::handleSysPluginAdd(DockItem *item, QString uuid)
{
    if (!item || m_sysPlugins.values().indexOf(uuid) != -1)
        return;

    m_sysPlugins.insert(item, uuid);

    if (uuid == SHUTDOWN_PLUGIN_ID)
    {
        if (m_sysPlugins.values().contains(DATETIME_PLUGIN_ID))
            emit itemInsert(sysPluginItem(DATETIME_PLUGIN_ID), item);
        else
            emit itemAppend(item);

        return;
    }

    if (uuid == SYSTRAY_PLUGIN_ID) {
        if (m_dockModeData->getDockMode() != Dock::FashionMode) {
            emit itemAppend(item);
        }
        else {
            emit itemInsert(sysPluginItem(SHUTDOWN_PLUGIN_ID), item);
        }
    }
    else {
        emit itemInsert(NULL, item);
    }
}

void DockPluginsManager::handleNormalPluginAdd(DockItem *item, QString uuid)
{
    if (!item || m_normalPlugins.values().indexOf(uuid) != -1)
        return;

    if (m_dockModeData->getDockMode() != Dock::FashionMode) {

        qDebug() << uuid << m_sysPlugins.values();

        if (m_sysPlugins.values().contains(SHUTDOWN_PLUGIN_ID))
            itemInsert(sysPluginItem(SHUTDOWN_PLUGIN_ID), item);
        else if (m_sysPlugins.values().contains(DATETIME_PLUGIN_ID))
            itemInsert(sysPluginItem(DATETIME_PLUGIN_ID), item);
        else
            itemInsert(nullptr, item);
    }
    else {
        //Normal plug placed in the far left on Fashion Mode
        emit itemAppend(item);
    }

    m_normalPlugins.insert(item, uuid);
}
