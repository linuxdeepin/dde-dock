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

DEFINES += PLUGIN_API_VERSION=1.0 NEW_DOCK_LAYOUT


RESOURCES += \
    dark.qrc \
    light.qrc

CONFIG += c++11 link_pkgconfig
PKGCONFIG += gtk+-2.0 x11 cairo xcb xcb-ewmh xcb-damage dtkbase dtkutil dtkwidget

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
    src/controller/dockmodedata.h \
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
    src/widgets/highlighteffect.h \
    src/widgets/reflectioneffect.h \
    src/mainwidget.h \
    src/controller/signalmanager.h \
    src/widgets/previewwindow.h \
    src/dbus/dbusdisplay.h \
    src/dbus/dbuslauncher.h \
    src/dbus/dbusdockentry.h \
    src/panel/dockpanel.h \
    src/widgets/app/dockappbg.h \
    src/widgets/app/dockappicon.h \
    src/widgets/app/dockappitem.h \
    src/widgets/app/dockapplayout.h \
    src/controller/apps/dockappmanager.h \
    src/widgets/dockitem.h \
    src/widgets/plugin/dockpluginitem.h \
    src/widgets/plugin/dockpluginssettingwindow.h \
    src/controller/plugins/dockpluginsmanager.h \
    src/widgets/movablelayout.h \
    src/widgets/plugin/dockpluginlayout.h \
    src/controller/old/appmanager.h \
    src/controller/old/pluginitemwrapper.h \
    src/controller/old/pluginproxy.h \
    src/controller/old/pluginssettingframe.h \
    src/controller/old/pluginmanager.h \
    src/controller/plugins/dockpluginitemwrapper.h \
    src/controller/plugins/dockpluginproxy.h \
    src/widgets/app/apppreview/apppreviewloader.h \
    src/widgets/app/apppreview/apppreviewloaderframe.h \
    src/widgets/app/apppreview/apppreviewscontainer.h \
    src/widgets/old/abstractdockitem.h \
    src/widgets/old/appbackground.h \
    src/widgets/old/appicon.h \
    src/widgets/old/appitem.h \
    src/widgets/old/docklayout.h \
    src/widgets/old/screenmask.h \
    src/widgets/launcher/docklauncheritem.h \
    src/widgets/old/launcheritem.h \
    src/controller/stylemanager.h

SOURCES += \
    libs/xcb_misc.cpp \
    src/controller/dockmodedata.cpp \
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
    src/widgets/highlighteffect.cpp \
    src/widgets/reflectioneffect.cpp \
    src/main.cpp \
    src/mainwidget.cpp \
    src/controller/signalmanager.cpp \
    src/widgets/previewwindow.cpp \
    src/dbus/dbusdisplay.cpp \
    src/dbus/dbuslauncher.cpp \
    src/dbus/dbusdockentry.cpp \
    src/panel/dockpanel.cpp \
    src/widgets/app/dockappbg.cpp \
    src/widgets/app/dockappicon.cpp \
    src/widgets/app/dockappitem.cpp \
    src/widgets/app/dockapplayout.cpp \
    src/controller/apps/dockappmanager.cpp \
    src/widgets/dockitem.cpp \
    src/widgets/plugin/dockpluginitem.cpp \
    src/widgets/plugin/dockpluginssettingwindow.cpp \
    src/controller/plugins/dockpluginsmanager.cpp \
    src/widgets/movablelayout.cpp \
    src/widgets/plugin/dockpluginlayout.cpp \
    src/controller/old/appmanager.cpp \
    src/controller/old/pluginitemwrapper.cpp \
    src/controller/old/pluginproxy.cpp \
    src/controller/old/pluginssettingframe.cpp \
    src/controller/old/pluginmanager.cpp \
    src/controller/plugins/dockpluginitemwrapper.cpp \
    src/controller/plugins/dockpluginproxy.cpp \
    src/widgets/app/apppreview/apppreviewloader.cpp \
    src/widgets/app/apppreview/apppreviewloaderframe.cpp \
    src/widgets/app/apppreview/apppreviewscontainer.cpp \
    src/widgets/old/abstractdockitem.cpp \
    src/widgets/old/appbackground.cpp \
    src/widgets/old/appicon.cpp \
    src/widgets/old/appitem.cpp \
    src/widgets/old/docklayout.cpp \
    src/widgets/old/screenmask.cpp \
    src/widgets/launcher/docklauncheritem.cpp \
    src/widgets/old/launcheritem.cpp \
    src/controller/stylemanager.cpp
