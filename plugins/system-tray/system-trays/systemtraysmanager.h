#ifndef SYSTEMTRAYSMANAGER_H
#define SYSTEMTRAYSMANAGER_H

#include "abstracttrayloader.h"
#include "abstracttraywidget.h"

#include <QObject>

class SystemTraysManager : public QObject
{
    Q_OBJECT

public:
    explicit SystemTraysManager(QObject *parent = nullptr);

Q_SIGNALS:
    void systemTrayWidgetAdded(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void systemTrayWidgetRemoved(const QString &itemKey);

public Q_SLOTS:
    void startLoad();

private:
    QList<AbstractTrayLoader *> m_loaderList;
};

#endif // SYSTEMTRAYSMANAGER_H
