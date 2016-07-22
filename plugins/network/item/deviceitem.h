#ifndef DEVICEITEM_H
#define DEVICEITEM_H

#include "networkmanager.h"

#include <QWidget>

class DeviceItem : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceItem(const NetworkDevice::NetworkType type, const QUuid &deviceUuid);

    const QUuid uuid() const;
    const NetworkDevice::NetworkType type() const;

protected:
    NetworkDevice::NetworkType m_type;
    QUuid m_deviceUuid;

    NetworkManager *m_networkManager;
};

#endif // DEVICEITEM_H
