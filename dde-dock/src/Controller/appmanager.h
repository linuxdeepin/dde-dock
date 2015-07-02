#ifndef APPMANAGER_H
#define APPMANAGER_H

#include <QObject>
#include "DBus/dbusentrymanager.h"
#include "DBus/dbusentryproxyer.h"

class AppManager : public QObject
{
    Q_OBJECT
public:
    explicit AppManager(QObject *parent = 0);

signals:

public slots:
};

#endif // APPMANAGER_H
