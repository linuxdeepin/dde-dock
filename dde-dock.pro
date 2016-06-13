QT       += core gui widgets dbus x11extras

TARGET = dde-dock
TEMPLATE = app
CONFIG += c++11 link_pkgconfig

PKGCONFIG += xcb-ewmh

SOURCES += main.cpp \
    window/mainwindow.cpp \
    xcb/xcb_misc.cpp \
    item/dockitem.cpp \
    panel/mainpanel.cpp \
    controller/dockitemcontroller.cpp \
    dbus/dbusdockentry.cpp \
    dbus/dbusdisplay.cpp \
    item/appitem.cpp \
    util/docksettings.cpp \
    item/placeholderitem.cpp \
    dbus/dbusclientmanager.cpp \
    dbus/dbusdock.cpp

HEADERS  += \
    window/mainwindow.h \
    xcb/xcb_misc.h \
    item/dockitem.h \
    panel/mainpanel.h \
    controller/dockitemcontroller.h \
    dbus/dbusdockentry.h \
    dbus/dbusdisplay.h \
    item/appitem.h \
    util/docksettings.h \
    item/placeholderitem.h \
    dbus/dbusclientmanager.h \
    dbus/dbusdock.h
