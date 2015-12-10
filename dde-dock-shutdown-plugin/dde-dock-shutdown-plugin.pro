include(../../common.pri)

QT       += core gui widgets

TARGET = dde-dock-shutdown-plugin
TEMPLATE = lib
CONFIG += plugin c++11
DESTDIR = $$_PRO_FILE_PWD_/../

SOURCES += \ 
    shutdownplugin.cpp

HEADERS += \ 
    shutdownplugin.h

DISTFILES += dde-dock-shutdown-plugin.json

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

QMAKE_MOC_OPTIONS += -I/usr/include/

RESOURCES += \
    resources.qrc
