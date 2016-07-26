#include "accesspointwidget.h"

#include <QHBoxLayout>
#include <QDebug>

AccessPointWidget::AccessPointWidget(const AccessPoint &ap)
    : QWidget(nullptr),

      m_ssid(new QLabel)
{

    m_ssid->setText(ap.ssid());
    m_ssid->setStyleSheet("color:white;");

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_ssid);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setLayout(centeralLayout);
}
