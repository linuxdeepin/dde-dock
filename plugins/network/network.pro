
include(../../interfaces/interfaces.pri)

QT              += widgets svg dbus
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       += dtkbase dtkwidget

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
    item/deviceitem.h \
    item/wirelessitem.h \
    item/applet/wirelessapplet.h \
    item/applet/devicecontrolwidget.h \
    item/applet/accesspoint.h

SOURCES += \
    networkplugin.cpp \
    item/wireditem.cpp \
    dbus/dbusnetwork.cpp \
    networkmanager.cpp \
    networkdevice.cpp \
    util/imageutil.cpp \
    item/deviceitem.cpp \
    item/wirelessitem.cpp \
    item/applet/wirelessapplet.cpp \
    item/applet/devicecontrolwidget.cpp \
    item/applet/accesspoint.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
    resources.qrc

