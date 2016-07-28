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

    void setActive(const bool active);

private:
    void setStrengthIcon(const int strength);

private:
    bool m_active;

    AccessPoint m_ap;
    QLabel *m_ssid;
    QLabel *m_securityIcon;
    QLabel *m_strengthIcon;
};

#endif // ACCESSPOINTWIDGET_H
