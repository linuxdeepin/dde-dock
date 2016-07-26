#include "devicecontrolwidget.h"

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

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_deviceName);
    centeralLayout->addWidget(m_switchBtn);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setLayout(centeralLayout);
    setFixedHeight(20);

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
