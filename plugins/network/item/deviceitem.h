#ifndef DEVICEITEM_H
#define DEVICEITEM_H

#include "networkmanager.h"

#include <QWidget>

class DeviceItem : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceItem(const QUuid &deviceUuid);

    const QUuid uuid() const;

    virtual NetworkDevice::NetworkType type() const = 0;
    virtual NetworkDevice::NetworkState state() const = 0;
    virtual const QString itemCommand() const;
    virtual QWidget *itemApplet();
    virtual QWidget *itemPopup();

protected:
    QSize sizeHint() const;

protected:
    QUuid m_deviceUuid;

    NetworkManager *m_networkManager;
};

#endif // DEVICEITEM_H
