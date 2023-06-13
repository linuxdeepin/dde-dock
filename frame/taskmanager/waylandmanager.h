// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WAYLANDMANAGER_H
#define WAYLANDMANAGER_H

#include "taskmanager/entry.h"
#include "windowinfok.h"

#include <QObject>
#include <QMap>
#include <QMutex>

class TaskManager;

// 管理wayland窗口
class WaylandManager : public QObject
{
    Q_OBJECT
public:
    explicit WaylandManager(TaskManager *_taskmanager, QObject *parent = nullptr);

    void registerWindow(const QString &objPath);
    void unRegisterWindow(const QString &objPath);

    WindowInfoK *findWindowById(uint activeWin);
    WindowInfoK *findWindowByXid(XWindow xid);
    WindowInfoK *findWindowByObjPath(QString objPath);
    void insertWindow(QString objPath, WindowInfoK *windowInfo);
    void deleteWindow(QString objPath);

private:
    TaskManager *m_taskmanager;
    QMap<QString, WindowInfoK *> m_kWinInfos;       // dbusObjectPath -> kwayland window Info
    QMap<XWindow, WindowInfoK *> m_windowInfoMap;
    QMutex m_mutex;
};

#endif // WAYLANDMANAGER_H
