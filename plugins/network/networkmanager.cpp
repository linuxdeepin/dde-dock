#include "networkmanager.h"

NetworkManager *NetworkManager::INSTANCE = nullptr;

NetworkManager *NetworkManager::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new NetworkManager(parent);

    return INSTANCE;
}

NetworkManager::NetworkManager(QObject *parent)
    : QObject(parent),

      m_networkInter(new DBusNetwork(this))
{
    qDebug() << m_networkInter->activeConnections();

    QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->activeConnections().toUtf8());
    qDebug() << doc;

    QJsonObject obj = doc.object();
    for (auto value : obj)
    {
        qDebug() << value.toObject().value("Uuid").toString();
    }

    qDebug() << QJsonDocument::fromJson(m_networkInter->devices().toUtf8());
}
