/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKPLUGINLAYOUT_H
#define DOCKPLUGINLAYOUT_H

#include "../movablelayout.h"
#include "../../dbus/dbusdisplay.h"
#include "../../controller/plugins/dockpluginsmanager.h"

class DockPluginLayout : public MovableLayout
{
    Q_OBJECT
public:
    explicit DockPluginLayout(QWidget *parent = 0);

    QSize sizeHint() const;
    void initAllPlugins();

signals:
    void needPreviewHide(bool immediately);
    void needPreviewShow(DockItem *item, QPoint pos);
    void needPreviewUpdate();
    void itemHoverableChange(bool v);
    void pluginsInitDone();

protected:
    void insertWidget(const int index, DockItem *widget);

private:
    void initPluginManager();
    DisplayRect getScreenRect();

private:
    DockPluginsManager *m_pluginManager;
    QTimer *m_pluginLoadFinishedTimer;
};

#endif // DOCKPLUGINLAYOUT_H
