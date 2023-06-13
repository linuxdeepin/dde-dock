// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "screenspliter_xcb.h"
#include "appitem.h"

#include <QX11Info>

#include <X11/X.h>
#include <X11/Xlib.h>
#include <xcb/xproto.h>

#define LEFT            1
#define RIGHT           2
#define TOP             3
#define BOTTOM          4
#define LEFTTOP         5
#define RIGHTTOP        6
#define LEFTBOTTOM      9
#define RIGHTBOTTOM     10
#define SPLITUNKNOW     0

static xcb_atom_t internAtom(const char *name, bool only_if_exist)
{
    if (!name || *name == 0)
        return XCB_NONE;

    xcb_connection_t *connection = QX11Info::connection();
    xcb_intern_atom_cookie_t cookie = xcb_intern_atom(connection, only_if_exist, strlen(name), name);
    xcb_intern_atom_reply_t *reply = xcb_intern_atom_reply(connection, cookie, 0);
    if (!reply)
        return XCB_NONE;

    xcb_atom_t atom = reply->atom;
    free(reply);

    return atom;
}

static QByteArray windowProperty(quint32 WId, xcb_atom_t propAtom, xcb_atom_t typeAtom, quint32 len)
{
    xcb_connection_t *conn = QX11Info::connection();
    xcb_get_property_cookie_t cookie = xcb_get_property(conn, false, WId, propAtom, typeAtom, 0, len);
    xcb_generic_error_t *err = nullptr;
    xcb_get_property_reply_t *reply = xcb_get_property_reply(conn, cookie, &err);

    QByteArray data;
    if (reply != nullptr) {
        int valueLen = xcb_get_property_value_length(reply);
        const char *buf = static_cast<const char *>(xcb_get_property_value(reply));
        data.append(buf, valueLen);
        free(reply);
    }

    if (err != nullptr) {
        free(err);
    }

    return data;
}

ScreenSpliter_Xcb::ScreenSpliter_Xcb(AppItem *appItem, QObject *parent)
    : ScreenSpliter(appItem, parent)
{
}

void ScreenSpliter_Xcb::startSplit(const QRect &rect)
{
    if (!suportSplitScreen())
        return;

    showSplitScreenEffect(rect, true);
}

bool ScreenSpliter_Xcb::split(ScreenSpliter::SplitDirection direction)
{
    if (!suportSplitScreen())
        return false;

    quint32 WId = appItem()->windowsInfos().keys().first();
    xcb_client_message_event_t xev;
    xev.response_type = XCB_CLIENT_MESSAGE;
    xev.type = internAtom("_DEEPIN_SPLIT_WINDOW", false);
    xev.window = WId;
    xev.format = 32;
    xev.data.data32[0] = direction_x11(direction);  // 1: 左分屏 2: 右分屏 5 左上 6 右上 9 左下 10 右下 15: 全屏
    xev.data.data32[1] = 1;                         // 1 进入预览 0 不进入预览

    xcb_send_event(QX11Info::connection(), false, QX11Info::appRootWindow(QX11Info::appScreen()),
                   SubstructureNotifyMask, (const char *)&xev);
    xcb_flush(QX11Info::connection());

    return true;
}

uint32_t ScreenSpliter_Xcb::direction_x11(ScreenSpliter::SplitDirection direction)
{
    static QMap<ScreenSpliter::SplitDirection, int> directionMapping = {
        { ScreenSpliter::Left, LEFT },
        { ScreenSpliter::Right, RIGHT },
        { ScreenSpliter::Top, TOP },
        { ScreenSpliter::Bottom, TOP },
        { ScreenSpliter::LeftTop, LEFTTOP },
        { ScreenSpliter::RightTop, RIGHTTOP },
        { ScreenSpliter::LeftBottom, LEFTBOTTOM },
        { ScreenSpliter::RightBottom, RIGHTBOTTOM }
    };

    return directionMapping.value(direction, SPLITUNKNOW);
}

void ScreenSpliter_Xcb::showSplitScreenEffect(const QRect &rect, bool visible)
{
    if (!suportSplitScreen())
        return;

    quint32 WId = appItem()->windowsInfos().keys().first();
    // 触发分屏的效果
    xcb_client_message_event_t xev;
    xev.response_type = XCB_CLIENT_MESSAGE;
    xev.type = internAtom("_DEEPIN_SPLIT_OUTLINE", false);
    xev.window = WId;
    xev.format = 32;
    xev.data.data32[0] = visible ? 1 : 0;                 // 1: 显示 0: 取消
    xev.data.data32[1] = rect.x();                        // X坐标
    xev.data.data32[2] = rect.y();                        // Y坐标
    xev.data.data32[3] = rect.width();                    // width
    xev.data.data32[4] = rect.height();                   // height

    xcb_send_event(QX11Info::connection(), false, QX11Info::appRootWindow(QX11Info::appScreen()),
                   SubstructureNotifyMask, (const char *)&xev);
    xcb_flush(QX11Info::connection());
}

bool ScreenSpliter_Xcb::suportSplitScreen()
{
    // 判断所有的窗口，只要有一个窗口支持分屏，就认为它支持分屏
    QList<quint32> winIds = appItem()->windowsInfos().keys();
    for (const quint32 &winId : winIds) {
        if (windowSupportSplit(winId))
            return true;
    }

    return false;
}

bool ScreenSpliter_Xcb::releaseSplit()
{
    showSplitScreenEffect(QRect(), false);
    return true;
}

bool ScreenSpliter_Xcb::windowSupportSplit(quint32 winId)
{
    xcb_atom_t propAtom = internAtom("_DEEPIN_NET_SUPPORTED", true);
    QByteArray data = windowProperty(winId, propAtom, XCB_ATOM_CARDINAL, 4);

    bool supported = false;
    if (const char *cdata = data.constData())
        supported = *(reinterpret_cast<const quint8 *>(cdata));

    return supported;
}
