// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef XCBUTILS_H
#define XCBUTILS_H

#include <xcb/xproto.h>
#include <xcb/xcb_ewmh.h>
#include <xcb/xcb_icccm.h>

#include <list>
#include <string>
#include <vector>
#include <map>

#define MAXLEN 0xffff
#define MAXALLOWEDACTIONLEN 256
#define ATOMNONE 0

typedef xcb_window_t XWindow ;
typedef xcb_atom_t XCBAtom;
typedef xcb_destroy_notify_event_t DestroyEvent;
typedef xcb_map_notify_event_t MapEvent;
typedef xcb_configure_notify_event_t ConfigureEvent;
typedef xcb_property_notify_event_t PropertyEvent;
typedef xcb_event_mask_t EventMask;

typedef struct {
    std::string instanceName;
    std::string className;
} WMClass;

typedef struct {
    int16_t x, y;
    uint16_t width, height;
} Geometry;

typedef struct {
    uint32_t flags;
    uint32_t functions;
    uint32_t decorations;
    int32_t inputMode;
    uint32_t status;
} MotifWMHints;

typedef struct {
  uint32_t width;   /** Icon width */
  uint32_t height;  /** Icon height */
  std::vector<uint32_t> data;   /** Rows, left to right and top to bottom of the CARDINAL ARGB */
} WMIcon;

typedef struct WindowFrameExtents {
    uint32_t Left;
    uint32_t Right;
    uint32_t Top;
    uint32_t Bottom;
    WindowFrameExtents(int left = 0, int right = 0, int top = 0, int bottom = 0): Left(left), Right(right), Top(top), Bottom(bottom) {}
    bool isNull() { return Left == 0 && Right == 0 && Top == 0 && Bottom == 0;}
} WindowFrameExtents;

// 缓存atom，减少X访问  TODO 加读写锁
class AtomCache {
public:
    AtomCache();

    XCBAtom getVal(std::string name);
    std::string getName(XCBAtom atom);
    void store(std::string name, XCBAtom value);

public:
    std::map<std::string, XCBAtom> m_atoms;
    std::map<XCBAtom, std::string> m_atomNames;
};

// XCB接口封装， 参考getCurrentWMDesktop
class XCBUtils
{
    XCBUtils();
    XCBUtils(const XCBUtils &other);
    XCBUtils & operator= (const XCBUtils &other);
    ~XCBUtils();

public:
    static XCBUtils *instance() {
        static XCBUtils instance;
        return &instance;
    }

    /************************* xcb method ***************************/
    // 分配XID
    XWindow allocId();

    // 刷新
    void flush();

    /************************* xpropto method ***************************/
    // 杀掉进程
    void killClientChecked(XWindow xid);

    // 获取属性reply, 返回值必须free
    xcb_get_property_reply_t *getPropertyValueReply(XWindow xid, XCBAtom property, XCBAtom type = XCB_ATOM_ATOM);

    // 获取属性
    void *getPropertyValue(XWindow xid, XCBAtom property, XCBAtom type = XCB_ATOM_ATOM);

    // 获取字符串属性
    std::string getUTF8PropertyStr(XWindow xid, XCBAtom property);

    // 获取名称对应的Atom
    XCBAtom getAtom(const char *name);

    // 获取Atom对应的名称
    std::string getAtomName(XCBAtom atom);

    // 获取窗口矩形
    Geometry getWindowGeometry(XWindow xid);

    // 判断当前窗口是否正常
    bool isGoodWindow(XWindow xid);

    // 获取窗口
    MotifWMHints getWindowMotifWMHints(XWindow xid);

    bool hasXEmbedInfo(XWindow xid);

    /************************* ewmh method ***************************/

    // 获取活动窗口 _NET_ACTIVE_WINDOW
    XWindow getActiveWindow();

    // 设置活动窗口 _NET_ACTIVE_WINDOW 属性
    void setActiveWindow(XWindow xid);

    // 改变活动窗口
    void changeActiveWindow(XWindow newActiveXid);

    // 重新排列窗口
    void restackWindow(XWindow xid);

    // 获取窗口列表 _NET_CLIENT_LIST
    std::list<XWindow> getClientList();

    // 获取窗口列表 _NET_CLIENT_LIST_STACKING
    std::list<XWindow> getClientListStacking();

    // 获取窗口状态 _NET_WM_STATE
    /*
        _NET_WM_STATE_MODAL, ATOM
        _NET_WM_STATE_STICKY, ATOM
        _NET_WM_STATE_MAXIMIZED_VERT, ATOM
        _NET_WM_STATE_MAXIMIZED_HORZ, ATOM
        _NET_WM_STATE_SHADED, ATOM
        _NET_WM_STATE_SKIP_TASKBAR, ATOM
        _NET_WM_STATE_SKIP_PAGER, ATOM
        _NET_WM_STATE_HIDDEN, ATOM
        _NET_WM_STATE_FULLSCREEN, ATOM
        _NET_WM_STATE_ABOVE, ATOM
        _NET_WM_STATE_BELOW, ATOM
        _NET_WM_STATE_DEMANDS_ATTENTION, ATOM
    */
    std::vector<XCBAtom> getWMState(XWindow xid);

    // 获取窗口类型 _NET_WM_WINDOW_TYPE
    // Rationale: This hint is intended to replace the MOTIF hints.
    // One of the objections to the MOTIF hints is that they are a purely visual description of the window decoration.
    // By describing the function of the window, the Window Manager can apply consistent decoration and behavior to windows of the same type.
    // Possible examples of behavior include keeping dock/panels on top or allowing pinnable menus / toolbars to only be hidden
    // when another window has focus
    /*
        _NET_WM_WINDOW_TYPE_DESKTOP, ATOM
        _NET_WM_WINDOW_TYPE_DOCK, ATOM
        _NET_WM_WINDOW_TYPE_TOOLBAR, ATOM
        _NET_WM_WINDOW_TYPE_MENU, ATOM
        _NET_WM_WINDOW_TYPE_UTILITY, ATOM
        _NET_WM_WINDOW_TYPE_SPLASH, ATOM
        _NET_WM_WINDOW_TYPE_DIALOG, ATOM
        _NET_WM_WINDOW_TYPE_DROPDOWN_MENU, ATOM
        _NET_WM_WINDOW_TYPE_POPUP_MENU, ATOM
        _NET_WM_WINDOW_TYPE_TOOLTIP, ATOM
        _NET_WM_WINDOW_TYPE_NOTIFICATION, ATOM
        _NET_WM_WINDOW_TYPE_COMBO, ATOM
        _NET_WM_WINDOW_TYPE_DND, ATOM
        _NET_WM_WINDOW_TYPE_NORMAL, ATOM
     * */
    std::vector<XCBAtom> getWMWindoType(XWindow xid);

    // 获取窗口许可动作 _NET_WM_ALLOWED_ACTIONS
    std::vector<XCBAtom> getWMAllowedActions(XWindow xid);

    // 设置窗口许可动作
    void setWMAllowedActions(XWindow xid, std::vector<XCBAtom> actions);

    // 获取窗口名称 _NET_WM_NAME
    std::string getWMName(XWindow xid);

    // 获取窗口所属进程 _NET_WM_PID
    uint32_t getWMPid(XWindow xid);

    // 获取窗口图标 _NET_WM_ICON_NAME
    std::string getWMIconName(XWindow xid);

    // 获取窗口图标信息 _NET_WM_ICON
    WMIcon getWMIcon(XWindow xid);

    // WM_CLIENT_LEADER
    XWindow getWMClientLeader(XWindow xid);

    // 关闭窗口 _NET_CLOSE_WINDOW
    void requestCloseWindow(XWindow xid, uint32_t timestamp);

    // 获取窗口对应桌面 _NET_WM_DESKTOP
    uint32_t getWMDesktop(XWindow xid);

    // 设置窗口当前桌面
    void setWMDesktop(XWindow xid, uint32_t desktop);

    // 设置当前桌面属性
    void setCurrentWMDesktop(uint32_t desktop);

    // 请求改变当前桌面
    void changeCurrentDesktop(uint32_t newDesktop, uint32_t timestamp);

    // 获取当前桌面 _NET_CURRENT_DESKTOP
    uint32_t getCurrentWMDesktop();


    /************************* icccm method ***************************/
    // The WM_TRANSIENT_FOR hint of the ICCCM allows clients to specify that a toplevel window may be closed before the client finishes.
    // A typical example of a transient window is a dialog.
    // Some dialogs can be open for a long time, while the user continues to work in the main window.
    // Other dialogs have to be closed before the user can continue to work in the main window
    XWindow getWMTransientFor(XWindow xid);

    uint32_t getWMUserTime(XWindow xid);

    int getWMUserTimeWindow(XWindow xid);

    // 获取窗口类型
    WMClass getWMClass(XWindow xid);

    // 最小化窗口
    void minimizeWindow(XWindow xid);

    // 最大化窗口
    void maxmizeWindow(XWindow xid);

    /************************* other method ***************************/
    // 获取窗口command
    std::vector<std::string> getWMCommand(XWindow xid);

    // 解析属性为UTF8格式字符串
    std::string getUTF8StrFromReply(xcb_get_property_reply_t *reply);

    // 解析属性为UTF8格式字符串字符数组
    std::vector<std::string> getUTF8StrsFromReply(xcb_get_property_reply_t *reply);

    // 获取根窗口
    XWindow getRootWindow();

    // 注册事件
    void registerEvents(XWindow xid, uint32_t eventMask);

private:
    XWindow getDecorativeWindow(XWindow xid);
    WindowFrameExtents getWindowFrameExtents(XWindow xid);

private:
    xcb_connection_t *m_connect;
    int m_screenNum;

    xcb_ewmh_connection_t m_ewmh;
    AtomCache m_atomCache;  // 和ewmh中Atom类型存在重复部分，扩张了自定义类型
};

#endif // XCBUTILS_H
