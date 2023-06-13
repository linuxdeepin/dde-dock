// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef ENTRIES_H
#define ENTRIES_H

#include "entry.h"
#include "constants.h"
#include "taskmanager/windowinfobase.h"

#include <QVector>
#include <QWeakPointer>
#include <qlist.h>

#define MAX_UNOPEN_RECENT_COUNT 3

class TaskManager;

// 所有应用管理类
class Entries
{
public:
    Entries(TaskManager *_taskmanager);

    const QList<Entry *> unDockedEntries() const;

    bool shouldInRecent();

    void removeLastRecent();
    void updateShowRecent();
    void updateEntriesMenu();
    void append(Entry *entry);
    void remove(Entry *entry);
    void moveEntryToLast(Entry *entry);
    void insert(Entry *entry, int index);
    void move(int oldIndex, int newIndex);
    void setDisplayMode(Dock::DisplayMode displayMode);
    void handleActiveWindowChanged(XWindow activeWindId);

    QString queryWindowIdentifyMethod(XWindow windowId);
    QStringList getEntryIDs();

    Entry *getByWindowPid(int pid);
    Entry *getByInnerId(QString innerId);
    Entry *getByWindowId(XWindow windowId);
    Entry *getByDesktopFilePath(const QString &filePath);
    Entry *getDockedEntryByDesktopFile(const QString &desktopFile);

    QList<Entry*> getEntries();
    QVector<Entry *> filterDockedEntries();

private:
    void insertCb(Entry *entry, int index);
    void removeCb(Entry *entry);

private:
    QList<Entry *> m_items;
    TaskManager *m_taskmanager;
};

#endif // ENTRIES_H
