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
#include "../frame/util/imageutil.h"
#include "../widgets/tipswidget.h"

#include <QHBoxLayout>
#include <QIcon>
#include <QApplication>

#include <DHiDPIHelper>
#include <DGuiApplicationHelper>
#include <DApplication>
#include <DLabel>

#define ICON_SIZE   24
#define APP_TITLE_SIZE 110

DGUI_USE_NAMESPACE
using namespace Dock;

const QPixmap getIconFromTheme(const QString &name, const QSize &size, const qreal ratio)
{
    QPixmap ret = QIcon::fromTheme(name, QIcon::fromTheme("application-x-desktop")).pixmap(size * ratio);
    ret.setDevicePixelRatio(ratio);

    return ret;
}

SinkInputWidget::SinkInputWidget(const QString &inputPath, QWidget *parent)
    : QWidget(parent)
    , m_inputInter(new DBusSinkInput(inputPath, this))
    , m_appBtn(new DIconButton(this))
    , m_volumeBtnMin(new DIconButton(this))
    , m_volumeIconMax(new QLabel(this))
    , m_volumeSlider(new VolumeSlider(this))
    , m_volumeLabel(new TipsWidget(this))
{
    const QString iconName = m_inputInter->icon();
    m_appBtn->setAccessibleName("app-" + iconName + "-icon");
    m_appBtn->setIcon(getIconFromTheme(iconName, QSize(ICON_SIZE, ICON_SIZE), devicePixelRatioF()));
    m_appBtn->setIconSize(QSize(ICON_SIZE, ICON_SIZE));
    m_appBtn->setFlat(true);

    DLabel *titleLabel = new DLabel;
    titleLabel->setForegroundRole(DPalette::TextTitle);
    titleLabel->setText(titleLabel->fontMetrics().elidedText(m_inputInter->name(), Qt::TextElideMode::ElideRight, APP_TITLE_SIZE));

    m_volumeBtnMin->setAccessibleName("volume-button");
    m_volumeBtnMin->setFixedSize(ICON_SIZE, ICON_SIZE);
    m_volumeBtnMin->setIconSize(QSize(ICON_SIZE, ICON_SIZE));
    m_volumeBtnMin->setIcon(DHiDPIHelper::loadNxPixmap("://audio-volume-low-symbolic.svg"));
    m_volumeBtnMin->setFlat(true);

    m_volumeIconMax->setFixedSize(ICON_SIZE, ICON_SIZE);

    m_volumeSlider->setAccessibleName("app-" + iconName + "-slider");
    m_volumeSlider->setMinimum(0);
    m_volumeSlider->setMaximum(1000);

    // 应用图标+名称
    QHBoxLayout *appLayout = new QHBoxLayout();
    appLayout->setAlignment(Qt::AlignLeft);
    appLayout->addWidget(m_appBtn);
    appLayout->addSpacing(10);
    appLayout->addWidget(titleLabel);
    appLayout->addStretch();
    appLayout->addWidget(m_volumeLabel, 0, Qt::AlignRight);
    appLayout->setSpacing(0);
    appLayout->setMargin(0);

    // 音量图标+slider
    QHBoxLayout *volumeCtrlLayout = new QHBoxLayout;
    volumeCtrlLayout->addSpacing(2);
    volumeCtrlLayout->addWidget(m_volumeBtnMin);
    volumeCtrlLayout->addSpacing(10);
    volumeCtrlLayout->addWidget(m_volumeSlider);
    volumeCtrlLayout->addSpacing(10);
    volumeCtrlLayout->addWidget(m_volumeIconMax);
    volumeCtrlLayout->setSpacing(0);
    volumeCtrlLayout->setMargin(0);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addLayout(appLayout);
    centralLayout->addSpacing(6);
    centralLayout->addLayout(volumeCtrlLayout);
    centralLayout->setSpacing(2);
    centralLayout->setMargin(0);

    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SinkInputWidget::setVolume);
    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SinkInputWidget::onVolumeChanged);
//    connect(m_volumeSlider, &VolumeSlider::requestPlaySoundEffect, this, &SinkInputWidget::onPlaySoundEffect);
    connect(m_volumeBtnMin, &DIconButton::clicked, this, &SinkInputWidget::setMute);
    connect(m_inputInter, &DBusSinkInput::MuteChanged, this, &SinkInputWidget::setMuteIcon);
    connect(m_inputInter, &DBusSinkInput::VolumeChanged, this, [ = ] {
        m_volumeSlider->setValue(m_inputInter->volume() * 1000);
        QString str = QString::number(int(m_inputInter->volume() * 100)) + '%';
        m_volumeLabel->setText(str);
        refreshIcon();
    });
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &SinkInputWidget::refreshIcon);
    connect(qApp, &DApplication::iconThemeChanged, this, &SinkInputWidget::refreshIcon);

    setLayout(centralLayout);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
    setFixedHeight(60);

    setMuteIcon();
    refreshIcon();
    onVolumeChanged();

    emit m_inputInter->VolumeChanged();
}

void SinkInputWidget::setVolume(const int value)
{
    m_inputInter->SetVolumeQueued(double(value) / 1000.0, false);
    refreshIcon();
}

void SinkInputWidget::setMute()
{
    m_inputInter->SetMuteQueued(!m_inputInter->mute());
}

void SinkInputWidget::setMuteIcon()
{
    if (m_inputInter->mute()) {
        const auto ratio = devicePixelRatioF();
        QString iconString = "audio-volume-muted-symbolic";
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
            iconString.append("-dark");
        }
        QPixmap muteIcon = QIcon::fromTheme(iconString).pixmap(ICON_SIZE * ratio, ICON_SIZE * ratio);
        muteIcon.setDevicePixelRatio(ratio);
        QPixmap appIconSource(getIconFromTheme(m_inputInter->icon(), QSize(ICON_SIZE, ICON_SIZE), devicePixelRatioF()));

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
        m_volumeBtnMin->setIcon(appIconSource);
    } else {
        m_volumeBtnMin->setIcon(getIconFromTheme(m_inputInter->icon(), QSize(ICON_SIZE, ICON_SIZE), devicePixelRatioF()));
    }

    refreshIcon();
}

void SinkInputWidget::onPlaySoundEffect()
{
    // set the mute property to false to play sound effects.
    m_inputInter->SetMuteQueued(false);
}

void SinkInputWidget::refreshIcon()
{
    if (!m_inputInter)
        return;

    QString iconLeft = QString(m_inputInter->mute() ? "audio-volume-muted-symbolic" : "audio-volume-low-symbolic");
    QString iconRight = QString("audio-volume-high-symbolic");

    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        iconLeft.append("-dark");
        iconRight.append("-dark");
    }

    const auto ratio = devicePixelRatioF();
    QPixmap ret = ImageUtil::loadSvg(iconRight, ":/", ICON_SIZE, ratio);
    m_volumeIconMax->setPixmap(ret);

    ret = ImageUtil::loadSvg(iconLeft, ":/", ICON_SIZE, ratio);
    m_volumeBtnMin->setIcon(ret);
}

void SinkInputWidget:: onVolumeChanged()
{
    QString str = QString::number(int(m_inputInter->volume() * 100)) + '%';
    m_volumeLabel->setText(str);
}
