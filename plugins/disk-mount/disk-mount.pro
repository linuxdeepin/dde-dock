
include(../../interfaces/interfaces.pri)

QT              += widgets svg dbus
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(disk-mount)
DESTDIR          = $$_PRO_FILE_PWD_/../
DISTFILES       += disk-mount.json

HEADERS += \
    diskmountplugin.h \
    dbus/dbusdiskmount.h \
    dbus/variant/diskinfo.h \
    diskcontrolwidget.h \
    diskpluginitem.h \
    imageutil.h

SOURCES += \
    diskmountplugin.cpp \
    dbus/dbusdiskmount.cpp \
    dbus/variant/diskinfo.cpp \
    diskcontrolwidget.cpp \
    diskpluginitem.cpp \
    imageutil.cpp

target.path = $${PREFIX}/lib/dde-dock/plugins/
INSTALLS += target

RESOURCES += \
    resources.qrc
