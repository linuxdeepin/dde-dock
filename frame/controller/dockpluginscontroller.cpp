// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "traypluginitem.h"

#include <QDebug>
#include <QDir>
#include <QDrag>

DockPluginsController::DockPluginsController(QObject *parent)
    : AbstractPluginsController(parent)
{
    setObjectName("DockPlugin");
}

void DockPluginsController::itemAdded(PluginsItemInterface *const itemInter, const QString &itemKey)
{
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();

    // check if same item added
    if (mPluginsMap.contains(itemInter))
        if (mPluginsMap[itemInter].contains(itemKey))
            return;

    // 取 plugin api
    QPluginLoader *pluginLoader = qobject_cast<QPluginLoader*>(mPluginsMap[itemInter].value("pluginloader"));
    if (!pluginLoader) {
        return;
    }
    const QJsonObject &meta = pluginLoader->metaData().value("MetaData").toObject();
    const QString &pluginApi = meta.value("api").toString();

    PluginsItem *item = nullptr;
    if (itemInter->pluginName() == "tray") {
        item = new TrayPluginItem(itemInter, itemKey, pluginApi);
        if (item->graphicsEffect()) {
            item->graphicsEffect()->setEnabled(false);
        }
        connect(static_cast<TrayPluginItem *>(item), &TrayPluginItem::trayVisableCountChanged,
                this, &DockPluginsController::trayVisableCountChanged, Qt::UniqueConnection);
    } else {
        item = new PluginsItem(itemInter, itemKey, pluginApi);
    }

    mPluginsMap[itemInter][itemKey] = item;

    emit pluginItemInserted(item);
}

void DockPluginsController::itemUpdate(PluginsItemInterface *const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->update();

    emit pluginItemUpdated(item);
}

void DockPluginsController::itemRemoved(PluginsItemInterface *const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(item);

    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();
    mPluginsMap[itemInter].remove(itemKey);

    // do not delete the itemWidget object(specified in the plugin interface)
    item->centralWidget()->setParent(nullptr);

    if (item->isDragging()) {
        QDrag::cancel();
    }

    // just delete our wrapper object(PluginsItem)
    item->deleteLater();
}

void DockPluginsController::requestWindowAutoHide(PluginsItemInterface *const itemInter, const QString &itemKey, const bool autoHide)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestWindowAutoHide(autoHide);
}

void DockPluginsController::requestRefreshWindowVisible(PluginsItemInterface *const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestRefreshWindowVisible();
}

void DockPluginsController::requestSetAppletVisible(PluginsItemInterface *const itemInter, const QString &itemKey, const bool visible)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    if (visible) {
        item->showPopupApplet(itemInter->itemPopupApplet(itemKey));
    } else {
        item->hidePopup();
    }
}

void DockPluginsController::startLoader()
{
    loadLocalPlugins();
    loadSystemPlugins();
}

void DockPluginsController::loadLocalPlugins()
{
    QString pluginsDir(QString("%1/.local/lib/dde-dock/plugins/").arg(QDir::homePath()));

    if (!QDir(pluginsDir).exists()) {
        return;
    }

    qDebug() << "using dock local plugins dir:" << pluginsDir;

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}

void DockPluginsController::loadSystemPlugins()
{
    QString pluginsDir(qApp->applicationDirPath() + "/../plugins");
#ifndef QT_DEBUG
    pluginsDir = "/usr/lib/dde-dock/plugins";
#endif
    qDebug() << "using dock plugins dir:" << pluginsDir;

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}
