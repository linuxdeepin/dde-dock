
include(../../interfaces/interfaces.pri)

QT              += widgets
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(system-tray)
DESTDIR          = $$_PRO_FILE_PWD_/../

HEADERS += \
    systemtrayplugin.h \
    dbus/dbustraymanager.h

SOURCES += \
    systemtrayplugin.cpp \
    dbus/dbustraymanager.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target
