#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include <QObject>

class DockSettings : public QObject
{
    Q_OBJECT

public:
    enum DockSide {
        Top,
        Bottom,
        Left,
        Right,
    };

public:
    explicit DockSettings(QObject *parent = 0);

    DockSide side() const;

signals:
    void dockSideChanged(const DockSide side) const;
};

#endif // DOCKSETTINGS_H
