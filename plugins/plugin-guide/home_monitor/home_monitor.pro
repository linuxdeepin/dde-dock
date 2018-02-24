
QT              += widgets
TEMPLATE         = lib
CONFIG          += plugin c++11

TARGET          = $$qtLibraryTarget(home_monitor)
DESTDIR          = $$_PRO_FILE_PWD_
DISTFILES       += home_monitor.json

HEADERS += \
    homemonitorplugin.h \
    informationwidget.h

SOURCES += \
    homemonitorplugin.cpp \
    informationwidget.cpp

isEmpty(PREFIX) {
    PREFIX = /usr
}

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

