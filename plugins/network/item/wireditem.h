#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "deviceitem.h"

#include <QWidget>
#include <QLabel>

class WiredItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WiredItem(const QUuid &deviceUuid);

    NetworkDevice::NetworkType type() const;
    NetworkDevice::NetworkState state() const;
    QWidget *itemApplet();

protected:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);

private:
    void reloadIcon();
    void activeConnectionChanged(const QUuid &uuid);

private:
    bool m_connected;
    QPixmap m_icon;

    QLabel *m_itemTips;
};

#endif // WIREDITEM_H
