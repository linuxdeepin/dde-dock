/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PLUGINPROXY_H
#define PLUGINPROXY_H

#include "widgets/old/abstractdockitem.h"
#include "interfaces/dockplugininterface.h"
#include "interfaces/dockpluginproxyinterface.h"

class QPluginLoader;
class PluginProxy : public QObject, public DockPluginProxyInterface
{
    Q_OBJECT
public:
    PluginProxy(QPluginLoader * loader, DockPluginInterface * plugin);
    ~PluginProxy();

    bool isSystemPlugin();
    DockPluginInterface * plugin();

    Dock::DockMode dockMode() Q_DECL_OVERRIDE;

    void itemAddedEvent(QString id) Q_DECL_OVERRIDE;
    void itemRemovedEvent(QString id) Q_DECL_OVERRIDE;
    void infoChangedEvent(DockPluginInterface::InfoType type, const QString &id) Q_DECL_OVERRIDE;

signals:
    void itemAdded(AbstractDockItem * item, QString uuid);
    void itemRemoved(AbstractDockItem * item, QString uuid);
    void titleChanged(const QString &id);
    void configurableChanged(QString id);
    void enabledChanged(QString id);

private:
    QMap<QString, AbstractDockItem*> m_items;

    QPluginLoader * m_loader;
    DockPluginInterface * m_plugin;

    void itemSizeChangedEvent(QString id);
    void appletSizeChangedEvent(QString id);
};

#endif // PLUGINPROXY_H
