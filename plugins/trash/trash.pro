
include(../../interfaces/interfaces.pri)

QT              += widgets svg
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(trash)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += trash.json

HEADERS += \
    trashplugin.h \
    trashwidget.h

SOURCES += \
    trashplugin.cpp \
    trashwidget.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
