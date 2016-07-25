#ifndef WIRELESSITEM_H
#define WIRELESSITEM_H

#include "deviceitem.h"

#include <QHash>

class WirelessItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WirelessItem(const QUuid &uuid);

    QWidget *itemApplet();

protected:
    void paintEvent(QPaintEvent *e);

private:
    const QPixmap icon(const QString &key);

private:
    QHash<QString, QPixmap> m_icons;
};

#endif // WIRELESSITEM_H
