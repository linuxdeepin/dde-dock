#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "deviceitem.h"

#include <QWidget>

class WiredItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WiredItem(const QUuid &deviceUuid);

protected:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    QSize sizeHint() const;

private:
    void reloadIcon();

private:
    QPixmap m_icon;
};

#endif // WIREDITEM_H
