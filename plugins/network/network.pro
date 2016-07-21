
include(../../interfaces/interfaces.pri)

QT              += widgets svg dbus
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(network)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += network.json

HEADERS += \
    networkplugin.h \
    wireditem.h \
    dbus/dbusnetwork.h \
    networkmanager.h \
    networkdevice.h

SOURCES += \
    networkplugin.cpp \
    wireditem.cpp \
    dbus/dbusnetwork.cpp \
    networkmanager.cpp \
    networkdevice.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

