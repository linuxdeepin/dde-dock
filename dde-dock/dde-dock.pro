#-------------------------------------------------
#
# Project created by QtCreator 2015-06-20T10:09:57
#
#-------------------------------------------------

QT       += core gui x11extras

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = dde-dock
TEMPLATE = app


SOURCES += main.cpp\
        mainwidget.cpp \
    Panel/panel.cpp \
    Widgets/appicon.cpp \
    Widgets/appbackground.cpp \
    Widgets/dockconstants.cpp \
    Widgets/dockmodel.cpp \
    Widgets/dockview.cpp \
    Widgets/dockitemdelegate.cpp \
    Widgets/appitem.cpp \
    Widgets/docklayout.cpp \
    Widgets/windowpreview.cpp \
    Widgets/dockitem.cpp \
    Widgets/screenmask.cpp

HEADERS  += mainwidget.h \
    Panel/panel.h \
    Widgets/appicon.h \
    Widgets/appbackground.h \
    Widgets/dockconstants.h \
    Widgets/dockmodel.h \
    Widgets/dockview.h \
    Widgets/dockitemdelegate.h \
    Widgets/appitem.h \
    Widgets/docklayout.h \
    Widgets/windowpreview.h \
    Widgets/dockitem.h \
    Widgets/screenmask.h \
    dockplugininterface.h \
    abstractdockitem.h

RESOURCES += \
    images.qrc

PKGCONFIG += gtk+-2.0 x11
CONFIG += c++11 link_pkgconfig
