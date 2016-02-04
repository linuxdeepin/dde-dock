/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKPLUGINPROXYINTERFACE_H
#define DOCKPLUGINPROXYINTERFACE_H

#include <QString>

#include "dockconstants.h"
#include "dockplugininterface.h"

class DockPluginProxyInterface
{
public:
    virtual Dock::DockMode dockMode() = 0;

    virtual void itemAddedEvent(QString id) = 0;
    virtual void itemRemovedEvent(QString id) = 0;
    virtual void infoChangedEvent(DockPluginInterface::InfoType type, const QString &id) = 0;
};

#endif // DOCKPLUGINPROXYINTERFACE_H
