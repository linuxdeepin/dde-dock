/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKPLUGININTERFACE_H
#define DOCKPLUGININTERFACE_H

#include <QPixmap>
#include <QObject>
#include <QStringList>

#include "dockconstants.h"

class DockPluginProxyInterface;
class DockPluginInterface
{
public:
    enum InfoType{
        ItemSize,//Q_DECL_DEPRECATED
        AppletSize,//Q_DECL_DEPRECATED
        Title,//Q_DECL_DEPRECATED
        CanDisable,//Q_DECL_DEPRECATED
        InfoTypeItemSize,
        InfoTypeAppletSize,
        InfoTypeTitle,
        InfoTypeEnable,
        InfoTypeConfigurable
    };

    virtual ~DockPluginInterface() {}

    virtual QString getPluginName() = 0;

    virtual void init(DockPluginProxyInterface *proxy) = 0;
    virtual void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) = 0;

    virtual QStringList ids() = 0;
    virtual QString getName(QString id) = 0;
    virtual QString getTitle(QString id) = 0;
    virtual QString getCommand(QString id) = 0;
    virtual QPixmap getIcon(QString id) {Q_UNUSED(id); return QPixmap("");}
    virtual bool configurable(const QString &id) = 0;
    virtual bool enabled(const QString &id) = 0;
    virtual void setEnabled(const QString &id, bool enabled) = 0;
    virtual QWidget * getItem(QString id) = 0;
    virtual QWidget * getApplet(QString id) = 0;
    virtual QString getMenuContent(QString id) = 0;
    virtual void invokeMenuItem(QString id, QString itemId, bool checked) = 0;
};

QT_BEGIN_NAMESPACE

#define DockPluginInterface_iid "org.deepin.Dock.PluginInterface"

Q_DECLARE_INTERFACE(DockPluginInterface, DockPluginInterface_iid)

QT_END_NAMESPACE

#endif // DOCKPLUGININTERFACE_H
