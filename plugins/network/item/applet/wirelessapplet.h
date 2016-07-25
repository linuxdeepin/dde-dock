#ifndef WIRELESSAPPLET_H
#define WIRELESSAPPLET_H

#include "devicecontrolwidget.h"
#include "../../dbus/dbusnetwork.h"

#include <QScrollArea>
#include <QVBoxLayout>

class WirelessApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessApplet(const QString &devicePath, QWidget *parent = 0);

private:
    void setDeviceInfo();

private slots:
    void APChanged(const QString &devPath, const QString &info);

private:
    const QString m_devicePath;

    QVBoxLayout *m_centeralLayout;
    QWidget *m_centeralWidget;
    DeviceControlWidget *m_controlPanel;
    DBusNetwork *m_networkInter;
};

#endif // WIRELESSAPPLET_H
