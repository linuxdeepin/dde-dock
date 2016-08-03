#include "soundapplet.h"
#include "componments/horizontalseparator.h"

#include <QLabel>
#include <QIcon>

#define WIDTH       200
#define ICON_SIZE   24

DWIDGET_USE_NAMESPACE

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent),

      m_centeralWidget(new QWidget),
      m_appControlWidget(new QWidget),
      m_volumeBtn(new DImageButton),
      m_volumeSlider(new VolumeSlider),

      m_audioInter(new DBusAudio(this)),
      m_defSinkInter(nullptr)
{
    QIcon::setThemeName("deepin");

    QLabel *deviceLabel = new QLabel;
    deviceLabel->setText(tr("Device"));
    deviceLabel->setStyleSheet("color:white;");

    QHBoxLayout *deviceLineLayout = new QHBoxLayout;
    deviceLineLayout->addWidget(deviceLabel);
    deviceLineLayout->addWidget(new HorizontalSeparator);
    deviceLineLayout->setMargin(0);
    deviceLineLayout->setSpacing(10);

    QHBoxLayout *volumeCtrlLayout = new QHBoxLayout;
    volumeCtrlLayout->addWidget(m_volumeBtn);
    volumeCtrlLayout->addWidget(m_volumeSlider);
    volumeCtrlLayout->setSpacing(0);
    volumeCtrlLayout->setMargin(0);

    QLabel *appLabel = new QLabel;
    appLabel->setText(tr("Application"));
    appLabel->setStyleSheet("color:white;");

    QHBoxLayout *appLineLayout = new QHBoxLayout;
    appLineLayout->addWidget(appLabel);
    appLineLayout->addWidget(new HorizontalSeparator);
    appLineLayout->setMargin(0);
    appLineLayout->setSpacing(10);

    QVBoxLayout *appLayout = new QVBoxLayout;
    appLayout->addLayout(appLineLayout);
    appLayout->setSpacing(0);
    appLayout->setMargin(0);

    m_volumeBtn->setFixedSize(ICON_SIZE, ICON_SIZE);
    m_volumeSlider->setMinimum(0);
    m_volumeSlider->setMaximum(1000);

    m_appControlWidget->setLayout(appLayout);

    m_centeralLayout = new QVBoxLayout;
    m_centeralLayout->addLayout(deviceLineLayout);
    m_centeralLayout->addLayout(volumeCtrlLayout);
    m_centeralLayout->addWidget(m_appControlWidget);

    m_centeralWidget->setLayout(m_centeralLayout);
    m_centeralWidget->setFixedWidth(WIDTH);

    setFixedWidth(WIDTH);
    setWidget(m_centeralWidget);
    setFrameStyle(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    connect(m_volumeBtn, &DImageButton::clicked, this, &SoundApplet::toggleMute);
    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SoundApplet::volumeSliderValueChanged);
    connect(this, static_cast<void (SoundApplet::*)(DBusSink*) const>(&SoundApplet::defaultSinkChanged), this, &SoundApplet::onVolumeChanged);

    QMetaObject::invokeMethod(this, "defaultSinkChanged", Qt::QueuedConnection);
}

void SoundApplet::defaultSinkChanged()
{
    delete m_defSinkInter;

    const QDBusObjectPath defSinkPath = m_audioInter->GetDefaultSink();
    m_defSinkInter = new DBusSink(defSinkPath.path(), this);

    connect(m_defSinkInter, &DBusSink::VolumeChanged, this, &SoundApplet::onVolumeChanged);
    connect(m_defSinkInter, &DBusSink::MuteChanged, this, &SoundApplet::onVolumeChanged);

    emit defaultSinkChanged(m_defSinkInter);
}

void SoundApplet::onVolumeChanged()
{
    const double volmue = m_defSinkInter->volume();
    const bool mute = m_defSinkInter->mute() || volmue < 0.001;

    m_volumeSlider->setValue(std::min(1000.0, volmue * 1000));

    QString volumeString;
    if (mute)
        volumeString = "muted";
    else if (volmue >= double(2)/3)
        volumeString = "high";
    else if (volmue >= double(1)/3)
        volumeString = "medium";
    else
        volumeString = "low";

    const QString iconString = QString("audio-volume-%1-symbolic").arg(volumeString);
    m_volumeBtn->setPixmap(QIcon::fromTheme(iconString).pixmap(ICON_SIZE, ICON_SIZE));
}

void SoundApplet::volumeSliderValueChanged()
{
    m_defSinkInter->SetVolume(double(m_volumeSlider->value()) / 1000, false);
}

void SoundApplet::toggleMute()
{
    m_defSinkInter->SetMute(!m_defSinkInter->mute());
}
