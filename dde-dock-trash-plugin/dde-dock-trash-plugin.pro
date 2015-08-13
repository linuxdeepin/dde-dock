#-------------------------------------------------
#
# Project created by QtCreator 2015-07-20T14:06:01
#
#-------------------------------------------------

QT       += core gui widgets dbus

TARGET = dde-dock-trash-plugin
TEMPLATE = lib
CONFIG += plugin c++11

INCLUDEPATH += ../dde-dock/src

SOURCES += \
    mainitem.cpp \
    trashplugin.cpp \
    dbus/dbusfileoperations.cpp \
    dbus/dbusfiletrashmonitor.cpp \
    dbus/dbustrashjob.cpp \
    dbus/dbusemptytrashjob.cpp \
    dbus/dbuslauncher.cpp

HEADERS += \
    mainitem.h \
    trashplugin.h \
    dbus/dbusfileoperations.h \
    dbus/dbusfiletrashmonitor.h \
    dbus/dbustrashjob.h \
    dbus/dbusemptytrashjob.h \
    dbus/dbuslauncher.h

DISTFILES += dde-dock-trash-plugin.json

unix {
    target.path = /usr/share/dde-dock/plugins/
    INSTALLS += target
}

RESOURCES +=
