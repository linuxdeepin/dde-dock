
include(../../interfaces/interfaces.pri)

QT              += widgets
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(shutdown)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += shutdown.json

HEADERS += \
    shutdownplugin.h

SOURCES += \
    shutdownplugin.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
    resources.qrc
