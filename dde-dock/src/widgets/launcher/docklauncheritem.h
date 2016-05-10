/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKLAUNCHERITEM_H
#define DOCKLAUNCHERITEM_H

#include <QWidget>

#include "../dockitem.h"
#include "../app/dockappicon.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"
#include "dbus/dbuslaunchercontroller.h"

class QProcess;

class DockLauncherItem : public DockItem
{
    Q_OBJECT
public:
    explicit DockLauncherItem(QWidget *parent = 0);
    ~DockLauncherItem();

    QString getItemId() Q_DECL_OVERRIDE {return "dde-launcher";}
    QString getTitle() Q_DECL_OVERRIDE { return tr("Launcher"); }
    QWidget * getApplet() Q_DECL_OVERRIDE { return NULL; }

protected:
    void enterEvent(QEvent *) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void mouseReleaseEvent(QMouseEvent *event) Q_DECL_OVERRIDE;

private slots:
    void slotMousePress(QMouseEvent *event);
    void slotMouseRelease(QMouseEvent *event);
    void updateIcon();

private:
    void changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode);
    void reanchorIcon();

private:
    DockAppIcon * m_appIcon;
    DBusLauncherController *m_launcherInter;
    QString m_menuInterfacePath = "";
    DockModeData * m_dockModeData = DockModeData::instance();
};

#endif // DOCKLAUNCHERITEM_H
