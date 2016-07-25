#ifndef WIRELESSAPPLET_H
#define WIRELESSAPPLET_H

#include <QScrollArea>

class WirelessApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessApplet(const QString &devicePath, QWidget *parent = 0);

private:
    const QString m_devicePath;
};

#endif // WIRELESSAPPLET_H
