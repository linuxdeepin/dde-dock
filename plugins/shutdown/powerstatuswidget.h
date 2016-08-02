#ifndef POWERSTATUSWIDGET_H
#define POWERSTATUSWIDGET_H

#include "dbus/dbuspower.h"

#include <QWidget>

class PowerStatusWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PowerStatusWidget(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);

private:
    QPixmap getBatteryIcon();

private:
    DBusPower *m_powerInter;
};

#endif // POWERSTATUSWIDGET_H
