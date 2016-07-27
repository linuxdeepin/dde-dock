#include "accesspointwidget.h"

#include <QHBoxLayout>
#include <QDebug>

AccessPointWidget::AccessPointWidget(const AccessPoint &ap)
    : QWidget(nullptr),

      m_ssid(new QLabel),
      m_securityIcon(new QLabel),
      m_strengthIcon(new QLabel)
{

    m_ssid->setText(ap.ssid());
    m_ssid->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_ssid->setAlignment(Qt::AlignVCenter | Qt::AlignLeft);
    m_ssid->setStyleSheet("color:white;");

    if (ap.secured())
        m_securityIcon->setPixmap(QPixmap(":/wireless/resources/wireless/security.svg"));
    else
        m_securityIcon->setPixmap(QPixmap(16, 16));

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_securityIcon);
    centeralLayout->addSpacing(5);
    centeralLayout->addWidget(m_strengthIcon);
    centeralLayout->addSpacing(10);
    centeralLayout->addWidget(m_ssid);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setStrengthIcon(ap.strength());
    setLayout(centeralLayout);
}

void AccessPointWidget::setStrengthIcon(const int strength)
{
    if (strength == 100)
        return m_strengthIcon->setPixmap(QPixmap(":/wireless/resources/wireless/wireless-8-symbolic.svg"));

    m_strengthIcon->setPixmap(QPixmap(QString(":/wireless/resources/wireless/wireless-%1-symbolic.svg").arg(strength / 10 & ~0x1)));
}
