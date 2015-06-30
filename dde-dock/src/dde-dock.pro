#-------------------------------------------------
#
# Project created by QtCreator 2015-06-20T10:09:57
#
#-------------------------------------------------

QT       += core gui x11extras

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = dde-dock
TEMPLATE = app


SOURCES += \
    main.cpp \
    mainwidget.cpp \
    Widgets/appbackground.cpp \
    Widgets/appicon.cpp \
    Widgets/appitem.cpp \
    Widgets/dockconstants.cpp \
    Widgets/dockitem.cpp \
    Widgets/dockitemdelegate.cpp \
    Widgets/docklayout.cpp \
    Widgets/dockmodel.cpp \
    Widgets/dockview.cpp \
    Widgets/screenmask.cpp \
    Widgets/windowpreview.cpp \
    Panel/panel.cpp

HEADERS  += \
    abstractdockitem.h \
    dockplugininterface.h \
    mainwidget.h \
    Widgets/appbackground.h \
    Widgets/appicon.h \
    Widgets/appitem.h \
    Widgets/dockconstants.h \
    Widgets/dockitem.h \
    Widgets/dockitemdelegate.h \
    Widgets/docklayout.h \
    Widgets/dockmodel.h \
    Widgets/dockview.h \
    Widgets/screenmask.h \
    Widgets/windowpreview.h \
    Panel/panel.h

RESOURCES += \
    images.qrc

PKGCONFIG += gtk+-2.0 x11
CONFIG += c++11 link_pkgconfig
