/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKMODEDATA_H
#define DOCKMODEDATA_H

#include <QObject>
#include <QStringList>
#include <QTimer>
#include <QDebug>
#include "dbus/dbusdocksetting.h"
#include "dbus/dbushidestatemanager.h"
#include "interfaces/dockconstants.h"

class DockModeData : public QObject
{
    Q_OBJECT
public:
    static DockModeData * instance();

    Dock::DockMode getDockMode();
    void setDockMode(Dock::DockMode value);
    Dock::HideMode getHideMode();
    void setHideMode(Dock::HideMode value);
    Dock::HideState getHideState();
    void setHideState(Dock::HideState value);

    int getDockHeight();
    int getItemHeight();
    int getNormalItemWidth();
    int getActivedItemWidth();
    int getAppItemSpacing();
    int getAppIconSize();

    int getAppletsItemHeight();
    int getAppletsItemWidth();
    int getAppletsItemSpacing();
    int getAppletsIconSize();

signals:
    void dockModeChanged(Dock::DockMode newMode,Dock::DockMode oldMode);
    void hideModeChanged(Dock::HideMode newMode,Dock::HideMode oldMode);

private:
    explicit DockModeData(QObject *parent = 0);

    void initDDS();
    void onDockModeChanged(int mode);
    void onHideModeChanged(int mode);

private:
    static DockModeData * m_dockModeData;

    Dock::DockMode m_currentMode = Dock::EfficientMode;
    Dock::HideMode m_hideMode = Dock::KeepShowing;
    Dock::HideState m_hideState = Dock::HideStateShown; //make sure the preview can be show

    DBusDockSetting *m_dockSetting = NULL;
    DBusHideStateManager *m_hideStateManager = new DBusHideStateManager(this);
};

#endif // DOCKMODEDATA_H
