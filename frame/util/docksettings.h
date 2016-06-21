#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"
#include "dbus/dbusdock.h"
#include "controller/dockitemcontroller.h"

#include <QObject>
#include <QSize>

using namespace Dock;

class DockSettings : public QObject
{
    Q_OBJECT

public:
    explicit DockSettings(QObject *parent = 0);

    Position position() const;
    const QSize mainWindowSize() const;

signals:
    void dataChanged() const;

public slots:
    void updateGeometry();

private slots:

private:
    Position m_position;
    QSize m_mainWindowSize;

    DBusDock *m_dockInter;
    DockItemController *m_itemController;
};

#endif // DOCKSETTINGS_H
