/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "soundapplet.h"
#include "sinkinputwidget.h"
#include "componments/horizontalseparator.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "util/utils.h"
#include <DGuiApplicationHelper>

#include <QLabel>
#include <QIcon>
#include <QScrollBar>
#include <DApplication>

#define WIDTH       200
#define MAX_HEIGHT  300
#define ICON_SIZE   24

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent)
    , m_centralWidget(new QWidget)
    , m_applicationTitle(new QWidget)
    , m_volumeBtn(new DImageButton)
    , m_volumeIconMax(new QLabel)
    , m_volumeSlider(new VolumeSlider)
    , m_soundShow(new Dock::TipsWidget)
    , m_audioInter(new DBusAudio(this))
    , m_defSinkInter(nullptr)
{
    //    QIcon::setThemeName("deepin");
    m_centralWidget->setAccessibleName("volumn-centralwidget");
    m_volumeBtn->setAccessibleName("volume-button");
    m_volumeIconMax->setAccessibleName("volume-iconmax");
    m_volumeSlider->setAccessibleName("volume-slider");
    m_soundShow->setAccessibleName("volume-soundtips");
    this->horizontalScrollBar()->setAccessibleName("volume-horizontalscrollbar");
    this->verticalScrollBar()->setAccessibleName("volume-verticalscrollbar");

    m_volumeIconMax->setFixedSize(ICON_SIZE, ICON_SIZE);

    m_soundShow->setText(QString("%1%").arg(0));

    Dock::TipsWidget *deviceLabel = new Dock::TipsWidget;
    deviceLabel->setText(tr("Device"));

    QHBoxLayout *deviceLayout =new QHBoxLayout;
    deviceLayout->addSpacing(2);
    deviceLayout->addWidget(deviceLabel,0, Qt::AlignLeft);
    deviceLayout->addWidget(m_soundShow,0,Qt::AlignRight);
    deviceLayout->setSpacing(0);
    deviceLayout->setMargin(0);

    QVBoxLayout *deviceLineLayout = new QVBoxLayout;
    deviceLineLayout->addLayout(deviceLayout);
    //    deviceLineLayout->addSpacing(12);
    deviceLineLayout->addWidget(new HorizontalSeparator);
    deviceLineLayout->setMargin(0);
    deviceLineLayout->setSpacing(10);

    QHBoxLayout *volumeCtrlLayout = new QHBoxLayout;
    volumeCtrlLayout->addSpacing(2);
    volumeCtrlLayout->addWidget(m_volumeBtn);
    volumeCtrlLayout->addSpacing(10);
    volumeCtrlLayout->addWidget(m_volumeSlider);
    volumeCtrlLayout->addSpacing(10);
    volumeCtrlLayout->addWidget(m_volumeIconMax);
    volumeCtrlLayout->setSpacing(0);
    volumeCtrlLayout->setMargin(0);

    Dock::TipsWidget *appLabel = new Dock::TipsWidget;
    appLabel->setText(tr("Application"));

    QVBoxLayout *appLineHLayout = new QVBoxLayout;
    appLineHLayout->addWidget(new HorizontalSeparator);
    appLineHLayout->addWidget(appLabel);
    appLineHLayout->setMargin(0);
    appLineHLayout->setSpacing(10);

    QVBoxLayout *appLineVLayout = new QVBoxLayout;
    appLineVLayout->addSpacing(10);
    appLineVLayout->addLayout(appLineHLayout);
    appLineVLayout->addSpacing(8);
    appLineVLayout->setSpacing(0);
    appLineVLayout->setMargin(0);

    m_applicationTitle->setLayout(appLineVLayout);
    m_applicationTitle->setAccessibleName("applicationtitle");

    m_volumeBtn->setFixedSize(ICON_SIZE, ICON_SIZE);
    m_volumeSlider->setMinimum(0);
    m_volumeSlider->setMaximum(m_audioInter->maxUIVolume() * 100.0f);

    m_centralLayout = new QVBoxLayout;
    m_centralLayout->addLayout(deviceLineLayout);
    m_centralLayout->addSpacing(8);
    m_centralLayout->addLayout(volumeCtrlLayout);
    m_centralLayout->addWidget(m_applicationTitle);

    m_centralWidget->setLayout(m_centralLayout);
    m_centralWidget->setFixedWidth(WIDTH);
    m_centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);

    setFixedWidth(WIDTH);
    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    connect(m_volumeBtn, &DImageButton::clicked, this, &SoundApplet::toggleMute);
    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SoundApplet::volumeSliderValueChanged);
    connect(m_volumeSlider, &VolumeSlider::requestPlaySoundEffect, this, &SoundApplet::onPlaySoundEffect);
    connect(m_audioInter, &DBusAudio::SinkInputsChanged, this, &SoundApplet::sinkInputsChanged);
    connect(m_audioInter, &DBusAudio::DefaultSinkChanged, this, static_cast<void (SoundApplet::*)()>(&SoundApplet::defaultSinkChanged));
    connect(m_audioInter, &DBusAudio::IncreaseVolumeChanged, this, &SoundApplet::increaseVolumeChanged);
    connect(this, static_cast<void (SoundApplet::*)(DBusSink *) const>(&SoundApplet::defaultSinkChanged), this, &SoundApplet::onVolumeChanged);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &SoundApplet::refreshIcon);
    connect(qApp, &DApplication::iconThemeChanged, this, &SoundApplet::refreshIcon);

    QMetaObject::invokeMethod(this, "defaultSinkChanged", Qt::QueuedConnection);
    QMetaObject::invokeMethod(this, "sinkInputsChanged", Qt::QueuedConnection);

    refreshIcon();
}

int SoundApplet::volumeValue() const
{
    return m_volumeSlider->value();
}

int SoundApplet::maxVolumeValue() const
{
    return m_volumeSlider->maximum();
}

VolumeSlider *SoundApplet::mainSlider()
{
    return m_volumeSlider;
}

void SoundApplet::defaultSinkChanged()
{
    delete m_defSinkInter;

    const QDBusObjectPath defSinkPath = m_audioInter->defaultSink();
    m_defSinkInter = new DBusSink(defSinkPath.path(), this);

    connect(m_defSinkInter, &DBusSink::VolumeChanged, this, &SoundApplet::onVolumeChanged);
    connect(m_defSinkInter, &DBusSink::MuteChanged, this, &SoundApplet::onVolumeChanged);

    emit defaultSinkChanged(m_defSinkInter);
}

void SoundApplet::onVolumeChanged()
{
    const float volume = m_defSinkInter->volume();

    m_volumeSlider->setValue(std::min(150.0f, volume * 100.0f));

    m_soundShow->setText(QString::number(volume * 100) + '%');
    emit volumeChanged(m_volumeSlider->value());
    refreshIcon();
}

void SoundApplet::volumeSliderValueChanged()
{
    m_defSinkInter->SetVolumeQueued(m_volumeSlider->value() / 100.0f, false);
}

void SoundApplet::sinkInputsChanged()
{
    m_centralWidget->setVisible(false);
    QVBoxLayout *appLayout = m_centralLayout;
    while (QLayoutItem *item = appLayout->takeAt(4)) {
        delete item->widget();
        delete item;
    }

    m_applicationTitle->setVisible(false);
    for (auto input : m_audioInter->sinkInputs()) {
        m_applicationTitle->setVisible(true);
        appLayout->addWidget(new HorizontalSeparator);

        SinkInputWidget *si = new SinkInputWidget(input.path());
        appLayout->addWidget(si);
    }

    const int contentHeight = m_centralWidget->sizeHint().height();
    m_centralWidget->setFixedHeight(contentHeight);
    m_centralWidget->setVisible(true);
    setFixedHeight(std::min(contentHeight, MAX_HEIGHT));
}

void SoundApplet::toggleMute()
{
    m_defSinkInter->SetMuteQueued(!m_defSinkInter->mute());
}

void SoundApplet::onPlaySoundEffect()
{

}

void SoundApplet::increaseVolumeChanged()
{
    m_volumeSlider->setMaximum(m_audioInter->maxUIVolume() * 100.0f);
}

void SoundApplet::refreshIcon()
{
    if (!m_defSinkInter)
        return;

    const bool mute = m_defSinkInter->mute();

    QString volumeString;

    if (mute) {
        volumeString = "muted";
    } else {
        volumeString = "low";
    }

    QString iconLeft = QString("audio-volume-%1-symbolic").arg(volumeString);
    QString iconRight = QString("audio-volume-high-symbolic");

    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        iconLeft.append("-dark");
        iconRight.append("-dark");
    }

    const auto ratio = devicePixelRatioF();
    QPixmap ret = ImageUtil::loadSvg(iconRight, ":/", ICON_SIZE, ratio);
    m_volumeIconMax->setPixmap(ret);

    ret = ImageUtil::loadSvg(iconLeft, ":/", ICON_SIZE, ratio);
    m_volumeBtn->setPixmap(ret);
}
