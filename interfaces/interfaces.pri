HEADERS += \
    $$PWD/pluginsiteminterface.h \
    $$PWD/constants.h \
    $$PWD/pluginproxyinterface.h

SOURCES += \

INCLUDEPATH += $$PWD

isEmpty(PREFIX)
{
    PREFIX = /usr
}
