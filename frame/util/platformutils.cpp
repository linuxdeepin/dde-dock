// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "platformutils.h"
#include "utils.h"

#include <QX11Info>

#include <X11/Xlib.h>

#define NORMAL_WINDOW_PROP_NAME "WM_CLASS"
#define WINE_WINDOW_PROP_NAME "__wine_prefix"
#define IS_WINE_WINDOW_BY_WM_CLASS "explorer.exe"

QString PlatformUtils::getAppNameForWindow(quint32 winId)
{
    // is normal application
    QString appName = getWindowProperty(winId, NORMAL_WINDOW_PROP_NAME);
    if (!appName.isEmpty() && appName != IS_WINE_WINDOW_BY_WM_CLASS)
        return appName;

    // is wine application
    appName = getWindowProperty(winId, WINE_WINDOW_PROP_NAME).split("/").last();
    if (!appName.isEmpty())
        return appName;

    // fallback to window id
    return QString::number(winId);
}

QString PlatformUtils::getWindowProperty(quint32 winId, QString propName)
{
    const auto display = Utils::IS_WAYLAND_DISPLAY ? XOpenDisplay(nullptr) : QX11Info::display();
    if (!display) {
        qWarning() << "QX11Info::display() is " << display;
        return QString();
    }

    Atom atom_prop = XInternAtom(display, propName.toLocal8Bit(), true);
    if (!atom_prop) {
        qDebug() << "Error: get window property failed, invalid property atom";
        return QString();
    }

    Atom actual_type_return;
    int actual_format_return;
    unsigned long nitems_return;
    unsigned long bytes_after_return;
    unsigned char *prop_return;

    int r = XGetWindowProperty(display, winId, atom_prop, 0, 100, false, AnyPropertyType,
            &actual_type_return, &actual_format_return, &nitems_return,
            &bytes_after_return, &prop_return);

    Q_UNUSED(r);
    if (Utils::IS_WAYLAND_DISPLAY)
        XCloseDisplay(display);

    return QString::fromLocal8Bit((char*)prop_return);
}
