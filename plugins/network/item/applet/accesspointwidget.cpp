#include "accesspointwidget.h"

#include <QHBoxLayout>
#include <QDebug>

AccessPointWidget::AccessPointWidget(const AccessPoint &ap)
    : QWidget(nullptr),

      m_active(false),
      m_ap(ap),
      m_ssidBtn(new QPushButton(this)),
      m_securityIcon(new QLabel),
      m_strengthIcon(new QLabel)
{
    m_ssidBtn->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_ssidBtn->setText(ap.ssid());
    m_ssidBtn->setObjectName("Ssid");

    if (ap.secured())
        m_securityIcon->setPixmap(QPixmap(":/wireless/resources/wireless/security.svg"));
    else
        m_securityIcon->setPixmap(QPixmap(16, 16));

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_securityIcon);
    centeralLayout->addSpacing(5);
    centeralLayout->addWidget(m_strengthIcon);
    centeralLayout->addSpacing(10);
    centeralLayout->addWidget(m_ssidBtn);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setStrengthIcon(ap.strength());
    setLayout(centeralLayout);
    setStyleSheet("AccessPointWidget #Ssid {"
                  "color:white;"
                  "text-align:left;"
                  "}"
                  "AccessPointWidget[active=true] #Ssid {"
                  "color:#2ca7f8;"
                  "}");

    connect(m_ssidBtn, &QPushButton::clicked, this, &AccessPointWidget::ssidClicked);
}

bool AccessPointWidget::active() const
{
    return m_active;
}

void AccessPointWidget::setActive(const bool active)
{
    if (m_active == active)
        return;

    m_active = active;
    setStyleSheet(styleSheet());
}

void AccessPointWidget::setStrengthIcon(const int strength)
{
    if (strength == 100)
        return m_strengthIcon->setPixmap(QPixmap(":/wireless/resources/wireless/wireless-8-symbolic.svg"));

    m_strengthIcon->setPixmap(QPixmap(QString(":/wireless/resources/wireless/wireless-%1-symbolic.svg").arg(strength / 10 & ~0x1)));
}

void AccessPointWidget::ssidClicked()
{
    emit requestActiveAP(QDBusObjectPath(m_ap.path()), m_ap.ssid());
}
