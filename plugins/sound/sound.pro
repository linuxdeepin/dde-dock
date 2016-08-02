
include(../../interfaces/interfaces.pri)

QT              += widgets svg
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       += dtkbase dtkwidget

TARGET          = $$qtLibraryTarget(sound)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += sound.json

HEADERS += \
    soundplugin.h \
    sounditem.h

SOURCES += \
    soundplugin.cpp \
    sounditem.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
