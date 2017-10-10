
include(../../interfaces/interfaces.pri)

QT              += widgets svg dbus
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       += gsettings-qt dtkwidget

TARGET          = $$qtLibraryTarget(sound)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += sound.json

HEADERS += \
    soundplugin.h \
    sounditem.h \
    soundapplet.h \
    dbus/dbusaudio.h \
    dbus/dbussink.h \
    componments/horizontalseparator.h \
    componments/volumeslider.h \
    dbus/dbussinkinput.h \
    sinkinputwidget.h

SOURCES += \
    soundplugin.cpp \
    sounditem.cpp \
    soundapplet.cpp \
    dbus/dbusaudio.cpp \
    dbus/dbussink.cpp \
    componments/horizontalseparator.cpp \
    componments/volumeslider.cpp \
    dbus/dbussinkinput.cpp \
    sinkinputwidget.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
    resources/resources.qrc
