// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "xcbutils.h"

#include <cstdint>
#include <utility>
#include <iostream>
#include <cstring>
#include <memory>
#include <algorithm>

#include <X11/Xlib.h>
#include <X11/extensions/XRes.h>

XCBUtils::XCBUtils()
{
    m_connect = xcb_connect(nullptr, &m_screenNum); // nullptr表示默认使用环境变量$DISPLAY获取屏幕
    if (xcb_connection_has_error(m_connect)) {
        std::cout << "XCBUtils: init xcb_connect error" << std::endl;
        return;
    }

    if (!xcb_ewmh_init_atoms_replies(&m_ewmh,
                                     xcb_ewmh_init_atoms(m_connect, &m_ewmh),   // 初始化Atom
                                     nullptr))
        std::cout << "XCBUtils: init ewmh  error" << std::endl;
}

XCBUtils::~XCBUtils()
{
    if (m_connect) {
        xcb_disconnect(m_connect);    // 关闭连接并释放
        m_connect = nullptr;
    }
}

XWindow XCBUtils::allocId()
{
    return xcb_generate_id(m_connect);
}

void XCBUtils::flush()
{
    xcb_flush(m_connect);
}

void XCBUtils::killClientChecked(XWindow xid)
{
    xcb_kill_client_checked(m_connect, xid);
}

xcb_get_property_reply_t *XCBUtils::getPropertyValueReply(XWindow xid, XCBAtom property, XCBAtom type)
{
    xcb_get_property_cookie_t cookie = xcb_get_property(m_connect,
                                                        0,
                                                        xid,
                                                        property,
                                                        type,
                                                        0,
                                                        MAXLEN);
    return xcb_get_property_reply(m_connect, cookie, nullptr);
}

void *XCBUtils::getPropertyValue(XWindow xid, XCBAtom property, XCBAtom type)
{
    void *value = nullptr;
    xcb_get_property_reply_t *reply = getPropertyValueReply(xid, property, type);
    if (reply) {
        if (xcb_get_property_value_length(reply) > 0) {
            value = xcb_get_property_value(reply);
        }
        free(reply);
    }
    return value;
}

std::string XCBUtils::getUTF8PropertyStr(XWindow xid, XCBAtom property)
{
    std::string ret;
    xcb_get_property_reply_t *reply = getPropertyValueReply(xid, property, m_ewmh.UTF8_STRING);
    if (reply) {
        ret = getUTF8StrFromReply(reply);
        free(reply);
    }
    return ret;
}

XCBAtom XCBUtils::getAtom(const char *name)
{
    XCBAtom ret = m_atomCache.getVal(name);
    if (ret == ATOMNONE) {
        xcb_intern_atom_cookie_t cookie = xcb_intern_atom(m_connect, false, strlen(name), name);
        std::shared_ptr<xcb_intern_atom_reply_t> reply(xcb_intern_atom_reply(m_connect, cookie, nullptr), [=](xcb_intern_atom_reply_t* reply){
            free(reply);}
        );
        if (reply) {
            m_atomCache.store(name, xcb_atom_t(reply->atom));
            ret = reply->atom;
        }
    }

    return ret;
}

std::string XCBUtils::getAtomName(XCBAtom atom)
{
    std::string ret = m_atomCache.getName(atom);
    if (ret.empty()) {
        xcb_get_atom_name_cookie_t cookie = xcb_get_atom_name(m_connect, atom);
        std::shared_ptr<xcb_get_atom_name_reply_t> reply(
            xcb_get_atom_name_reply(m_connect, cookie, nullptr),
            [=](xcb_get_atom_name_reply_t* reply) {free(reply);});
        if (reply) {
            char *name = xcb_get_atom_name_name(reply.get());
            if (name) {
                m_atomCache.store(name, atom);
                ret = name;
            }
        }
    }

    return ret;
}

Geometry XCBUtils::getWindowGeometry(XWindow xid)
{
    xcb_get_geometry_cookie_t cookie = xcb_get_geometry(m_connect, xcb_drawable_t(xid));
    std::shared_ptr<xcb_get_geometry_reply_t> reply(
        xcb_get_geometry_reply(m_connect, cookie, nullptr),
        [=](xcb_get_geometry_reply_t* reply){free(reply);}
    );
    if (!reply) {
        std::cout << xid << " getWindowGeometry err" << std::endl;
        return Geometry();
    }

    Geometry ret;
    ret.x = reply->x;
    ret.y = reply->y;
    ret.width = reply->width;
    ret.height = reply->height;

    const xcb_setup_t *xcbSetup = xcb_get_setup(m_connect);
    if (!xcbSetup)
        return Geometry();

    xcb_screen_iterator_t xcbScreenIterator = xcb_setup_roots_iterator(xcbSetup);
    std::shared_ptr<xcb_translate_coordinates_reply_t> translateReply(
        xcb_translate_coordinates_reply(m_connect, 
            xcb_translate_coordinates(m_connect, xid, xcbScreenIterator.data->root, 0, 0),
            nullptr),
        [=](xcb_translate_coordinates_reply_t* translateReply){free(translateReply);});

    if (translateReply) {
        ret.x = translateReply->dst_x;
        ret.y = translateReply->dst_y;
    }

    XWindow dWin = getDecorativeWindow(xid);
    reply.reset(xcb_get_geometry_reply(m_connect, xcb_get_geometry(m_connect, xcb_drawable_t(dWin)), nullptr),
                [=](xcb_get_geometry_reply_t* reply){free(reply);});
    if (!reply)
        return ret;

    if (reply->x == ret.x && reply->y == ret.y) {
        // 无标题的窗口，比如deepin-editor, dconf-editor等
        WindowFrameExtents windowFrameRect = getWindowFrameExtents(xid);
        if (!windowFrameRect.isNull()) {
            int x = ret.x + windowFrameRect.Left;
            int y = ret.y + windowFrameRect.Top;
            int width = ret.width - (windowFrameRect.Left + windowFrameRect.Right);
            int height = ret.height - (windowFrameRect.Top + windowFrameRect.Bottom);
            ret.x = x;
            ret.y = y;
            ret.width = width;
            ret.height = height;
        }
    }

    return ret;
}

XWindow XCBUtils::getDecorativeWindow(XWindow xid)
{
    XWindow winId = xid;
    for (int i = 0; i < 10; i++) {
        xcb_query_tree_cookie_t cookie = xcb_query_tree(m_connect, winId);
        std::shared_ptr<xcb_query_tree_reply_t> reply(
            xcb_query_tree_reply(m_connect, cookie, nullptr),
            [=](xcb_query_tree_reply_t* reply){free(reply);}
        );
        if (!reply) return 0;
        if (reply->root == reply->parent) return winId;

        winId = reply->parent;
    }

    return 0;
}

WindowFrameExtents XCBUtils::getWindowFrameExtents(XWindow xid)
{
    xcb_atom_t perp = getAtom("_NET_FRAME_EXTENTS");
    xcb_get_property_cookie_t cookie = xcb_get_property(m_connect, false, xid, perp, XCB_ATOM_CARDINAL, 0, 4);
    std::shared_ptr<xcb_get_property_reply_t> reply(
        xcb_get_property_reply(m_connect, cookie, nullptr),
        [=](xcb_get_property_reply_t* reply){free(reply);}
    );
    if (!reply || reply->format == 0) {
        perp = getAtom("_GTK_FRAME_EXTENTS");
        cookie = xcb_get_property(m_connect, false, xid, perp, XCB_ATOM_CARDINAL, 0, 4);
        reply.reset(xcb_get_property_reply(m_connect, cookie, nullptr), [=](xcb_get_property_reply_t* reply){free(reply);});
        if (!reply)
            return WindowFrameExtents();
    }

    if (reply->format != 32 || reply->value_len != 4) {
        return WindowFrameExtents();
    }
    uint32_t *data = static_cast<uint32_t *>(xcb_get_property_value(reply.get()));

    if (!data)
        return WindowFrameExtents();

    WindowFrameExtents winFrame(data[0], data[1], data[2], data[3]);
    return winFrame;
}

XWindow XCBUtils::getActiveWindow()
{
    XWindow ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_active_window(&m_ewmh, m_screenNum);
    if (!xcb_ewmh_get_active_window_reply(&m_ewmh, cookie, &ret, nullptr)) {
        std::cout << "getActiveWindow error" << std::endl;
    }

    return ret;
}

void XCBUtils::setActiveWindow(XWindow xid)
{
    xcb_ewmh_set_active_window(&m_ewmh, m_screenNum, xid);
}

void XCBUtils::changeActiveWindow(XWindow newActiveXid)
{
    xcb_ewmh_request_change_active_window(&m_ewmh, m_screenNum, newActiveXid, XCB_EWMH_CLIENT_SOURCE_TYPE_OTHER, XCB_CURRENT_TIME, XCB_WINDOW_NONE);
    flush();
}

void XCBUtils::restackWindow(XWindow xid)
{
    xcb_ewmh_request_restack_window(&m_ewmh, m_screenNum, xid, 0, XCB_STACK_MODE_ABOVE);
}

std::list<XWindow> XCBUtils::getClientList()
{
    std::list<XWindow> ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_client_list(&m_ewmh, m_screenNum);
    xcb_ewmh_get_windows_reply_t reply;
    if (xcb_ewmh_get_client_list_reply(&m_ewmh, cookie, &reply, nullptr)) {
        for (uint32_t i = 0; i < reply.windows_len; i++) {
            ret.push_back(reply.windows[i]);
        }

        xcb_ewmh_get_windows_reply_wipe(&reply);
    } else {
        std::cout << "getClientList error" << std::endl;
    }

    return ret;
}

std::list<XWindow> XCBUtils::getClientListStacking()
{
    std::list<XWindow> ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_client_list_stacking(&m_ewmh, m_screenNum);
    xcb_ewmh_get_windows_reply_t reply;
    if (xcb_ewmh_get_client_list_stacking_reply(&m_ewmh, cookie, &reply, nullptr)) {
        for (uint32_t i = 0; i < reply.windows_len; i++) {
            ret.push_back(reply.windows[i]);
        }

        xcb_ewmh_get_windows_reply_wipe(&reply);
    } else {
        std::cout << "getClientListStacking error" << std::endl;
    }

    return ret;
}

std::vector<XCBAtom> XCBUtils::getWMState(XWindow xid)
{
    std::vector<XCBAtom> ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_state(&m_ewmh, xid);
    xcb_ewmh_get_atoms_reply_t reply; // a list of Atom
    if (xcb_ewmh_get_wm_state_reply(&m_ewmh, cookie, &reply, nullptr)) {
        for (uint32_t i = 0; i < reply.atoms_len; i++) {
            ret.push_back(reply.atoms[i]);
        }

        xcb_ewmh_get_atoms_reply_wipe(&reply);
    } else {
        std::cout << xid << " getWMState error" << std::endl;
    }

    return ret;
}

std::vector<XCBAtom> XCBUtils::getWMWindoType(XWindow xid)
{
    std::vector<XCBAtom> ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_window_type(&m_ewmh, xid);
    xcb_ewmh_get_atoms_reply_t reply; // a list of Atom
    if (xcb_ewmh_get_wm_window_type_reply(&m_ewmh, cookie, &reply, nullptr)) {
        for (uint32_t i = 0; i < reply.atoms_len; i++) {
            ret.push_back(reply.atoms[i]);
        }

        xcb_ewmh_get_atoms_reply_wipe(&reply);
    } else {
        std::cout << xid << " getWMWindoType error" << std::endl;
    }

    return ret;
}

std::vector<XCBAtom> XCBUtils::getWMAllowedActions(XWindow xid)
{
    std::vector<XCBAtom> ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_allowed_actions(&m_ewmh, xid);
    xcb_ewmh_get_atoms_reply_t reply;   // a list of Atoms
    if (xcb_ewmh_get_wm_allowed_actions_reply(&m_ewmh, cookie, &reply, nullptr)) {
        for (uint32_t i = 0; i < reply.atoms_len; i++) {
            ret.push_back(reply.atoms[i]);
        }

        xcb_ewmh_get_atoms_reply_wipe(&reply);
    } else {
        std::cout << xid << " getWMAllowedActions error" << std::endl;
    }

    return ret;
}

void XCBUtils::setWMAllowedActions(XWindow xid, std::vector<XCBAtom> actions)
{
    XCBAtom list[MAXALLOWEDACTIONLEN] {0};
    for (size_t i = 0; i < actions.size(); i++) {
        list[i] = actions[i];
    }

    xcb_ewmh_set_wm_allowed_actions(&m_ewmh, xid, actions.size(), list);
}

std::string XCBUtils::getWMName(XWindow xid)
{
    std::string ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_name(&m_ewmh, xid);
    xcb_ewmh_get_utf8_strings_reply_t reply;
    if (xcb_ewmh_get_wm_name_reply(&m_ewmh, cookie, &reply, nullptr)) {
        ret.assign(reply.strings, reply.strings_len);
        // 释放utf8_strings_reply分配的内存
        xcb_ewmh_get_utf8_strings_reply_wipe(&reply);
    } else {
        std::cout << xid << " getWMName error" << std::endl;
    }

    return ret;
}

uint32_t XCBUtils::getWMPid(XWindow xid)
{
    // NOTE(black_desk): code copy from https://gitlab.gnome.org/GNOME/metacity/-/merge_requests/13/diffs

    XResClientIdSpec spec = {
        .client = xid,
        .mask = XRES_CLIENT_ID_PID_MASK,
    };

    std::shared_ptr<Display> dpy = {
        XOpenDisplay(nullptr),
        [](Display *p){ XCloseDisplay(p); },
    };

    long num_ids;
    XResClientIdValue *client_ids;
    XResQueryClientIds(dpy.get(),
                       1,
                       &spec,
                       &num_ids,
                       &client_ids);

    pid_t pid = -1;
    for (long i = 0; i < num_ids; i++) {
        if (client_ids[i].spec.mask == XRES_CLIENT_ID_PID_MASK) {
            pid = XResGetClientPid(&client_ids[i]);
            break;
        }
    }

    XResClientIdsDestroy(num_ids, client_ids);
    return pid;
}

std::string XCBUtils::getWMIconName(XWindow xid)
{
    std::string ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_icon_name(&m_ewmh, xid);
    xcb_ewmh_get_utf8_strings_reply_t reply;
    if (!xcb_ewmh_get_wm_icon_name_reply(&m_ewmh, cookie, &reply, nullptr)) {
        std::cout << xid << " getWMIconName error" << std::endl;
    }

    ret.assign(reply.strings);

    return ret;
}

WMIcon XCBUtils::getWMIcon(XWindow xid)
{
    WMIcon wmIcon{};
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_icon(&m_ewmh, xid);
    xcb_ewmh_get_wm_icon_reply_t reply;
    xcb_generic_error_t* error;

    auto ret = xcb_ewmh_get_wm_icon_reply(&m_ewmh, cookie, &reply, &error);

    if (error) {
        std::cout << "failed to get wm icon" << error->error_code;
        std::free(error);
        return wmIcon;
    }

    if (ret) {
        auto fcn = [](xcb_ewmh_wm_icon_iterator_t it) {
            // https://specifications.freedesktop.org/wm-spec/wm-spec-1.3.html#idm45582154990752
            // The first two cardinals are width, height. Data is in rows, left to right and top to bottom
            // Two cardinals means width and heighr, not offset.
            const auto size = it.width * it.height;
            std::vector<uint32_t> ret(size);
            // data数据是按行从左至右，从上至下排列
            uint32_t *data = it.data;
            if (!data) {
                return ret;
            }

            std::copy_n(data, size, ret.begin());
            return ret;
        };

        // 获取icon中size最大的图标
        xcb_ewmh_wm_icon_iterator_t iter = xcb_ewmh_get_wm_icon_iterator(&reply);
        xcb_ewmh_wm_icon_iterator_t wmIconIt{0, 0, nullptr};
        for (; iter.rem; xcb_ewmh_get_wm_icon_next(&iter)) {
            const uint32_t size = iter.width * iter.height;
            if (size > 0 && size > wmIconIt.width * wmIconIt.height) {
                wmIconIt = iter;
            }
        }

        wmIcon = WMIcon{wmIconIt.width, wmIconIt.height, fcn(wmIconIt)};

        xcb_ewmh_get_wm_icon_reply_wipe(&reply); // clear
    }

    return wmIcon;
}

XWindow XCBUtils::getWMClientLeader(XWindow xid)
{
    XWindow ret = 0;
    XCBAtom atom = getAtom("WM_CLIENT_LEADER");
    void *value = getPropertyValue(xid, atom, XCB_ATOM_INTEGER);
    if (value) {
        ret = *(XWindow*)(value);
    }
    return ret;
}

void XCBUtils::requestCloseWindow(XWindow xid, uint32_t timestamp)
{
    xcb_ewmh_request_close_window(&m_ewmh, m_screenNum, xid, timestamp, XCB_EWMH_CLIENT_SOURCE_TYPE_OTHER);
}

uint32_t XCBUtils::getWMDesktop(XWindow xid)
{
    uint32_t ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_desktop(&m_ewmh, xid);
    if (!xcb_ewmh_get_wm_desktop_reply(&m_ewmh, cookie, &ret, nullptr)) {
        std::cout << xid << " getWMDesktop error" << std::endl;
    }

    return ret;
}

void XCBUtils::setWMDesktop(XWindow xid, uint32_t desktop)
{
    xcb_ewmh_set_wm_desktop(&m_ewmh, xid, desktop);
}

void XCBUtils::setCurrentWMDesktop(uint32_t desktop)
{
    xcb_ewmh_set_current_desktop(&m_ewmh, m_screenNum, desktop);
}

void XCBUtils::changeCurrentDesktop(uint32_t newDesktop, uint32_t timestamp)
{
    xcb_ewmh_request_change_current_desktop(&m_ewmh, m_screenNum, newDesktop, timestamp);
}

uint32_t XCBUtils::getCurrentWMDesktop()
{
    uint32_t ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_current_desktop(&m_ewmh, m_screenNum);
    if (!xcb_ewmh_get_current_desktop_reply(&m_ewmh, cookie, &ret, nullptr)) {
        std::cout << "getCurrentWMDesktop error" << std::endl;
    }

    return ret;
}

bool XCBUtils::isGoodWindow(XWindow xid)
{
    bool ret = false;
    xcb_get_geometry_cookie_t cookie = xcb_get_geometry(m_connect, xid);
    xcb_generic_error_t **errStore = nullptr;
    xcb_get_geometry_reply_t *reply = xcb_get_geometry_reply(m_connect, cookie, errStore);
    if (reply) {
        // 正常获取窗口geometry则判定为good
        if (!errStore) {
            ret = true;
        } else {
            free(errStore);
        }

        free(reply);
    }
    return ret;
}

// TODO XCB下无_MOTIF_WM_HINTS属性
MotifWMHints XCBUtils::getWindowMotifWMHints(XWindow xid)
{
    XCBAtom atomWmHints = getAtom("_MOTIF_WM_HINTS");
    xcb_get_property_cookie_t cookie = xcb_get_property(m_connect, false, xid, atomWmHints, atomWmHints, 0, 5);
    std::unique_ptr<xcb_get_property_reply_t> reply(xcb_get_property_reply(m_connect, cookie, nullptr));
    if (!reply || reply->format != 32 || reply->value_len != 5)
        return MotifWMHints{0, 0, 0, 0, 0};

    uint32_t *data = static_cast<uint32_t *>(xcb_get_property_value(reply.get()));
    MotifWMHints ret;
    ret.flags = data[0];
    ret.functions = data[1];
    ret.decorations = data[2];
    ret.inputMode = data[3];
    ret.status = data[4];
    return ret;
}

bool XCBUtils::hasXEmbedInfo(XWindow xid)
{
    //XCBAtom atom = getAtom("_XEMBED_INFO");

    return false;
}

XWindow XCBUtils::getWMTransientFor(XWindow xid)
{
    XWindow ret;
    xcb_get_property_cookie_t cookie = xcb_icccm_get_wm_transient_for(m_connect, xid);
    if (!xcb_icccm_get_wm_transient_for_reply(m_connect, cookie, &ret, nullptr)) {
        std::cout << xid << " getWMTransientFor error" << std::endl;
    }

    return ret;
}

uint32_t XCBUtils::getWMUserTime(XWindow xid)
{
    uint32_t ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_user_time(&m_ewmh, xid);
    if (!xcb_ewmh_get_wm_user_time_reply(&m_ewmh, cookie, &ret, nullptr)) {
        std::cout << xid << " getWMUserTime error" << std::endl;
    }

    return ret;
}

int XCBUtils::getWMUserTimeWindow(XWindow xid)
{
    XCBAtom ret;
    xcb_get_property_cookie_t cookie = xcb_ewmh_get_wm_user_time_window(&m_ewmh, xid);
    if (!xcb_ewmh_get_wm_user_time_window_reply(&m_ewmh, cookie, &ret, nullptr)) {
        std::cout << xid << " getWMUserTimeWindow error" << std::endl;
    }

    return ret;
}

WMClass XCBUtils::getWMClass(XWindow xid)
{
    WMClass ret;
    xcb_get_property_cookie_t cookie = xcb_icccm_get_wm_class(m_connect, xid);
    xcb_icccm_get_wm_class_reply_t reply;
    reply.instance_name = nullptr;
    reply.class_name = nullptr;
    xcb_icccm_get_wm_class_reply(m_connect, cookie, &reply, nullptr);   // 返回值为0不一定表示失败， 故不做返回值判断
    if (reply.class_name)
        ret.className.assign(reply.class_name);

    if (reply.instance_name)
        ret.instanceName.assign(reply.instance_name);

    if (reply.class_name || reply.instance_name) {
        xcb_icccm_get_wm_class_reply_wipe(&reply);
    }

    return ret;
}

void XCBUtils::minimizeWindow(XWindow xid)
{
    uint32_t data[2];
    data[0] = XCB_ICCCM_WM_STATE_ICONIC;
    data[1] = XCB_NONE;
    xcb_ewmh_send_client_message(m_connect, xid, getRootWindow(),getAtom("WM_CHANGE_STATE"), 2, data);
    flush();
}

void XCBUtils::maxmizeWindow(XWindow xid)
{
    xcb_ewmh_request_change_wm_state(&m_ewmh
                                     , m_screenNum
                                     , xid
                                     , XCB_EWMH_WM_STATE_ADD
                                     , getAtom("_NET_WM_STATE_MAXIMIZED_VERT")
                                     , getAtom("_NET_WM_STATE_MAXIMIZED_HORZ")
                                     , XCB_EWMH_CLIENT_SOURCE_TYPE_OTHER);
}

// TODO
std::vector<std::string> XCBUtils::getWMCommand(XWindow xid)
{
    std::vector<std::string> ret;
    xcb_get_property_reply_t *reply = getPropertyValueReply(xid, XCB_ATOM_WM_COMMAND, m_ewmh.UTF8_STRING);
    if (reply) {
        ret = getUTF8StrsFromReply(reply);
        free(reply);
    }

    return ret;
}

std::string XCBUtils::getUTF8StrFromReply(xcb_get_property_reply_t *reply)
{
    std::string ret;
    if (!reply || reply->format != 8) {
        return ret;
    }

    char data[12] = {0};
    for (uint32_t i=0; i < reply->value_len; i++) {
        data[i] = char(reply->pad0[i]);
    }
    ret.assign(data);
    return ret;
}

std::vector<std::string> XCBUtils::getUTF8StrsFromReply(xcb_get_property_reply_t *reply)
{
    std::vector<std::string> ret;
    if (!reply) {
        return ret;
    }

    if (reply->format != 8) {
        return ret;
    }


    // 字符串拆分
    uint32_t start = 0;
    for (uint32_t i=0; i < reply->value_len; i++) {
        if (reply->pad0[i] == 0) {
            char data[12] = {0};
            int count = 0;
            for (uint32_t j=start; j < i; j++)
                data[count++] = char(reply->pad0[j]);

            data[count] = 0;
            ret.push_back(data);
        }
    }

    return ret;
}

XWindow XCBUtils::getRootWindow()
{
    XWindow rootWindow = 0;
    /* Get the first screen */
    xcb_screen_t *screen = xcb_setup_roots_iterator(xcb_get_setup(m_connect)).data;
    if (screen) {
        rootWindow = screen->root;
    }

    std::cout << "getRootWinodw: " << rootWindow << std::endl;
    return rootWindow;
}

void XCBUtils::registerEvents(XWindow xid, uint32_t eventMask)
{
    uint32_t value[1] = {eventMask};
    xcb_void_cookie_t cookie = xcb_change_window_attributes_checked(m_connect,
                                                                    xid,
                                                                    XCB_CW_EVENT_MASK,
                                                                    &value);
    flush();

    xcb_generic_error_t *error = xcb_request_check(m_connect, cookie);
    if (error != nullptr) {
        std::cout << "window " << xid << "registerEvents error" << std::endl;
    }
}


AtomCache::AtomCache()
{
}

XCBAtom AtomCache::getVal(std::string name)
{
    XCBAtom atom = ATOMNONE;
    auto search = m_atoms.find(name);
    if (search != m_atoms.end()) {
        atom = search->second;
    }

    return atom;
}

std::string AtomCache::getName(XCBAtom atom)
{
    std::string ret;
    auto search = m_atomNames.find(atom);
    if (search != m_atomNames.end()) {
        ret = search->second;
    }

    return ret;
}

void AtomCache::store(std::string name, XCBAtom value)
{
    m_atoms[name] = value;
    m_atomNames[value] = name;
}
