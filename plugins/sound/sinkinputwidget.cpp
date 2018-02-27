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

#include "sinkinputwidget.h"

#include <QHBoxLayout>
#include <QIcon>
#include <QApplication>
#include <DHiDPIHelper>

DWIDGET_USE_NAMESPACE

const QPixmap getIconFromTheme(const QString &name, const QSize &size)
{
    const auto ratio = qApp->devicePixelRatio();
    QPixmap ret = QIcon::fromTheme(name, QIcon::fromTheme("application-x-desktop")).pixmap(size * ratio);
    ret.setDevicePixelRatio(ratio);

    return ret;
}

SinkInputWidget::SinkInputWidget(const QString &inputPath, QWidget *parent)
    : QWidget(parent),

      m_inputInter(new DBusSinkInput(inputPath, this)),

      m_volumeIcon(new DImageButton),
      m_volumeSlider(new VolumeSlider)
{
    const QString iconName = m_inputInter->icon();
    m_volumeIcon->setAccessibleName("app-" + iconName + "-icon");
    m_volumeIcon->setPixmap(getIconFromTheme(iconName, QSize(24, 24)));
    m_volumeSlider->setAccessibleName("app-" + iconName + "-slider");

    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->addWidget(m_volumeIcon);
    centralLayout->addSpacing(10);
    centralLayout->addWidget(m_volumeSlider);
    centralLayout->setSpacing(2);
    centralLayout->setMargin(0);

    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SinkInputWidget::setVolume);
    connect(m_volumeSlider, &VolumeSlider::requestPlaySoundEffect, this, &SinkInputWidget::onPlaySoundEffect);
    connect(m_volumeIcon, &DImageButton::clicked, this, &SinkInputWidget::setMute);
    connect(m_inputInter, &DBusSinkInput::MuteChanged, this, &SinkInputWidget::setMuteIcon);
    connect(m_inputInter, &DBusSinkInput::VolumeChanged, this, [=] { m_volumeSlider->setValue(m_inputInter->volume() * 1000); });

    setLayout(centralLayout);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
    setFixedHeight(30);

    setMuteIcon();

    emit m_inputInter->VolumeChanged();
}

void SinkInputWidget::setVolume(const int value)
{
    m_inputInter->SetVolumeQueued(double(value) / 1000.0, false);
}

void SinkInputWidget::setMute()
{
    m_inputInter->SetMuteQueued(!m_inputInter->mute());
}

void SinkInputWidget::setMuteIcon()
{
    if (m_inputInter->mute()) {
        const auto ratio = devicePixelRatioF();
        QPixmap muteIcon = DHiDPIHelper::loadNxPixmap(":/icons/image/audio-volume-muted-symbolic.svg");
        QPixmap appIconSource(getIconFromTheme(m_inputInter->icon(), QSize(24, 24)));

        QPixmap temp(appIconSource.size());
        temp.fill(Qt::transparent);
        temp.setDevicePixelRatio(ratio);
        QPainter p1(&temp);
        p1.drawPixmap(0, 0, appIconSource);
        p1.setCompositionMode(QPainter::CompositionMode_DestinationIn);
        p1.fillRect(temp.rect(), QColor(0, 0, 0, 40));
        p1.end();
        appIconSource = temp;

        QPainter p(&appIconSource);
        p.drawPixmap(0, 0, muteIcon);

        appIconSource.setDevicePixelRatio(ratio);
        m_volumeIcon->setPixmap(appIconSource);
    } else {
        m_volumeIcon->setPixmap(getIconFromTheme(m_inputInter->icon(), QSize(24, 24)));
    }
}

void SinkInputWidget::onPlaySoundEffect()
{
    // set the mute property to false to play sound effects.
    m_inputInter->SetMuteQueued(false);
}
