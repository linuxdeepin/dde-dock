#-------------------------------------------------
#
# Project created by QtCreator 2015-06-20T10:09:57
#
#-------------------------------------------------

QT       += core gui x11extras dbus svg

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = dde-dock
TEMPLATE = app
INCLUDEPATH += ./src ./libs

DEFINES += PLUGIN_API_VERSION=1.0 #NEW_DOCK_LAYOUT


RESOURCES += \
    dark.qrc \
    light.qrc

PKGCONFIG += gtk+-2.0 x11 cairo xcb xcb-ewmh xcb-damage dui
CONFIG += c++11 link_pkgconfig

include (../cutelogger/cutelogger.pri)

target.path = /usr/bin/

headers.files += src/interfaces/dockconstants.h \
    src/interfaces/dockplugininterface.h \
    src/interfaces/dockpluginproxyinterface.h
headers.path = /usr/include/dde-dock

dbus_service.files += com.deepin.dde.dock.service
dbus_service.path = /usr/share/dbus-1/services

INSTALLS += dbus_service headers target

HEADERS += \
    src/interfaces/dockconstants.h \
    src/interfaces/dockplugininterface.h \
    src/interfaces/dockpluginproxyinterface.h \
    libs/xcb_misc.h \
    src/controller/appmanager.h \
    src/controller/dockmodedata.h \
    src/controller/plugins/dockpluginmanager.h \
    src/controller/plugins/dockpluginproxy.h \
    src/controller/plugins/pluginitemwrapper.h \
    src/controller/plugins/pluginssettingframe.h \
    src/dbus/dbusclientmanager.h \
    src/dbus/dbusdockedappmanager.h \
    src/dbus/dbusdocksetting.h \
    src/dbus/dbusentrymanager.h \
    src/dbus/dbushidestatemanager.h \
    src/dbus/dbusmenu.h \
    src/dbus/dbusmenumanager.h \
    src/dbus/dbuspanelmanager.h \
    src/panel/panel.h \
    src/panel/panelmenu.h \
    src/widgets/abstractdockitem.h \
    src/widgets/appbackground.h \
    src/widgets/appicon.h \
    src/widgets/appitem.h \
    src/widgets/docklayout.h \
    src/widgets/highlighteffect.h \
    src/widgets/launcheritem.h \
    src/widgets/reflectioneffect.h \
    src/widgets/screenmask.h \
    src/mainwidget.h \
    src/controller/signalmanager.h \
    src/widgets/apppreview/apppreviewloader.h \
    src/widgets/apppreview/apppreviewscontainer.h \
    src/widgets/apppreview/apppreviewloaderframe.h \
    src/widgets/previewwindow.h \
    src/dbus/dbusdisplay.h \
    src/dbus/dbuslauncher.h \
    src/dbus/dbusdockentry.h \
    src/panel/dockpanel.h \
    src/widgets/app/dockappbg.h \
    src/widgets/app/dockappicon.h \
    src/widgets/app/dockappitem.h \
    src/widgets/app/dockapplayout.h \
    src/widgets/app/dockitem.h \
    src/widgets/app/movablelayout.h \
    src/controller/apps/dockappmanager.h

SOURCES += \
    libs/xcb_misc.cpp \
    src/controller/appmanager.cpp \
    src/controller/dockmodedata.cpp \
    src/controller/plugins/dockpluginmanager.cpp \
    src/controller/plugins/dockpluginproxy.cpp \
    src/controller/plugins/pluginitemwrapper.cpp \
    src/controller/plugins/pluginssettingframe.cpp \
    src/dbus/dbusclientmanager.cpp \
    src/dbus/dbusdockedappmanager.cpp \
    src/dbus/dbusdocksetting.cpp \
    src/dbus/dbusentrymanager.cpp \
    src/dbus/dbushidestatemanager.cpp \
    src/dbus/dbusmenu.cpp \
    src/dbus/dbusmenumanager.cpp \
    src/dbus/dbuspanelmanager.cpp \
    src/panel/panel.cpp \
    src/panel/panelmenu.cpp \
    src/widgets/abstractdockitem.cpp \
    src/widgets/appbackground.cpp \
    src/widgets/appicon.cpp \
    src/widgets/appitem.cpp \
    src/widgets/docklayout.cpp \
    src/widgets/highlighteffect.cpp \
    src/widgets/launcheritem.cpp \
    src/widgets/reflectioneffect.cpp \
    src/widgets/screenmask.cpp \
    src/main.cpp \
    src/mainwidget.cpp \
    src/controller/signalmanager.cpp \
    src/widgets/apppreview/apppreviewloader.cpp \
    src/widgets/apppreview/apppreviewscontainer.cpp \
    src/widgets/apppreview/apppreviewloaderframe.cpp \
    src/widgets/previewwindow.cpp \
    src/dbus/dbusdisplay.cpp \
    src/dbus/dbuslauncher.cpp \
    src/dbus/dbusdockentry.cpp \
    src/panel/dockpanel.cpp \
    src/widgets/app/dockappbg.cpp \
    src/widgets/app/dockappicon.cpp \
    src/widgets/app/dockappitem.cpp \
    src/widgets/app/dockapplayout.cpp \
    src/widgets/app/dockitem.cpp \
    src/widgets/app/movablelayout.cpp \
    src/controller/apps/dockappmanager.cpp
