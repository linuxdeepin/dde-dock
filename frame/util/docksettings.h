#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include <QObject>
#include <QSize>

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
    const QSize mainWindowSize() const;

public slots:
    void updateGeometry();

signals:
    void dataChanged() const;

private:
    QSize m_mainWindowSize;
};

#endif // DOCKSETTINGS_H
