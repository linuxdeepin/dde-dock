/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QIcon>

#include "trashplugin.h"

TrashPlugin::TrashPlugin()
{
    QIcon::setThemeName("deepin");

    m_item = new MainItem();
    connect(this, &TrashPlugin::menuItemInvoked, m_item, &MainItem::emptyTrash);
}

void TrashPlugin::init(DockPluginProxyInterface *proxy)
{
    m_proxy = proxy;

    setMode(proxy->dockMode());
}

QString TrashPlugin::getPluginName()
{
    return tr("Trash");
}

QStringList TrashPlugin::ids()
{
    return QStringList(m_id);
}

QString TrashPlugin::getName(QString)
{
    return getPluginName();
}

QString TrashPlugin::getTitle(QString)
{
    return getPluginName();
}

QString TrashPlugin::getCommand(QString)
{
    return "";
}

bool TrashPlugin::configurable(const QString &)
{
    return false;
}

bool TrashPlugin::enabled(const QString &)
{
    return true;
}

void TrashPlugin::setEnabled(const QString &, bool)
{

}

QWidget * TrashPlugin::getItem(QString)
{
    return m_item;
}

QWidget * TrashPlugin::getApplet(QString)
{
    return NULL;
}

void TrashPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode != oldMode)
        setMode(newMode);
}

QString TrashPlugin::getMenuContent(QString)
{
    QJsonObject contentObj;

    QJsonArray items;

    items.append(createMenuItem("clear_trash", tr("Empty")));

    contentObj.insert("items", items);

    return QString(QJsonDocument(contentObj).toJson());
}

void TrashPlugin::invokeMenuItem(QString, QString itemId, bool checked)
{
    qWarning() << "Menu check:" << itemId << checked;
    emit menuItemInvoked();
}

// private methods
void TrashPlugin::setMode(Dock::DockMode mode)
{
    m_mode = mode;

    if (mode == Dock::FashionMode)
        m_proxy->itemAddedEvent(m_id);
    else{
        m_proxy->itemRemovedEvent(m_id);
        m_item->setParent(NULL);
    }
}

QJsonObject TrashPlugin::createMenuItem(QString itemId, QString itemName, bool checkable, bool checked)
{
    QJsonObject itemObj;

    itemObj.insert("itemId", itemId);
    itemObj.insert("itemText", itemName);
    itemObj.insert("itemIcon", "");
    itemObj.insert("itemIconHover", "");
    itemObj.insert("itemIconInactive", "");
    itemObj.insert("itemExtra", "");
    itemObj.insert("isActive", m_trashMonitor->ItemCount() > 0);
    itemObj.insert("isCheckable", checkable);
    itemObj.insert("checked", checked);
    itemObj.insert("itemSubMenu", QJsonObject());

    return itemObj;
}

TrashPlugin::~TrashPlugin()
{

}

#if QT_VERSION < 0x050000
Q_EXPORT_PLUGIN2(dde-dock-trash-plugin, TrashPlugin)
#endif // QT_VERSION < 0x050000
