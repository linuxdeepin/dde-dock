#ifndef ACCESSPOINTWIDGET_H
#define ACCESSPOINTWIDGET_H

#include "accesspoint.h"

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QDBusObjectPath>

#include <dimagebutton.h>

class AccessPointWidget : public QWidget
{
    Q_OBJECT
    Q_PROPERTY(bool active READ active WRITE setActive DESIGNABLE true)

public:
    explicit AccessPointWidget(const AccessPoint &ap);

    bool active() const;
    void setActive(const bool active);

signals:
    void requestActiveAP(const QDBusObjectPath &apPath, const QString &ssid) const;
    void requestDeactiveAP(const AccessPoint &ap) const;

private:
    void setStrengthIcon(const int strength);

private slots:
    void ssidClicked();
    void disconnectBtnClicked();

private:
    bool m_active;

    AccessPoint m_ap;
    QPushButton *m_ssidBtn;
    Dtk::Widget::DImageButton *m_disconnectBtn;
    QLabel *m_securityIcon;
    QLabel *m_strengthIcon;
};

#endif // ACCESSPOINTWIDGET_H
