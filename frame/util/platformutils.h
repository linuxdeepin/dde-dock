#ifndef PLATFORMUTILS_H
#define PLATFORMUTILS_H

#include <QObject>

class PlatformUtils
{
public:
    static QString getAppNameForWindow(quint32 winId);

private:
    static QString getWindowProperty(quint32 winId, QString propName);
};

#endif // PLATFORMUTILS_H
