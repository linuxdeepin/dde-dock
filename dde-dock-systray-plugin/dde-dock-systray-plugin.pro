#-------------------------------------------------
#
# Project created by QtCreator 2015-06-29T20:08:12
#
#-------------------------------------------------

QT       += core gui dbus x11extras

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

PKGCONFIG += xcb xcb-image xcb-composite dtkbase dtkwidget
CONFIG += plugin c++11 link_pkgconfig

TARGET = $$qtLibraryTarget(dde-dock-systray-plugin)
TEMPLATE = lib
INCLUDEPATH += ../dde-dock/src/

SOURCES += systrayplugin.cpp \
    dbustraymanager.cpp \
    compositetrayitem.cpp \
    trayicon.cpp \
    ../dde-dock/src/dbus/dbusentrymanager.cpp

HEADERS  += systrayplugin.h \
    dbustraymanager.h \
    compositetrayitem.h \
    trayicon.h \
    ../dde-dock/src/dbus/dbusentrymanager.h

RESOURCES += images.qrc
DISTFILES += dde-dock-systray-plugin.json

target.path = /usr/lib/dde-dock/plugins/
INSTALLS += target
