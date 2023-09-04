// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "x11manager.h"
#include "taskmanager.h"
#include "common.h"

#include <QDebug>
#include <QTimer>

/*
 *  使用Xlib监听X Events
 *  使用XCB接口与X进行交互
 * */

#include <ctype.h>
#include <X11/Xos.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xproto.h>

#define XCB XCBUtils::instance()

X11Manager::X11Manager(TaskManager *_taskmanager, QObject *parent)
    : QObject(parent)
    , m_taskmanager(_taskmanager)
    , m_mutex(QMutex(QMutex::NonRecursive))
    , m_listenXEvent(true)
{
    m_rootWindow = XCB->getRootWindow();
}

void X11Manager::listenXEventUseXlib()
{

    Display *dpy;
    int screen;
    char *displayname = nullptr;
    Window w;
    XSetWindowAttributes attr;
    XWindowAttributes wattr;

    dpy = XOpenDisplay (displayname);
    if (!dpy) {
        exit (1);
    }

    screen = DefaultScreen (dpy);
    w = RootWindow(dpy, screen);

    const struct {
        const char *name;
        long mask;
    } events[] = {
    { "keyboard", KeyPressMask | KeyReleaseMask | KeymapStateMask },
    { "mouse", ButtonPressMask | ButtonReleaseMask | EnterWindowMask |
                LeaveWindowMask | PointerMotionMask | Button1MotionMask |
                Button2MotionMask | Button3MotionMask | Button4MotionMask |
                Button5MotionMask | ButtonMotionMask },
    { "button", ButtonPressMask | ButtonReleaseMask },
    { "expose", ExposureMask },
    { "visibility", VisibilityChangeMask },
    { "structure", StructureNotifyMask },
    { "substructure", SubstructureNotifyMask | SubstructureRedirectMask },
    { "focus", FocusChangeMask },
    { "property", PropertyChangeMask },
    { "colormap", ColormapChangeMask },
    { "owner_grab_button", OwnerGrabButtonMask },
    { nullptr, 0 }
};

    long mask = 0;
    for (int i = 0; events[i].name; i++)
        mask |= events[i].mask;

    attr.event_mask = mask;

    XGetWindowAttributes(dpy, w, &wattr);

    attr.event_mask &= ~SubstructureRedirectMask;
    XSelectInput(dpy, w, attr.event_mask);

    while (m_listenXEvent) {
        XEvent event;
        XNextEvent (dpy, &event);

        switch (event.type) {
        case DestroyNotify: {
            XDestroyWindowEvent *eD = (XDestroyWindowEvent *)(&event);
            // qDebug() <<  "DestroyNotify windowId=" << eD->window;

            handleDestroyNotifyEvent(XWindow(eD->window));
            break;
        }
        case MapNotify: {
            XMapEvent *eM = (XMapEvent *)(&event);
            // qDebug() << "MapNotify windowId=" << eM->window;

            handleMapNotifyEvent(XWindow(eM->window));
            break;
        }
        case ConfigureNotify: {
            XConfigureEvent *eC = (XConfigureEvent *)(&event);
            // qDebug() << "ConfigureNotify windowId=" << eC->window;

            handleConfigureNotifyEvent(XWindow(eC->window), eC->x, eC->y, eC->width, eC->height);
            break;
        }
        case PropertyNotify: {
            XPropertyEvent *eP = (XPropertyEvent *)(&event);
            // qDebug() << "PropertyNotify windowId=" << eP->window;

            handlePropertyNotifyEvent(XWindow(eP->window), XCBAtom(eP->atom));
            break;
        }
        case UnmapNotify: {
            // 当松开鼠标的时候会触发该事件，在松开鼠标的时候，需要检测当前窗口是否符合智能隐藏的条件，因此在此处加上该功能
            // 如果不加上该处理，那么就会出现将窗口从任务栏下方移动到屏幕中央的时候，任务栏不隐藏
            handleActiveWindowChangedX();
            break;
        }
        default:
            //qDebug() << "unused event type " << event.type;
            break;
        }
    }

    XCloseDisplay (dpy);
}

void X11Manager::listenXEventUseXCB()
{
    /*
    xcb_get_window_attributes_cookie_t cookie = xcb_get_window_attributes(XCB->getConnect(), XCB->getRootWindow());
    xcb_get_window_attributes_reply_t *reply = xcb_get_window_attributes_reply(XCB->getConnect(), cookie, NULL);
    if (reply) {
        uint32_t valueMask = reply->your_event_mask;
        valueMask &= ~XCB_CW_OVERRIDE_REDIRECT;
        uint32_t mask[2] = {0};
        mask[0] = valueMask;
        //xcb_change_window_attributes(XCB->getConnect(), XCB->getRootWindow(), valueMask, mask);

        free(reply);
    }

    xcb_generic_event_t *event;
    while ( (event = xcb_wait_for_event (XCB->getConnect())) ) {
        eventHandler(event->response_type & ~0x80, event);
    }
    */
}

/**
 * @brief X11Manager::registerWindow 注册X11窗口
 * @param xid
 * @return
 */
WindowInfoX *X11Manager::registerWindow(XWindow xid)
{
    qInfo() << "registWindow: windowId=" << xid;
    WindowInfoX *ret = nullptr;
    do {
        if (m_windowInfoMap.find(xid) != m_windowInfoMap.end()) {
            ret = m_windowInfoMap[xid];
            break;
        }

        WindowInfoX *winInfo = new WindowInfoX(xid);
        if (!winInfo)
            break;

        listenWindowXEvent(winInfo);
        m_windowInfoMap[xid] = winInfo;
        ret = winInfo;
    } while (0);

    return ret;
}

// 取消注册X11窗口
void X11Manager::unregisterWindow(XWindow xid)
{
    qInfo() << "unregisterWindow: windowId=" << xid;
    if (m_windowInfoMap.find(xid) != m_windowInfoMap.end()) {
        m_windowInfoMap.remove(xid);
    }
}

WindowInfoX *X11Manager::findWindowByXid(XWindow xid)
{
    WindowInfoX *ret = nullptr;
    if (m_windowInfoMap.find(xid) != m_windowInfoMap.end())
        ret = m_windowInfoMap[xid];

    return ret;
}

void X11Manager::handleClientListChanged()
{
    QSet<XWindow> newClientList, oldClientList, addClientList, rmClientList;
    for (auto atom : XCB->getClientList())
        newClientList.insert(atom);

    for (auto atom : m_taskmanager->getClientList())
        oldClientList.insert(atom);

    addClientList = newClientList - oldClientList;
    rmClientList = oldClientList - newClientList;
    m_taskmanager->setClientList(newClientList.values());

    // 处理新增窗口
    for (auto xid : addClientList) {
        WindowInfoX *info = registerWindow(xid);
        if (!XCB->isGoodWindow(xid))
            continue;

        uint32_t pid = XCB->getWMPid(xid);
        WMClass wmClass = XCB->getWMClass(xid);
        QString wmName(XCB->getWMName(xid).c_str());
        if (pid != 0 || (wmClass.className.size() > 0 && wmClass.instanceName.size() > 0)
                || wmName.size() > 0 || XCB->getWMCommand(xid).size() > 0) {

            if (info) {
                Q_EMIT requestAttachOrDetachWindow(info);
            }
        }
    }

    // 处理需要移除的窗口
    for (auto xid : rmClientList) {
        WindowInfoX *info = m_windowInfoMap[xid];
        if (info) {
            m_taskmanager->detachWindow(info);
            unregisterWindow(xid);
        } else {
            // no window
            auto entry = m_taskmanager->getEntryByWindowId(xid);
            if (entry && !m_taskmanager->isDocked(entry->getFileName())) {
                m_taskmanager->removeAppEntry(entry);
            }
        }
    }
}

void X11Manager::handleActiveWindowChangedX()
{
    XWindow active = XCB->getActiveWindow();
    WindowInfoX *info = findWindowByXid(active);
    if (info) {
        Q_EMIT requestHandleActiveWindowChange(info);
    }
}

void X11Manager::listenRootWindowXEvent()
{
    uint32_t eventMask = EventMask::XCB_EVENT_MASK_PROPERTY_CHANGE | XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY;
    XCB->registerEvents(m_rootWindow, eventMask);
    handleActiveWindowChangedX();
    handleClientListChanged();
}

/**
 * @brief X11Manager::listenWindowXEvent 监听窗口事件
 * @param winInfo
 */
void X11Manager::listenWindowXEvent(WindowInfoX *winInfo)
{
    uint32_t eventMask = EventMask::XCB_EVENT_MASK_PROPERTY_CHANGE | EventMask::XCB_EVENT_MASK_STRUCTURE_NOTIFY | EventMask::XCB_EVENT_MASK_VISIBILITY_CHANGE;
    XCB->registerEvents(winInfo->getXid(), eventMask);
}

void X11Manager::handleRootWindowPropertyNotifyEvent(XCBAtom atom)
{
    if (atom == XCB->getAtom("_NET_CLIENT_LIST")) {
        // 窗口列表改变
        handleClientListChanged();
    } else if (atom == XCB->getAtom("_NET_ACTIVE_WINDOW")) {
        // 活动窗口改变
        handleActiveWindowChangedX();
    } else if (atom == XCB->getAtom("_NET_SHOWING_DESKTOP")) {
        // 更新任务栏隐藏状态
        Q_EMIT requestUpdateHideState(false);
    }
}

// destory event
void X11Manager::handleDestroyNotifyEvent(XWindow xid)
{
    WindowInfoX *winInfo = findWindowByXid(xid);
    if (!winInfo)
        return;

    m_taskmanager->detachWindow(winInfo);
    unregisterWindow(xid);
}

// map event
void X11Manager::handleMapNotifyEvent(XWindow xid)
{
    WindowInfoX *winInfo = registerWindow(xid);
    if (!winInfo)
        return;

    // TODO QTimer不能在非主线程执行，使用单独线程开发定时器处理非主线程类似定时任务
    //QTimer::singleShot(2 * 1000, this, [=] {
    qInfo() << "handleMapNotifyEvent: pass 2s, now call idnetifyWindow, windowId=" << winInfo->getXid();
    QString innerId;
    AppInfo *appInfo = m_taskmanager->identifyWindow(winInfo, innerId);
    m_taskmanager->markAppLaunched(appInfo);
    //});
}

// config changed event 检测窗口大小调整和重绘应用，触发智能隐藏更新
void X11Manager::handleConfigureNotifyEvent(XWindow xid, int x, int y, int width, int height)
{
    WindowInfoX *winInfo = findWindowByXid(xid);
    if (!winInfo || m_taskmanager->getDockHideMode() != HideMode::SmartHide)
        return;

    WMClass wmClass = winInfo->getWMClass();
    if (wmClass.className.c_str() == frontendWindowWmClass)
        return;     // ignore frontend window ConfigureNotify event

    Q_EMIT requestUpdateHideState(winInfo->isGeometryChanged(x, y, width, height));
}

// property changed event
void X11Manager::handlePropertyNotifyEvent(XWindow xid, XCBAtom atom)
{
    if (xid == m_rootWindow) {
        handleRootWindowPropertyNotifyEvent(atom);
        return;
    }

    WindowInfoX *winInfo = findWindowByXid(xid);
    if (!winInfo)
        return;

    QString newInnerId;
    bool needAttachOrDetach = false;
    if (atom == XCB->getAtom("_NET_WM_STATE")) {
        winInfo->updateWmState();
        needAttachOrDetach = true;
    } else if (atom == XCB->getAtom("_GTK_APPLICATION_ID")) {
        QString gtkAppId;
        winInfo->setGtkAppId(gtkAppId);
        newInnerId = winInfo->genInnerId(winInfo);
    } else if (atom == XCB->getAtom("_NET_WM_PID")) {
        winInfo->updateProcessInfo();
        newInnerId = winInfo->genInnerId(winInfo);
    } else if (atom == XCB->getAtom("_NET_WM_NAME")) {
        winInfo->updateWmName();
        newInnerId = winInfo->genInnerId(winInfo);
    } else if (atom == XCB->getAtom("_NET_WM_ICON")) {
        winInfo->updateIcon();
    } else if (atom == XCB->getAtom("_NET_WM_ALLOWED_ACTIONS")) {
        winInfo->updateWmAllowedActions();
    } else if (atom == XCB->getAtom("_MOTIF_WM_HINTS")) {
        winInfo->updateMotifWmHints();
    } else if (atom == XCB_ATOM_WM_CLASS) {
        winInfo->updateWmClass();
        newInnerId = winInfo->genInnerId(winInfo);
        needAttachOrDetach = true;
    } else if (atom == XCB->getAtom("_XEMBED_INFO")) {
        winInfo->updateHasXEmbedInfo();
        needAttachOrDetach = true;
    } else if (atom == XCB->getAtom("_NET_WM_WINDOW_TYPE")) {
        winInfo->updateWmWindowType();
        needAttachOrDetach = true;
    } else if (atom == XCB_ATOM_WM_TRANSIENT_FOR) {
        winInfo->updateHasWmTransientFor();
        needAttachOrDetach = true;
    }

    if (!newInnerId.isEmpty() && winInfo->getUpdateCalled() && winInfo->getInnerId() != newInnerId) {
        // winInfo.innerId changed
        m_taskmanager->detachWindow(winInfo);
        winInfo->setInnerId(newInnerId);
        needAttachOrDetach = true;
    }

    if (needAttachOrDetach && winInfo) {
        Q_EMIT requestAttachOrDetachWindow(winInfo);
    }

    Entry *entry = m_taskmanager->getEntryByWindowId(xid);
    if (!entry)
        return;

    if (atom == XCB->getAtom("_NET_WM_STATE")) {
        // entry->updateExportWindowInfos();
    } else if (atom == XCB->getAtom("_NET_WM_ICON")) {
        if (entry->getCurrentWindowInfo() == winInfo) {
            entry->updateIcon();
        }
    } else if (atom == XCB->getAtom("_NET_WM_NAME")) {
        if (entry->getCurrentWindowInfo() == winInfo) {
            entry->updateName();
        }
        // entry->updateExportWindowInfos();
    } else if (atom == XCB->getAtom("_NET_WM_ALLOWED_ACTIONS")) {
        entry->updateMenu();
    }
}

void X11Manager::eventHandler(uint8_t type, void *event)
{
    qInfo() << "eventHandler" << "type = " << type;
    switch (type) {
    case XCB_MAP_NOTIFY:    // 17   注册新窗口
        qInfo() << "eventHandler: XCB_MAP_NOTIFY";
        break;
    case XCB_DESTROY_NOTIFY:    // 19   销毁窗口
        qInfo() << "eventHandler: XCB_DESTROY_NOTIFY";
        break;
    case XCB_CONFIGURE_NOTIFY:  // 22   窗口变化
        qInfo() << "eventHandler: XCB_CONFIGURE_NOTIFY";
        break;
    case XCB_PROPERTY_NOTIFY:   // 28   窗口属性改变
        qInfo() << "eventHandler: XCB_PROPERTY_NOTIFY";
        break;
    }
}

void X11Manager::addWindowLastConfigureEvent(XWindow xid, ConfigureEvent *event)
{
    delWindowLastConfigureEvent(xid);

    QMutexLocker locker(&m_mutex);
    QTimer *timer = new QTimer();
    timer->setInterval(configureNotifyDelay);
    m_windowLastConfigureEventMap[xid] = QPair(event, timer);
}

QPair<ConfigureEvent *, QTimer *> X11Manager::getWindowLastConfigureEvent(XWindow xid)
{
    QPair<ConfigureEvent *, QTimer *> ret;
    QMutexLocker locker(&m_mutex);
    if (m_windowLastConfigureEventMap.find(xid) != m_windowLastConfigureEventMap.end())
        ret = m_windowLastConfigureEventMap[xid];

    return ret;
}

void X11Manager::delWindowLastConfigureEvent(XWindow xid)
{
    QMutexLocker locker(&m_mutex);
    if (m_windowLastConfigureEventMap.find(xid) != m_windowLastConfigureEventMap.end()) {
        QPair<ConfigureEvent*, QTimer*> item = m_windowLastConfigureEventMap[xid];
        m_windowLastConfigureEventMap.remove(xid);
        delete item.first;
        item.second->deleteLater();
    }
}
