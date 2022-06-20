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
#include "volumewidget.h"
#include "brightnessmodel.h"
#include "imageutil.h"
#include "volumemodel.h"
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

DGUI_USE_NAMESPACE

#define ICON_SIZE 24
#define BACKSIZE 36

VolumeWidget::VolumeWidget(VolumeModel *model, QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_model(model)
    , m_sliderContainer(new SliderContainer(this))
    , m_defaultSink(m_model->defaultSink())
{
    initUi();
    initConnection();
}

VolumeWidget::~VolumeWidget()
{
}

void VolumeWidget::initUi()
{
    if (m_defaultSink)
        m_sliderContainer->slider()->setValue(m_defaultSink->volume());

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(17, 0, 12, 0);
    mainLayout->addWidget(m_sliderContainer);

    QPixmap leftPixmap = ImageUtil::loadSvg(leftIcon(), QSize(ICON_SIZE, ICON_SIZE));
    QPixmap rightPixmap = ImageUtil::loadSvg(rightIcon(), QSize(ICON_SIZE, ICON_SIZE));
    m_sliderContainer->updateSlider(SliderContainer::IconPosition::LeftIcon, { leftPixmap.size(), QSize(), leftPixmap, 12});
    m_sliderContainer->updateSlider(SliderContainer::IconPosition::RightIcon, { rightPixmap.size(), QSize(BACKSIZE, BACKSIZE), rightPixmap, 12});

    SliderProxyStyle *proxy = new SliderProxyStyle;
    proxy->setParent(m_sliderContainer->slider());
    m_sliderContainer->slider()->setStyle(proxy);

    bool existActiveOutputDevice = m_model->existActiveOutputDevice();
    setEnabled(existActiveOutputDevice);
}

void VolumeWidget::initConnection()
{
    auto setCtrlVolumeValue = [this](int volume) {
        m_sliderContainer->blockSignals(true);
        m_sliderContainer->slider()->setValue(volume);
        m_sliderContainer->blockSignals(false);
    };
    if (m_defaultSink)
        connect(m_defaultSink, &AudioSink::volumeChanged, this, setCtrlVolumeValue);

    connect(m_model, &VolumeModel::defaultSinkChanged, this, [ this, setCtrlVolumeValue ](AudioSink *sink) {
        if (m_defaultSink)
            disconnect(m_defaultSink);

        m_defaultSink = sink;
        if (sink) {
            setCtrlVolumeValue(sink->volume());
            connect(m_defaultSink, &AudioSink::volumeChanged, this, setCtrlVolumeValue);
        }
    });

    connect(m_sliderContainer->slider(), &QSlider::valueChanged, this, [ this ](int value) {
        AudioSink *sink = m_model->defaultSink();
        if (sink)
            sink->setVolume(value, true);
    });

    connect(m_model, &VolumeModel::muteChanged, this, [ this ] {
        m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, QIcon(leftIcon()));
    });

    connect(m_sliderContainer, &SliderContainer::iconClicked, this, [ this ](const SliderContainer::IconPosition icon) {
        switch (icon) {
        case SliderContainer::IconPosition::LeftIcon: {
            if (m_model->existActiveOutputDevice())
                m_model->setMute(!m_model->isMute());
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

void VolumeWidget::showEvent(QShowEvent *event)
{
    DBlurEffectWidget::showEvent(event);
    Q_EMIT visibleChanged(true);
}

void VolumeWidget::hideEvent(QHideEvent *event)
{
    DBlurEffectWidget::hideEvent(event);
    Q_EMIT visibleChanged(false);
}

const QString VolumeWidget::leftIcon()
{
    bool existActiveOutputDevice = m_model->existActiveOutputDevice();
    const bool mute = existActiveOutputDevice ? m_model->isMute() : true;
    if (mute)
        return QString(":/icons/resources/audio-volume-muted-dark");

    return QString(":/icons/resources/volume");
}

const QString VolumeWidget::rightIcon()
{
    return QString(":/icons/resources/broadcast");
}
