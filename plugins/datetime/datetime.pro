
QT              += widgets
TEMPLATE         = lib
CONFIG          += plugin c++11 link_pkgconfig
PKGCONFIG       +=

TARGET          = $$qtLibraryTarget(datetime)
DESTDIR         = $$_PRO_FILE_PWD_/../

HEADERS += \
    datetimeplugin.h \
    datetimeitem.h

SOURCES += \
    datetimeplugin.cpp \
    datetimeitem.cpp

include(../../interfaces/interfaces.pri)

INCLUDEPATH += "../../frame/item"
