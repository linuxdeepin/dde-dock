HEADERS += \
    $$PWD/pluginsiteminterface.h \
    $$PWD/constants.h \
    $$PWD/pluginproxyinterface.h

INCLUDEPATH += $$PWD

isEmpty(PREFIX)
{
    PREFIX = /usr
}
