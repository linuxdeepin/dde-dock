#-------------------------------------------------
#
# Project created by QtCreator 2015-06-20T10:09:57
#
#-------------------------------------------------

QT       += core gui x11extras

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = dde-dock
TEMPLATE = app
INCLUDEPATH += src/

SOURCES += \
    src/main.cpp \
    src/mainwidget.cpp \
    src/Widgets/appbackground.cpp \
    src/Widgets/appicon.cpp \
    src/Widgets/docklayout.cpp \
    src/Widgets/screenmask.cpp \
    src/Widgets/windowpreview.cpp \
    src/Panel/panel.cpp \
    src/Widgets/appitem.cpp \
    src/systraymanager.cpp \
    src/Panel/panelmenu.cpp \
    src/Controller/dockmodedata.cpp \
    src/Controller/dockconstants.cpp

HEADERS  += \
    src/abstractdockitem.h \
    src/dockplugininterface.h \
    src/mainwidget.h \
    src/Widgets/appbackground.h \
    src/Widgets/appicon.h \
    src/Widgets/docklayout.h \
    src/Widgets/screenmask.h \
    src/Widgets/windowpreview.h \
    src/Panel/panel.h \
    src/Widgets/appitem.h \
    src/systraymanager.h \
    src/Panel/panelmenu.h \
    src/Controller/dockmodedata.h \
    src/Controller/dockconstants.h

RESOURCES += \
    images.qrc \
    qss.qrc

PKGCONFIG += gtk+-2.0 x11
CONFIG += c++11 link_pkgconfig
