#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "deviceitem.h"

#include <QWidget>
#include <QLabel>
#include <QTimer>

class WiredItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WiredItem(const QUuid &deviceUuid);

    NetworkDevice::NetworkType type() const;
    NetworkDevice::NetworkState state() const;
    QWidget *itemPopup();
    const QString itemCommand() const override;

protected:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);

private slots:
    void refreshIcon();
    void reloadIcon();
    void activeConnectionChanged(const QUuid &uuid);
    void deviceStateChanged(const NetworkDevice &device);

private:
    bool m_connected;
    QPixmap m_icon;

    QLabel *m_itemTips;
    QTimer *m_delayTimer;
};

#endif // WIREDITEM_H
