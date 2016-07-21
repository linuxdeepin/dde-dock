#ifndef NETWORKMANAGER_H
#define NETWORKMANAGER_H

#include "dbus/dbusnetwork.h"

#include <QJsonObject>
#include <QJsonDocument>
#include <QJsonArray>

class NetworkManager : public QObject
{
    Q_OBJECT

public:
    static NetworkManager *instance(QObject *parent = nullptr);

private:
    explicit NetworkManager(QObject *parent = 0);

private:
    DBusNetwork *m_networkInter;

    static NetworkManager *INSTANCE;
};

#endif // NETWORKMANAGER_H
