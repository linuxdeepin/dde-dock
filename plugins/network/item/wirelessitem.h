#ifndef WIRELESSITEM_H
#define WIRELESSITEM_H

#include "constants.h"

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
    void resizeEvent(QResizeEvent *e);

private:
    const QPixmap iconPix(const Dock::DisplayMode displayMode, const int size);
    const QPixmap backgroundPix(const int size);

private:
    QHash<QString, QPixmap> m_icons;
};

#endif // WIRELESSITEM_H
