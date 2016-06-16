#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"

#include <QObject>
#include <QSize>

class DockSettings : public QObject
{
    Q_OBJECT

public:
    explicit DockSettings(QObject *parent = 0);

    DockSide side() const;
    const QSize mainWindowSize() const;

public slots:
    void updateGeometry();

signals:
    void dataChanged() const;

private:
    QSize m_mainWindowSize;
};

#endif // DOCKSETTINGS_H
