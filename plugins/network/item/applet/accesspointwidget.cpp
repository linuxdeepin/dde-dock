#include "accesspointwidget.h"
#include "horizontalseperator.h"

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
    {
        QPixmap pixmap(16, 16);
        pixmap.fill(Qt::transparent);
        m_securityIcon->setPixmap(pixmap);
    }

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_securityIcon);
    infoLayout->addSpacing(5);
    infoLayout->addWidget(m_strengthIcon);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_ssidBtn);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(15, 0, 0, 0);

    HorizontalSeperator *seperator = new HorizontalSeperator;
    seperator->setFixedHeight(1);

    QVBoxLayout *centeralLayout = new QVBoxLayout;
    centeralLayout->addWidget(seperator);
    centeralLayout->addLayout(infoLayout);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setStrengthIcon(ap.strength());
    setLayout(centeralLayout);
    setStyleSheet("AccessPointWidget #Ssid {"
                  "color:white;"
                  "background-color:transparent;"
                  "border:none;"
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
