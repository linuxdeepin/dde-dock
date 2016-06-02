#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    DockItemController *instance(QObject *parent);

signals:
    void dockItemCountChanged() const;

private:
    explicit DockItemController(QObject *parent = 0);

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
