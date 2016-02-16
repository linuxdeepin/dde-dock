/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef APPMANAGER_H
#define APPMANAGER_H

#include "dbus/dbusentrymanager.h"
#include "dbus/dbusdockedappmanager.h"
#include "widgets/old/appitem.h"
#include "widgets/old/launcheritem.h"

class AppManager : public QObject
{
    Q_OBJECT
public:
    explicit AppManager(QObject *parent = 0);
    void initEntries();

    void setDockingItemId(const QString &dockingItemId);

signals:
    void requestSort();
    void entryAdded(AbstractDockItem *item, bool delayShow);
    void entryAppend(AbstractDockItem *item, bool delayShow);
    void entryRemoved(const QString &id);

private:
    void sortItemList();    //Sort and append item to dock
    void onEntryRemoved(const QString &id);
    void onEntryAdded(const QDBusObjectPath &path);

private:
    QStringList m_ids;
    QString m_dockingItemId = "";   //drop to the dock one
    DBusEntryManager *m_entryManager = NULL;
    DBusDockedAppManager *m_dockAppManager = new DBusDockedAppManager(this);
    QMap<QString, AbstractDockItem *> m_initItemList; //Juse for initialization <id, item>
};

#endif // APPMANAGER_H
