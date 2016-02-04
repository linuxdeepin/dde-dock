/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QPluginLoader>

#include "dockpluginproxy.h"
#include "dockpluginitemwrapper.h"
#include "controller/dockmodedata.h"

DockPluginProxy::DockPluginProxy(QPluginLoader * loader, DockPluginInterface * plugin) :
    QObject(),
    m_loader(loader),
    m_plugin(plugin)
{
}

DockPluginProxy::~DockPluginProxy()
{
    foreach (QString id, m_items.keys()) {
        emit itemRemoved(m_items.take(id), id);
    }
    m_items.clear();

    m_loader->unload();
    m_loader->deleteLater();

    qDebug() << "Plugin unloaded: " << m_loader->fileName();
}

DockPluginInterface * DockPluginProxy::plugin()
{
    return m_plugin;
}

Dock::DockMode DockPluginProxy::dockMode()
{
    return DockModeData::instance()->getDockMode();
}

bool DockPluginProxy::isSystemPlugin()
{
    return m_loader->metaData()["MetaData"].toObject()["sys_plugin"].toBool();
}

void DockPluginProxy::itemAddedEvent(QString id)
{
    if (m_plugin->getItem(id)) {
        qDebug() << "Item added on plugin " << m_plugin->getPluginName() << id;

        DockItem * item = new DockPluginItemWrapper(m_plugin, id);
        m_items[id] = item;

        emit itemAdded(item, id);
    }
}

void DockPluginProxy::itemRemovedEvent(QString id)
{
    qDebug() << "Item removed on plugin " << m_plugin->getPluginName() << id;

    DockItem * item = m_items.value(id);
    if (item) {
        item->hidePreview(true);
        m_items.take(id);

        emit itemRemoved(item, id);
    }
}

void DockPluginProxy::infoChangedEvent(DockPluginInterface::InfoType type, const QString &id)
{
    switch (type) {
    case DockPluginInterface::InfoTypeItemSize:
    case DockPluginInterface::ItemSize: //Q_DECL_DEPRECATED
        itemSizeChangedEvent(id);
        break;
    case DockPluginInterface::InfoTypeAppletSize:
    case DockPluginInterface::AppletSize:   //Q_DECL_DEPRECATED
        appletSizeChangedEvent(id);
        break;
    case DockPluginInterface::InfoTypeTitle:
        emit titleChanged(id);
        break;
    case DockPluginInterface::InfoTypeConfigurable:
        emit configurableChanged(id);
        break;
    case DockPluginInterface::InfoTypeEnable:
        emit enabledChanged(id);
        break;
    default:
        break;
    }
}

void DockPluginProxy::itemSizeChangedEvent(QString id)
{

    DockItem * item = m_items.value(id);
    if (item) {
        qWarning() << "Item size changed on plugin " << m_plugin->getPluginName() << id;
        item->adjustSize();
    }
}

void DockPluginProxy::appletSizeChangedEvent(QString id)
{
    DockItem * item = m_items.value(id);
    if (item) {
        qWarning() << "Applet size changed on plugin " << m_plugin->getPluginName() << id;
        item->needPreviewUpdate();
    }
}
