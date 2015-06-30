#-------------------------------------------------
#
# Project created by QtCreator 2015-06-29T20:08:12
#
#-------------------------------------------------

QT       += core gui

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

INCLUDEPATH += ../dde-dock/src/
CONFIG += plugin c++11

TARGET = $$qtLibraryTarget(dock-systray-plugin)
TEMPLATE = lib

SOURCES += systrayplugin.cpp \
    docktrayitem.cpp

HEADERS  += systrayplugin.h \
    docktrayitem.h
