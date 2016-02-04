/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef SHUTDOWNPLUGIN_H
#define SHUTDOWNPLUGIN_H

#include <QObject>
#include <QLabel>

#include "interfaces/dockplugininterface.h"
#include "interfaces/dockpluginproxyinterface.h"

class ShutdownPlugin : public QObject, public DockPluginInterface
{
    Q_OBJECT
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "dde-dock-shutdown-plugin.json")
    Q_INTERFACES(DockPluginInterface)

public:
    explicit ShutdownPlugin(QObject *parent = 0);

    void init(DockPluginProxyInterface *proxy) Q_DECL_OVERRIDE;
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) Q_DECL_OVERRIDE;
    void invokeMenuItem(QString id, QString itemId, bool checked) Q_DECL_OVERRIDE;
    void setEnabled(const QString &id, bool enabled) Q_DECL_OVERRIDE;
    bool configurable(const QString &id) Q_DECL_OVERRIDE;
    bool enabled(const QString &id) Q_DECL_OVERRIDE;
    QPixmap getIcon(QString id) Q_DECL_OVERRIDE;
    QString getPluginName() Q_DECL_OVERRIDE;
    QString getName(QString id) Q_DECL_OVERRIDE;
    QString getTitle(QString id) Q_DECL_OVERRIDE;
    QString getCommand(QString id) Q_DECL_OVERRIDE;
    QString getMenuContent(QString id) Q_DECL_OVERRIDE;
    QStringList ids() Q_DECL_OVERRIDE;
    QWidget *getItem(QString id) Q_DECL_OVERRIDE;
    QWidget *getApplet(QString id) Q_DECL_OVERRIDE;

private:
    DockPluginProxyInterface *m_proxy;

    QLabel *m_mainWidget;
};

#endif // SHUTDOWNPLUGIN_H
