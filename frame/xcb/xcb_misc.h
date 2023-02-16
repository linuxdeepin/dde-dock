// Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef XCB_MISC_H
#define XCB_MISC_H

#include <QtCore>

#include <xcb/xcb_ewmh.h>

class XcbMisc
{

public:
    enum Orientation {
        OrientationLeft,
        OrientationRight,
        OrientationTop,
        OrientationBottom
    };

    enum WindowType {
        Dock,
        Desktop
    };

    virtual ~XcbMisc();

    static XcbMisc * instance();

    void set_window_type(xcb_window_t winId, WindowType winType);
    void clear_strut_partial(xcb_window_t winId);
    void set_strut_partial(xcb_window_t winId, Orientation orientation, uint strut, uint start, uint end);
    void set_window_icon_geometry(xcb_window_t winId, QRect geo);

private:
    XcbMisc();

    xcb_ewmh_connection_t m_ewmh_connection;
};

#endif // XCB_MISC_H
