#include "accesspointwidget.h"
#include "horizontalseperator.h"

#include <QHBoxLayout>
#include <QDebug>

DWIDGET_USE_NAMESPACE

AccessPointWidget::AccessPointWidget(const AccessPoint &ap)
    : QWidget(nullptr),

      m_active(false),
      m_ap(ap),
      m_ssidBtn(new QPushButton(this)),
      m_disconnectBtn(new DImageButton(this)),
      m_securityIcon(new QLabel),
      m_strengthIcon(new QLabel)
{
    m_ssidBtn->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_ssidBtn->setText(ap.ssid());
    m_ssidBtn->setObjectName("Ssid");

    m_disconnectBtn->setVisible(false);
    m_disconnectBtn->setNormalPic(":/wireless/resources/wireless/selected.png");
    m_disconnectBtn->setHoverPic(":/wireless/resources/wireless/disconnect.png");
    m_disconnectBtn->setPressPic(":/wireless/resources/wireless/disconnect_pressed.png");

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
    infoLayout->addWidget(m_disconnectBtn);
    infoLayout->addSpacing(20);
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
//                  "color:#2ca7f8;"
                  "}");

    connect(m_ssidBtn, &QPushButton::clicked, this, &AccessPointWidget::ssidClicked);
    connect(m_disconnectBtn, &DImageButton::clicked, this, &AccessPointWidget::disconnectBtnClicked);
}

bool AccessPointWidget::active() const
{
    return m_active;
}

void AccessPointWidget::setActive(const bool active)
{
    if (m_active == active)
        return;

    m_disconnectBtn->setVisible(active);
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

void AccessPointWidget::disconnectBtnClicked()
{
    emit requestDeactiveAP(m_ap);

    setActive(false);
}
