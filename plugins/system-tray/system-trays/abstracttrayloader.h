#ifndef ABSTRACTTRAYLOADER_H
#define ABSTRACTTRAYLOADER_H

#include "abstracttraywidget.h"

#include <QDBusConnectionInterface>
#include <QObject>

class AbstractTrayLoader : public QObject
{
    Q_OBJECT
public:
    explicit AbstractTrayLoader(const QString &waitService, QObject *parent = nullptr);

Q_SIGNALS:
    void systemTrayAdded(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void systemTrayRemoved(const QString &itemKey);

public Q_SLOTS:
    virtual void load() = 0;

public:
    inline bool waitService() { return !m_waitingService.isEmpty(); }
    bool serviceExist();
    void waitServiceForLoad();

private Q_SLOTS:
    void onServiceOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);

private:
    QDBusConnectionInterface *m_dbusDaemonInterface;

    QString m_waitingService;
};

#endif // ABSTRACTTRAYLOADER_H
