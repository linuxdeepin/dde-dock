#ifndef WIRELESSITEM_H
#define WIRELESSITEM_H

#include "constants.h"

#include "deviceitem.h"
#include "applet/wirelessapplet.h"

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
    const QPixmap cachedPix(const QString &key, const int size);

private slots:
    void init();

private:
    QHash<QString, QPixmap> m_icons;

    WirelessApplet *m_applet;
};

#endif // WIRELESSITEM_H
