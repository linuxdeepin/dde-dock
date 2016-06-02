QT       += core gui widgets dbus x11extras

TARGET = dde-dock
TEMPLATE = app
CONFIG += c++11 link_pkgconfig

PKGCONFIG += xcb-ewmh

SOURCES += main.cpp \
    window/mainwindow.cpp \
    dbus/dbusdockentrymanager.cpp \
    xcb/xcb_misc.cpp \
    item/dockitem.cpp \
    panel/mainpanel.cpp \
    controller/dockitemcontroller.cpp

HEADERS  += \
    window/mainwindow.h \
    dbus/dbusdockentrymanager.h \
    xcb/xcb_misc.h \
    item/dockitem.h \
    panel/mainpanel.h \
    controller/dockitemcontroller.h
