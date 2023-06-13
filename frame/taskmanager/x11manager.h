// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef X11MANAGER_H
#define X11MANAGER_H

#include "windowinfox.h"
#include "xcbutils.h"

#include <QObject>
#include <QMap>
#include <QMutex>
#include <QTimer>

class TaskManager;

class X11Manager : public QObject
{
    Q_OBJECT
public:
    explicit X11Manager(TaskManager *_taskmanager, QObject *parent = nullptr);

    WindowInfoX *findWindowByXid(XWindow xid);
    WindowInfoX *registerWindow(XWindow xid);
    void unregisterWindow(XWindow xid);

    void handleClientListChanged();
    void handleActiveWindowChangedX();
    void listenRootWindowXEvent();
    void listenWindowXEvent(WindowInfoX *winInfo);

    void handleRootWindowPropertyNotifyEvent(XCBAtom atom);
    void handleDestroyNotifyEvent(XWindow xid);
    void handleMapNotifyEvent(XWindow xid);
    void handleConfigureNotifyEvent(XWindow xid, int x, int y, int width, int height);
    void handlePropertyNotifyEvent(XWindow xid, XCBAtom atom);

    void eventHandler(uint8_t type, void *event);
    void listenWindowEvent(WindowInfoX *winInfo);
    void listenXEventUseXlib();
    void listenXEventUseXCB();

Q_SIGNALS:
    void requestUpdateHideState(bool delay);
    void requestHandleActiveWindowChange(WindowInfoBase *info);
    void requestAttachOrDetachWindow(WindowInfoBase *info);

private:
    void addWindowLastConfigureEvent(XWindow xid, ConfigureEvent* event);
    QPair<ConfigureEvent*, QTimer*> getWindowLastConfigureEvent(XWindow xid);
    void delWindowLastConfigureEvent(XWindow xid);

private:
    QMap<XWindow, WindowInfoX *> m_windowInfoMap;
    TaskManager *m_taskmanager;
    QMap<XWindow, QPair<ConfigureEvent*, QTimer*>> m_windowLastConfigureEventMap; // 手动回收ConfigureEvent和QTimer
    QMutex m_mutex;
    XWindow m_rootWindow;                                                         // 根窗口
    bool m_listenXEvent;                                                          // 监听X事件
};

#endif // X11MANAGER_H
