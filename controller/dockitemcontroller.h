#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include "dbus/dbusdockentrymanager.h"

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    static DockItemController *instance(QObject *parent);

signals:
    void dockItemCountChanged() const;

private:
    explicit DockItemController(QObject *parent = 0);

    DBusDockEntryManager *m_entryManager;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
