#include "sinkinputwidget.h"

#include <QHBoxLayout>
#include <QIcon>

DWIDGET_USE_NAMESPACE

SinkInputWidget::SinkInputWidget(const QString &inputPath, QWidget *parent)
    : QWidget(parent),

      m_inputInter(new DBusSinkInput(inputPath, this)),

      m_volumeIcon(new DImageButton),
      m_volumeSlider(new VolumeSlider)
{
    const QString iconName = m_inputInter->icon();
    m_volumeIcon->setAccessibleName("app-" + iconName + "-icon");
    m_volumeIcon->setPixmap(QIcon::fromTheme(iconName).pixmap(24, 24));
    m_volumeSlider->setAccessibleName("app-" + iconName + "-slider");
    m_volumeSlider->setValue(m_inputInter->volume() * 1000);

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_volumeIcon);
    centeralLayout->addSpacing(10);
    centeralLayout->addWidget(m_volumeSlider);
    centeralLayout->setSpacing(2);
    centeralLayout->setMargin(0);

    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SinkInputWidget::setVolume);

    setLayout(centeralLayout);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
    setFixedHeight(30);
}

void SinkInputWidget::setVolume(const int value)
{
    m_inputInter->SetVolume(double(value) / 1000.0, false);
}
