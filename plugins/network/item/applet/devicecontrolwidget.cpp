#include "devicecontrolwidget.h"
#include "horizontalseperator.h"

#include <QHBoxLayout>
#include <QDebug>

DWIDGET_USE_NAMESPACE

DeviceControlWidget::DeviceControlWidget(QWidget *parent)
    : QWidget(parent)
{
    m_deviceName = new QLabel;
    m_deviceName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_deviceName->setStyleSheet("color:white;");

    m_switchBtn = new DSwitchButton;

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_deviceName);
    infoLayout->addWidget(m_switchBtn);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(15, 0, 5, 0);

    m_seperator = new HorizontalSeperator;
    m_seperator->setFixedHeight(1);
    m_seperator->setColor(0.1);

    QVBoxLayout *centeralLayout = new QVBoxLayout;
    centeralLayout->addStretch();
    centeralLayout->addLayout(infoLayout);
    centeralLayout->addStretch();
    centeralLayout->addWidget(m_seperator);
    centeralLayout->setMargin(0);
    centeralLayout->setSpacing(0);

    setLayout(centeralLayout);
    setFixedHeight(30);

    connect(m_switchBtn, &DSwitchButton::checkedChanged, this, &DeviceControlWidget::deviceEnableChanged);
}

void DeviceControlWidget::setDeviceName(const QString &name)
{
    m_deviceName->setText(name);
}

void DeviceControlWidget::setDeviceEnabled(const bool enable)
{
    m_switchBtn->setChecked(enable);
}

void DeviceControlWidget::setSeperatorVisible(const bool visible)
{
    m_seperator->setVisible(visible);
}
