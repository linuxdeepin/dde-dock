#ifndef ACCESSPOINT_H
#define ACCESSPOINT_H

#include <QObject>

class AccessPoint : public QObject
{
    Q_OBJECT

public:
    explicit AccessPoint(QObject *parent = 0);
};

#endif // ACCESSPOINT_H
