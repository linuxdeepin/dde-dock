#ifndef ACCESSPOINTWIDGET_H
#define ACCESSPOINTWIDGET_H

#include "accesspoint.h"

#include <QWidget>
#include <QLabel>

class AccessPointWidget : public QWidget
{
    Q_OBJECT

public:
    explicit AccessPointWidget(const AccessPoint &ap);

private:
    QLabel *m_ssid;
};

#endif // ACCESSPOINTWIDGET_H
