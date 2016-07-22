
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
    item/wireditem.h \
    dbus/dbusnetwork.h \
    networkmanager.h \
    networkdevice.h \
    util/imageutil.h \
    item/deviceitem.h

SOURCES += \
    networkplugin.cpp \
    item/wireditem.cpp \
    dbus/dbusnetwork.cpp \
    networkmanager.cpp \
    networkdevice.cpp \
    util/imageutil.cpp \
    item/deviceitem.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
    resources.qrc

