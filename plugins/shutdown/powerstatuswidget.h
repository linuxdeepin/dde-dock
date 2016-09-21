#ifndef POWERSTATUSWIDGET_H
#define POWERSTATUSWIDGET_H

#include "dbus/dbuspower.h"

#include <QWidget>

#define POWER_KEY       "power"

class PowerStatusWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PowerStatusWidget(QWidget *parent = 0);

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    QPixmap getBatteryIcon();

private:
    DBusPower *m_powerInter;
};

#endif // POWERSTATUSWIDGET_H
