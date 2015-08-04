#-------------------------------------------------
#
# Project created by QtCreator 2015-06-29T20:08:12
#
#-------------------------------------------------

QT       += core gui dbus

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += plugin c++11

TARGET = $$qtLibraryTarget(dock-systray-plugin)
TEMPLATE = lib

SOURCES += systrayplugin.cpp \
    dbustraymanager.cpp \
    compositetrayitem.cpp \
    trayicon.cpp

HEADERS  += systrayplugin.h \
    dbustraymanager.h \
    compositetrayitem.h \
    trayicon.h

RESOURCES += images.qrc
DISTFILES += dde-dock-systray-plugin.json

target.path = /usr/share/dde-dock/plugins/
INSTALLS += target
