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
    m_volumeIcon->setPixmap(QIcon::fromTheme(m_inputInter->icon()).pixmap(24, 24));
    m_volumeSlider->setValue(m_inputInter->volume() * 1000);

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_volumeIcon);
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
