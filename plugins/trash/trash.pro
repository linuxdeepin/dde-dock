
include(../../interfaces/interfaces.pri)

QT              += widgets svg dbus
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       += dtkwidget dframeworkdbus

TARGET          = $$qtLibraryTarget(trash)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += trash.json

HEADERS += \
    trashplugin.h \
    trashwidget.h \
    popupcontrolwidget.h

SOURCES += \
    trashplugin.cpp \
    trashwidget.cpp \
    popupcontrolwidget.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
