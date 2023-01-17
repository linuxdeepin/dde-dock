/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "soundwidget.h"
#include "imageutil.h"
#include "imageutil.h"
#include "slidercontainer.h"

#include <DGuiApplicationHelper>

#include <QDBusConnectionInterface>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDBusPendingReply>
#include <QDebug>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QLabel>
#include <QEvent>
#include <QHBoxLayout>
#include <QMetaMethod>
#include <QPainter>

DGUI_USE_NAMESPACE

#define ICON_SIZE 18
#define BACKSIZE 36

SoundWidget::SoundWidget(QWidget *parent)
    : QWidget(parent)
    , m_dbusAudio(new DBusAudio("org.deepin.dde.Audio1", "/org/deepin/dde/Audio1", QDBusConnection::sessionBus(), this))
    , m_sliderContainer(new SliderContainer(this))
    , m_defaultSink(new DBusSink("org.deepin.dde.Audio1", m_dbusAudio->defaultSink().path(), QDBusConnection::sessionBus(), this))
{
    initUi();
    initConnection();
}

SoundWidget::~SoundWidget()
{
}

void SoundWidget::initUi()
{
    if (m_defaultSink)
        m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100);

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(17, 0, 12, 0);
    mainLayout->addWidget(m_sliderContainer);

    onThemeTypeChanged();
    m_sliderContainer->setRange(0, std::round(m_dbusAudio->maxUIVolume() * 100.00));
    m_sliderContainer->setPageStep(2);

    SliderProxyStyle *proxy = new SliderProxyStyle;
    m_sliderContainer->setSliderProxyStyle(proxy);

    setEnabled(existActiveOutputDevice());
}

void SoundWidget::initConnection()
{
    connect(m_defaultSink, &DBusSink::VolumeChanged, this, [ this ](double value) { m_sliderContainer->updateSliderValue(std::round(value * 100.00)); });
    connect(m_defaultSink, &DBusSink::MuteChanged, this, [ = ] { m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100); });

    connect(m_dbusAudio, &DBusAudio::DefaultSinkChanged, this, [ this ](const QDBusObjectPath &value) {
        if (m_defaultSink)
            delete m_defaultSink;

        m_defaultSink = new DBusSink("org.deepin.dde.Audio1", value.path(), QDBusConnection::sessionBus(), this);
        connect(m_defaultSink, &DBusSink::VolumeChanged, this, [ this ](double value) { m_sliderContainer->updateSliderValue(std::round(value * 100.00)); });
        connect(m_defaultSink, &DBusSink::MuteChanged, this, [ = ] { m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100); });

        m_sliderContainer->updateSliderValue(std::round(m_defaultSink->volume() * 100.00));
    });

    connect(m_dbusAudio, &DBusAudio::MaxUIVolumeChanged, this, [ = ] (double maxValue) {
        m_sliderContainer->setRange(0, std::round(maxValue * 100.00));
    });

    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        m_defaultSink->SetVolume(value * 0.01, true);
    });

    connect(m_defaultSink, &DBusSink::MuteChanged, this, [ this ] {
        m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, QIcon(leftIcon()));
    });

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &SoundWidget::onThemeTypeChanged);

    connect(m_sliderContainer, &SliderContainer::iconClicked, this, [ this ](const SliderContainer::IconPosition icon) {
        switch (icon) {
        case SliderContainer::IconPosition::LeftIcon: {
            if (existActiveOutputDevice())
                m_defaultSink->SetMute(!m_defaultSink->mute());
            break;
        }
        case SliderContainer::IconPosition::RightIcon: {
            // 弹出音量选择对话框
            Q_EMIT rightIconClick();
            break;
        }
        }
    });
}

const QString SoundWidget::leftIcon()
{
    const bool mute = existActiveOutputDevice() ? m_defaultSink->mute() : true;
    if (mute) {
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
            return QString(":/audio-volume-muted-symbolic-dark.svg");

        return QString(":/audio-volume-muted-symbolic.svg");
    }

    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
        return QString(":/audio-volume-medium-symbolic-dark.svg");

    return QString(":/audio-volume-medium-symbolic.svg");
}

const QString SoundWidget::rightIcon()
{
    return QString(":/icons/resources/broadcast.svg");
}

void SoundWidget::convertThemePixmap(QPixmap &pixmap)
{
    // 图片是黑色的，如果当前主题为白色主题，则无需转换
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
        return;

    // 如果是黑色主题，则转换成白色图像
    QPainter painter(&pixmap);
    painter.setCompositionMode(QPainter::CompositionMode_SourceIn);
    painter.fillRect(pixmap.rect(), Qt::white);
    painter.end();
}

/** 判断是否存在未禁用的声音输出设备
 * @brief SoundApplet::existActiveOutputDevice
 * @return 存�