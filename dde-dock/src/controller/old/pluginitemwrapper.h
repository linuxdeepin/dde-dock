/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PLUGINITEMWRAPPER_H
#define PLUGINITEMWRAPPER_H

#include <QJsonObject>

#include "widgets/old/abstractdockitem.h"
#include "interfaces/dockplugininterface.h"
#include "dbus/dbusdisplay.h"

class QMouseEvent;
class PluginItemWrapper : public AbstractDockItem
{
    Q_OBJECT
public:
    PluginItemWrapper(DockPluginInterface *plugin, QString id, QWidget * parent = 0);
    virtual ~PluginItemWrapper();

    QString id() const;

    QString getTitle() Q_DECL_OVERRIDE;
    QWidget * getApplet() Q_DECL_OVERRIDE;

    QString getMenuContent() Q_DECL_OVERRIDE;
    void invokeMenuItem(QString itemId, bool checked) Q_DECL_OVERRIDE;

protected:
    void enterEvent(QEvent * event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent * event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent * event) Q_DECL_OVERRIDE;
//    void mouseReleaseEvent(QMouseEvent * event) Q_DECL_OVERRIDE;

private:
    DBusDisplay *m_display = NULL;
    QWidget *m_pluginItemContents = NULL;
    DockPluginInterface * m_plugin;
    QString m_id;

    QJsonObject createMenuItem(QString itemId, QString itemName, bool checkable, bool checked);
    const int DOCK_PREVIEW_MARGIN = 7;
};

#endif // PLUGINITEMWRAPPER_H
