HEADERS += \
    $$PWD/pluginsiteminterface.h \
    $$PWD/constants.h

INCLUDEPATH += $$PWD

isEmpty(PREFIX)
{
    PREFIX = /usr
}
