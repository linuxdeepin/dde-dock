#include "devicecontrolwidget.h"
#include "horizontalseperator.h"
#include "refreshbutton.h"

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

    RefreshButton *refreshBtn = new RefreshButton;
    refreshBtn->setVisible(m_switchBtn->checked());

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_deviceName);
    infoLayout->addWidget(refreshBtn);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_switchBtn);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(15, 0, 5, 0);

//    m_seperator = new HorizontalSeperator;
//    m_seperator->setFixedHeight(1);
//    m_seperator->setColor(Qt::black);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addStretch();
    centralLayout->addLayout(infoLayout);
    centralLayout->addStretch();
//    centralLayout->addWidget(m_seperator);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedHeight(30);

    connect(m_switchBtn, &DSwitchButton::checkedChanged, this, &DeviceControlWidget::deviceEnableChanged);
    connect(m_switchBtn, &DSwitchButton::checkedChanged, refreshBtn, &RefreshButton::setVisible);
    connect(refreshBtn, &RefreshButton::clicked, this, &DeviceControlWidget::requestRefresh);
}

void DeviceControlWidget::setDeviceName(const QString &name)
{
    m_deviceName->setText(name);
}

void DeviceControlWidget::setDeviceEnabled(const bool enable)
{
    m_switchBtn->setChecked(enable);
}

//void DeviceControlWidget::setSeperatorVisible(const bool visible)
//{
//    m_seperator->setVisible(visible);
//}
