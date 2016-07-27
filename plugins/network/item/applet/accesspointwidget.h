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
    void setStrengthIcon(const int strength);

private:
    QLabel *m_ssid;
    QLabel *m_securityIcon;
    QLabel *m_strengthIcon;
};

#endif // ACCESSPOINTWIDGET_H
