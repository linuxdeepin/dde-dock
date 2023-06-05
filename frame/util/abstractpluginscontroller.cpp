// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "abstractpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "utils.h"
#include "pluginmanagerinterface.h"

#include <DNotifySender>
#include <DSysInfo>

#include <QDebug>
#include <QDir>
#include <QMapIterator>

static const QStringList CompatiblePluginApiList {
    "1.1.1",
    "1.2",
    "1.2.1",
    "1.2.2",
    DOCK_PLUGIN_API_VERSION
};

AbstractPluginsController::AbstractPluginsController(QObject *parent)
    : QObject(parent)
    , m_pluginManager(nullptr)
{
    qApp->installEventFilter(this);
}

AbstractPluginsController::~AbstractPluginsController()
{
    for (auto inter : m_pluginsMap.keys()) {
        delete m_pluginsMap.value(inter).value("pluginloader");
        m_pluginsMap[inter]["pluginloader"] = nullptr;
        m_pluginsMap.remove(inter);
        delete inter;
        inter = nullptr;
    }
}

void AbstractPluginsController::startLoader(PluginLoader *loader)
{
    connect(loader, &PluginLoader::finished, loader, &PluginLoader::deleteLater, Qt::QueuedConnection);
    connect(loader, &PluginLoader::pluginFounded, this, [ = ](const QString &pluginFile) {
        QPair<QString, PluginsItemInterface *> pair;
        pair.first = pluginFile;
        pair.second = nullptr;
        m_pluginLoadMap.insert(pair, false);
    });
    connect(loader, &PluginLoader::pluginFounded, this, &AbstractPluginsController::loadPlugin, Qt::QueuedConnection);

    int delay = Utils::SettingValue("com.deepin.dde.dock", "/com/deepin/dde/dock/", "delay-plugins-time", 0).toInt();
    QTimer::singleShot(delay, loader, [ = ] { loader->start(QThread::LowestPriority); });
}

void AbstractPluginsController::displayModeChanged()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->displayModeChanged(displayMode);
}

void AbstractPluginsController::positionChanged()
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->positionChanged(position);
}

void AbstractPluginsController::loadPlugin(const QString &pluginFile)
{
    QPluginLoader *pluginLoader = new QPluginLoader(pluginFile, this);
    const QJsonObject &meta = pluginLoader->metaData().value("MetaData").toObject();
    const QString &pluginApi = meta.value("api").toString();
    bool pluginIsValid = true;
    if (pluginApi.isEmpty() || !CompatiblePluginApiList.contains(pluginApi)) {
        qDebug() << objectName()
                 << "plugin api version not matched! expect versions:" << CompatiblePluginApiList
                 << ", got version:" << pluginApi
                 << ", the plugin file is:" << pluginFile;

        pluginIsValid = false;
    }

    PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
    if (!interface) {
        qDebug() << objectName() << "load plugin failed!!!" << pluginLoader->errorString() << pluginFile;

        pluginLoader->unload();
        pluginLoader->deleteLater();

        pluginIsValid = false;
    }

    if (!pluginIsValid) {
        for (auto &pair : m_pluginLoadMap.keys()) {
            if (pair.first == pluginFile) {
                m_pluginLoadMap.remove(pair);
            }
        }
        QString notifyMessage(tr("The plugin %1 is not compatible with the system."));
        Dtk::Core::DUtil::DNotifySender(notifyMessage.arg(QFileInfo(pluginFile).fileName())).appIcon("dialog-warning").call();
        return;
    }

    QMapIterator<QPair<QString, PluginsItemInterface *>, bool> it(m_pluginLoadMap);
    while (it.hasNext()) {
        it.next();
        if (it.key().first == pluginFile) {
            m_pluginLoadMap.remove(it.key());
            QPair<QString, PluginsItemInterface *> newPair;
            newPair.first = pluginFile;
            newPair.second = interface;
            m_pluginLoadMap.insert(newPair, false);
            break;
        }
    }

    // 保存 PluginLoader 对象指针
    QMap<QString, QObject *> interfaceData;
    interfaceData["pluginloader"] = pluginLoader;
    m_pluginsMap.insert(interface, interfaceData);

    PluginManagerInterface * pluginManager = dynamic_cast<PluginManagerInterface *>(interface);
    if (pluginManager) {
        m_pluginManager = pluginManager;
        connect(m_pluginManager, &PluginManagerInterface::pluginLoadFinished, this, &AbstractPluginsController::pluginLoaderFinished);
    }

    // NOTE(justforlxz): 插件的所有初始化工作都在init函数中进行，
    // loadPlugin函数是按队列执行的，initPlugin函数会有可能导致
    // 函数执行被阻塞。
    QTimer::singleShot(1, this, [ = ] {
        initPlugin(interface);
    });
}

void AbstractPluginsController::initPlugin(PluginsItemInterface *interface)
{
    if (!interface)
        return;

    qDebug() << objectName() << "init plugin: " << interface->pluginName();
    interface->init(this);

    for (const auto &pair : m_pluginLoadMap.keys()) {
        if (pair.second == interface)
            m_pluginLoadMap.insert(pair, true);
    }

    qDebug() << objectName() << "init plugin finished: " << interface->pluginName();

}

bool AbstractPluginsController::eventFilter(QObject *object, QEvent *event)
{
    if (object != qApp || event->type() != QEvent::DynamicPropertyChange)
        return false;

    QDynamicPropertyChangeEvent *const dpce = static_cast<QDynamicPropertyChangeEvent *>(event);
    const QString propertyName = dpce->propertyName();

    if (propertyName == PROP_POSITION)
        positionChanged();
    else if (propertyName == PROP_DISPLAY_MODE)
        displayModeChanged();

    return false;
}

PluginManagerInterface *AbstractPluginsController::pluginManager() const
{
    return m_pluginManager;
}
