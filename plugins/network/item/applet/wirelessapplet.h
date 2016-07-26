#ifndef WIRELESSAPPLET_H
#define WIRELESSAPPLET_H

#include "devicecontrolwidget.h"
#include "accesspoint.h"
#include "../../dbus/dbusnetwork.h"

#include <QScrollArea>
#include <QVBoxLayout>
#include <QList>
#include <QTimer>

class WirelessApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessApplet(const QString &devicePath, QWidget *parent = 0);

private:
    void setDeviceInfo();
    void loadAPList();

private slots:
    void init();
    void APPropertiesChanged(const QString &devPath, const QString &info);
    void updateAPList();
    void deviceEnableChanged(const bool enable);

private:
    const QString m_devicePath;

    QList<AccessPoint> m_apList;

    QTimer *m_updateAPTimer;

    QVBoxLayout *m_centeralLayout;
    QWidget *m_centeralWidget;
    DeviceControlWidget *m_controlPanel;
    DBusNetwork *m_networkInter;
};

#endif // WIRELESSAPPLET_H
