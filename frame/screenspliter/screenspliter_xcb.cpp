/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
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

ScreenSpliter_Xcb::ScreenSpliter_Xcb(AppItem *appItem, DockEntryInter *entryInter, QObject *parent)
    : ScreenSpliter(appItem, entryInter, parent)
    , m_isSplitCreateWindow(false)
{
    connect(entryInter, &DockEntryInter::WindowInfosChanged,
            this, &ScreenSpliter_Xcb::onUpdateWindowInfo, Qt::QueuedConnection);
}

void ScreenSpliter_Xcb::startSplit(const QRect &rect)
{
    if (!openWindow()) {
        m_effectRect = rect;
        return;
    }
    showSplitScreenEffect(rect, true);
}

bool ScreenSpliter_Xcb::split(ScreenSpliter::SplitDirection direction)
{
    if (!openWindow())
        return false;

    // 如果当前的应用不支持分屏，也无需分屏，检查分屏的时候至少需要一个窗口，因此这里写在打开窗口之后
    quint32 WId = splittingWindowWId();
    if (WId == 0) {
        // 如果当前存在主动打开的窗口，那么就关闭当前主动打开的窗口
        if (m_isSplitCreateWindow) {
            entryInter()->ForceQuit();
            m_isSplitCreateWindow = false;
        }
        return false;
    }

    xcb_client_message_event_t xev;

    xev.response_type = XCB_CLIENT_MESSAGE;
    xev.type = internAtom("_DEEPIN_SPLIT_WINDOW", false);
    xev.window = WId;
    xev.format = 32;
    xev.data.data32[0] = direction_x11(direction); // 1: 左分屏 2: 右分屏 5 左上 6 右上 9 左下 10 右下 15: 全屏
    xev.data.data32[1] = 1;         // 1 进入预览 0 不进入预览

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

bool ScreenSpliter_Xcb::openWindow()
{
    // 查看当前应用是否有打开的窗口，如果没有，则先打开一个窗口
    const WindowInfoMap windowlist = entryInter()->windowInfos();
    if (!windowlist.isEmpty())
        return true;

    if (!m_isSplitCreateWindow) {
        // 如果当前没有打开窗口，且未执行打开操作
        entryInter()->Activate(QX11Info::getTimestamp());
        m_isSplitCreateWindow = true;
    }

    return false;
}

void ScreenSpliter_Xcb::showSplitScreenEffect(const QRect &rect, bool visible)
{
    quint32 WId = splittingWindowWId();
    if (WId == 0)
        return;

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

void ScreenSpliter_Xcb::onUpdateWindowInfo(const WindowInfoMap &info)
{
    // 如果打开的是第一个窗口，且这个打开的窗口是通过拖动二分屏的方式打开，且当前是结束拖拽
    // 并且不支持分屏那么这个窗口就需要关闭
    if (!appItem()->isDragging()) {
        releaseSplit();
    } else if (!m_effectRect.isEmpty() && info.size() > 0) {
        // 只有当需要触发分屏效果的时候，发现当前没有窗口，则记录当前分屏的区域，保存在m_effectRect中
        // 在新增窗口的时候，如果返现m_effectRect有值，则重新触发分屏，并且清空m_effectRect，防止再次打开窗口的时候再次触发分屏效果
        showSplitScreenEffect(m_effectRect, true);
        m_effectRect.setRect(0, 0, 0, 0);
    }
}

bool ScreenSpliter_Xcb::suportSplitScreen()
{
    // 如果当前的窗口的数量为0，则不知道它是否支持分屏，则始终让其返回true,然后打开窗口，因为窗口打开后，
    // 要过一段事件才能收到信号，等收到信号后才知道它是否支持分屏，在窗口显示后会根据当前是否请求过执行分屏操作
    // 来决定是否执行分屏的操作
    if (entryInter()->windowInfos().size() == 0)
        return true;

    return (splittingWindowWId() != 0);
}

bool ScreenSpliter_Xcb::releaseSplit()
{
    showSplitScreenEffect(QRect(), false);
    if (!m_isSplitCreateWindow)
        return false;

    if (!entryInter()->windowInfos().isEmpty() && splittingWindowWId() == 0) {
        // 释放后，如果当前的窗口是通过验证是否支持二分屏的方式来新建的窗口(m_isSplitCreateWindow == true)
        // 并且存在打开的窗口（也有可能不存在打开的窗口，打开的窗口最后才出来，时机上不好控制，所以这种情况
        // 在updateWindowInfos函数里面做了处理），并且打开的窗口不支持二分屏，则此时关闭新打开的窗口
        entryInter()->ForceQuit();
    }

    m_isSplitCreateWindow = false;
    return true;
}

quint32 ScreenSpliter_Xcb::splittingWindowWId()
{
    WindowInfoMap windowsInfo = entryInter()->windowInfos();
    if (windowsInfo.size() == 0)
        return 0;

    quint32 WId = windowsInfo.keys().first();
    xcb_atom_t propAtom = internAtom("_DEEPIN_NET_SUPPORTED", true);
    QByteArray data = windowProperty(WId, propAtom, XCB_ATOM_CARDINAL, 4);

    bool supported = false;
    if (const char *cdata = data.constData())
        supported = *(reinterpret_cast<const quint8 *>(cdata));

    return supported ? WId : 0;
}
