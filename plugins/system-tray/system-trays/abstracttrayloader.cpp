#include "abstracttrayloader.h"

#include <QDebug>

AbstractTrayLoader::AbstractTrayLoader(const QString &waitService, QObject *parent)
    : QObject(parent),
      m_dbusDaemonInterface(QDBusConnection::sessionBus().interface()),
      m_waitingService(waitService)
{
}

bool AbstractTrayLoader::serviceExist()
{
    bool exist = m_dbusDaemonInterface->isServiceRegistered(m_waitingService).value();

    if (!exist) {
        qDebug() << m_waitingService << "daemon has not started";
    }

    return exist;
}

void AbstractTrayLoader::waitServiceForLoad()
{
    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &AbstractTrayLoader::onServiceOwnerChanged, Qt::UniqueConnection);
}

void AbstractTrayLoader::onServiceOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (m_waitingService.isEmpty() || newOwner.isEmpty()) {
        return;
    }

    if (m_waitingService == name) {
        qDebug() << m_waitingService << "daemon started, load tray and disconnect";
        load();
        disconnect(m_dbusDaemonInterface);
        m_waitingService = QString();
    }
}
