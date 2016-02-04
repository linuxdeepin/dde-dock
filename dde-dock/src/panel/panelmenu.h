/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PANELMENU_H
#define PANELMENU_H

#include <QWidget>
#include <QLabel>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QDebug>

#include "dbus/dbusmenumanager.h"
#include "dbus/dbusmenu.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"

class PanelMenu : public QObject
{
    Q_OBJECT
public:
    enum OperationType {
        ToFashionMode,
        ToEfficientMode,
        ToClassicMode,
        ToKeepShowing,
        ToKeepHidden,
        ToSmartHide,
        ToPluginSetting
    };
    enum MenuGroup{
        DockModeGroup,
        HideModeGroup
    };

    static PanelMenu * instance();

    void showMenu(int x,int y);

signals:
    void settingPlugin();
    void menuItemInvoked();

private:
    explicit PanelMenu(QObject *parent = 0);

    void changeToFashionMode();
    void changeToEfficientMode();
    void changeToClassicMode();
    void changeToKeepShowing();
    void changeToKeepHidden();
    void changeToSmartHide();

    void onItemInvoked(const QString &itemId, bool result);

    QJsonObject createItemObj(const QString &itemName, OperationType type);
    QJsonObject createRadioItemObj(const QString &itemName, OperationType type, MenuGroup group, bool check);

private:
    static PanelMenu * m_panelMenu;
    QString m_menuInterfacePath = "";
    DBusDockSetting m_dockSetting;
    DBusMenuManager *m_menuManager = NULL;
    DockModeData *m_dockModeData = DockModeData::instance();

};

#endif // PANELMENU_H
