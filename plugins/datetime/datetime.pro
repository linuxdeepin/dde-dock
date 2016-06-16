
include(../../interfaces/interfaces.pri)

QT              += widgets
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(datetime)
DESTDIR          = $$_PRO_FILE_PWD_/../

HEADERS += \
    datetimeplugin.h

SOURCES += \
    datetimeplugin.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target
